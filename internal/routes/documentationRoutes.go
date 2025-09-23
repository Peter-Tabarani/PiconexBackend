package routes

import (
	"database/sql"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/handlers"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"
	"github.com/gorilla/mux"
)

func RegisterDocumentationRoutes(router *mux.Router, db *sql.DB) {
	documentationRouter := router.PathPrefix("/documentation").Subrouter()
	documentationRouter.Use(utils.WithCORS, utils.AuthMiddleware)

	documentationRouter.Handle("",
		utils.RollMiddleware(map[string][]string{
			"GET": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetDocumentations(db, w, r)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("GET", "OPTIONS")

	documentationRouter.Handle("/{activity_id}",
		utils.RollMiddleware(map[string][]string{
			"GET": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetDocumentationByID(db, w, r)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("GET", "OPTIONS")
}
