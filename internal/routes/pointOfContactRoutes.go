package routes

import (
	"database/sql"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/handlers"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"

	"github.com/gorilla/mux"
)

func RegisterPointOfContactRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/point-of-contact", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.GetPointsOfContact(db, w, r)
		case "POST":
			handlers.CreatePointOfContact(db, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})).Methods("GET", "POST", "OPTIONS")

	router.HandleFunc("/point-of-contact/{activity_id}", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.GetPointOfContactByID(db, w, r)
		case "PUT":
			handlers.UpdatePointOfContact(db, w, r)
		case "DELETE":
			handlers.DeletePointOfContact(db, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})).Methods("GET", "PUT", "DELETE", "OPTIONS")

	router.HandleFunc("/point-of-contact/admin/{id}/date/{date}", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetPointsOfContactByAdminIDAndDate(db, w, r)
	})).Methods("GET", "OPTIONS")

	router.HandleFunc("/future-point-of-contact/student/{student_id}/admin/{admin_id}", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetFuturePointsOfContactByStudentIDAndAdminID(db, w, r)
	})).Methods("GET", "OPTIONS")

}
