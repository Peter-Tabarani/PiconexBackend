package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"

	"github.com/gorilla/mux"
)

func GetSpecificDocumentationByStudentID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	studentID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT
			ac.activity_id, sd.id, sd.doc_type, ac.date, ac.time, d.file
		FROM specific_documentation sd
		JOIN documentation d ON sd.activity_id = d.activity_id
		JOIN activity ac ON ac.activity_id = sd.activity_id
		WHERE sd.id = ?
	`

	rows, err := db.Query(query, studentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var docs []models.Specific_Documentation

	for rows.Next() {
		var doc models.Specific_Documentation
		if err := rows.Scan(
			&doc.Activity_ID, &doc.ID, &doc.DocType, &doc.Date, &doc.Time, &doc.File,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		docs = append(docs, doc)
	}

	if len(docs) == 0 {
		http.Error(w, "No specific documentation found for the student", http.StatusNotFound)
		return
	}

	jsonBytes, err := json.MarshalIndent(docs, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
