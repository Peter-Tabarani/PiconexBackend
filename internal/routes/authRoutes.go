package routes

import (
	"database/sql"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/handlers"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"
	"github.com/gorilla/mux"
)

func RegisterAuthRoutes(router *mux.Router, db *sql.DB) {
	publicAuth := router.PathPrefix("/").Subrouter()
	publicAuth.Use(utils.WithCORS)

	publicAuth.HandleFunc("/signup/student", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.SignupHandler(db, w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}).Methods("POST", "OPTIONS")

	publicAuth.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.LoginHandler(db, w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}).Methods("POST", "OPTIONS")

	protectedAuth := router.PathPrefix("/").Subrouter()
	protectedAuth.Use(utils.WithCORS, utils.AuthMiddleware)

	protectedAuth.Handle("/signup",
		utils.RollMiddleware(map[string][]string{
			"POST": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				handlers.AdminSignupStudentHandler(db, w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("POST", "OPTIONS")
}
