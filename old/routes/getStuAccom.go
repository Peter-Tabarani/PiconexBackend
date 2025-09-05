package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

func GetStuAccom(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, accommodation_id
		FROM stu_accom
	`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var stuAccomList []models.StudentAccommodation

	for rows.Next() {
		var sa models.StudentAccommodation
		if err := rows.Scan(&sa.ID, &sa.AccommodationID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		stuAccomList = append(stuAccomList, sa)
	}

	jsonBytes, err := json.MarshalIndent(stuAccomList, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
