package routes

import (
	"database/sql"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/handlers"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"

	"github.com/gorilla/mux"
)

func RegisterSpecificDocumentationRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/specific-documentation", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.GetSpecificDocumentations(db, w, r)
		case "POST":
			handlers.CreateSpecificDocumentation(db, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})).Methods("GET", "POST", "OPTIONS")

	router.HandleFunc("/specific-documentation/{activity_id}", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.GetSpecificDocumentationByID(db, w, r)
		case "PUT":
			handlers.UpdateSpecificDocumentation(db, w, r)
		case "DELETE":
			handlers.DeleteSpecificDocumentation(db, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})).Methods("GET", "PUT", "DELETE", "OPTIONS")

	router.HandleFunc("/specific-documentation/student/{id}", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetSpecificDocumentationByStudentID(db, w, r)
	})).Methods("GET", "OPTIONS")
}
