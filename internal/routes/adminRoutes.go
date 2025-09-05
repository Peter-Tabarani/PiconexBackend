package routes

import (
	"database/sql"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/handlers"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"

	"github.com/gorilla/mux"
)

func RegisterAdminRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/admin", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.GetAdmins(db, w, r)
		case "POST":
			handlers.CreateAdmin(db, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})).Methods("GET", "POST", "OPTIONS")

	router.HandleFunc("/admin/{id}", utils.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.GetAdminByID(db, w, r)
		case "PUT":
			handlers.UpdateAdmin(db, w, r)
		case "DELETE":
			handlers.DeleteAdmin(db, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})).Methods("GET", "PUT", "DELETE", "OPTIONS")
}
