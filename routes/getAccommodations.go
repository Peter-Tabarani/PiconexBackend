package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/models"

	_ "github.com/go-sql-driver/mysql"
)

func GetAccommodations(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT
			am.accommodation_id, am.name, am.description
		FROM accommodation am
	`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var accommodations []models.Accommodation

	for rows.Next() {
		var am models.Accommodation
		err := rows.Scan(
			&am.Accommodation_ID, &am.Name, &am.Description,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		accommodations = append(accommodations, am)
	}

	jsonBytes, err := json.MarshalIndent(accommodations, "", "    ") // Pretty print with 4 spaces indent
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
