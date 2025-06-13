package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/models"

	_ "github.com/go-sql-driver/mysql"
)

func GetPinned(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT admin_id, student_id
		FROM pinned
	`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pinnedList []models.Pinned

	for rows.Next() {
		var p models.Pinned
		if err := rows.Scan(&p.AdminID, &p.StudentID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pinnedList = append(pinnedList, p)
	}

	jsonBytes, err := json.MarshalIndent(pinnedList, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
