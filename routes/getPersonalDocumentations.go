package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/models"

	_ "github.com/go-sql-driver/mysql"
)

func GetPersonalDocumentations(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT
		ac.activity_id, a.id, ac.date, ac.time, d.file
	FROM personal_documentation pd
	JOIN activity ac ON pd.activity_id = ac.activity_id
	JOIN admin a ON pd.id = a.id
	JOIN documentation d ON d.activity_id = pd.activity_id
	`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var personal_documentations []models.Personal_Documentation

	for rows.Next() {
		var pd models.Personal_Documentation
		err := rows.Scan(
			&pd.Activity_ID, &pd.ID, &pd.Date, &pd.Time, &pd.File,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		personal_documentations = append(personal_documentations, pd)
	}

	jsonBytes, err := json.MarshalIndent(personal_documentations, "", "    ") // Pretty print with 4 spaces indent
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
