package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"
	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func GetStudents(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not GET
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// All data being selected for this GET command
	query := `
		SELECT
			p.id, p.first_name, p.preferred_name, p.middle_name, p.last_name,
			p.email, p.phone_number, p.pronouns, p.sex, p.gender,
			p.birthday, p.address, p.city, p.state, p.zip_code, p.country, s.year,
			s.start_year, s.planned_grad_year, s.housing, s.dining
		FROM student s
		JOIN person p ON s.id = p.id
	`

	// Executes written SQL
	rows, err := db.Query(query)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to obtain students")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	// Creates an empty slice to obtain results
	students := make([]models.Student, 0)

	// Reads each row returned by the database
	for rows.Next() {
		var s models.Student
		if err := rows.Scan(
			&s.ID, &s.FirstName, &s.PreferredName, &s.MiddleName, &s.LastName,
			&s.Email, &s.PhoneNumber, &s.Pronouns, &s.Sex, &s.Gender,
			&s.Birthday, &s.Address, &s.City, &s.State, &s.ZipCode, &s.Country,
			&s.Year, &s.StartYear, &s.PlannedGradYear, &s.Housing, &s.Dining,
		); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to scan student")
			log.Println("Row scan error:", err)
			return
		}
		students = append(students, s)
	}

	// Checks for errors during iteration such as network interruptions and driver errors
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Writes JSON response & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, students)
}

func GetStudentByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not GET
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing student ID")
		return
	}

	// Converts the "id" string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid student ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// SQL query to select a single student
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

	// Empty variable for student struct
	var s models.Student

	// Executes query
	err = db.QueryRow(query, id).Scan(
		&s.ID,
		&s.FirstName,
		&s.PreferredName,
		&s.MiddleName,
		&s.LastName,
		&s.Email,
		&s.PhoneNumber,
		&s.Pronouns,
		&s.Sex,
		&s.Gender,
		&s.Birthday,
		&s.Address,
		&s.City,
		&s.State,
		&s.ZipCode,
		&s.Country,
		&s.Year,
		&s.StartYear,
		&s.PlannedGradYear,
		&s.Housing,
		&s.Dining,
	)

	// Error message if no rows are found
	if err == sql.ErrNoRows {
		utils.WriteError(w, http.StatusNotFound, "Student not found")
		return
	} else if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch student")
		log.Println("DB query error:", err)
		return
	}

	// Writes JSON response & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, s)
}

func GetStudentsByName(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not GET
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extracts path variables from the request
	vars := mux.Vars(r)
	raw, ok := vars["name"]
	if !ok || len(strings.TrimSpace(raw)) == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Missing or invalid name")
		return
	}

	// Split the name by spaces
	words := strings.Fields(raw)

	// Build a condition group for each word and join them with AND
	var conditions []string
	var args []interface{}

	for _, word := range words {
		word = "%" + strings.ToLower(word) + "%"
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

	// SQL query to select students matching the name
	query := `
		SELECT
			s.id, p.first_name, p.preferred_name, p.middle_name, p.last_name,
			p.email, p.phone_number, p.pronouns, p.sex, p.gender, p.birthday,
			p.address, p.city, p.state, p.zip_code, p.country,
			s.year, s.start_year, s.planned_grad_year, s.housing, s.dining
		FROM student s
		JOIN person p ON s.id = p.id
		WHERE ` + whereClause

	// Executes written SQL
	rows, err := db.Query(query, args...)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to obtain students")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	// Creates an empty slice to obtain results
	students := make([]models.Student, 0)

	// Reads each row returned by the database
	for rows.Next() {
		var s models.Student
		if err := rows.Scan(
			&s.ID, &s.FirstName, &s.PreferredName, &s.MiddleName, &s.LastName,
			&s.Email, &s.PhoneNumber, &s.Pronouns, &s.Sex, &s.Gender,
			&s.Birthday, &s.Address, &s.City, &s.State, &s.ZipCode, &s.Country,
			&s.Year, &s.StartYear, &s.PlannedGradYear, &s.Housing, &s.Dining,
		); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to scan student")
			log.Println("Row scan error:", err)
			return
		}
		students = append(students, s)
	}

	// Error message if no students were found
	if len(students) == 0 {
		utils.WriteError(w, http.StatusNotFound, "No student found")
		return
	}

	// Writes JSON response & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, students)
}

func CreateStudent(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not POST
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Empty variables for student struct
	var s models.Student

	// Decodes JSON body from the request into "s" variable
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents extra unexpected fields
	if err := decoder.Decode(&s); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		log.Println("JSON decode error:", err)
		return
	}

	// TECH DEBT: Validates required fields

	// Executes SQL to insert into person table
	res, err := db.ExecContext(r.Context(),
		`INSERT INTO person (
			first_name, preferred_name, middle_name, last_name, email,
			phone_number, pronouns, sex, gender, birthday,
			address, city, state, zip_code, country
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		s.FirstName, s.PreferredName, s.MiddleName, s.LastName,
		s.Email, s.PhoneNumber, s.Pronouns, s.Sex, s.Gender, s.Birthday,
		s.Address, s.City, s.State, s.ZipCode, s.Country,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert into person")
		log.Println("DB insert error:", err)
		return
	}

	// Gets the last inserted person ID
	lastID, err := res.LastInsertId()

	// Error message if LastInsertId fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get inserted person ID")
		log.Println("LastInsertId error:", err)
		return
	}

	// Executes SQL to insert into student table
	_, err = db.ExecContext(r.Context(),
		`INSERT INTO student (id, year, start_year, planned_grad_year, housing, dining)
		VALUES (?, ?, ?, ?, ?, ?)`,
		lastID, s.Year, s.StartYear, s.PlannedGradYear, s.Housing, s.Dining,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert into student")
		log.Println("DB insert error:", err)
		return
	}

	// Writes JSON response including the new ID & sends a HTTP 201 response code
	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message":   "Student created successfully",
		"studentId": lastID,
	})
}

func DeleteStudent(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not DELETE
	if r.Method != http.MethodDelete {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Converts the "id" string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid student ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// Executes SQL to delete from student
	res, err := db.ExecContext(r.Context(),
		"DELETE FROM student WHERE id = ?",
		id,
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete student")
		log.Println("DB delete error:", err)
		return
	}

	// Gets the number of rows affected by the delete
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get rows affected")
		log.Println("RowsAffected error:", err)
		return
	}

	// Error message if no rows were deleted
	if rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, "Student not found")
		return
	}

	// Writes JSON response & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Student deleted successfully",
	})
}

func UpdateStudent(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Only allow PUT method
	if r.Method != http.MethodPut {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extracts path variables from the request
	vars := mux.Vars(r)
	id := vars["id"]

	// Empty variable for student struct
	var s models.Student

	// Decode JSON directly into a temporary struct
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&s); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		log.Println("JSON decode error:", err)
		return
	}

	// Execute direct SQL update
	_, err := db.ExecContext(r.Context(),
		`UPDATE person
		 SET first_name = ?, preferred_name = ?, middle_name = ?, last_name = ?,
		     email = ?, phone_number = ?, pronouns = ?, sex = ?, gender = ?,
		     birthday = ?, address = ?, city = ?, state = ?, zip_code = ?, country = ?
		 WHERE id = ?`,
		s.FirstName, s.PreferredName, s.MiddleName, s.LastName,
		s.Email, s.PhoneNumber, s.Pronouns, s.Sex, s.Gender,
		s.Birthday, s.Address, s.City, s.State, s.ZipCode, s.Country,
		id,
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update student")
		log.Println("DB update error:", err)
		return
	}

	// Optionally update student table fields
	res, err := db.ExecContext(r.Context(),
		`UPDATE student
		 SET year = ?, start_year = ?, planned_grad_year = ?, housing = ?, dining = ?
		 WHERE id = ?`,
		s.Year, s.StartYear, s.PlannedGradYear, s.Housing, s.Dining,
		id,
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update student details")
		log.Println("DB update error:", err)
		return
	}

	// Gets the number of rows affected by the update
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get rows affected")
		log.Println("RowsAffected error:", err)
		return
	}

	// Error message if no rows were updated
	if rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, "Student not found")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Student updated successfully",
	})
}
