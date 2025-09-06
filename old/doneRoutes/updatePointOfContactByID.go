package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UpdatePOCRequest struct {
	EventDate string `json:"event_date"` // "YYYY-MM-DD"
	EventTime string `json:"event_time"` // "HH:MM:SS"
	EventType string `json:"event_type"` // "trad", "al", or "in"
	StudentID *int   `json:"student_id,omitempty"`
	AdminID   *int   `json:"admin_id,omitempty"`
}

func UpdatePointOfContactByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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
