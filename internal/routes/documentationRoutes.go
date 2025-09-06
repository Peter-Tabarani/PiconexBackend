package routes

import (
	"database/sql"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/handlers"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"
	"github.com/gorilla/mux"
)

// RegisterPersonRoutes registers all person endpoints to the router
func RegisterDocumentationRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/documentation", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetDocumentations(db, w, r)
	})).Methods("GET", "OPTIONS")

	router.HandleFunc("/documentation/{id}", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetDocumentationByID(db, w, r)
	})).Methods("GET", "OPTIONS")
}
