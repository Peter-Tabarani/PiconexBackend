package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"

	"github.com/gorilla/mux"
)

func GetPersonalDocumentationByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["activity_id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid activity ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT
			pd.activity_id,
			pd.id,
			a.date,
			a.time,
			d.file
		FROM personal_documentation pd
		JOIN documentation d ON pd.activity_id = d.activity_id
		JOIN activity a ON pd.activity_id = a.activity_id
		WHERE pd.activity_id = ?
	`

	var pd models.Personal_Documentation
	err = db.QueryRow(query, id).Scan(&pd.Activity_ID, &pd.ID, &pd.Date, &pd.Time, &pd.File)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Personal documentation not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	jsonBytes, err := json.MarshalIndent(pd, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
