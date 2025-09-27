package routes

import (
	"database/sql"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/handlers"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"

	"github.com/gorilla/mux"
)

func RegisterPersonalDocumentationRoutes(router *mux.Router, db *sql.DB) {
	pdRouter := router.PathPrefix("/personal-documentation").Subrouter()
	pdRouter.Use(utils.WithCORS, utils.AuthMiddleware)

	pdRouter.Handle("",
		utils.RollMiddleware(map[string][]string{
			"GET":  {"admin"},
			"POST": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetPersonalDocumentations(db, w, r)
			case http.MethodPost:
				handlers.CreatePersonalDocumentation(db, w, r)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("GET", "POST", "OPTIONS")

	pdRouter.Handle("/{activity_id}",
		utils.RollMiddleware(map[string][]string{
			"GET":    {"admin"},
			"PUT":    {"admin"},
			"DELETE": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetPersonalDocumentationByID(db, w, r)
			case http.MethodPut:
				handlers.UpdatePersonalDocumentation(db, w, r)
			case http.MethodDelete:
				handlers.DeletePersonalDocumentation(db, w, r)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("GET", "PUT", "DELETE", "OPTIONS")
}
