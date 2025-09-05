package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	//trying to push
	"github.com/Peter-Tabarani/PiconexBackend/internal/models"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

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
