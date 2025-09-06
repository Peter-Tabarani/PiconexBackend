package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"

	"github.com/gorilla/mux"
)

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
