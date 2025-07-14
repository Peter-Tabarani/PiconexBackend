package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/routes"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var db *sql.DB

func main() {
	// 1. Connect to MySQL
	dsn := "peter:pjaplmTabs7!@tcp(127.0.0.1:3306)/piconexdb"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("❌ Error opening database:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("❌ Error connecting to database:", err)
	}
	fmt.Println("✅ Connected to MySQL successfully!")

	// 2. Build router
	router := mux.NewRouter()

	// Define your routes
	router.HandleFunc("/person", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetPersons(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/student", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetStudents(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetAdmins(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/activity", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetActivities(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/documentation", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetDocumentations(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/personal-documentation", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetPersonalDocumentations(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/specific-documentation", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetSpecificDocumentations(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/disability", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetDisabilities(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/accommodation", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetAccommodations(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/point-of-contact", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetPointOfContact(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/stu-dis", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetStuDis(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/stu-accom", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetStuAccom(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/pinned", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetPinned(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/pocadmin", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetPocAdmin(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/person/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetPersonByID(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/student/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetStudentByID(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/admin/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetAdminByID(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/activity/{activity_id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetActivityByID(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/documentation/{activity_id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetDocumentationByID(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/personal-documentation/{activity_id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetPersonalDocumentationByID(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/specific-documentation/{activity_id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetSpecificDocumentationByID(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/point-of-contact/{activity_id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetPointOfContactByID(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/disability/{disability_id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetDisabilityByID(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/accommodation/{accommodation_id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetAccommodationByID(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/pinned/admin/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetPinnedByAdminID(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/student/name/{name}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetStudentByName(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/activity/date/{date}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetActivitiesByDate(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/activity/student/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetActivitiesByStudentID(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/specific-documentation/student/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetSpecificDocumentationByStudentID(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/point-of-contact/admin/{id}/date/{date}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetPointOfContactByAdminIDAndDate(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/future-point-of-contact/student/{student_id}/admin/{admin_id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetFuturePointOfContactByStudentIDAndAdminID(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/disability/student/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetDisabilitiesByStudentID(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/accommodation/student/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}
		routes.GetAccommodationsByStudentID(db, w, r)
	}).Methods("GET", "OPTIONS")

	// 3. CORS middleware
	headersOk := handlers.AllowedHeaders([]string{
		"X-Requested-With", "Content-Type", "Authorization", "ngrok-skip-browser-warning",
	})
	originsOk := handlers.AllowedOrigins([]string{
		"*", // You can replace this with "http://localhost:5173" in production
	})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})

	// Apply CORS middleware directly on the router
	corsRouter := handlers.CORS(originsOk, headersOk, methodsOk)(router)

	// 4. Start server with CORS enabled
	http.ListenAndServe("0.0.0.0:8080", corsRouter)
}
