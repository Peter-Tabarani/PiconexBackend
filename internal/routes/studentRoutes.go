package routes

import (
	"database/sql"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/handlers"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"

	"github.com/gorilla/mux"
)

func RegisterStudentRoutes(router *mux.Router, db *sql.DB) {
	studentRouter := router.PathPrefix("/student").Subrouter()
	studentRouter.Use(utils.WithCORS, utils.AuthMiddleware)

	studentRouter.Handle("",
		utils.RollMiddleware(map[string][]string{
			"GET":  {"student", "admin"},
			"POST": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetStudents(db, w, r)
			case http.MethodPost:
				handlers.CreateStudent(db, w, r)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("GET", "POST", "OPTIONS")

	studentRouter.Handle("/{id}",
		utils.RollMiddleware(map[string][]string{
			"GET":    {"student", "admin"},
			"PUT":    {"admin"},
			"DELETE": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetStudentByID(db, w, r)
			case http.MethodPut:
				handlers.UpdateStudent(db, w, r)
			case http.MethodDelete:
				handlers.DeleteStudent(db, w, r)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("GET", "PUT", "DELETE", "OPTIONS")

	studentRouter.Handle("/name/{name}",
		utils.RollMiddleware(map[string][]string{
			"GET": {"student", "admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetStudentsByName(db, w, r)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("GET", "OPTIONS")
}
