package routes

import (
	"database/sql"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/handlers"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"

	"github.com/gorilla/mux"
)

func RegisterActivityRoutes(router *mux.Router, db *sql.DB) {
	activityRouter := router.PathPrefix("/activity").Subrouter()
	activityRouter.Use(utils.WithCORS, utils.AuthMiddleware)

	activityRouter.Handle("",
		utils.RollMiddleware(map[string][]string{
			"GET": {"student", "admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetActivities(db, w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("GET", "OPTIONS")

	activityRouter.Handle("/{activity_id}",
		utils.RollMiddleware(map[string][]string{
			"GET": {"student", "admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetActivityByID(db, w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("GET", "OPTIONS")

	activityRouter.Handle("/date/{date}",
		utils.RollMiddleware(map[string][]string{
			"GET": {"student", "admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetActivitiesByDate(db, w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("GET", "OPTIONS")
}
