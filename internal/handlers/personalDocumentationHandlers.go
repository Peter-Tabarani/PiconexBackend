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

func GetPersonalDocumentations(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// All data being selected for this GET command
	query := `
		SELECT
			pd.personal_documentation_id, pd.admin_id, a.activity_datetime, d.file
		FROM personal_documentation pd
		JOIN activity a ON pd.personal_documentation_id = a.activity_id
		JOIN documentation d ON pd.personal_documentation_id = d.documentation_id
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query)

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
		// Parses the current data into fields of "a" variable
		if err := rows.Scan(&pd.PersonalDocumentationID, &pd.AdminID, &pd.ActivityDateTime, &pd.File); err != nil {
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
		SELECT pd.personal_documentation_id, pd.admin_id, a.activity_datetime, d.file
		FROM personal_documentation pd
		JOIN activity a ON pd.personal_documentation_id = a.activity_id
		JOIN documentation d ON pd.personal_documentation_id = d.documentation_id
		WHERE pd.personal_documentation_id = ?
	`

	// Empty variable for personal_documentation struct
	var pd models.PersonalDocumentation

	// Executes query
	err = db.QueryRowContext(r.Context(), query, personalDocumentationID).Scan(&pd.PersonalDocumentationID, &pd.AdminID, &pd.ActivityDateTime, &pd.File)

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
	if pd.AdminID == 0 || len(pd.File) == 0 {
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
		"INSERT INTO documentation (documentation_id, file) VALUES (?, ?)",
		lastID, pd.File,
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
	if pd.AdminID == 0 || len(pd.File) == 0 {
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

	// Executes written SQL to update the documentation data
	_, err = tx.ExecContext(r.Context(),
		"UPDATE documentation SET file=? WHERE documentation_id=?",
		pd.File, personalDocumentationID,
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

	// Start transaction
	tx, err := db.BeginTx(r.Context(), nil)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to begin transaction")
		log.Println("BeginTx error:", err)
		return
	}
	defer tx.Rollback()

	// Executes written SQL to delete the personal documentation
	res, err := tx.ExecContext(r.Context(),
		"DELETE FROM personal_documentation WHERE personal_documentation_id = ?", personalDocumentationID)

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
		utils.WriteError(w, http.StatusNotFound, "Personal documentation not found")
		return
	}

	// Executes written SQL to delete the documentation
	res, err = tx.ExecContext(r.Context(),
		"DELETE FROM documentation WHERE documentation_id = ?", personalDocumentationID)

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
		"DELETE FROM activity WHERE activity_id = ?", personalDocumentationID)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete activity")
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
		utils.WriteError(w, http.StatusNotFound, "Activity not found")
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to commit transaction")
		log.Println("Transaction commit error:", err)
		return
	}

	// Writes JSON response confirming deletion & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Personal documentation deleted successfully",
	})
}
