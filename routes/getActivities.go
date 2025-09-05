package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

func GetActivities(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT
			ac.activity_id, ac.date, ac.time
		FROM activity ac
	`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var activities []models.Activity

	for rows.Next() {
		var ac models.Activity
		err := rows.Scan(
			&ac.Activity_ID, &ac.Date, &ac.Time,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		activities = append(activities, ac)
	}

	jsonBytes, err := json.MarshalIndent(activities, "", "    ") // Pretty print with 4 spaces indent
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
