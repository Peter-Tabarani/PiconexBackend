package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

func GetPocAdmin(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	query := `
        SELECT activity_id, admin_id
        FROM poc_adm
    `

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pocAdmins []models.PocAdmin

	for rows.Next() {
		var pa models.PocAdmin
		if err := rows.Scan(&pa.ActivityID, &pa.AdminID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pocAdmins = append(pocAdmins, pa)
	}

	jsonBytes, err := json.MarshalIndent(pocAdmins, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
