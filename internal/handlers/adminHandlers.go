package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func GetAdmins(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not GET
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// All data being selected for this GET command
	query := `
		SELECT
			p.id, p.first_name, p.preferred_name, p.middle_name, p.last_name,
			p.email, p.phone_number, p.pronouns, p.sex, p.gender,
			p.birthday, p.address, p.city, p.state, p.zip_code, p.country, a.title
		FROM admin a
		JOIN person p ON a.id = p.id
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
			&a.ID, &a.FirstName, &a.PreferredName, &a.MiddleName, &a.LastName,
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
	// Error message for any request that is not GET
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing admin ID")
		return
	}

	// Converts the "id" string to an integer
	adminID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid admin ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// All data being selected for this GET command
	query := `
		SELECT
			p.id, p.first_name, p.preferred_name, p.middle_name, p.last_name,
			p.email, p.phone_number, p.pronouns, p.sex, p.gender,
			p.birthday, p.address, p.city, p.state, p.zip_code, p.country, a.title
		FROM admin a
		JOIN person p ON a.id = p.id
		WHERE a.id = ?
	`

	// Empty variable for admin struct
	var a models.Admin

	// Executes written SQL and retrieves only one row
	err = db.QueryRowContext(r.Context(), query, adminID).Scan(
		&a.ID, &a.FirstName, &a.PreferredName, &a.MiddleName, &a.LastName,
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
	// Error message for any request that is not POST
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
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
	if a.FirstName == "" || a.LastName == "" || a.Email == "" || a.Title == "" {
		utils.WriteError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Executes written SQL to insert a new person
	res, err := db.ExecContext(r.Context(),
		`INSERT INTO person (
			first_name, preferred_name, middle_name, last_name,
			email, phone_number, pronouns, sex, gender,
			birthday, address, city, state, zip_code, country
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		a.FirstName, a.PreferredName, a.MiddleName, a.LastName,
		a.Email, a.PhoneNumber, a.Pronouns, a.Sex, a.Gender,
		a.Birthday, a.Address, a.City, a.State, a.ZipCode, a.Country,
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert admin")
		log.Println("DB insert error:", err)
		return
	}

	// Gets the ID of the newly inserted person
	lastID, err := res.LastInsertId()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get last insert ID")
		log.Println("LastInsertId error:", err)
		return
	}

	// Executes written SQL to insert into admin table
	_, err = db.ExecContext(r.Context(),
		"INSERT INTO admin (id, title) VALUES (?, ?)",
		lastID, a.Title,
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert admin title")
		log.Println("DB insert error:", err)
		return
	}

	// Writes JSON response including the new ID & sends a HTTP 201 response code
	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Admin created successfully",
		"adminId": lastID,
	})
}

func DeleteAdmin(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not DELETE
	if r.Method != http.MethodDelete {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Converts the "id" string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid admin ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// Executes written SQL to delete the admin
	res, err := db.ExecContext(r.Context(), "DELETE FROM person WHERE id = ?", id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete admin")
		log.Println("DB delete error:", err)
		return
	}

	// Gets the number of rows affected by the delete
	rowsAffected, err := res.RowsAffected()
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

	// Writes JSON response confirming deletion & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Admin deleted successfully",
	})
}

func UpdateAdmin(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not PUT
	if r.Method != http.MethodPut {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Converts the "id" string to an integer
	id, err := strconv.Atoi(idStr)
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
	if a.FirstName == "" || a.LastName == "" || a.Email == "" || a.Title == "" {
		utils.WriteError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Executes written SQL to update the person data
	_, err = db.ExecContext(r.Context(),
		`UPDATE person SET
			first_name=?, preferred_name=?, middle_name=?, last_name=?,
			email=?, phone_number=?, pronouns=?, sex=?, gender=?,
			birthday=?, address=?, city=?, state=?, zip_code=?, country=?
		WHERE id=?`,
		a.FirstName, a.PreferredName, a.MiddleName, a.LastName,
		a.Email, a.PhoneNumber, a.Pronouns, a.Sex, a.Gender,
		a.Birthday, a.Address, a.City, a.State, a.ZipCode, a.Country,
		id,
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update admin")
		log.Println("DB update error:", err)
		return
	}

	// Executes written SQL to update the admin title
	res, err := db.ExecContext(r.Context(),
		"UPDATE admin SET title=? WHERE id=?",
		a.Title, id,
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update admin title")
		log.Println("DB update error:", err)
		return
	}

	// Gets the number of rows affected by the update
	rowsAffected, err := res.RowsAffected()
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

	// Writes JSON response confirming update & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Admin updated successfully",
	})
}
