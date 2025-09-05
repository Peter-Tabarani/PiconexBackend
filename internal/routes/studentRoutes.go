package routes

import (
	"database/sql"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/handlers"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"

	"github.com/gorilla/mux"
)

func RegisterStudentRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/student", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.GetStudents(db, w, r)
		case "POST":
			handlers.CreateStudent(db, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})).Methods("GET", "POST", "OPTIONS")

	router.HandleFunc("/student/{id}", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.GetStudentByID(db, w, r)
		case "PUT":
			handlers.UpdateStudent(db, w, r)
		case "DELETE":
			handlers.DeleteStudent(db, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})).Methods("GET", "PUT", "DELETE", "OPTIONS")

	router.HandleFunc("/student/name/{name}", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetStudentsByName(db, w, r)
	})).Methods("GET", "OPTIONS")
}
