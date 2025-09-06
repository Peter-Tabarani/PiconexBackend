package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

func GetPointOfContact(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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
