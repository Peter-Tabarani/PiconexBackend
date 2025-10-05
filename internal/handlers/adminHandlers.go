package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"
	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func GetAdmins(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// All data being selected for this GET command
	query := `
		SELECT
			p.person_id, p.first_name, p.preferred_name, p.middle_name, p.last_name,
			p.email, p.phone_number, p.pronouns, p.sex, p.gender,
			p.birthday, p.address, p.city, p.state, p.zip_code, p.country, a.title
		FROM admin a
		JOIN person p ON a.admin_id = p.person_id
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query)

	// Error message if QueryContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to obtain admins")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	// Creates an empty slice to obtain results
	admins := make([]models.Admin, 0)

	// Reads each row returned by the database
	for rows.Next() {
		var a models.Admin
		// Parses the current data into fields of "a" variable
		if err := rows.Scan(
			&a.AdminID, &a.FirstName, &a.PreferredName, &a.MiddleName, &a.LastName,
			&a.Email, &a.PhoneNumber, &a.Pronouns, &a.Sex, &a.Gender,
			&a.Birthday, &a.Address, &a.City, &a.State, &a.ZipCode, &a.Country,
			&a.Title,
		); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse admins")
			log.Println("Row scan error:", err)
			return
		}
		// Adds the obtained data to the slice
		admins = append(admins, a)
	}

	// Checks for errors during iteration such as network interruptions and driver errors
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Writes the slice as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, admins)
}

func GetAdminByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr, ok := vars["admin_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing admin ID")
		return
	}

	// Converts the "admin_id" string to an integer
	adminID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid admin ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// All data being selected for this GET command
	query := `
		SELECT
			a.admin_id, p.first_name, p.preferred_name, p.middle_name, p.last_name,
			p.email, p.phone_number, p.pronouns, p.sex, p.gender,
			p.birthday, p.address, p.city, p.state, p.zip_code, p.country, a.title
		FROM admin a
		JOIN person p ON a.admin_id = p.person_id
		WHERE a.admin_id = ?
	`

	// Empty variable for admin struct
	var a models.Admin

	// Executes written SQL and retrieves only one row
	err = db.QueryRowContext(r.Context(), query, adminID).Scan(
		&a.AdminID, &a.FirstName, &a.PreferredName, &a.MiddleName, &a.LastName,
		&a.Email, &a.PhoneNumber, &a.Pronouns, &a.Sex, &a.Gender,
		&a.Birthday, &a.Address, &a.City, &a.State, &a.ZipCode, &a.Country,
		&a.Title,
	)

	// Error message if no rows are found
	if err == sql.ErrNoRows {
		utils.WriteError(w, http.StatusNotFound, "Admin not found")
		return
		// Error message if QueryRowContext or scan fails
	} else if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch admin")
		log.Println("DB query error:", err)
		return
	}

	// Writes the struct as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, a)
}

func CreateAdmin(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Decodes JSON body from the request into "a" variable
	type CreateAdminRequest struct {
		models.Admin
		Password string `json:"password"`
	}

	// Decodes JSON body from the request into "a" variable
	var a CreateAdminRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents extra unexpected fields
	if err := decoder.Decode(&a); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		log.Println("JSON decode error:", err)
		return
	}

	// Validates required fields
	if a.FirstName == "" || a.LastName == "" || a.Email == "" || a.PhoneNumber == "" ||
		a.Sex == "" || a.Birthday == "" || a.Address == "" || a.City == "" ||
		a.Country == "" || a.Title == "" || a.Password == "" {
		utils.WriteError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Start transaction
	tx, err := db.BeginTx(r.Context(), nil)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to begin transaction")
		log.Println("BeginTx error:", err)
		return
	}
	defer tx.Rollback()

	// Executes SQL to insert into person table
	res, err := tx.ExecContext(r.Context(),
		`INSERT INTO person (
		first_name, preferred_name, middle_name, last_name,
		email, phone_number, pronouns, sex, gender,
		birthday, address, city, state, zip_code, country
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		a.FirstName, a.PreferredName, a.MiddleName, a.LastName,
		a.Email, a.PhoneNumber, a.Pronouns, a.Sex, a.Gender,
		a.Birthday, a.Address, a.City, a.State, a.ZipCode, a.Country,
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

	// Executes SQL to insert into admin table
	_, err = tx.ExecContext(r.Context(),
		"INSERT INTO admin (admin_id, title) VALUES (?, ?)",
		lastID, a.Title,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert admin title")
		log.Println("DB insert error:", err)
		return
	}

	// Hashes password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to hash password")
		log.Println("Password hashing error:", err)
		return
	}

	// Adds hash and role to the users table
	_, err = tx.ExecContext(r.Context(),
		`INSERT INTO users (id, password_hash, role) VALUES (?, ?, ?)`,
		lastID, string(hashedPassword), "admin",
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create admin login")
		log.Println("DB insert error:", err)
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to commit transaction")
		log.Println("Transaction commit error:", err)
		return
	}

	// Writes JSON response including the new ID & sends a HTTP 201 response code
	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Admin created successfully",
		"adminId": lastID,
	})
}

func UpdateAdmin(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr, ok := vars["admin_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing admin ID")
		return
	}

	// Converts the "admin_id" string to an integer
	adminID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid admin ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// Empty variable for admin struct
	var a models.Admin

	// Decodes JSON body from the request into "a" variable
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents extra unexpected fields
	if err := decoder.Decode(&a); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		log.Println("JSON decode error:", err)
		return
	}

	// Validates required fields
	if a.FirstName == "" || a.LastName == "" || a.Email == "" || a.PhoneNumber == "" ||
		a.Sex == "" || a.Birthday == "" || a.Address == "" || a.City == "" ||
		a.Country == "" || a.Title == "" {
		utils.WriteError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Start transaction
	tx, err := db.BeginTx(r.Context(), nil)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to begin transaction")
		log.Println("BeginTx error:", err)
		return
	}
	defer tx.Rollback()

	// Executes written SQL to update the person data
	_, err = tx.ExecContext(r.Context(),
		`UPDATE person SET
			first_name=?, preferred_name=?, middle_name=?, last_name=?,
			email=?, phone_number=?, pronouns=?, sex=?, gender=?,
			birthday=?, address=?, city=?, state=?, zip_code=?, country=?
		WHERE person_id=?`,
		a.FirstName, a.PreferredName, a.MiddleName, a.LastName,
		a.Email, a.PhoneNumber, a.Pronouns, a.Sex, a.Gender,
		a.Birthday, a.Address, a.City, a.State, a.ZipCode, a.Country,
		adminID,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update admin")
		log.Println("DB update error:", err)
		return
	}

	// Executes written SQL to update the admin title
	res, err := tx.ExecContext(r.Context(),
		"UPDATE admin SET title=? WHERE admin_id=?",
		a.Title, adminID,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update admin title")
		log.Println("DB update error:", err)
		return
	}

	// Gets the number of rows affected by the update
	rowsAffected, err := res.RowsAffected()

	// Error message if RowsAffected fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get rows affected")
		log.Println("RowsAffected error:", err)
		return
	}

	// Error message if no rows were updated
	if rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, "Admin not found")
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to commit transaction")
		log.Println("Transaction commit error:", err)
		return
	}

	// Writes JSON response confirming update & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Admin updated successfully",
	})
}

func DeleteAdmin(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr, ok := vars["admin_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing admin ID")
		return
	}

	// Converts the "admin_id" string to an integer
	adminID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid admin ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// Start transaction
	tx, err := db.BeginTx(r.Context(), nil)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to begin transaction")
		log.Println("BeginTx error:", err)
		return
	}
	defer tx.Rollback()

	// Executes SQL to delete from admin
	res, err := tx.ExecContext(r.Context(),
		"DELETE FROM admin WHERE admin_id = ?",
		adminID,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete admin")
		log.Println("DB delete error:", err)
		return
	}

	// Gets the number of rows affected
	rowsAffected, err := res.RowsAffected()

	// Error message if RowsAffected fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get rows affected")
		log.Println("RowsAffected error:", err)
		return
	}

	// Error message if no rows were deleted
	if rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, "Admin not found")
		return
	}

	// Executes SQL to delete from person
	res, err = tx.ExecContext(r.Context(), "DELETE FROM person WHERE person_id = ?", adminID)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete person")
		log.Println("DB delete person error:", err)
		return
	}

	// Gets the number of rows affected for person
	rowsAffected, err = res.RowsAffected()

	// Error message if RowsAffected fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get rows affected for person")
		log.Println("RowsAffected person error:", err)
		return
	}

	// Error message if no rows were deleted
	if rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, "Person not found")
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to commit transaction")
		log.Println("Transaction commit error:", err)
		return
	}

	// Writes JSON response confirming deletion
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Admin deleted successfully",
	})
}
