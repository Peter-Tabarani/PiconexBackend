package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Local struct for login request body
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decode JSON request into "req" variable
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Variables for DB values
	var userID int
	var passwordHash, role string

	// Look up user by email in users + person tables
	err := db.QueryRow(`
		SELECT u.id, u.password_hash, u.role
		FROM users u
		JOIN person p ON p.id = u.id
		WHERE p.email = ?`, req.Email,
	).Scan(&userID, &passwordHash, &role)

	// Return unauthorized if not found
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Compare provided password with stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := utils.CreateJWT(userID, role)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	// Return token in JSON response
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func SignupHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Local struct for request
	type AdminSignupStudentRequest struct {
		ID       int    `json:"id"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decodes JSON body from the request into "req" variable
	var req AdminSignupStudentRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents extra unexpected fields
	if err := decoder.Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		log.Println("JSON decode error:", err)
		return
	}

	// Hashes password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to hash password")
		log.Println("Password hashing error:", err)
		return
	}

	// Adds hash and role to the users table
	_, err = db.ExecContext(r.Context(),
		`INSERT INTO users (id, password_hash, role) VALUES (?, ?, ?)`,
		req.ID, string(hashedPassword), "student",
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create student login")
		log.Println("DB insert error:", err)
		return
	}

	// Writes JSON response including the new ID & sends a HTTP 201 response code
	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message":   "Student signup completed successfully",
		"studentId": req.ID,
	})
}

func SignupStudentHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Empty variables for student struct
	type CreateStudentRequest struct {
		models.Student
		Password string `json:"password"`
	}

	// Decodes JSON body from the request into "req" variable
	var req CreateStudentRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents extra unexpected fields
	if err := decoder.Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		log.Println("JSON decode error:", err)
		return
	}

	// TECH DEBT: Validates required fields

	// Executes SQL to insert into person table
	res, err := db.ExecContext(r.Context(),
		`INSERT INTO person (
			first_name, preferred_name, middle_name, last_name, email,
			phone_number, pronouns, sex, gender, birthday,
			address, city, state, zip_code, country
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		req.FirstName, req.PreferredName, req.MiddleName, req.LastName,
		req.Email, req.PhoneNumber, req.Pronouns, req.Sex, req.Gender, req.Birthday,
		req.Address, req.City, req.State, req.ZipCode, req.Country,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert into person")
		log.Println("DB insert error:", err)
		return
	}

	// Gets the last inserted person ID
	lastID, err := res.LastInsertId()

	// Error message if LastInsertId fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get inserted person ID")
		log.Println("LastInsertId error:", err)
		return
	}

	// Executes SQL to insert into student table
	_, err = db.ExecContext(r.Context(),
		`INSERT INTO student (id, year, start_year, planned_grad_year, housing, dining)
		VALUES (?, ?, ?, ?, ?, ?)`,
		lastID, req.Year, req.StartYear, req.PlannedGradYear, req.Housing, req.Dining,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert into student")
		log.Println("DB insert error:", err)
		return
	}

	// Hashes password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to hash password")
		log.Println("Password hashing error:", err)
		return
	}

	// Adds hash and role to the users table
	_, err = db.ExecContext(r.Context(),
		`INSERT INTO users (id, password_hash, role) VALUES (?, ?, ?)`,
		lastID, string(hashedPassword), "student",
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create student login")
		log.Println("DB insert error:", err)
		return
	}

	// Writes JSON response including the new ID & sends a HTTP 201 response code
	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message":   "Student signup and creation completed successfully",
		"studentId": lastID,
	})
}
