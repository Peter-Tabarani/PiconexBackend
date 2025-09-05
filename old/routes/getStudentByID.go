package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func GetStudentByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT
			s.id,
			p.first_name, p.preferred_name, p.middle_name, p.last_name,
			p.email, p.phone_number, p.pronouns, p.sex, p.gender,
			p.birthday, p.address, p.city, p.state, p.zip_code, p.country,
			s.year, s.start_year, s.planned_grad_year, s.housing, s.dining
		FROM student s
		JOIN person p ON s.id = p.id
		WHERE s.id = ?
	`

	var student models.Student
	err = db.QueryRow(query, id).Scan(
		&student.ID,
		&student.FirstName,
		&student.PreferredName,
		&student.MiddleName,
		&student.LastName,
		&student.Email,
		&student.PhoneNumber,
		&student.Pronouns,
		&student.Sex,
		&student.Gender,
		&student.Birthday,
		&student.Address,
		&student.City,
		&student.State,
		&student.ZipCode,
		&student.Country,
		&student.Year,
		&student.StartYear,
		&student.PlannedGradYear,
		&student.Housing,
		&student.Dining,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Student not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	jsonBytes, err := json.MarshalIndent(student, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
