package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Peter-Tabarani/PiconexBackend/models"

	"github.com/gorilla/mux"
)

func GetStudentByName(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	raw := vars["name"]
	words := strings.Fields(raw) // Split by space

	if len(words) == 0 {
		http.Error(w, "Invalid name", http.StatusBadRequest)
		return
	}

	// Build a condition group for each word and join them with AND
	var conditions []string
	var args []interface{}

	for _, word := range words {
		word = "%" + strings.ToLower(word) + "%" // Match anywhere, not just prefix
		group := `(
			LOWER(p.first_name) LIKE ? OR
			LOWER(p.last_name) LIKE ? OR
			LOWER(p.preferred_name) LIKE ? OR
			LOWER(p.middle_name) LIKE ?
		)`
		conditions = append(conditions, group)
		for i := 0; i < 4; i++ {
			args = append(args, word)
		}
	}

	whereClause := strings.Join(conditions, " AND ")

	query := `
		SELECT 
			s.id, p.first_name, p.preferred_name, p.middle_name, p.last_name,
			p.email, p.phone_number, p.pronouns, p.sex, p.gender, p.birthday,
			p.address, p.city, p.state, p.zip_code, p.country,
			s.year, s.start_year, s.planned_grad_year, s.housing, s.dining
		FROM student s
		JOIN person p ON s.id = p.id
		WHERE ` + whereClause

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var students []models.Student

	for rows.Next() {
		var s models.Student
		err := rows.Scan(
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

	if len(students) == 0 {
		http.Error(w, "No student found", http.StatusNotFound)
		return
	}

	jsonBytes, err := json.MarshalIndent(students, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
