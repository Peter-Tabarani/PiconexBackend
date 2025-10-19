package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"
	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func GetSpecificDocumentations(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extracts optional query parameter from the request
	studentIDStr := r.URL.Query().Get("student_id")

	// Base SQL query for retrieving specific documentation
	query := `
		SELECT
			sd.specific_documentation_id, sd.student_id, sd.doc_type, a.activity_datetime, d.file
		FROM specific_documentation sd
		JOIN activity a ON sd.specific_documentation_id = a.activity_id
		JOIN documentation d ON sd.specific_documentation_id = d.documentation_id
	`

	args := []any{}

	// Optional filter by student_id
	if studentIDStr != "" {
		// Converts the "student_id" string to an integer
		studentID, err := strconv.Atoi(studentIDStr)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "Invalid student ID")
			log.Println("Invalid ID parse error:", err)
			return
		}
		query += " WHERE sd.student_id = ?"
		args = append(args, studentID)
	}

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query, args...)
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
		if err := rows.Scan(&sd.SpecificDocumentationID, &sd.StudentID, &sd.DocType, &sd.ActivityDateTime, &sd.File); err != nil {
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
	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr, ok := vars["specific_documentation_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing specific documentation ID")
		return
	}

	// Converts the "specific_documentation_id" string to an integer
	specificDocumentationID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid specific documentation ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// SQL query to select a single specific_documentation
	query := `
		SELECT sd.specific_documentation_id, sd.student_id, sd.doc_type, a.activity_datetime, d.file
		FROM specific_documentation sd
		JOIN activity a ON sd.specific_documentation_id = a.activity_id
		JOIN documentation d ON sd.specific_documentation_id = d.documentation_id
		WHERE sd.specific_documentation_id = ?
	`

	// Empty variable for specific_documentation struct
	var sd models.SpecificDocumentation

	// Executes query
	err = db.QueryRowContext(r.Context(), query, specificDocumentationID).Scan(&sd.SpecificDocumentationID, &sd.StudentID, &sd.DocType, &sd.ActivityDateTime, &sd.File)

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

func CreateSpecificDocumentation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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

	// Automatically set activity_datetime to now
	sd.ActivityDateTime = time.Now()

	// Validates required fields
	if sd.StudentID == 0 || sd.DocType == "" || len(sd.File) == 0 {
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

	// Executes SQL to insert into activity table
	res, err := tx.ExecContext(r.Context(),
		"INSERT INTO activity (activity_datetime) VALUES (?)",
		sd.ActivityDateTime,
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
	_, err = tx.ExecContext(r.Context(),
		"INSERT INTO documentation (documentation_id, file) VALUES (?, ?)",
		lastID, sd.File,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert documentation")
		log.Println("DB insert error:", err)
		return
	}

	// Inserts into specific_documentation table
	_, err = tx.ExecContext(r.Context(),
		"INSERT INTO specific_documentation (specific_documentation_id, student_id, doc_type) VALUES (?, ?, ?)",
		lastID, sd.StudentID, sd.DocType,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert specific documentation")
		log.Println("DB insert error:", err)
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to commit transaction")
		log.Println("Transaction commit error:", err)
		return
	}

	// Writes JSON response & sends a HTTP 201 response code
	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message":                   "Specific documentation created successfully",
		"specific_documentation_id": lastID,
	})
}

func UpdateSpecificDocumentation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr, ok := vars["specific_documentation_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing specific documentation ID")
		return
	}

	// Converts the "specific_documentation_id" string to an integer
	specificDocumentationID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid specific documentation ID")
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

	// Automatically set activity_datetime to now
	sd.ActivityDateTime = time.Now()

	// Validates required fields
	if sd.StudentID == 0 || sd.DocType == "" || len(sd.File) == 0 {
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

	// Executes written SQL to update the activity data
	_, err = tx.ExecContext(r.Context(),
		"UPDATE activity SET activity_datetime=? WHERE activity_id=?",
		sd.ActivityDateTime, specificDocumentationID,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update activity")
		log.Println("DB update error:", err)
		return
	}

	// Executes written SQL to update the documentation data
	_, err = tx.ExecContext(r.Context(),
		"UPDATE documentation SET file=? WHERE documentation_id=?",
		sd.File, specificDocumentationID,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update documentation")
		log.Println("DB update error:", err)
		return
	}

	// Executes written SQL to update the specific documentation data
	res, err := tx.ExecContext(r.Context(),
		"UPDATE specific_documentation SET doc_type=?, student_id=? WHERE specific_documentation_id=?",
		sd.DocType, sd.StudentID, specificDocumentationID,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update specific documentation")
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
		utils.WriteError(w, http.StatusNotFound, "Specific documentation not found")
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to commit transaction")
		log.Println("Transaction commit error:", err)
		return
	}

	// Writes JSON response & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Specific documentation updated successfully",
	})
}

func DeleteSpecificDocumentation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr, ok := vars["specific_documentation_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing specific documentation ID")
		return
	}

	// Converts the "specific_documentation_id" string to an integer
	specificDocumentationID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid specific documentation ID")
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

	// Executes written SQL to delete the specific documentation
	res, err := tx.ExecContext(r.Context(),
		"DELETE FROM specific_documentation WHERE specific_documentation_id = ?", specificDocumentationID)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete specific documentation")
		log.Println("DB delete error:", err)
		return
	}

	// Gets the number of rows affected by the delete
	rowsAffected, err := res.RowsAffected()

	// Error message if RowsAffected fails
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

	// Executes written SQL to delete the documentation
	res, err = tx.ExecContext(r.Context(),
		"DELETE FROM documentation WHERE documentation_id = ?", specificDocumentationID)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete documentation")
		log.Println("DB delete error:", err)
		return
	}

	// Gets the number of rows affected by the delete
	rowsAffected, err = res.RowsAffected()

	// Error message if RowsAffected fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get rows affected")
		log.Println("RowsAffected error:", err)
		return
	}

	// Error message if no rows were deleted
	if rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, "Documentation not found")
		return
	}

	// Executes written SQL to delete the activity
	res, err = tx.ExecContext(r.Context(),
		"DELETE FROM activity WHERE activity_id = ?", specificDocumentationID)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete activity")
		log.Println("DB delete activity error:", err)
		return
	}

	// Gets the number of rows affected by the delete
	rowsAffected, err = res.RowsAffected()

	// Error message if RowsAffected fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get rows affected")
		log.Println("RowsAffected error:", err)
		return
	}

	// Error message if no rows were deleted
	if rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, "Activity not found")
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to commit transaction")
		log.Println("Transaction commit error:", err)
		return
	}

	// Writes JSON response & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Specific documentation deleted successfully",
	})
}

func DeleteSpecificDocumentationByStudentID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extract student_id from route params
	vars := mux.Vars(r)
	studentIDStr, ok := vars["student_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing student ID")
		return
	}

	// Converts the "student_id" string to an integer
	studentID, err := strconv.Atoi(studentIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid student ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// Begin a transaction (not strictly required for single multi-table DELETE, but safer)
	tx, err := db.BeginTx(r.Context(), nil)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to begin transaction")
		log.Println("BeginTx error:", err)
		return
	}
	defer tx.Rollback()

	// Multi-table delete query:
	// Deletes from specific_documentation, documentation, and activity in one go
	query := `
		DELETE sd, d, a
		FROM specific_documentation sd
		JOIN documentation d ON d.documentation_id = sd.specific_documentation_id
		JOIN activity a ON a.activity_id = d.documentation_id
		WHERE sd.student_id = ?;
	`

	// Executes written SQL to delete the documentation
	res, err := tx.ExecContext(r.Context(), query, studentID)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete specific documentation")
		log.Println("Delete query error:", err)
		return
	}

	// Gets the number of rows affected by the delete
	rowsAffected, err := res.RowsAffected()

	// Error message if RowsAffected fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get rows affected")
		log.Println("RowsAffected error:", err)
		return
	}

	// Error message if no rows were deleted
	if rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, "No specific documentation found for this student")
		return
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to commit transaction")
		log.Println("Transaction commit error:", err)
		return
	}

	// Respond with success
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message":       "All specific documentation for student " + studentIDStr + " deleted successfully",
		"rows_affected": rowsAffected / 3, // Each specific documentation involves 3 rows deleted
	})
}
