package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"

	"github.com/gorilla/mux"
)

func GetDocumentations(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// All data being selected for this GET command
	query := `
		SELECT ac.activity_id, ac.activity_datetime, d.file
		FROM documentation d
		JOIN activity ac ON d.activity_id = ac.activity_id
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query)

	// Error message if QueryContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to obtain documentations")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	// Creates an empty slice to obtain results
	documentations := make([]models.Documentation, 0)

	// Reads each row returned by the database
	for rows.Next() {
		var d models.Documentation
		// Parses the current data into fields of "d" variable
		if err := rows.Scan(&d.ActivityID, &d.ActivityDateTime, &d.File); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse documentations")
			log.Println("Row scan error:", err)
			return
		}

		// Adds the obtained data to the slice
		documentations = append(documentations, d)
	}

	// Checks for errors during iteration such as network interruptions and driver errors
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Writes the slice as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, documentations)
}

func GetDocumentationByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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

	// All data being selected for this GET command
	query := `
		SELECT d.activity_id, a.activity_datetime, d.file
		FROM documentation d
		JOIN activity a ON d.activity_id = a.activity_id
		WHERE d.activity_id = ?
	`

	// Empty variable for documentation struct
	var d models.Documentation

	// Executes written SQL and retrieves only one row
	err = db.QueryRowContext(r.Context(), query, activityID).Scan(&d.ActivityID, &d.ActivityDateTime, &d.File)

	// Error message if no rows are found
	if err == sql.ErrNoRows {
		utils.WriteError(w, http.StatusNotFound, "Documentation not found")
		return
		// Error message if QueryRowContext or scan fails
	} else if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch documentation")
		log.Println("DB query error:", err)
		return
	}

	// Writes the struct as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, d)
}
