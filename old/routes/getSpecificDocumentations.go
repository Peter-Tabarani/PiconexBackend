package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

func GetSpecificDocumentations(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT
		ac.activity_id, s.id, sd.doc_type, ac.date, ac.time, d.file
	FROM specific_documentation sd
	JOIN activity ac ON sd.activity_id = ac.activity_id
	JOIN student s ON sd.id = s.id
	JOIN documentation d ON d.activity_id = sd.activity_id
	`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var specific_documentations []models.Specific_Documentation

	for rows.Next() {
		var sd models.Specific_Documentation
		err := rows.Scan(
			&sd.Activity_ID, &sd.ID, &sd.DocType, &sd.Date, &sd.Time, &sd.File,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		specific_documentations = append(specific_documentations, sd)
	}

	jsonBytes, err := json.MarshalIndent(specific_documentations, "", "    ") // Pretty print with 4 spaces indent
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
