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

func GetPointsOfContact(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// All data being selected for this GET command
	query := `
		SELECT
    		poc.point_of_contact_id,
    		a.activity_datetime,
    		poc.event_datetime,
    		poc.duration,
    		poc.event_type,
    		poc.student_id
		FROM point_of_contact poc
		JOIN activity a ON poc.point_of_contact_id = a.activity_id
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query)

	// Error message if QueryContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to obtain points of contact")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	// Creates an empty slice to obtain results
	pointsOfContact := make([]models.PointOfContact, 0)

	// Reads each row returned by the database
	for rows.Next() {
		var poc models.PointOfContact
		// Parses the current data into fields of "poc" variable
		if err := rows.Scan(
			&poc.PointOfContactID,
			&poc.ActivityDateTime,
			&poc.EventDateTime,
			&poc.Duration,
			&poc.EventType,
			&poc.StudentID,
		); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse points of contact")
			log.Println("Row scan error:", err)
			return
		}

		// Adds the obtained data to the slice
		pointsOfContact = append(pointsOfContact, poc)
	}

	// Checks for errors during iteration such as network interruptions and driver errors
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Writes the slice as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, pointsOfContact)
}

func GetPointOfContactByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extracts path variables from the request
	vars := mux.Vars(r)
	pointOfContactIDStr, ok := vars["point_of_contact_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing point of contact ID")
		return
	}

	// Converts the "point_of_contact_id" string to an integer
	pointOfContactID, err := strconv.Atoi(pointOfContactIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid point of contact ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// All data being selected for this GET command
	query := `
		SELECT
    		poc.point_of_contact_id,
    		a.activity_datetime,
    		poc.event_datetime,
    		poc.duration,
    		poc.event_type,
    		poc.student_id
		FROM point_of_contact poc
		JOIN activity a ON poc.point_of_contact_id = a.activity_id
		WHERE poc.point_of_contact_id = ?
	`

	// Empty variable for PointOfContact struct
	var poc models.PointOfContact

	// Executes written SQL and retrieves only one row
	err = db.QueryRowContext(r.Context(), query, pointOfContactID).Scan(
		&poc.PointOfContactID,
		&poc.ActivityDateTime,
		&poc.EventDateTime,
		&poc.Duration,
		&poc.EventType,
		&poc.StudentID,
	)

	// Error message if no rows are found
	if err == sql.ErrNoRows {
		utils.WriteError(w, http.StatusNotFound, "Point of Contact not found")
		return
		// Error message if QueryRowContext or scan fails
	} else if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch point of contact")
		log.Println("DB query error:", err)
		return
	}

	// Writes the struct as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, poc)
}

func GetFuturePointsOfContactByAdminID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extracts path variables from the request
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
		log.Println("Invalid admin ID parse error:", err)
		return
	}

	currentDate := time.Now().Format("2006-01-02") // MySQL DATE format

	// Query: join point_of_contact -> activity -> poc_adm
	query := `
		SELECT
    		poc.point_of_contact_id,
    		a.activity_datetime,
    		poc.event_datetime,
    		poc.duration,
    		poc.event_type,
    		poc.student_id
		FROM point_of_contact poc
		JOIN activity a ON poc.point_of_contact_id = a.activity_id
		JOIN poc_admin pa ON poc.point_of_contact_id = pa.point_of_contact_id
		WHERE pa.admin_id = ? AND poc.event_datetime > ?
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query, adminID, currentDate)

	// Error message if QueryContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to obtain future points of contact")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	// Creates an empty slice to obtain results
	pointsOfContact := make([]models.PointOfContact, 0)

	// Reads each row returned by the database
	for rows.Next() {
		var poc models.PointOfContact
		// Parses the current data into fields of "poc" variable
		if err := rows.Scan(
			&poc.PointOfContactID,
			&poc.ActivityDateTime,
			&poc.EventDateTime,
			&poc.Duration,
			&poc.EventType,
			&poc.StudentID,
		); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse points of contact")
			log.Println("Row scan error:", err)
			return
		}

		// Adds the obtained data to the slice
		pointsOfContact = append(pointsOfContact, poc)
	}

	// Checks for errors during iteration
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Writes the slice as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, pointsOfContact)
}

func GetPastPointsOfContactByStudentIDAndAdminID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extracts path variables from the request
	vars := mux.Vars(r)
	studentIDStr, ok1 := vars["student_id"]
	adminIDStr, ok2 := vars["admin_id"]
	if !ok1 || !ok2 {
		utils.WriteError(w, http.StatusBadRequest, "Missing student ID or admin ID")
		return
	}

	// Converts the "student_id" and "admin_id" strings to integers
	studentID, err := strconv.Atoi(studentIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid student ID")
		log.Println("Invalid student ID parse error:", err)
		return
	}

	adminID, err := strconv.Atoi(adminIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid admin ID")
		log.Println("Invalid admin ID parse error:", err)
		return
	}

	currentDate := time.Now().Format("2006-01-02") // MySQL DATE format

	// Query: join point_of_contact -> activity -> poc_adm
	query := `
		SELECT
    		poc.point_of_contact_id,
    		a.activity_datetime,
    		poc.event_datetime,
    		poc.duration,
    		poc.event_type,
    		poc.student_id
		FROM point_of_contact poc
		JOIN activity a ON poc.point_of_contact_id = a.activity_id
		JOIN poc_admin pa ON poc.point_of_contact_id = pa.point_of_contact_id
		WHERE poc.student_id = ? AND pa.admin_id = ? AND poc.event_datetime < ?
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query, studentID, adminID, currentDate)

	// Error message if QueryContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to obtain future points of contact")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	// Creates an empty slice to obtain results
	pointsOfContact := make([]models.PointOfContact, 0)

	// Reads each row returned by the database
	for rows.Next() {
		var poc models.PointOfContact
		// Parses the current data into fields of "poc" variable
		if err := rows.Scan(
			&poc.PointOfContactID,
			&poc.ActivityDateTime,
			&poc.EventDateTime,
			&poc.Duration,
			&poc.EventType,
			&poc.StudentID,
		); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse points of contact")
			log.Println("Row scan error:", err)
			return
		}

		// Adds the obtained data to the slice
		pointsOfContact = append(pointsOfContact, poc)
	}

	// Checks for errors during iteration
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Writes the slice as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, pointsOfContact)
}

func GetFuturePointsOfContactByStudentID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extracts path variables from the request
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
		log.Println("Invalid student ID parse error:", err)
		return
	}

	currentDate := time.Now().Format("2006-01-02") // MySQL DATE format

	// Query: join point_of_contact -> activity -> poc_adm
	query := `
		SELECT
			poc.point_of_contact_id,
			a.activity_datetime,
			poc.event_datetime,
			poc.duration,
			poc.event_type,
			poc.student_id
		FROM point_of_contact poc
		JOIN activity a ON poc.point_of_contact_id = a.activity_id
		WHERE poc.student_id = ? AND poc.event_datetime > ?
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query, studentID, currentDate)

	// Error message if QueryContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to obtain future points of contact")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	// Creates an empty slice to obtain results
	pointsOfContact := make([]models.PointOfContact, 0)

	// Reads each row returned by the database
	for rows.Next() {
		var poc models.PointOfContact
		// Parses the current data into fields of "poc" variable
		if err := rows.Scan(
			&poc.PointOfContactID,
			&poc.ActivityDateTime,
			&poc.EventDateTime,
			&poc.Duration,
			&poc.EventType,
			&poc.StudentID,
		); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse points of contact")
			log.Println("Row scan error:", err)
			return
		}

		// Adds the obtained data to the slice
		pointsOfContact = append(pointsOfContact, poc)
	}

	// Checks for errors during iteration
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Writes the slice as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, pointsOfContact)
}

func GetFuturePointsOfContactByStudentIDAndAdminID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extracts path variables from the request
	vars := mux.Vars(r)
	studentIDStr, ok1 := vars["student_id"]
	adminIDStr, ok2 := vars["admin_id"]
	if !ok1 || !ok2 {
		utils.WriteError(w, http.StatusBadRequest, "Missing student ID or admin ID")
		return
	}

	// Converts the "student_id" and "admin_id" strings to integers
	studentID, err := strconv.Atoi(studentIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid student ID")
		log.Println("Invalid student ID parse error:", err)
		return
	}

	adminID, err := strconv.Atoi(adminIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid admin ID")
		log.Println("Invalid admin ID parse error:", err)
		return
	}

	currentDate := time.Now().Format("2006-01-02") // MySQL DATE format

	// Query: join point_of_contact -> activity -> poc_adm
	query := `
		SELECT
    		poc.point_of_contact_id,
    		a.activity_datetime,
    		poc.event_datetime,
    		poc.duration,
    		poc.event_type,
    		poc.student_id
		FROM point_of_contact poc
		JOIN activity a ON poc.point_of_contact_id = a.activity_id
		JOIN poc_admin pa ON poc.point_of_contact_id = pa.point_of_contact_id
		WHERE poc.student_id = ? AND pa.admin_id = ? AND poc.event_datetime > ?
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query, studentID, adminID, currentDate)

	// Error message if QueryContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to obtain future points of contact")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	// Creates an empty slice to obtain results
	pointsOfContact := make([]models.PointOfContact, 0)

	// Reads each row returned by the database
	for rows.Next() {
		var poc models.PointOfContact
		// Parses the current data into fields of "poc" variable
		if err := rows.Scan(
			&poc.PointOfContactID,
			&poc.ActivityDateTime,
			&poc.EventDateTime,
			&poc.Duration,
			&poc.EventType,
			&poc.StudentID,
		); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse points of contact")
			log.Println("Row scan error:", err)
			return
		}

		// Adds the obtained data to the slice
		pointsOfContact = append(pointsOfContact, poc)
	}

	// Checks for errors during iteration
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Writes the slice as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, pointsOfContact)
}

func CreatePointOfContact(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Empty variable for PointOfContact struct
	var poc models.PointOfContact

	// Decodes JSON body from the request into "poc" variable
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents extra unexpected fields
	if err := decoder.Decode(&poc); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		log.Println("JSON decode error:", err)
		return
	}

	// Automatically set activity_datetime to now
	poc.ActivityDateTime = time.Now()

	// Validates required fields
	if poc.StudentID == 0 || poc.Duration == 0 || poc.EventType == "" || poc.EventDateTime.IsZero() {
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

	// Executes written SQL to insert a new activity
	res, err := tx.ExecContext(r.Context(),
		`INSERT INTO activity (activity_datetime) VALUES (?)`,
		poc.ActivityDateTime,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert activity")
		log.Println("DB insert activity error:", err)
		return
	}

	// Gets the ID of the newly inserted accommodation
	lastID, err := res.LastInsertId()

	// Error message if LastInsertId fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get last insert ID for activity")
		log.Println("LastInsertId error:", err)
		return
	}

	// Executes written SQL to insert a new point of contact
	_, err = tx.ExecContext(r.Context(),
		`INSERT INTO point_of_contact (point_of_contact_id, event_datetime, duration, event_type, student_id) VALUES (?, ?, ?, ?, ?)`,
		lastID, poc.EventDateTime, poc.Duration, poc.EventType, poc.StudentID,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert point of contact")
		log.Println("DB insert point_of_contact error:", err)
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to commit transaction")
		log.Println("Transaction commit error:", err)
		return
	}

	// Writes JSON response including the new activity_id & sends a HTTP 201 response code
	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message":             "Point of Contact created successfully",
		"point_of_contact_id": lastID,
	})
}

func UpdatePointOfContact(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extracts path variables from the request
	vars := mux.Vars(r)
	pointOfContactIDStr, ok := vars["point_of_contact_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing activity ID")
		return
	}

	// Converts the "point_of_contact_id" string to an integer
	pointOfContactID, err := strconv.Atoi(pointOfContactIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid point of contact ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// Empty variable for PointOfContact struct
	var poc models.PointOfContact

	// Decodes JSON body from the request into "poc" variable
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents extra unexpected fields
	if err := decoder.Decode(&poc); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		log.Println("JSON decode error:", err)
		return
	}

	// Automatically set activity_datetime to now
	poc.ActivityDateTime = time.Now()

	// Validates required fields
	if poc.StudentID == 0 || poc.Duration == 0 || poc.EventType == "" || poc.EventDateTime.IsZero() {
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

	// Updates the activity table first
	_, err = tx.ExecContext(r.Context(),
		`UPDATE activity SET activity_datetime=? WHERE activity_id=?`,
		poc.ActivityDateTime, pointOfContactID)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update activity")
		log.Println("DB update error:", err)
		return
	}

	// Updates the point_of_contact table
	res, err := tx.ExecContext(r.Context(),
		`UPDATE point_of_contact SET event_datetime=?, duration=?, event_type=?, student_id=? WHERE point_of_contact_id=?`,
		poc.EventDateTime, poc.Duration, poc.EventType, poc.StudentID, pointOfContactID)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update point of contact")
		log.Println("DB update point_of_contact error:", err)
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
		utils.WriteError(w, http.StatusNotFound, "Point of Contact not found")
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
		"message": "Point of Contact updated successfully",
	})
}

func DeletePointOfContact(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extracts path variables from the request
	vars := mux.Vars(r)
	pointOfContactIDStr, ok := vars["point_of_contact_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing point of contact ID")
		return
	}

	// Converts the "point_of_contact_id" string to an integer
	pointOfContactID, err := strconv.Atoi(pointOfContactIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid point of contact ID")
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

	// Executes SQL to delete from point of contact
	res, err := tx.ExecContext(r.Context(),
		"DELETE FROM point_of_contact WHERE point_of_contact_id = ?",
		pointOfContactID,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete point of contact")
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
		utils.WriteError(w, http.StatusNotFound, "Point of contact not found")
		return
	}

	// Executes written SQL to delete the point of contact
	res, err = tx.ExecContext(r.Context(),
		"DELETE FROM activity WHERE activity_id = ?", pointOfContactID)

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
		"message": "Point of Contact deleted successfully",
	})
}
