package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"
	"github.com/gorilla/mux"
)

func GetSpecificDocumentationByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	activityIDStr := vars["activity_id"]
	activityID, err := strconv.Atoi(activityIDStr)
	if err != nil {
		http.Error(w, "Invalid activity ID", http.StatusBadRequest)
		return
	}

	query := `
        SELECT
            sd.activity_id,
            sd.id,
            sd.doc_type,
            a.date,
            a.time,
            d.file
        FROM specific_documentation sd
        JOIN activity a ON sd.activity_id = a.activity_id
        JOIN documentation d ON sd.activity_id = d.activity_id
        WHERE sd.activity_id = ?
    `

	row := db.QueryRow(query, activityID)

	var sd models.Specific_Documentation

	err = row.Scan(
		&sd.Activity_ID,
		&sd.ID,
		&sd.DocType,
		&sd.Date,
		&sd.Time,
		&sd.File,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No documentation found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.MarshalIndent(sd, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
