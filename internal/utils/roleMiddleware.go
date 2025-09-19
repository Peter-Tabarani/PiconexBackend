package utils

import (
	"net/http"
)

func RollMiddleware(methodRoles map[string][]string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(RoleKey).(string)
		if !ok || role == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		allowedRoles := methodRoles[r.Method]
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				next.ServeHTTP(w, r)
				return
			}
		}

		http.Error(w, "Forbidden: insufficient role", http.StatusForbidden)
	})
}
