package routes

import (
	"database/sql"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/handlers"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"
	"github.com/gorilla/mux"
)

func RegisterActivityRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/activity", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetActivities(db, w, r)
	})).Methods("GET", "OPTIONS")

	router.HandleFunc("/activity/{activity_id}", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetActivityByID(db, w, r)
	})).Methods("GET", "OPTIONS")

	router.HandleFunc("/activity/date/{date}", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetActivitiesByDate(db, w, r)
	})).Methods("GET", "OPTIONS")

	router.HandleFunc("/activity/student/{id}", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetActivitiesByStudentID(db, w, r)
	})).Methods("GET", "OPTIONS")
}
