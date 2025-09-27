package routes

import (
	"database/sql"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/handlers"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"
	"github.com/gorilla/mux"
)

func RegisterPersonRoutes(router *mux.Router, db *sql.DB) {
	personRouter := router.PathPrefix("/person").Subrouter()
	personRouter.Use(utils.WithCORS, utils.AuthMiddleware)

	personRouter.Handle("",
		utils.RollMiddleware(map[string][]string{
			"GET": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetPersons(db, w, r)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("GET", "OPTIONS")

	personRouter.Handle("/{id}",
		utils.RollMiddleware(map[string][]string{
			"GET": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetPersonByID(db, w, r)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("GET", "OPTIONS")
}
