package routes

import (
	"database/sql"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/handlers"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"

	"github.com/gorilla/mux"
)

func RegisterRelationshipRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/pinned", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.GetPinned(db, w, r)
		case "POST":
			handlers.CreatePinned(db, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})).Methods("GET", "POST", "OPTIONS")

	router.HandleFunc("/pinned/{admin_id}/{student_id}", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		handlers.DeletePinned(db, w, r)
	})).Methods("GET", "OPTIONS")

	router.HandleFunc("/pinned/admin/{id}", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetPinnedByAdminID(db, w, r)
	})).Methods("GET", "OPTIONS")
}
