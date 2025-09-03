package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"context"
	"time"

	"github.com/Peter-Tabarani/PiconexBackend/routes"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var db *sql.DB

func main() {
	// 1. Connect to MySQL
	dsn := "piconex:pjaplmTabs7!@tcp(178.156.189.138:3306)/piconexdb"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("❌ Error opening database:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("❌ Error connecting to database:", err)
	}

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

		switch r.Method {
		case "GET":
			routes.GetStudents(db, w, r)
		case "POST":
			routes.CreateStudent(db, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}).Methods("GET", "POST", "PUT", "OPTIONS")
	router.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}

		switch r.Method {
		case "GET":
			routes.GetAdmins(db, w, r)
		case "POST":
			routes.CreateAdmin(db, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}).Methods("GET", "POST", "OPTIONS")
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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		switch r.Method {
		case "GET":
			routes.GetPersonalDocumentations(db, w, r)
		case "POST":
			routes.CreatePersonalDocumentation(db, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}).Methods("GET", "POST", "OPTIONS")
	router.HandleFunc("/specific-documentation", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}

		switch r.Method {
		case "GET":
			routes.GetSpecificDocumentations(db, w, r)
		case "POST":
			routes.CreateSpecificDocumentation(db, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}).Methods("GET", "POST", "OPTIONS")
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

		switch r.Method {
		case "GET":
			routes.GetAccommodations(db, w, r)
		case "POST":
			routes.CreateAccommodation(db, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}).Methods("GET", "POST", "OPTIONS")
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
			w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}

		switch r.Method {
		case "GET":
			routes.GetStudentByID(db, w, r)
		case "PUT":
			routes.UpdateStudentByID(db, w, r)
		case "DELETE":
			routes.DeleteStudentByID(db, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}).Methods("GET", "PUT", "DELETE", "OPTIONS")
	router.HandleFunc("/admin/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}

		switch r.Method {
		case "GET":
			routes.GetAdminByID(db, w, r)
		case "PUT":
			routes.UpdateAdminByID(db, w, r)
		case "DELETE":
			routes.DeleteAdminByID(db, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}).Methods("GET", "PUT", "DELETE", "OPTIONS")
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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		switch r.Method {
		case "GET":
			routes.GetPersonalDocumentationByID(db, w, r)
		case "PUT":
			routes.UpdatePersonalDocumentationByID(db, w, r)
		case "DELETE":
			routes.DeletePersonalDocumentationByID(db, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}).Methods("GET", "PUT", "DELETE", "OPTIONS")
	router.HandleFunc("/specific-documentation/{activity_id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}

		switch r.Method {
		case "GET":
			routes.GetSpecificDocumentationByID(db, w, r)
		case "PUT":
			routes.UpdateSpecificDocumentationByID(db, w, r)
		case "DELETE":
			routes.DeleteSpecificDocumentationByID(db, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}).Methods("GET", "PUT", "DELETE", "OPTIONS")
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
	router.HandleFunc("/accommodation/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")
			w.WriteHeader(http.StatusOK)
			return
		}

		switch r.Method {
		case "GET":
			routes.GetAccommodationByID(db, w, r)
		case "PUT":
			routes.UpdateAccommodationByID(db, w, r)
		case "DELETE":
			routes.DeleteAccommodationByID(db, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}).Methods("GET", "PUT", "DELETE", "OPTIONS")
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
	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: corsRouter,
	}

	// Channel to listen for SIGTERM / SIGINT
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Run server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	log.Println("✅ Server started on :8080")

	// Wait for termination signal
	<-stop
	log.Println("🛑 Shutting down gracefully...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited properly")

}
