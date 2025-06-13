package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Peter-Tabarani/PiconexBackend/models"

	"github.com/gorilla/mux"
)

func GetDisabilitiesByStudentID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	studentID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT 
			d.disability_id, d.name, d.description
		FROM stu_dis sd
		JOIN disability d ON sd.disability_id = d.disability_id
		WHERE sd.id = ?
	`

	rows, err := db.Query(query, studentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var disabilities []models.Disability

	for rows.Next() {
		var d models.Disability
		if err := rows.Scan(&d.Disability_ID, &d.Name, &d.Description); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		disabilities = append(disabilities, d)
	}

	if len(disabilities) == 0 {
		http.Error(w, "No disabilities found for the student", http.StatusNotFound)
		return
	}

	jsonBytes, err := json.MarshalIndent(disabilities, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
