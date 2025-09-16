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

func GetSpecificDocumentations(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not GET
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// All data being selected for this GET command
	query := `
		SELECT
			sd.activity_id, sd.id, sd.doc_type, a.date, a.time, d.file
		FROM specific_documentation sd
		JOIN activity a ON sd.activity_id = a.activity_id
		JOIN documentation d ON sd.activity_id = d.activity_id
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query)

	// Error message if QueryContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to obtain specific documentations")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	// Creates an empty slice to obtain results
	specificDocumentation := make([]models.SpecificDocumentation, 0)

	// Reads each row returned by the database
	for rows.Next() {
		var sd models.SpecificDocumentation
		// Parses the current data into fields of "sd" variable
		if err := rows.Scan(&sd.Activity_ID, &sd.ID, &sd.DocType, &sd.Date, &sd.Time, &sd.File); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to scan specific documentation")
			log.Println("Row scan error:", err)
			return
		}

		// Adds the obtained data to the slice
		specificDocumentation = append(specificDocumentation, sd)
	}

	// Checks for errors during iteration such as network interruptions and driver errors
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Writes JSON response & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, specificDocumentation)
}

func GetSpecificDocumentationByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not GET
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr, ok := vars["activity_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing activity ID")
		return
	}

	// Converts the "activity_id" string to an integer
	activityID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid activity ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// SQL query to select a single specific_documentation
	query := `
		SELECT sd.activity_id, sd.id, sd.doc_type, a.date, a.time, d.file
		FROM specific_documentation sd
		JOIN activity a ON sd.activity_id = a.activity_id
		JOIN documentation d ON sd.activity_id = d.activity_id
		WHERE sd.activity_id = ?
	`

	// Empty variable for specific_documentation struct
	var sd models.SpecificDocumentation

	// Executes query
	err = db.QueryRow(query, activityID).Scan(&sd.Activity_ID, &sd.ID, &sd.DocType, &sd.Date, &sd.Time, &sd.File)

	// Error message if no rows are found
	if err == sql.ErrNoRows {
		utils.WriteError(w, http.StatusNotFound, "Specific documentation not found")
		return
		// Error message if QueryRowContext or scan fails
	} else if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch specific documentation")
		log.Println("DB query error:", err)
		return
	}

	// Writes JSON response & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, sd)
}

func GetSpecificDocumentationByStudentID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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

	// SQL query to select all specific_documentation for a student
	query := `
		SELECT sd.activity_id, sd.id, sd.doc_type, a.date, a.time, d.file
		FROM specific_documentation sd
		JOIN activity a ON sd.activity_id = a.activity_id
		JOIN documentation d ON sd.activity_id = d.activity_id
		WHERE sd.id = ?
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query, studentID)

	// Error message if QueryContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to obtain specific documentations")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	// Creates an empty slice to obtain results
	specificDocumentation := make([]models.SpecificDocumentation, 0)

	// Reads each row returned by the database
	for rows.Next() {
		var sd models.SpecificDocumentation
		// Parses the current data into fields of "sd" variable
		if err := rows.Scan(&sd.Activity_ID, &sd.ID, &sd.DocType, &sd.Date, &sd.Time, &sd.File); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to scan specific documentation")
			log.Println("Row scan error:", err)
			return
		}
		// Adds the obtained data to the slice
		specificDocumentation = append(specificDocumentation, sd)
	}

	// Checks for errors during iteration such as network interruptions and driver errors
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Writes JSON response & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, specificDocumentation)
}

func CreateSpecificDocumentation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not POST
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Empty variable for specific_documentation struct
	var sd models.SpecificDocumentation

	// Decodes JSON body from the request into "sd" variable
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents extra unexpected fields
	if err := decoder.Decode(&sd); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		log.Println("JSON decode error:", err)
		return
	}

	// TECH DEBT: Validates required fields

	// Executes SQL to insert into activity table
	res, err := db.ExecContext(r.Context(),
		"INSERT INTO activity (date, time) VALUES (?, ?)",
		sd.Date, sd.Time,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert activity")
		log.Println("DB insert error:", err)
		return
	}

	// Gets the last inserted activity_id
	lastID, err := res.LastInsertId()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get inserted activity ID")
		log.Println("LastInsertId error:", err)
		return
	}

	// Inserts into documentation table
	_, err = db.ExecContext(r.Context(),
		"INSERT INTO documentation (activity_id, file) VALUES (?, ?)",
		lastID, sd.File,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert documentation")
		log.Println("DB insert error:", err)
		return
	}

	// Inserts into specific_documentation table
	_, err = db.ExecContext(r.Context(),
		"INSERT INTO specific_documentation (activity_id, id, doc_type) VALUES (?, ?, ?)",
		lastID, sd.ID, sd.DocType,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert specific documentation")
		log.Println("DB insert error:", err)
		return
	}

	// Writes JSON response & sends a HTTP 201 response code
	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message":     "Specific documentation created successfully",
		"activity_id": lastID,
	})
}

func DeleteSpecificDocumentation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not DELETE
	if r.Method != http.MethodDelete {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr, ok := vars["activity_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing activity ID")
		return
	}

	// Converts the "activity_id" string to an integer
	activityID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid activity ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// Executes SQL to delete from activity (will cascade to documentation and specific_documentation)
	res, err := db.ExecContext(r.Context(),
		"DELETE FROM activity WHERE activity_id=?",
		activityID,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete specific documentation")
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
		utils.WriteError(w, http.StatusNotFound, "Specific documentation not found")
		return
	}

	// Writes JSON response & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Specific documentation deleted successfully",
	})
}

func UpdateSpecificDocumentation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not PUT
	if r.Method != http.MethodPut {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr, ok := vars["activity_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing activity ID")
		return
	}

	// Converts the "activity_id" string to an integer
	activityID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid activity ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// Empty variable for specific_documentation struct
	var sd models.SpecificDocumentation

	// Decodes JSON body from the request into "sd" variable
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents extra unexpected fields
	if err := decoder.Decode(&sd); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		log.Println("JSON decode error:", err)
		return
	}

	// Executes written SQL to update the activity data
	_, err = db.ExecContext(r.Context(),
		"UPDATE activity SET date=?, time=? WHERE activity_id=?",
		sd.Date, sd.Time, activityID,
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update activity")
		log.Println("DB update error:", err)
		return
	}

	// Executes written SQL to update the documentation data
	_, err = db.ExecContext(r.Context(),
		"UPDATE documentation SET file=? WHERE activity_id=?",
		sd.File, activityID,
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update documentation")
		log.Println("DB update error:", err)
		return
	}

	// Executes written SQL to update the specific documentation data
	res, err := db.ExecContext(r.Context(),
		"UPDATE specific_documentation SET doc_type=?, id=? WHERE activity_id=?",
		sd.DocType, sd.ID, activityID,
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update specific documentation")
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
		utils.WriteError(w, http.StatusNotFound, "Specific documentation not found")
		return
	}

	// Writes JSON response & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Specific documentation updated successfully",
	})
}
