package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"
	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func GetActivities(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// All data being selected for this GET command
	query := `
		SELECT
			activity_id, date, time
		FROM activity
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query)

	// Error message if QueryContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to obtain activities")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	// Creates an empty slice to obtain results
	activities := make([]models.Activity, 0)

	// Reads each row returned by the database
	for rows.Next() {
		var a models.Activity
		// Parses the current data into fields of "a" variable
		if err := rows.Scan(&a.Activity_ID, &a.Date, &a.Time); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse activities")
			log.Println("Row scan error:", err)
			return
		}

		// Adds the obtained data to the slice
		activities = append(activities, a)
	}

	// Checks for errors during iteration such as network interruptions and driver errors
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Writes the slice as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, activities)
}

func GetActivityByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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
        SELECT activity_id, date, time
        FROM activity
        WHERE activity_id = ?
    `

	// Empty variable for activity struct
	var a models.Activity

	// Executes written SQL and retrieves only one row
	err = db.QueryRowContext(r.Context(), query, activityID).Scan(
		&a.Activity_ID, &a.Date, &a.Time,
	)

	// Error message if no rows are found
	if err == sql.ErrNoRows {
		utils.WriteError(w, http.StatusNotFound, "Activity not found")
		return
		// Error message if QueryRowContext or scan fails
	} else if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch activity")
		log.Println("DB query error:", err)
		return
	}

	// Writes the struct as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, a)
}

func GetActivitiesByDate(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extracts path variables from the request
	vars := mux.Vars(r)
	date := vars["date"]
	if date == "" {
		utils.WriteError(w, http.StatusBadRequest, "Date is required")
		return
	}

	// All data being selected for this GET command
	query := `
		SELECT activity_id, date, time
		FROM activity
		WHERE date = ?
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query, date)

	// Error message if QueryContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to obtain activities")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	// Creates an empty slice to obtain results
	activities := make([]models.Activity, 0)

	// Reads each row returned by the database
	for rows.Next() {
		var a models.Activity
		// Parses the current data into fields of "a" variable
		if err := rows.Scan(&a.Activity_ID, &a.Date, &a.Time); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse activities")
			log.Println("Row scan error:", err)
			return
		}

		// Adds the obtained data to the slice
		activities = append(activities, a)
	}

	// Checks for errors during iteration such as network interruptions and driver errors
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Writes the slice as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, activities)
}
