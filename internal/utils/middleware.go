package utils

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type contextKey string

// Context keys for storing authenticated user data
const (
	UserIDKey contextKey = "userID"
	RoleKey   contextKey = "role"
)

const SuperKey = "superkey"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract and validate "Authorization: Bearer <token>" header
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			WriteError(w, http.StatusUnauthorized, "Missing token")
			log.Println("Auth error: missing token")
			return
		}
		tokenString := authHeader[7:]

		// 🔑 Superkey bypass
		if tokenString == SuperKey {
			ctx := context.WithValue(r.Context(), UserIDKey, -1) // -1 means "system"
			ctx = context.WithValue(ctx, RoleKey, "superadmin")  // special role
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Parse JWT claims
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtSecret, nil
		})
		if err != nil || !token.Valid {
			WriteError(w, http.StatusUnauthorized, "Invalid token")
			log.Println("Auth error: invalid token: ", err)
			return
		}

		// Store user ID and role in context for downstream use
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, RoleKey, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RollMiddleware(methodRoles map[string][]string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract role from context
		role, ok := r.Context().Value(RoleKey).(string)
		if !ok || role == "" {
			WriteError(w, http.StatusUnauthorized, "Unauthorized")
			log.Println("Role middleware error: missing role in context")
			return
		}

		// Check if role is allowed for this HTTP method
		allowedRoles := methodRoles[r.Method]
		for _, allowedRole := range allowedRoles {
			if role == "superadmin" || role == allowedRole {
				next.ServeHTTP(w, r)
				return
			}
		}

		// Role not permitted
		WriteError(w, http.StatusForbidden, "Forbidden: insufficient role")
		log.Printf("Role middleware error: role %q not allowed for %s\n", role, r.Method)
	})
}

func OwnershipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract role from context
		role, ok := r.Context().Value(RoleKey).(string)
		if !ok {
			WriteError(w, http.StatusUnauthorized, "Unauthorized")
			log.Println("Ownership error: missing role in context")
			return
		}

		// Extract user ID from context
		userID, ok := r.Context().Value(UserIDKey).(int)
		if !ok {
			WriteError(w, http.StatusUnauthorized, "Unauthorized")
			log.Println("Ownership error: missing user ID in context")
			return
		}

		// 🔑 Admins and superadmins always allowed
		if role == "admin" || role == "superadmin" {
			next.ServeHTTP(w, r)
			return
		}

		// Extract requested resource ID from route variables
		vars := mux.Vars(r)
		idStr := vars["id"]

		id, err := strconv.Atoi(idStr)
		if err != nil {
			WriteError(w, http.StatusBadRequest, "Invalid student ID")
			log.Println("Invalid ID parse error:", err)
			return
		}

		// Students can only access their own ID
		if role == "student" && userID != id {
			WriteError(w, http.StatusForbidden, "Forbidden: not owner")
			log.Printf("Ownership error: student %d tried to access student %d\n", userID, id)
			return
		}

		// Allowed → pass request through
		next.ServeHTTP(w, r)
	})
}

func ResourceOwnershipMiddleware(db *sql.DB, table, idColumn, studentColumn string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, _ := r.Context().Value(RoleKey).(string)
		userID, _ := r.Context().Value(UserIDKey).(int)

		// Admins always allowed
		if role == "admin" || role == "superadmin" {
			next.ServeHTTP(w, r)
			return
		}

		if role == "student" {
			vars := mux.Vars(r)
			idStr := vars[idColumn]
			id, err := strconv.Atoi(idStr)
			if err != nil {
				WriteError(w, http.StatusBadRequest, "Invalid ID")
				log.Println("Ownership error: invalid ID parse:", err)
				return
			}

			// Check ownership in DB
			var ownerID int
			query := fmt.Sprintf("SELECT %s FROM %s WHERE %s = ?", studentColumn, table, idColumn)
			err = db.QueryRowContext(r.Context(), query, id).Scan(&ownerID)
			if err != nil {
				if err == sql.ErrNoRows {
					WriteError(w, http.StatusNotFound, "Resource not found")
					log.Println("DB query error:", err)
				} else {
					WriteError(w, http.StatusInternalServerError, "Failed to verify ownership")
					log.Println("Ownership DB error:", err)
				}
				return
			}

			if ownerID != userID {
				WriteError(w, http.StatusForbidden, "Forbidden: not owner")
				log.Printf("Ownership error: student %d tried to access %s=%d in %s\n", userID, idColumn, id, table)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func ResourceCreateOwnershipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, _ := r.Context().Value(RoleKey).(string)
		userID, _ := r.Context().Value(UserIDKey).(int)

		// Admins always allowed
		if role == "admin" || role == "superadmin" {
			next.ServeHTTP(w, r)
			return
		}

		// Only enforce for students
		if role == "student" {
			// Decode a copy of the JSON body into a map for validation
			var payload map[string]interface{}
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				WriteError(w, http.StatusBadRequest, "Invalid request body")
				log.Println("Ownership create error: failed to read body:", err)
				return
			}

			// Reset body so the next decoder in handler still works
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			if err := json.Unmarshal(bodyBytes, &payload); err != nil {
				WriteError(w, http.StatusBadRequest, "Invalid JSON body")
				log.Println("Ownership create error: JSON unmarshal failed:", err)
				return
			}

			// Check for student_id field
			if sid, ok := payload["id"].(float64); ok {
				if int(sid) != userID {
					WriteError(w, http.StatusForbidden, "You can only create records for yourself")
					log.Printf("Ownership create error: student %d tried to create record for student %d\n", userID, int(sid))
					return
				}
			} else {
				// Optionally enforce that student_id must be present
				// Or automatically set it to userID if your schema allows
				WriteError(w, http.StatusBadRequest, "Missing student ID")
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
