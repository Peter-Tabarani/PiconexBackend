package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"

	"github.com/gorilla/mux"
)

func GetAccommodationsByStudentID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	studentID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT
			a.accommodation_id, a.name, a.description
		FROM stu_accom sa
		JOIN accommodation a ON sa.accommodation_id = a.accommodation_id
		WHERE sa.id = ?
	`

	rows, err := db.Query(query, studentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var accommodations []models.Accommodation

	for rows.Next() {
		var a models.Accommodation
		if err := rows.Scan(&a.Accommodation_ID, &a.Name, &a.Description); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		accommodations = append(accommodations, a)
	}

	if len(accommodations) == 0 {
		http.Error(w, "No accommodations found for the student", http.StatusNotFound)
		return
	}

	jsonBytes, err := json.MarshalIndent(accommodations, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
