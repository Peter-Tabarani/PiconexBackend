package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"
	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func GetStudents(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT
			p.id, p.first_name, p.preferred_name, p.middle_name, p.last_name,
			p.email, p.phone_number, p.pronouns, p.sex, p.gender,
			p.birthday, p.address, p.city, p.state, p.zip_code, p.country, s.year,
			s.start_year, s.planned_grad_year, s.housing, s.dining
		FROM student s
		JOIN person p ON s.id = p.id
	`

	rows, err := db.Query(query)
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
			&s.Email, &s.PhoneNumber, &s.Pronouns, &s.Sex, &s.Gender,
			&s.Birthday, &s.Address, &s.City, &s.State, &s.ZipCode, &s.Country,
			&s.Year, &s.StartYear, &s.PlannedGradYear, &s.Housing, &s.Dining,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		students = append(students, s)
	}

	jsonBytes, err := json.MarshalIndent(students, "", "    ") // Pretty print with 4 spaces indent
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

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

func GetStudentsByName(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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

type CreateStudentRequest struct {
	Person  models.Person  `json:"person"`
	Student models.Student `json:"student"`
}

func CreateStudent(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	personQuery := `
		INSERT INTO person (
			first_name, preferred_name, middle_name, last_name, email,
			phone_number, pronouns, sex, gender, birthday,
			address, city, state, zip_code, country
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := tx.Exec(personQuery,
		req.Person.FirstName, req.Person.PreferredName, req.Person.MiddleName, req.Person.LastName,
		req.Person.Email, req.Person.PhoneNumber, req.Person.Pronouns, req.Person.Sex,
		req.Person.Gender, req.Person.Birthday, req.Person.Address, req.Person.City,
		req.Person.State, req.Person.ZipCode, req.Person.Country,
	)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to insert into person: "+err.Error(), http.StatusInternalServerError)
		return
	}

	personID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to retrieve person ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	studentQuery := `
		INSERT INTO student (
			id, year, start_year, planned_grad_year, housing, dining
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err = tx.Exec(studentQuery,
		personID, req.Student.Year, req.Student.StartYear, req.Student.PlannedGradYear,
		req.Student.Housing, req.Student.Dining,
	)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to insert into student: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Student created successfully",
		"studentId": personID,
	})
}

func DeleteStudent(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	// Step 1: Nullify student_id in point_of_contact (ON DELETE SET NULL doesn't cascade)
	_, err = db.Exec(`UPDATE point_of_contact SET student_id = NULL WHERE student_id = ?`, id)
	if err != nil {
		http.Error(w, "Failed to update point_of_contact: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 2: Delete from student (cascades through specific_documentation, stu_accom, stu_dis, person, etc.)
	_, err = db.Exec(`DELETE FROM student WHERE id = ?`, id)
	if err != nil {
		http.Error(w, "Failed to delete student: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Student deleted successfully"})
}

func UpdateStudent(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var s models.Student
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Update only fields that were sent (optional: dynamic SQL builder)
	_, err = db.Exec(`
		UPDATE person
		SET gender = ?
		WHERE id = ?`, s.Gender, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Student updated successfully")
}
