package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/models"

	_ "github.com/go-sql-driver/mysql"
)

func GetStuDis(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, disability_id
		FROM stu_dis
	`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var stuDisList []models.StudentDisability

	for rows.Next() {
		var sd models.StudentDisability
		if err := rows.Scan(&sd.ID, &sd.DisabilityID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		stuDisList = append(stuDisList, sd)
	}

	jsonBytes, err := json.MarshalIndent(stuDisList, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
