package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"
	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func GetPersonalDocumentations(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extracts optional query parameter from the request
	adminIDStr := r.URL.Query().Get("admin_id")

	// Base SQL query for retrieving personal documentation
	query := `
		SELECT
			pd.personal_documentation_id,
			pd.admin_id,
			a.activity_datetime,
			d.file_name,
			d.file_path,
			d.mime_type,
			d.size_bytes,
			d.uploaded_by
		FROM personal_documentation pd
		JOIN activity a ON pd.personal_documentation_id = a.activity_id
		JOIN documentation d ON pd.personal_documentation_id = d.documentation_id
	`

	args := []any{}

	// Optional filter by admin_id
	if adminIDStr != "" {
		// Converts the "admin_id" string to an integer
		adminID, err := strconv.Atoi(adminIDStr)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "Invalid admin ID")
			log.Println("Invalid ID parse error:", err)
			return
		}
		query += " WHERE pd.admin_id = ?"
		args = append(args, adminID)
	}

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query, args...)
	// Error message if QueryContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to obtain personal documentation")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	// Creates an empty slice to obtain results
	personalDocumentation := make([]models.PersonalDocumentation, 0)

	// Reads each row returned by the database
	for rows.Next() {
		var pd models.PersonalDocumentation
		// Parses the current data into fields of "pd" variable
		if err := rows.Scan(
			&pd.PersonalDocumentationID,
			&pd.AdminID,
			&pd.ActivityDateTime,
			&pd.FileName,
			&pd.FilePath,
			&pd.MimeType,
			&pd.SizeBytes,
			&pd.UploadedBy,
		); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to scan personal documentation")
			log.Println("Row scan error:", err)
			return
		}
		// Adds the obtained data to the slice
		personalDocumentation = append(personalDocumentation, pd)
	}

	// Checks for errors during iteration such as network interruptions and driver errors
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Writes JSON response & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, personalDocumentation)
}

func GetPersonalDocumentationByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr, ok := vars["personal_documentation_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing personal documentation ID")
		return
	}

	// Converts the "personal_documentation_id" string to an integer
	personalDocumentationID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid personal documentation ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// SQL query to select a single personal_documentation
	query := `
		SELECT
			pd.personal_documentation_id,
			pd.admin_id,
			a.activity_datetime,
			d.file_name,
			d.file_path,
			d.mime_type,
			d.size_bytes,
			d.uploaded_by
		FROM personal_documentation pd
		JOIN activity a ON pd.personal_documentation_id = a.activity_id
		JOIN documentation d ON pd.personal_documentation_id = d.documentation_id
		WHERE pd.personal_documentation_id = ?
	`

	// Empty variable for personal_documentation struct
	var pd models.PersonalDocumentation

	// Executes query
	err = db.QueryRowContext(r.Context(), query, personalDocumentationID).Scan(
		&pd.PersonalDocumentationID,
		&pd.AdminID,
		&pd.ActivityDateTime,
		&pd.FileName,
		&pd.FilePath,
		&pd.MimeType,
		&pd.SizeBytes,
		&pd.UploadedBy,
	)
	// Error message if no rows are found
	if err == sql.ErrNoRows {
		utils.WriteError(w, http.StatusNotFound, "Personal documentation not found")
		return
		// Error message if QueryRowContext or scan fails
	} else if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch personal documentation")
		log.Println("DB query error:", err)
		return
	}

	// Writes JSON response & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, pd)
}

func DownloadPersonalDocumentation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extracts personal_documentation_id from URL path parameters
	vars := mux.Vars(r)
	idStr := vars["personal_documentation_id"]

	// Converts "personal_documentation_id" string to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid personal_documentation_id")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// SQL query to retrieve the full documentation record
	query := `
		SELECT
			pd.personal_documentation_id,
			pd.admin_id,
			a.activity_datetime,
			d.file_name,
			d.file_path,
			d.mime_type,
			d.size_bytes,
			d.uploaded_by
		FROM personal_documentation pd
		JOIN activity a ON pd.personal_documentation_id = a.activity_id
		JOIN documentation d ON pd.personal_documentation_id = d.documentation_id
		WHERE pd.personal_documentation_id = ?
	`

	// Creates an empty struct to store result
	var pd models.PersonalDocumentation

	// Executes the SQL query
	err = db.QueryRowContext(r.Context(), query, id).Scan(
		&pd.PersonalDocumentationID,
		&pd.AdminID,
		&pd.ActivityDateTime,
		&pd.FileName,
		&pd.FilePath,
		&pd.MimeType,
		&pd.SizeBytes,
		&pd.UploadedBy,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteError(w, http.StatusNotFound, "File not found")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to obtain documentation info")
		log.Println("DB query error:", err)
		return
	}

	// Clean and resolve full file path
	fullPath := filepath.Clean(pd.FilePath)

	// Sets headers for file download
	w.Header().Set("Content-Type", pd.MimeType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", pd.FileName))

	// Streams the file to the HTTP response
	http.ServeFile(w, r, fullPath)
}

func CreatePersonalDocumentation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Empty variable for personal_documentation struct
	var pd models.PersonalDocumentation

	// Decodes JSON body from the request into "pd" variable
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents extra unexpected fields
	if err := decoder.Decode(&pd); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		log.Println("JSON decode error:", err)
		return
	}

	// Automatically set activity_datetime to now
	pd.ActivityDateTime = time.Now()

	// Validates required fields
	if pd.AdminID == 0 || pd.FileName == "" || pd.FilePath == "" || pd.MimeType == "" || pd.SizeBytes == 0 {
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
		pd.ActivityDateTime,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert activity")
		log.Println("DB insert error:", err)
		return
	}

	// Gets the last inserted activity_id
	lastID, err := res.LastInsertId()

	// Error message if LastInsertId fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get inserted activity ID")
		log.Println("LastInsertId error:", err)
		return
	}

	// Inserts into documentation table
	_, err = tx.ExecContext(r.Context(),
		`INSERT INTO documentation (
			documentation_id,
			file_name,
			file_path,
			mime_type,
			size_bytes,
			uploaded_by
		) VALUES (?, ?, ?, ?, ?, ?)`,
		lastID,
		pd.FileName,
		pd.FilePath,
		pd.MimeType,
		pd.SizeBytes,
		pd.UploadedBy,
	)
	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert documentation")
		log.Println("DB insert error:", err)
		return
	}

	// Inserts into personal_documentation table
	_, err = tx.ExecContext(r.Context(),
		"INSERT INTO personal_documentation (personal_documentation_id, admin_id) VALUES (?, ?)",
		lastID, pd.AdminID,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert personal documentation")
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
		"message":                   "Personal documentation created successfully",
		"personal_documentation_id": lastID,
	})
}

func UpdatePersonalDocumentation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr, ok := vars["personal_documentation_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing personal documentation ID")
		return
	}

	// Converts the "personal_documentation_id" string to an integer
	personalDocumentationID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid personal documentation ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// Empty variable for personal_documentation struct
	var pd models.PersonalDocumentation

	// Decodes JSON body from the request into "pd" variable
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents extra unexpected fields
	if err := decoder.Decode(&pd); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		log.Println("JSON decode error:", err)
		return
	}

	// Automatically set activity_datetime to now
	pd.ActivityDateTime = time.Now()

	// Validates required fields
	if pd.AdminID == 0 || pd.FileName == "" || pd.FilePath == "" || pd.MimeType == "" || pd.SizeBytes == 0 {
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
		pd.ActivityDateTime, personalDocumentationID,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update activity")
		log.Println("DB update error:", err)
		return
	}

	_, err = tx.ExecContext(r.Context(),
		`UPDATE documentation
		SET file_name=?, file_path=?, mime_type=?, size_bytes=?, uploaded_by=?
		WHERE documentation_id=?`,
		pd.FileName, pd.FilePath, pd.MimeType, pd.SizeBytes, pd.UploadedBy, personalDocumentationID,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update documentation")
		log.Println("DB update error:", err)
		return
	}

	// Executes written SQL to update the personal documentation data
	res, err := tx.ExecContext(r.Context(),
		"UPDATE personal_documentation SET admin_id=? WHERE personal_documentation_id=?",
		pd.AdminID, personalDocumentationID,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update personal documentation")
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
		utils.WriteError(w, http.StatusNotFound, "Personal documentation not found")
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
		"message": "Personal documentation updated successfully",
	})
}

func DeletePersonalDocumentation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr, ok := vars["personal_documentation_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing personal documentation ID")
		return
	}

	// Converts the "personal_documentation_id" string to an integer
	personalDocumentationID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid personal documentation ID")
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
	// Deletes from personal_documentation, documentation, and activity in one go
	query := `
		DELETE pd, d, a
		FROM personal_documentation pd
		JOIN documentation d ON d.documentation_id = pd.personal_documentation_id
		JOIN activity a ON a.activity_id = d.documentation_id
		WHERE pd.personal_documentation_id = ?;
	`

	// Executes written SQL to delete the documentation
	res, err := tx.ExecContext(r.Context(), query, personalDocumentationID)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete personal documentation")
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
		utils.WriteError(w, http.StatusNotFound, "No personal documentation found for this ID")
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to commit transaction")
		log.Println("Transaction commit error:", err)
		return
	}

	// Respond with success
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message":       "Personal documentation " + idStr + " deleted successfully",
		"rows_affected": rowsAffected / 3, // Each personal documentation involves 3 rows deleted
	})
}

func DeletePersonalDocumentationByAdminID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extract admin_id from route params
	vars := mux.Vars(r)
	adminIDStr, ok := vars["admin_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing admin ID")
		return
	}

	// Converts the "admin_id" string to an integer
	adminID, err := strconv.Atoi(adminIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid admin ID")
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
	// Deletes from personal_documentation, documentation, and activity in one go
	query := `
		DELETE pd, d, a
		FROM personal_documentation pd
		JOIN documentation d ON d.documentation_id = pd.personal_documentation_id
		JOIN activity a ON a.activity_id = d.documentation_id
		WHERE pd.admin_id = ?;
	`

	// Executes written SQL to delete the documentation
	res, err := tx.ExecContext(r.Context(), query, adminID)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete personal documentation")
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
		utils.WriteError(w, http.StatusNotFound, "No personal documentation found for this admin")
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
		"message":       "All personal documentation for admin " + adminIDStr + " deleted successfully",
		"rows_affected": rowsAffected / 3, // Each personal documentation involves 3 rows deleted
	})
}
