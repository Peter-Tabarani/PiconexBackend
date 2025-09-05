package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"

	"github.com/gorilla/mux"
)

func GetDocumentationByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["activity_id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid activity ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT d.activity_id, a.date, a.time, d.file
		FROM documentation d
		JOIN activity a ON d.activity_id = a.activity_id
		WHERE d.activity_id = ?
	`

	var doc models.Documentation
	err = db.QueryRow(query, id).Scan(&doc.Activity_ID, &doc.Date, &doc.Time, &doc.File)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Documentation not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	jsonBytes, err := json.MarshalIndent(doc, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
