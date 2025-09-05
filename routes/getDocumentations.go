package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

func GetDocumentations(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT
			ac.activity_id, ac.date, ac.time, d.file
		FROM documentation d
		JOIN activity ac ON d.activity_id = ac.activity_id
	`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var documentations []models.Documentation

	for rows.Next() {
		var d models.Documentation
		err := rows.Scan(
			&d.Activity_ID, &d.Date, &d.Time, &d.File,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		documentations = append(documentations, d)
	}

	jsonBytes, err := json.MarshalIndent(documentations, "", "    ") // Pretty print with 4 spaces indent
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
