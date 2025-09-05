package routes

import (
	"database/sql"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/handlers"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"
	"github.com/gorilla/mux"
)

// RegisterPersonRoutes registers all person endpoints to the router
func RegisterPersonRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/person", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetPersons(db, w, r)
	})).Methods("GET", "OPTIONS")

	router.HandleFunc("/person/{id}", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetPersonByID(db, w, r)
	})).Methods("GET", "OPTIONS")
}
