package routes

import (
	"database/sql"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/handlers"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"

	"github.com/gorilla/mux"
)

func RegisterSpecificDocumentationRoutes(router *mux.Router, db *sql.DB) {
	sdRouter := router.PathPrefix("/specific-documentation").Subrouter()
	sdRouter.Use(utils.WithCORS, utils.AuthMiddleware)

	sdRouter.Handle(
		"",
		utils.RollMiddleware(map[string][]string{
			"GET":  {"admin"},
			"POST": {"student", "admin"},
		}, utils.ResourceCreateOwnershipMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetSpecificDocumentations(db, w, r)
			case http.MethodPost:
				handlers.CreateSpecificDocumentation(db, w, r)
			default:
				utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
		}))),
	).Methods("GET", "POST", "OPTIONS")

	sdRouter.Handle(
		"/{activity_id}",
		utils.RollMiddleware(map[string][]string{
			"GET":    {"student", "admin"},
			"PUT":    {"student", "admin"},
			"DELETE": {"student", "admin"},
		}, utils.ResourceOwnershipMiddleware(
			db,
			"specific_documentation",
			"activity_id",
			"id",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					handlers.GetSpecificDocumentationByID(db, w, r)
				case http.MethodPut:
					handlers.UpdateSpecificDocumentation(db, w, r)
				case http.MethodDelete:
					handlers.DeleteSpecificDocumentation(db, w, r)
				default:
					utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
				}
			}),
		)),
	).Methods("GET", "PUT", "DELETE", "OPTIONS")

	sdRouter.Handle(
		"/student/{id}",
		utils.RollMiddleware(
			map[string][]string{
				"GET": {"student", "admin"},
			},
			utils.OwnershipMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlers.GetSpecificDocumentationByStudentID(db, w, r)
			})),
		),
	).Methods("GET", "OPTIONS")
}
