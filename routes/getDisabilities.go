package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/models"

	_ "github.com/go-sql-driver/mysql"
)

func GetDisabilities(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT
			ds.disability_id, ds.name, ds.description
		FROM disability ds
	`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var disabilities []models.Disability

	for rows.Next() {
		var ds models.Disability
		err := rows.Scan(
			&ds.Disability_ID, &ds.Name, &ds.Description,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		disabilities = append(disabilities, ds)
	}

	jsonBytes, err := json.MarshalIndent(disabilities, "", "    ") // Pretty print with 4 spaces indent
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
