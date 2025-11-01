package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
			sd.specific_documentation_id,
			sd.student_id,
			sd.doc_type,
			a.activity_datetime,
			d.file_name,
			d.file_path,
			d.mime_type,
			d.size_bytes,
			d.uploaded_by
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
		if err := rows.Scan(
			&sd.SpecificDocumentationID,
			&sd.StudentID,
			&sd.DocType,
			&sd.ActivityDateTime,
			&sd.FileName,
			&sd.FilePath,
			&sd.MimeType,
			&sd.SizeBytes,
			&sd.UploadedBy,
		); err != nil {
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
		SELECT
			sd.specific_documentation_id,
			sd.student_id,
			sd.doc_type,
			a.activity_datetime,
			d.file_name,
			d.file_path,
			d.mime_type,
			d.size_bytes,
			d.uploaded_by
		FROM specific_documentation sd
		JOIN activity a ON sd.specific_documentation_id = a.activity_id
		JOIN documentation d ON sd.specific_documentation_id = d.documentation_id
		WHERE sd.specific_documentation_id = ?
	`

	// Empty variable for specific_documentation struct
	var sd models.SpecificDocumentation

	// Executes query
	err = db.QueryRowContext(r.Context(), query, specificDocumentationID).Scan(
		&sd.SpecificDocumentationID,
		&sd.StudentID,
		&sd.DocType,
		&sd.ActivityDateTime,
		&sd.FileName,
		&sd.FilePath,
		&sd.MimeType,
		&sd.SizeBytes,
		&sd.UploadedBy,
	)
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
	// Parses multipart form data from the request with a maximum upload size of 20MB
	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to parse form data")
		log.Println("Form parse error:", err)
		return
	}

	// Extracts "student_id" and "doc_type" fields from the multipart form
	studentIDStr := r.FormValue("student_id")
	docType := r.FormValue("doc_type")

	// Validates required form fields
	if studentIDStr == "" || docType == "" {
		utils.WriteError(w, http.StatusBadRequest, "Missing student_id or doc_type")
		return
	}

	// Converts "student_id" string to an integer
	studentID, err := strconv.Atoi(studentIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid student ID")
		log.Println("Invalid student ID parse error:", err)
		return
	}

	// Retrieves the uploaded file from the form
	file, header, err := r.FormFile("file")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Missing file in request")
		log.Println("Form file error:", err)
		return
	}
	defer file.Close()

	// Defines file storage directory and constructs a unique filename
	dstDir := "/home/piconex/database/files/specific"
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to ensure specific folder")
		log.Println("MkdirAll error:", err)
		return
	}
	fullPath := filepath.Join(dstDir, header.Filename)

	// Creates a new file at the destination path
	dst, err := os.Create(fullPath)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create file on server")
		log.Println("File create error:", err)
		return
	}
	defer dst.Close()

	// Copies the uploaded file content into the newly created file
	sizeBytes, err := io.Copy(dst, file)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to save uploaded file")
		log.Println("File write error:", err)
		return
	}

	// Detects the file's MIME type from the uploaded header
	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Begins a new database transaction
	tx, err := db.BeginTx(r.Context(), nil)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to begin transaction")
		log.Println("BeginTx error:", err)
		return
	}
	defer tx.Rollback()

	// Inserts a new activity record with the current timestamp
	now := time.Now()
	res, err := tx.ExecContext(r.Context(), "INSERT INTO activity (activity_datetime) VALUES (?)", now)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert activity")
		log.Println("Insert activity error:", err)
		return
	}

	// Retrieves the automatically generated activity_id from the database
	activityID, err := res.LastInsertId()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to retrieve inserted activity ID")
		log.Println("LastInsertId error:", err)
		return
	}

	// Inserts a new record into the documentation table with file metadata
	_, err = tx.ExecContext(r.Context(),
		`INSERT INTO documentation (documentation_id, file_name, file_path, mime_type, size_bytes, uploaded_by)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		activityID, header.Filename, fullPath, mimeType, sizeBytes, 5, // uploaded_by temporarily set to 5
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert documentation metadata")
		log.Println("Insert documentation error:", err)
		return
	}

	// Inserts a new record into the specific_documentation table linking the student and doc_type
	_, err = tx.ExecContext(r.Context(),
		`INSERT INTO specific_documentation (specific_documentation_id, student_id, doc_type)
		 VALUES (?, ?, ?)`,
		activityID, studentID, docType,
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert specific documentation entry")
		log.Println("Insert specific documentation error:", err)
		return
	}

	// Commits the transaction to finalize the database changes
	if err := tx.Commit(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to commit transaction")
		log.Println("Transaction commit error:", err)
		return
	}

	// Writes JSON response & sends a HTTP 201 response code
	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Specific documentation uploaded successfully",
		"id":      activityID,
		"path":    fullPath,
		"size":    sizeBytes,
	})
}

func DownloadSpecificDocumentation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extracts the specific_documentation_id from URL path parameters
	vars := mux.Vars(r)
	idStr := vars["specific_documentation_id"]

	// Converts "specific_documentation_id" string to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid specific_documentation_id")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// SQL query to retrieve the full documentation record
	query := `
		SELECT
			sd.specific_documentation_id,
			sd.student_id,
			a.activity_datetime,
			d.file_path,
			d.file_name,
			d.mime_type,
			d.size_bytes,
			d.uploaded_by,
			sd.doc_type
		FROM specific_documentation sd
		JOIN activity a ON sd.specific_documentation_id = a.activity_id
		JOIN documentation d ON sd.specific_documentation_id = d.documentation_id
		WHERE sd.specific_documentation_id = ?
	`

	// Creates an empty struct to store result
	var sd models.SpecificDocumentation

	// Executes the SQL query
	err = db.QueryRowContext(r.Context(), query, id).Scan(
		&sd.SpecificDocumentationID,
		&sd.StudentID,
		&sd.ActivityDateTime,
		&sd.FilePath,
		&sd.FileName,
		&sd.MimeType,
		&sd.SizeBytes,
		&sd.UploadedBy,
		&sd.DocType,
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
	fullPath := filepath.Clean(sd.FilePath)

	// Sets appropriate headers for file download
	w.Header().Set("Content-Type", sd.MimeType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", sd.FileName))

	// Streams the file to the HTTP response
	http.ServeFile(w, r, fullPath)
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
	if sd.StudentID == 0 || sd.DocType == "" || sd.FileName == "" || sd.FilePath == "" || sd.MimeType == "" || sd.SizeBytes == 0 {
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
		`UPDATE documentation
		SET file_name=?, file_path=?, mime_type=?, size_bytes=?, uploaded_by=?
		WHERE documentation_id=?`,
		sd.FileName, sd.FilePath, sd.MimeType, sd.SizeBytes, sd.UploadedBy, specificDocumentationID,
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

	// Retrieve file path before deleting from DB
	var filePath string
	err = db.QueryRowContext(r.Context(),
		`SELECT d.file_path
		 FROM documentation d
		 JOIN specific_documentation sd ON d.documentation_id = sd.specific_documentation_id
		 WHERE sd.specific_documentation_id = ?`,
		specificDocumentationID,
	).Scan(&filePath)

	// Handles missing or invalid file path case
	if err == sql.ErrNoRows {
		utils.WriteError(w, http.StatusNotFound, "No file found for this documentation ID")
		return
	} else if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to retrieve file path")
		log.Println("File path retrieval error:", err)
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
		WHERE sd.specific_documentation_id = ?;
	`

	// Executes written SQL to delete the specific documentation
	res, err := tx.ExecContext(r.Context(), query, specificDocumentationID)

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
		utils.WriteError(w, http.StatusNotFound, "No specific documentation found for this ID")
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to commit transaction")
		log.Println("Transaction commit error:", err)
		return
	}

	// Delete the physical file (after DB commit)
	if filePath != "" {
		if err := os.Remove(filePath); err != nil {
			// Non-fatal: log the issue but still return success
			log.Println("Failed to delete file from disk:", filePath, "Error:", err)
		}
	}

	// Respond with success JSON
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message":       "Specific documentation " + idStr + " deleted successfully",
		"file_deleted":  filePath,
		"rows_affected": rowsAffected / 3,
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

	// Retrieve all file paths before deleting from DB
	rows, err := db.QueryContext(r.Context(), `
		SELECT d.file_path
		FROM documentation d
		JOIN specific_documentation sd ON d.documentation_id = sd.specific_documentation_id
		WHERE sd.student_id = ?;
	`, studentID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to retrieve file paths")
		log.Println("File path retrieval error:", err)
		return
	}
	defer rows.Close()

	var filePaths []string
	for rows.Next() {
		var fp string
		if err := rows.Scan(&fp); err == nil && fp != "" {
			filePaths = append(filePaths, fp)
		}
	}
	_ = filePaths
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Error scanning file paths")
		log.Println("Rows scan error:", err)
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

	// Delete physical files (after DB commit)
	for _, path := range filePaths {
		if err := os.Remove(path); err != nil {
			log.Println("Failed to delete file:", path, "Error:", err)
		}
	}

	// Respond with success
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message":       "All specific documentation for student " + studentIDStr + " deleted successfully",
		"rows_affected": rowsAffected / 3,
		"files_deleted": len(filePaths),
	})
}
