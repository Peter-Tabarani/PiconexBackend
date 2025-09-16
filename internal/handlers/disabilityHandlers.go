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
)

func GetDisabilities(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not GET
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// All data being selected for this GET command
	query := `SELECT disability_id, name, description FROM disability`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to obtain disabilities")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	// Creates an empty slice to obtain results
	disabilities := make([]models.Disability, 0)

	// Reads each row returned by the database
	for rows.Next() {
		var d models.Disability
		// Parses the current data into fields of "d" variable
		if err := rows.Scan(&d.Disability_ID, &d.Name, &d.Description); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse disabilities")
			log.Println("Row scan error:", err)
			return
		}

		// Adds the obtained data to the slice
		disabilities = append(disabilities, d)
	}

	// Checks for errors during iteration such as network interruptions and driver errors
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Writes the slice as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, disabilities)
}

func GetDisabilityByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not GET
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr, ok := vars["disability_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing disability ID")
		return
	}

	// Converts the "disability_id" string to an integer
	disabilityID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid disability ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// All data being selected for this GET command
	query := `SELECT disability_id, name, description FROM disability WHERE disability_id = ?`

	// Empty variable for disability struct
	var d models.Disability

	// Executes written SQL and retrieves only one row
	err = db.QueryRowContext(r.Context(), query, disabilityID).Scan(&d.Disability_ID, &d.Name, &d.Description)

	// Error message if no rows are found
	if err == sql.ErrNoRows {
		utils.WriteError(w, http.StatusNotFound, "Disability not found")
		return
		// Error message if QueryRowContext or scan fails
	} else if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch disability")
		log.Println("DB query error:", err)
		return
	}

	// Writes the struct as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, d)
}

func GetDisabilitiesByStudentID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not GET
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing student ID")
		return
	}

	// Converts the "id" string to an integer
	studentID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid student ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// All data being selected for this GET command
	query := `
		SELECT d.disability_id, d.name, d.description
		FROM stu_dis sd
		JOIN disability d ON sd.disability_id = d.disability_id
		WHERE sd.id = ?
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query, studentID)

	// Error message if QueryContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to obtain disabilities for student")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	// Creates an empty slice to obtain results
	disabilities := make([]models.Disability, 0)

	// Reads each row returned by the database
	for rows.Next() {
		var d models.Disability
		// Parses the current data into fields of "d" variable
		if err := rows.Scan(&d.Disability_ID, &d.Name, &d.Description); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse disabilities")
			log.Println("Row scan error:", err)
			return
		}

		// Adds the obtained data to the slice
		disabilities = append(disabilities, d)
	}

	// Checks for errors during iteration such as network interruptions and driver errors
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Writes the slice as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, disabilities)
}

func CreateDisability(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not POST
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Empty variable for disability struct
	var d models.Disability

	// Decodes JSON body from the request into "d" variable
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents extra unexpected fields
	if err := decoder.Decode(&d); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		log.Println("JSON decode error:", err)
		return
	}

	// Validates required fields
	if d.Name == "" {
		utils.WriteError(w, http.StatusBadRequest, "Missing required fields: name or description")
		return
	}

	// Executes written SQL to insert a new disability
	res, err := db.ExecContext(r.Context(),
		"INSERT INTO disability (name, description) VALUES (?, ?)",
		d.Name, d.Description,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert disability")
		log.Println("DB insert error:", err)
		return
	}

	// Gets the ID of the newly inserted disability
	lastID, err := res.LastInsertId()

	// Error message if LastInsertId fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get last insert ID")
		log.Println("LastInsertId error:", err)
		return
	}

	// Writes JSON response including the new ID & sends a HTTP 201 response code
	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message":       "Disability created successfully",
		"disability_id": lastID,
	})
}

func DeleteDisability(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not DELETE
	if r.Method != http.MethodDelete {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr := vars["disability_id"]

	// Converts the "disability_id" string to an integer
	disabilityID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid disability ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// Executes written SQL to delete any student references first
	_, err = db.ExecContext(r.Context(), "DELETE FROM stu_dis WHERE disability_id = ?", disabilityID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete student references")
		log.Println("DB delete error:", err)
		return
	}

	// Executes written SQL to delete the disability
	res, err := db.ExecContext(r.Context(), "DELETE FROM disability WHERE disability_id = ?", disabilityID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete disability")
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
		utils.WriteError(w, http.StatusNotFound, "Disability not found")
		return
	}

	// Writes JSON response confirming deletion & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Disability deleted successfully",
	})
}

func UpdateDisability(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not PUT
	if r.Method != http.MethodPut {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr := vars["disability_id"]

	// Converts the "disability_id" string to an integer
	disabilityID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid disability ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// Empty variable for disability struct
	var d models.Disability

	// Decodes JSON body from the request into "d" variable
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents extra unexpected fields
	if err := decoder.Decode(&d); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		log.Println("JSON decode error:", err)
		return
	}

	// Validates required fields
	if d.Name == "" || d.Description == "" {
		utils.WriteError(w, http.StatusBadRequest, "Missing required fields: name or description")
		return
	}

	// Executes written SQL to update the disability
	res, err := db.ExecContext(r.Context(),
		"UPDATE disability SET name = ?, description = ? WHERE disability_id = ?",
		d.Name, d.Description, disabilityID,
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update disability")
		log.Println("DB update error:", err)
		return
	}

	// Gets the number of rows affected by the update
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to check update result")
		log.Println("RowsAffected error:", err)
		return
	}

	// Error message if no rows were updated
	if rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, "Disability not found")
		return
	}

	// Writes JSON response confirming update & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Disability updated successfully",
	})
}
