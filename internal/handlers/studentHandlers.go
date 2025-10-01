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
	// All data being selected for this GET command
	query := `
		SELECT
			p.person_id, p.first_name, p.preferred_name, p.middle_name, p.last_name,
			p.email, p.phone_number, p.pronouns, p.sex, p.gender,
			p.birthday, p.address, p.city, p.state, p.zip_code, p.country, s.year,
			s.start_year, s.planned_grad_year, s.housing, s.dining
		FROM student s
		JOIN person p ON s.student_id = p.person_id
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query)

	// Error message if QueryContext fails
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
		// Parses the current data into fields of "s" variable
		if err := rows.Scan(
			&s.StudentID, &s.FirstName, &s.PreferredName, &s.MiddleName, &s.LastName,
			&s.Email, &s.PhoneNumber, &s.Pronouns, &s.Sex, &s.Gender,
			&s.Birthday, &s.Address, &s.City, &s.State, &s.ZipCode, &s.Country,
			&s.Year, &s.StartYear, &s.PlannedGradYear, &s.Housing, &s.Dining,
		); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to scan student")
			log.Println("Row scan error:", err)
			return
		}

		// Adds the obtained data to the slice
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
	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr, ok := vars["student_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing student ID")
		return
	}

	// Converts the "student_id" string to an integer
	studentID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid student ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// SQL query to select a single student
	query := `
		SELECT
			s.student_id,
			p.first_name, p.preferred_name, p.middle_name, p.last_name,
			p.email, p.phone_number, p.pronouns, p.sex, p.gender,
			p.birthday, p.address, p.city, p.state, p.zip_code, p.country,
			s.year, s.start_year, s.planned_grad_year, s.housing, s.dining
		FROM student s
		JOIN person p ON s.student_id = p.person_id
		WHERE s.student_id = ?
	`

	// Empty variable for student struct
	var s models.Student

	// Executes query
	err = db.QueryRowContext(r.Context(), query, studentID).Scan(
		&s.StudentID,
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
		// Error message if QueryRowContext or scan fails
	} else if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch student")
		log.Println("DB query error:", err)
		return
	}

	// Writes JSON response & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, s)
}

func GetStudentsByName(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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
			s.student_id, p.first_name, p.preferred_name, p.middle_name, p.last_name,
			p.email, p.phone_number, p.pronouns, p.sex, p.gender, p.birthday,
			p.address, p.city, p.state, p.zip_code, p.country,
			s.year, s.start_year, s.planned_grad_year, s.housing, s.dining
		FROM student s
		JOIN person p ON s.student_id = p.person_id
		WHERE ` + whereClause

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query, args...)

	// Error message if QueryContext fails
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
		// Parses the current data into fields of "s" variable
		if err := rows.Scan(
			&s.StudentID, &s.FirstName, &s.PreferredName, &s.MiddleName, &s.LastName,
			&s.Email, &s.PhoneNumber, &s.Pronouns, &s.Sex, &s.Gender,
			&s.Birthday, &s.Address, &s.City, &s.State, &s.ZipCode, &s.Country,
			&s.Year, &s.StartYear, &s.PlannedGradYear, &s.Housing, &s.Dining,
		); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to scan student")
			log.Println("Row scan error:", err)
			return
		}

		// Adds the obtained data to the slice
		students = append(students, s)
	}

	// Writes JSON response & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, students)
}

func CreateStudent(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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

	// Validates required fields
	if s.FirstName == "" || s.LastName == "" || s.Email == "" || s.PhoneNumber == "" ||
		s.Sex == "" || s.Birthday == "" || s.Address == "" || s.City == "" ||
		s.Country == "" || s.Year == "" || s.StartYear == 0 || s.PlannedGradYear == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Start transaction
	tx, err := db.BeginTx(r.Context(), nil)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to begin transaction")
		log.Println("BeginTx error:", err)
		return
	}
	defer tx.Rollback()

	// Executes SQL to insert into person table
	res, err := tx.ExecContext(r.Context(),
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
	_, err = tx.ExecContext(r.Context(),
		`INSERT INTO student (student_id, year, start_year, planned_grad_year, housing, dining)
		VALUES (?, ?, ?, ?, ?, ?)`,
		lastID, s.Year, s.StartYear, s.PlannedGradYear, s.Housing, s.Dining,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert into student")
		log.Println("DB insert error:", err)
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to commit transaction")
		log.Println("Transaction commit error:", err)
		return
	}

	// Writes JSON response including the new ID & sends a HTTP 201 response code
	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message":   "Student created successfully",
		"studentId": lastID,
	})
}

func UpdateStudent(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr, ok := vars["student_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing student ID")
		return
	}

	// Converts the "student_id" string to an integer
	studentID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid student ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

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

	// Validates required fields
	if s.FirstName == "" || s.LastName == "" || s.Email == "" || s.PhoneNumber == "" ||
		s.Sex == "" || s.Birthday == "" || s.Address == "" || s.City == "" ||
		s.Country == "" || s.Year == "" || s.StartYear == 0 || s.PlannedGradYear == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Start transaction
	tx, err := db.BeginTx(r.Context(), nil)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to begin transaction")
		log.Println("BeginTx error:", err)
		return
	}
	defer tx.Rollback()

	// Execute direct SQL update
	_, err = tx.ExecContext(r.Context(),
		`UPDATE person
		 SET first_name = ?, preferred_name = ?, middle_name = ?, last_name = ?,
		     email = ?, phone_number = ?, pronouns = ?, sex = ?, gender = ?,
		     birthday = ?, address = ?, city = ?, state = ?, zip_code = ?, country = ?
		 WHERE person_id = ?`,
		s.FirstName, s.PreferredName, s.MiddleName, s.LastName,
		s.Email, s.PhoneNumber, s.Pronouns, s.Sex, s.Gender,
		s.Birthday, s.Address, s.City, s.State, s.ZipCode, s.Country,
		studentID,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update student")
		log.Println("DB update error:", err)
		return
	}

	// Optionally update student table fields
	res, err := tx.ExecContext(r.Context(),
		`UPDATE student
		 SET year = ?, start_year = ?, planned_grad_year = ?, housing = ?, dining = ?
		 WHERE student_id = ?`,
		s.Year, s.StartYear, s.PlannedGradYear, s.Housing, s.Dining,
		studentID,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update student details")
		log.Println("DB update error:", err)
		return
	}

	// Gets the number of rows affected by the update
	rowsAffected, err := res.RowsAffected()

	// Error message if RowsAffected fails
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

	// Commit transaction
	if err := tx.Commit(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to commit transaction")
		log.Println("Transaction commit error:", err)
		return
	}

	// Writes JSON response confirming update & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Student updated successfully",
	})
}

func DeleteStudent(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr, ok := vars["student_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing student ID")
		return
	}

	// Converts the "student_id" string to an integer
	studentID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid student ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// Start transaction
	tx, err := db.BeginTx(r.Context(), nil)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to begin transaction")
		log.Println("BeginTx error:", err)
		return
	}
	defer tx.Rollback()

	// Executes SQL to delete from student
	res, err := tx.ExecContext(r.Context(),
		"DELETE FROM student WHERE student_id = ?",
		studentID,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete student")
		log.Println("DB delete error:", err)
		return
	}

	// Gets the number of rows affected by the delete
	rowsAffected, err := res.RowsAffected()

	// Error message if RowsAffected fails
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

	// Executes SQL to delete from person
	res, err = tx.ExecContext(r.Context(), "DELETE FROM person WHERE person_id = ?", studentID)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete person")
		log.Println("DB delete person error:", err)
		return
	}

	// Gets the number of rows affected by the delete
	rowsAffected, err = res.RowsAffected()

	// Error message if RowsAffected fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get rows affected for person")
		log.Println("RowsAffected person error:", err)
		return
	}

	// Error message if no rows were deleted
	if rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, "Person not found")
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to commit transaction")
		log.Println("Transaction commit error:", err)
		return
	}

	// Writes JSON response confirming deletion
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Student deleted successfully",
	})
}
