package routes

import (
	"database/sql"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/handlers"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"

	"github.com/gorilla/mux"
)

func RegisterAccommodationRoutes(router *mux.Router, db *sql.DB) {
	accommodationRouter := router.PathPrefix("/accommodation").Subrouter()
	accommodationRouter.Use(utils.WithCORS, utils.AuthMiddleware)

	accommodationRouter.Handle("",
		utils.RollMiddleware(map[string][]string{
			"GET":  {"student", "admin"},
			"POST": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetAccommodations(db, w, r)
			case http.MethodPost:
				handlers.CreateAccommodation(db, w, r)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("GET", "POST", "OPTIONS")

	accommodationRouter.Handle("/{accommodation_id}",
		utils.RollMiddleware(map[string][]string{
			"GET":    {"student", "admin"},
			"PUT":    {"admin"},
			"DELETE": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetAccommodationByID(db, w, r)
			case http.MethodPut:
				handlers.UpdateAccommodation(db, w, r)
			case http.MethodDelete:
				handlers.DeleteAccommodation(db, w, r)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("GET", "PUT", "DELETE", "OPTIONS")

	accommodationRouter.Handle("/student/{id}",
		utils.RollMiddleware(map[string][]string{
			"GET": {"student", "admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlers.GetAccommodationsByStudentID(db, w, r)
		})),
	).Methods("GET", "OPTIONS")
}
