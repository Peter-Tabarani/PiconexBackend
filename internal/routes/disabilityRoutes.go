package routes

import (
	"database/sql"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/handlers"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"

	"github.com/gorilla/mux"
)

func RegisterDisabilityRoutes(router *mux.Router, db *sql.DB) {
	disabilityRouter := router.PathPrefix("/disability").Subrouter()
	disabilityRouter.Use(utils.WithCORS, utils.AuthMiddleware)

	disabilityRouter.Handle(
		"",
		utils.RollMiddleware(map[string][]string{
			"GET":  {"student", "admin"},
			"POST": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetDisabilities(db, w, r)
			case http.MethodPost:
				handlers.CreateDisability(db, w, r)
			default:
				utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
		})),
	).Methods("GET", "POST", "OPTIONS")

	disabilityRouter.Handle(
		"/{disability_id}",
		utils.RollMiddleware(map[string][]string{
			"GET":    {"student", "admin"},
			"PUT":    {"admin"},
			"DELETE": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetDisabilityByID(db, w, r)
			case http.MethodPut:
				handlers.UpdateDisability(db, w, r)
			case http.MethodDelete:
				handlers.DeleteDisability(db, w, r)
			default:
				utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
		})),
	).Methods("GET", "PUT", "DELETE", "OPTIONS")

	disabilityRouter.Handle(
		"/student/{student_id}",
		utils.RollMiddleware(
			map[string][]string{
				"GET": {"student", "admin"},
			},
			utils.OwnershipMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlers.GetDisabilitiesByStudentID(db, w, r)
			})),
		),
	).Methods("GET", "OPTIONS")
}
