package routes

import (
	"database/sql"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/handlers"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"

	"github.com/gorilla/mux"
)

func RegisterAdminRoutes(router *mux.Router, db *sql.DB) {
	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.Use(utils.WithCORS, utils.AuthMiddleware)

	adminRouter.Handle("",
		utils.RollMiddleware(map[string][]string{
			"GET":  {"admin"},
			"POST": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetAdmins(db, w, r)
			case http.MethodPost:
				handlers.CreateAdmin(db, w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("GET", "POST", "OPTIONS")

	adminRouter.Handle("/{id}",
		utils.RollMiddleware(map[string][]string{
			"GET":    {"admin"},
			"PUT":    {"admin"},
			"DELETE": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetAdminByID(db, w, r)
			case http.MethodPut:
				handlers.UpdateAdmin(db, w, r)
			case http.MethodDelete:
				handlers.DeleteAdmin(db, w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("GET", "PUT", "DELETE", "OPTIONS")
}
