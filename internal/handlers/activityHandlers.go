package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"
	"github.com/gorilla/mux"

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

func GetActivityByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["activity_id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid activity ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT activity_id, date, time
		FROM activity
		WHERE activity_id = ?
	`

	var activity models.Activity
	err = db.QueryRow(query, id).Scan(
		&activity.Activity_ID,
		&activity.Date,
		&activity.Time,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Activity not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	jsonBytes, err := json.MarshalIndent(activity, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

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

func GetActivitiesByStudentID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	studentID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT
			a.activity_id, a.date, a.time
		FROM specific_documentation sd
		JOIN activity a ON sd.activity_id = a.activity_id
		WHERE sd.id = ?
	`

	rows, err := db.Query(query, studentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var activities []models.Activity

	for rows.Next() {
		var a models.Activity
		if err := rows.Scan(&a.Activity_ID, &a.Date, &a.Time); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		activities = append(activities, a)
	}

	if len(activities) == 0 {
		http.Error(w, "No activities found for the student", http.StatusNotFound)
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
