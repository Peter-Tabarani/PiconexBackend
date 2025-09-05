package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"

	"github.com/gorilla/mux"
)

func GetPointOfContactByAdminIDAndDate(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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
