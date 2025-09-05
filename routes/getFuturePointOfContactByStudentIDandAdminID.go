package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"

	"github.com/gorilla/mux"
)

func GetFuturePointOfContactByStudentIDAndAdminID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentIDStr := vars["student_id"]
	adminIDStr := vars["admin_id"]

	studentID, err := strconv.Atoi(studentIDStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	adminID, err := strconv.Atoi(adminIDStr)
	if err != nil {
		http.Error(w, "Invalid admin ID", http.StatusBadRequest)
		return
	}

	currentDate := time.Now().Format("2006-01-02") // MySQL DATE format

	query := `
		SELECT activity_id, event_date, event_time, event_type, student_id, admin_id
		FROM point_of_contact
		WHERE student_id = ? AND admin_id = ? AND event_date > ?
	`

	rows, err := db.Query(query, studentID, adminID, currentDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var meetings []models.PointOfContact

	for rows.Next() {
		var m models.PointOfContact
		if err := rows.Scan(
			&m.Activity_ID,
			&m.Event_Date,
			&m.Event_Time,
			&m.Event_Type,
			&m.Student_ID,
			&m.Admin_ID,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		meetings = append(meetings, m)
	}

	if len(meetings) == 0 {
		http.Error(w, "No future meetings found", http.StatusNotFound)
		return
	}

	jsonBytes, err := json.MarshalIndent(meetings, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
