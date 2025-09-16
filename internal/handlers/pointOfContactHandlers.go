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
	// Error message for any request that is not GET
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// All data being selected for this GET command
	query := `
		SELECT
			a.activity_id, a.date, a.time,
			poc.event_date, poc.event_time, poc.event_type,
			poc.id
		FROM point_of_contact poc
		JOIN activity a ON poc.activity_id = a.activity_id
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query)
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
			&poc.Activity_ID, &poc.Date, &poc.Time,
			&poc.Event_Date, &poc.Event_Time, &poc.Event_Type,
			&poc.ID,
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
	// Error message for any request that is not GET
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extracts path variables from the request
	vars := mux.Vars(r)
	activityIDStr, ok := vars["activity_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing activity ID")
		return
	}

	// Converts the "activity_id" string to an integer
	activityID, err := strconv.Atoi(activityIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid activity ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// All data being selected for this GET command
	query := `
		SELECT
			a.activity_id, a.date, a.time,
			poc.event_date, poc.event_time, poc.event_type,
			poc.id
		FROM point_of_contact poc
		JOIN activity a ON poc.activity_id = a.activity_id
		WHERE poc.activity_id = ?
	`

	// Empty variable for PointOfContact struct
	var poc models.PointOfContact

	// Executes written SQL and retrieves only one row
	err = db.QueryRowContext(r.Context(), query, activityID).Scan(
		&poc.Activity_ID, &poc.Date, &poc.Time,
		&poc.Event_Date, &poc.Event_Time, &poc.Event_Type,
		&poc.ID,
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

func GetPointsOfContactByAdminIDAndDate(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not GET
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	date, dateOk := vars["date"]
	if !ok || !dateOk {
		utils.WriteError(w, http.StatusBadRequest, "Missing admin ID or date")
		return
	}

	// Converts the "id" string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid admin ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// Query: join point_of_contact -> activity and poc_adm
	query := `
		SELECT
			a.activity_id, a.date, a.time,
			poc.event_date, poc.event_time, poc.event_type,
			poc.id
		FROM point_of_contact poc
		JOIN activity a ON poc.activity_id = a.activity_id
		JOIN poc_adm pa ON poc.activity_id = pa.activity_id
		WHERE pa.id = ? AND poc.event_date = ?
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query, id, date)
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
		if err := rows.Scan(
			&poc.Activity_ID, &poc.Date, &poc.Time,
			&poc.Event_Date, &poc.Event_Time, &poc.Event_Type,
			&poc.ID,
		); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse points of contact")
			log.Println("Row scan error:", err)
			return
		}
		pointsOfContact = append(pointsOfContact, poc)
	}

	// Checks for errors during iteration
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Error message if no rows were found
	if len(pointsOfContact) == 0 {
		utils.WriteError(w, http.StatusNotFound, "No point of contact records found")
		return
	}

	// Writes the slice as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, pointsOfContact)
}

func GetFuturePointsOfContactByStudentIDAndAdminID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not GET
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

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
			a.activity_id, a.date, a.time,
			poc.event_date, poc.event_time, poc.event_type,
			poc.id
		FROM point_of_contact poc
		JOIN activity a ON poc.activity_id = a.activity_id
		JOIN poc_adm pa ON poc.activity_id = pa.activity_id
		WHERE poc.id = ? AND pa.id = ? AND poc.event_date > ?
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query, studentID, adminID, currentDate)
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
		if err := rows.Scan(
			&poc.Activity_ID, &poc.Date, &poc.Time,
			&poc.Event_Date, &poc.Event_Time, &poc.Event_Type,
			&poc.ID,
		); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse future points of contact")
			log.Println("Row scan error:", err)
			return
		}
		pointsOfContact = append(pointsOfContact, poc)
	}

	// Checks for errors during iteration
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Error message if no rows were found
	if len(pointsOfContact) == 0 {
		utils.WriteError(w, http.StatusNotFound, "No future points of contact found")
		return
	}

	// Writes the slice as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, pointsOfContact)
}

func CreatePointOfContact(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not POST
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
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

	// Validates required fields
	if poc.Event_Date == "" || poc.Event_Time == "" || poc.Event_Type == "" {
		utils.WriteError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	// First, insert into activity table
	res, err := db.ExecContext(r.Context(),
		`INSERT INTO activity (date, time) VALUES (?, ?)`,
		poc.Event_Date, poc.Event_Time,
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert activity")
		log.Println("DB insert activity error:", err)
		return
	}

	// Get the generated activity_id
	activityID, err := res.LastInsertId()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get last insert ID for activity")
		log.Println("LastInsertId error:", err)
		return
	}

	// Then, insert into point_of_contact table using the new activity_id
	_, err = db.ExecContext(r.Context(),
		`INSERT INTO point_of_contact (activity_id, event_date, event_time, event_type, id) VALUES (?, ?, ?, ?, ?)`,
		activityID, poc.Event_Date, poc.Event_Time, poc.Event_Type, poc.ID,
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert point of contact")
		log.Println("DB insert point_of_contact error:", err)
		return
	}

	// Writes JSON response including the new activity_id & sends a HTTP 201 response code
	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message":     "Point of Contact created successfully",
		"activity_id": activityID,
	})
}

func DeletePointOfContact(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not DELETE
	if r.Method != http.MethodDelete {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extracts path variables from the request
	vars := mux.Vars(r)
	activityIDStr, ok := vars["activity_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing activity ID")
		return
	}

	// Converts the "activity_id" string to an integer
	activityID, err := strconv.Atoi(activityIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid activity ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// Deletes the point_of_contact first (child table)
	res, err := db.ExecContext(r.Context(),
		"DELETE FROM point_of_contact WHERE activity_id = ?", activityID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete point of contact")
		log.Println("DB delete error:", err)
		return
	}

	// Check if any row was deleted
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get rows affected")
		log.Println("RowsAffected error:", err)
		return
	}
	if rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, "Point of Contact not found")
		return
	}

	// Deletes the activity row (parent table)
	_, err = db.ExecContext(r.Context(),
		"DELETE FROM activity WHERE activity_id = ?", activityID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete activity")
		log.Println("DB delete activity error:", err)
		return
	}

	// Writes JSON response confirming deletion & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Point of Contact deleted successfully",
	})
}

func UpdatePointOfContact(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not PUT
	if r.Method != http.MethodPut {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extracts path variables from the request
	vars := mux.Vars(r)
	activityIDStr, ok := vars["activity_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing activity ID")
		return
	}

	// Converts the "activity_id" string to an integer
	activityID, err := strconv.Atoi(activityIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid activity ID")
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

	// Validates required fields
	if poc.Event_Date == "" || poc.Event_Time == "" || poc.Event_Type == "" {
		utils.WriteError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Updates the activity table first
	_, err = db.ExecContext(r.Context(),
		`UPDATE activity SET date=?, time=? WHERE activity_id=?`,
		poc.Event_Date, poc.Event_Time, activityID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update activity")
		log.Println("DB update activity error:", err)
		return
	}

	// Updates the point_of_contact table
	res, err := db.ExecContext(r.Context(),
		`UPDATE point_of_contact SET event_date=?, event_time=?, event_type=?, id=? WHERE activity_id=?`,
		poc.Event_Date, poc.Event_Time, poc.Event_Type, poc.ID, activityID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update point of contact")
		log.Println("DB update point_of_contact error:", err)
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
		utils.WriteError(w, http.StatusNotFound, "Point of Contact not found")
		return
	}

	// Writes JSON response confirming update & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Point of Contact updated successfully",
	})
}
