package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

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
