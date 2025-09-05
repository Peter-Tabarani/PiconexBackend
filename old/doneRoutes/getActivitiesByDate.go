package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"

	"github.com/gorilla/mux"
)

func GetActivitiesByDate(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date := vars["date"] // expects format like "2025-06-01"

	query := `
		SELECT activity_id, date, time
		FROM activity
		WHERE date = ?
	`

	rows, err := db.Query(query, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var activities []models.Activity

	for rows.Next() {
		var a models.Activity
		err := rows.Scan(&a.Activity_ID, &a.Date, &a.Time)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		activities = append(activities, a)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(activities) == 0 {
		http.Error(w, "No activities found for the specified date", http.StatusNotFound)
		return
	}

	jsonBytes, err := json.MarshalIndent(activities, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
