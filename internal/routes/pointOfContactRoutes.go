package routes

import (
	"database/sql"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/handlers"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"

	"github.com/gorilla/mux"
)

func RegisterPointOfContactRoutes(router *mux.Router, db *sql.DB) {
	pocRouter := router.PathPrefix("/point-of-contact").Subrouter()
	pocRouter.Use(utils.WithCORS, utils.AuthMiddleware)

	pocRouter.Handle("",
		utils.RollMiddleware(map[string][]string{
			"GET":  {"admin"},
			"POST": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetPointsOfContact(db, w, r)
			case http.MethodPost:
				handlers.CreatePointOfContact(db, w, r)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("GET", "POST", "OPTIONS")

	pocRouter.Handle("/{activity_id}",
		utils.RollMiddleware(map[string][]string{
			"GET":    {"admin"},
			"PUT":    {"admin"},
			"DELETE": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetPointOfContactByID(db, w, r)
			case http.MethodPut:
				handlers.UpdatePointOfContact(db, w, r)
			case http.MethodDelete:
				handlers.DeletePointOfContact(db, w, r)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("GET", "PUT", "DELETE", "OPTIONS")

	pocRouter.Handle("/admin/{id}/date/{date}",
		utils.RollMiddleware(map[string][]string{
			"GET": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetPointsOfContactByAdminIDAndDate(db, w, r)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("GET", "OPTIONS")

	pocRouter.Handle("/future/student/{student_id}/admin/{admin_id}",
		utils.RollMiddleware(map[string][]string{
			"GET": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetFuturePointsOfContactByStudentIDAndAdminID(db, w, r)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("GET", "OPTIONS")
}
