package routes

import (
	"database/sql"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/handlers"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"

	"github.com/gorilla/mux"
)

func RegisterPointOfContactRoutes(router *mux.Router, db *sql.DB) {
	pocRouter := router.PathPrefix("/point-of-contact").Subrouter()
	pocRouter.Use(utils.WithCORS, utils.AuthMiddleware)

	pocRouter.Handle(
		"",
		utils.RollMiddleware(map[string][]string{
			"GET":  {"admin"},
			"POST": {"student", "admin"},
		}, utils.ResourceCreateOwnershipMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetPointsOfContact(db, w, r)
			case http.MethodPost:
				handlers.CreatePointOfContact(db, w, r)
			default:
				utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
		}))),
	).Methods("GET", "POST", "OPTIONS")

	pocRouter.Handle(
		"/past",
		utils.RollMiddleware(map[string][]string{
			"GET": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetPastPointsOfContact(db, w, r)
			default:
				utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
		})),
	).Methods("GET", "OPTIONS")

	pocRouter.Handle(
		"/future",
		utils.RollMiddleware(map[string][]string{
			"GET": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetFuturePointsOfContact(db, w, r)
			default:
				utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
		})),
	).Methods("GET", "OPTIONS")

	pocRouter.Handle(
		"/{point_of_contact_id}",
		utils.RollMiddleware(map[string][]string{
			"GET":    {"student", "admin"},
			"PUT":    {"student", "admin"},
			"DELETE": {"student", "admin"},
		}, utils.ResourceOwnershipMiddleware(
			db,
			"point_of_contact",
			"point_of_contact_id",
			"student_id",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					handlers.GetPointOfContactByID(db, w, r)
				case http.MethodPut:
					handlers.UpdatePointOfContact(db, w, r)
				case http.MethodDelete:
					handlers.DeletePointOfContact(db, w, r)
				default:
					utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
				}
			}),
		)),
	).Methods("GET", "PUT", "DELETE", "OPTIONS")
}
