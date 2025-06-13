package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Peter-Tabarani/PiconexBackend/models"

	"github.com/gorilla/mux"
)

func GetPinnedByAdminID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	adminIDStr := vars["id"]
	adminID, err := strconv.Atoi(adminIDStr)
	if err != nil {
		http.Error(w, "Invalid admin ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT
			p.admin_id, s.id, pe.first_name, pe.preferred_name, pe.middle_name, pe.last_name,
			pe.email, pe.phone_number, pe.pronouns, pe.sex, pe.gender, pe.birthday,
			pe.address, pe.city, pe.state, pe.zip_code, pe.country,
			s.year, s.start_year, s.planned_grad_year, s.housing, s.dining
		FROM pinned p
		JOIN student s ON p.student_id = s.id
		JOIN person pe ON s.id = pe.id
		WHERE p.admin_id = ?
	`

	rows, err := db.Query(query, adminID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var students []models.Student

	for rows.Next() {
		var s models.Student
		var discardAdminID int

		err := rows.Scan(
			&discardAdminID,
			&s.ID, &s.FirstName, &s.PreferredName, &s.MiddleName, &s.LastName,
			&s.Email, &s.PhoneNumber, &s.Pronouns, &s.Sex, &s.Gender, &s.Birthday,
			&s.Address, &s.City, &s.State, &s.ZipCode, &s.Country,
			&s.Year, &s.StartYear, &s.PlannedGradYear, &s.Housing, &s.Dining,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		students = append(students, s)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.MarshalIndent(students, "", "    ") // Pretty print
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
