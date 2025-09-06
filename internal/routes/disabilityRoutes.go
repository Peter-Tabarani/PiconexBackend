package routes

import (
	"database/sql"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/handlers"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"

	"github.com/gorilla/mux"
)

func RegisterDisabilityRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/disability", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.GetDisabilities(db, w, r)
		case "POST":
			handlers.CreateDisability(db, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})).Methods("GET", "POST", "OPTIONS")

	router.HandleFunc("/disability/{id}", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.GetDisabilityByID(db, w, r)
		case "PUT":
			handlers.UpdateDisability(db, w, r)
		case "DELETE":
			handlers.DeleteDisability(db, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})).Methods("GET", "PUT", "DELETE", "OPTIONS")

	router.HandleFunc("/disability/student/{id}", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetDisabilitiesByStudentID(db, w, r)
	})).Methods("GET", "OPTIONS")

}
