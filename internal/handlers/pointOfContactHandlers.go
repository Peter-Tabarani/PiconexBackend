package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"
	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func GetPointsOfContact(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT
			poc.activity_id,
			poc.event_date,
			poc.event_time,
			poc.event_type,
			poc.student_id,
			poc.admin_id
		FROM point_of_contact poc
	`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pointOfContacts []models.PointOfContact

	for rows.Next() {
		var poc models.PointOfContact
		var studentID sql.NullInt64
		var adminID sql.NullInt64

		err := rows.Scan(
			&poc.Activity_ID, &poc.Event_Date, &poc.Event_Time, &poc.Event_Type,
			&studentID, &adminID,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 🧪 Debug print to terminal/log

		if studentID.Valid {
			temp := int(studentID.Int64)
			poc.Student_ID = &temp
		}
		if adminID.Valid {
			temp := int(adminID.Int64)
			poc.Admin_ID = &temp
		}

		pointOfContacts = append(pointOfContacts, poc)
	}

	jsonBytes, err := json.MarshalIndent(pointOfContacts, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func GetPointOfContactByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	activityIDStr := vars["activity_id"]
	activityID, err := strconv.Atoi(activityIDStr)
	if err != nil {
		http.Error(w, "Invalid activity ID", http.StatusBadRequest)
		return
	}

	query := `
        SELECT
            poc.activity_id,
            poc.event_date,
            poc.event_time,
            poc.event_type,
            ps.student_id,
            pa.admin_id
        FROM point_of_contact poc
        LEFT JOIN poc_stu ps ON poc.activity_id = ps.activity_id
        LEFT JOIN poc_adm pa ON poc.activity_id = pa.activity_id
        WHERE poc.activity_id = ?
    `

	row := db.QueryRow(query, activityID)

	var poc models.PointOfContact
	var studentID sql.NullInt64
	var adminID sql.NullInt64

	err = row.Scan(
		&poc.Activity_ID,
		&poc.Event_Date,
		&poc.Event_Time,
		&poc.Event_Type,
		&studentID,
		&adminID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Point of Contact not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if studentID.Valid {
		poc.Student_ID = new(int)
		*poc.Student_ID = int(studentID.Int64)
	} else {
		poc.Student_ID = nil
	}

	if adminID.Valid {
		poc.Admin_ID = new(int)
		*poc.Admin_ID = int(adminID.Int64)
	} else {
		poc.Admin_ID = nil
	}

	jsonBytes, err := json.MarshalIndent(poc, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func GetPointsOfContactByAdminIDAndDate(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	date := vars["date"]

	adminID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid admin ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT activity_id, event_date, event_time, event_type, student_id, admin_id
		FROM point_of_contact
		WHERE admin_id = ? AND event_date = ?
	`

	rows, err := db.Query(query, adminID, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var contacts []models.PointOfContact

	for rows.Next() {
		var poc models.PointOfContact
		if err := rows.Scan(
			&poc.Activity_ID,
			&poc.Event_Date,
			&poc.Event_Time,
			&poc.Event_Type,
			&poc.Student_ID,
			&poc.Admin_ID,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		contacts = append(contacts, poc)
	}

	if len(contacts) == 0 {
		http.Error(w, "No point of contact records found", http.StatusNotFound)
		return
	}

	jsonBytes, err := json.MarshalIndent(contacts, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func GetFuturePointsOfContactByStudentIDAndAdminID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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

type CreatePOCRequest struct {
	ActivityID int    `json:"activity_id"`
	EventDate  string `json:"event_date"` // format: YYYY-MM-DD
	EventTime  string `json:"event_time"` // format: HH:MM:SS
	EventType  string `json:"event_type"` // "trad", "al", or "in"
	StudentID  *int   `json:"student_id,omitempty"`
	AdminID    *int   `json:"admin_id,omitempty"`
}

func CreatePointOfContact(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreatePOCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(`
		INSERT INTO point_of_contact (activity_id, event_date, event_time, event_type, student_id, admin_id)
		VALUES (?, ?, ?, ?, ?, ?)`,
		req.ActivityID, req.EventDate, req.EventTime, req.EventType, req.StudentID, req.AdminID,
	)
	if err != nil {
		http.Error(w, "Failed to insert point_of_contact: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Point of Contact created successfully",
	})
}

func DeletePointOfContact(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")

	vars := mux.Vars(r)
	activityIDStr := vars["activity_id"]

	activityID, err := strconv.Atoi(activityIDStr)
	if err != nil {
		http.Error(w, "Invalid activity ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`DELETE FROM point_of_contact WHERE activity_id = ?`, activityID)
	if err != nil {
		http.Error(w, "Failed to delete point_of_contact: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Point of Contact deleted successfully"})
}

type UpdatePOCRequest struct {
	EventDate string `json:"event_date"` // "YYYY-MM-DD"
	EventTime string `json:"event_time"` // "HH:MM:SS"
	EventType string `json:"event_type"` // "trad", "al", or "in"
	StudentID *int   `json:"student_id,omitempty"`
	AdminID   *int   `json:"admin_id,omitempty"`
}

func UpdatePointOfContact(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	activityIDStr := vars["activity_id"]
	activityID, err := strconv.Atoi(activityIDStr)
	if err != nil {
		http.Error(w, "Invalid activity ID", http.StatusBadRequest)
		return
	}

	var req UpdatePOCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`
		UPDATE point_of_contact 
		SET event_date=?, event_time=?, event_type=?, student_id=?, admin_id=?
		WHERE activity_id=?`,
		req.EventDate, req.EventTime, req.EventType, req.StudentID, req.AdminID, activityID,
	)
	if err != nil {
		http.Error(w, "Failed to update point_of_contact: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Point of Contact updated successfully"})
}
