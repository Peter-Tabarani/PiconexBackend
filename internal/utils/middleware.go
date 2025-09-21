package utils

import (
	"context"
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
			if role == allowedRole {
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
			WriteError(w, http.StatusForbidden, "Forbidden")
			log.Printf("Ownership error: student %d tried to access student %d\n", userID, id)
			return
		}

		// Allowed → pass request through
		next.ServeHTTP(w, r)
	})
}
