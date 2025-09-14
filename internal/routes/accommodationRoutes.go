package routes

import (
	"database/sql"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/handlers"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"

	"github.com/gorilla/mux"
)

func RegisterAccommodationRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/accommodation", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.GetAccommodations(db, w, r)
		case "POST":
			handlers.CreateAccommodation(db, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})).Methods("GET", "POST", "OPTIONS")

	router.HandleFunc("/accommodation/{accommodation_id}", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.GetAccommodationByID(db, w, r)
		case "PUT":
			handlers.UpdateAccommodation(db, w, r)
		case "DELETE":
			handlers.DeleteAccommodation(db, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})).Methods("GET", "PUT", "DELETE", "OPTIONS")

	router.HandleFunc("/accommodation/student/{accommodation_id}", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetAccommodationsByStudentID(db, w, r)
	})).Methods("GET", "OPTIONS")

}
