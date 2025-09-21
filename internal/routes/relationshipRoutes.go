package routes

import (
	"database/sql"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/handlers"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"

	"github.com/gorilla/mux"
)

func RegisterRelationshipRoutes(router *mux.Router, db *sql.DB) {
	pinnedRouter := router.PathPrefix("/pinned").Subrouter()
	pinnedRouter.Use(utils.WithCORS, utils.AuthMiddleware)

	pinnedRouter.Handle("",
		utils.RollMiddleware(map[string][]string{
			"GET":  {"admin"},
			"POST": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetPinned(db, w, r)
			case http.MethodPost:
				handlers.CreatePinned(db, w, r)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("GET", "POST", "OPTIONS")

	pinnedRouter.Handle("/{admin_id}/{student_id}",
		utils.RollMiddleware(map[string][]string{
			"DELETE": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodDelete:
				handlers.DeletePinned(db, w, r)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("DELETE", "OPTIONS")

	pinnedRouter.Handle("/admin/{id}",
		utils.RollMiddleware(map[string][]string{
			"GET": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetPinnedByAdminID(db, w, r)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("GET", "OPTIONS")

	stuAccomRouter := router.PathPrefix("/stu-accom").Subrouter()
	stuAccomRouter.Use(utils.WithCORS, utils.AuthMiddleware)

	stuAccomRouter.Handle("",
		utils.RollMiddleware(map[string][]string{
			"GET":  {"admin"},
			"POST": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetStuAccom(db, w, r)
			case http.MethodPost:
				handlers.CreateStuAccom(db, w, r)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("GET", "POST", "OPTIONS")

	stuAccomRouter.Handle("/{id}/{accommodation_id}",
		utils.RollMiddleware(map[string][]string{
			"DELETE": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodDelete:
				handlers.DeleteStuAccom(db, w, r)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("DELETE", "OPTIONS")

	stuDisRouter := router.PathPrefix("/stu-dis").Subrouter()
	stuDisRouter.Use(utils.WithCORS, utils.AuthMiddleware)

	stuDisRouter.Handle("",
		utils.RollMiddleware(map[string][]string{
			"GET":  {"admin"},
			"POST": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetStuDis(db, w, r)
			case http.MethodPost:
				handlers.CreateStuDis(db, w, r)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("GET", "POST", "OPTIONS")

	stuDisRouter.Handle("/{id}/{disability_id}",
		utils.RollMiddleware(map[string][]string{
			"DELETE": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodDelete:
				handlers.DeleteStuDis(db, w, r)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("DELETE", "OPTIONS")

	pocAdminRouter := router.PathPrefix("/poc-admin").Subrouter()
	pocAdminRouter.Use(utils.WithCORS, utils.AuthMiddleware)

	pocAdminRouter.Handle("",
		utils.RollMiddleware(map[string][]string{
			"GET":  {"admin"},
			"POST": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetPocAdmin(db, w, r)
			case http.MethodPost:
				handlers.CreatePocAdmin(db, w, r)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("GET", "POST", "OPTIONS")

	pocAdminRouter.Handle("/{activity_id}/{id}",
		utils.RollMiddleware(map[string][]string{
			"DELETE": {"admin"},
		}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodDelete:
				handlers.DeletePocAdmin(db, w, r)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods("DELETE", "OPTIONS")
}
