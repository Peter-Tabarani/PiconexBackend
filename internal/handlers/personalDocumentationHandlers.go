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

func GetPersonalDocumentations(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not GET
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// All data being selected for this GET command
	query := `
		SELECT
			pd.activity_id, pd.id, a.date, a.time, d.file
		FROM personal_documentation pd
		JOIN activity a ON pd.activity_id = a.activity_id
		JOIN documentation d ON pd.activity_id = d.activity_id
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
		if err := rows.Scan(&pd.Activity_ID, &pd.ID, &pd.Date, &pd.Time, &pd.File); err != nil {
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
	// Error message for any request that is not GET
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr, ok := vars["activity_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing admin ID")
		return
	}

	// Converts the "activity_id" string to an integer
	activityID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid activity ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// SQL query to select a single personal_documentation
	query := `
		SELECT pd.activity_id, pd.id, a.date, a.time, d.file
		FROM personal_documentation pd
		JOIN activity a ON pd.activity_id = a.activity_id
		JOIN documentation d ON pd.activity_id = d.activity_id
		WHERE pd.activity_id = ?
	`

	// Empty variable for personal_documentation struct
	var pd models.PersonalDocumentation

	// Executes query
	err = db.QueryRow(query, activityID).Scan(&pd.Activity_ID, &pd.ID, &pd.Date, &pd.Time, &pd.File)

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
	// Error message for any request that is not POST
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
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

	// TECH DEBT: Validate required fields

	// Executes SQL to insert into activity table
	res, err := db.ExecContext(r.Context(),
		"INSERT INTO activity (date, time) VALUES (?, ?)",
		pd.Date, pd.Time,
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
	_, err = db.ExecContext(r.Context(),
		"INSERT INTO documentation (activity_id, file) VALUES (?, ?)",
		lastID, pd.File,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert documentation")
		log.Println("DB insert error:", err)
		return
	}

	// Inserts into personal_documentation table
	_, err = db.ExecContext(r.Context(),
		"INSERT INTO personal_documentation (activity_id, id) VALUES (?, ?)",
		lastID, pd.ID,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert personal documentation")
		log.Println("DB insert error:", err)
		return
	}

	// Writes JSON response & sends a HTTP 201 response code
	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message":     "Personal documentation created successfully",
		"activity_id": lastID,
	})
}

func DeletePersonalDocumentation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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

	// Executes SQL to delete from activity (will cascade)
	res, err := db.ExecContext(r.Context(),
		"DELETE FROM activity WHERE activity_id=?",
		activityID,
	)

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

	// Writes JSON response confirming deletion & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Personal documentation deleted successfully",
	})
}

func UpdatePersonalDocumentation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not PUT
	if r.Method != http.MethodPut {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr := vars["activity_id"]

	// Converts the "activity_id" string to an integer
	activityID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid activity ID")
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

	// Executes written SQL to update the activity data
	_, err = db.ExecContext(r.Context(),
		"UPDATE activity SET date=?, time=? WHERE activity_id=?",
		pd.Date, pd.Time, activityID,
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update activity")
		log.Println("DB update error:", err)
		return
	}

	// Executes written SQL to update the documentation data
	_, err = db.ExecContext(r.Context(),
		"UPDATE documentation SET file=? WHERE activity_id=?",
		pd.File, activityID,
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update documentation")
		log.Println("DB update error:", err)
		return
	}

	// Executes written SQL to update the personal documentation data
	res, err := db.ExecContext(r.Context(),
		"UPDATE personal_documentation SET id=? WHERE activity_id=?",
		pd.ID, activityID,
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update personal documentation")
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
		utils.WriteError(w, http.StatusNotFound, "Personal documentation not found")
		return
	}

	// Writes JSON response & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Personal documentation updated successfully",
	})
}
