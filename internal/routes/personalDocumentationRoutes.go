package routes

import (
	"database/sql"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/handlers"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"

	"github.com/gorilla/mux"
)

func RegisterPersonalDocumentationRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/personal-documentation", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.GetPersonalDocumentations(db, w, r)
		case "POST":
			handlers.CreatePersonalDocumentation(db, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})).Methods("GET", "POST", "OPTIONS")

	router.HandleFunc("/personal-documentation/{activity_id}", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.GetPersonalDocumentationByID(db, w, r)
		case "PUT":
			handlers.UpdatePersonalDocumentation(db, w, r)
		case "DELETE":
			handlers.DeletePersonalDocumentation(db, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})).Methods("GET", "PUT", "DELETE", "OPTIONS")
}
