package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func GetPinned(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// All data being selected for this GET command
	query := `
		SELECT admin_id, student_id
		FROM pinned
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query)

	// Error message if QueryContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to obtain pinned records")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	// Creates an empty slice to obtain results
	pinnedList := make([]models.Pinned, 0)

	// Reads each row returned by the database
	for rows.Next() {
		var p models.Pinned
		// Parses the current data into fields of "p" variable
		if err := rows.Scan(&p.AdminID, &p.StudentID); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse pinned record")
			log.Println("Row scan error:", err)
			return
		}
		// Adds the obtained data to the slice
		pinnedList = append(pinnedList, p)
	}

	// Checks for errors during iteration such as network interruptions and driver errors
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Writes the slice as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, pinnedList)
}

func GetPin(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extracts path variables from the request
	vars := mux.Vars(r)
	adminIDStr, ok := vars["admin_id"]
	studentIDStr, ok2 := vars["student_id"]
	if !ok || !ok2 {
		utils.WriteError(w, http.StatusBadRequest, "Missing admin ID or student ID")
		return
	}

	// Converts the "admin_id" string to an integer
	adminID, err := strconv.Atoi(adminIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid admin ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// Converts the "student_id" string to an integer
	studentID, err := strconv.Atoi(studentIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid student ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	query := `
		SELECT 1
		FROM pinned
		WHERE admin_id = ? AND student_id = ?
		LIMIT 1
	`
	// Executes written SQL
	var exists bool
	err = db.QueryRowContext(r.Context(), query, adminID, studentID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			// Pin not found → false
			json.NewEncoder(w).Encode(false)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Database query error")
		log.Println("DB query error:", err)
		return
	}

	// If a row was found, exists = true
	json.NewEncoder(w).Encode(true)

}

func GetPinnedByAdminID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr, ok := vars["admin_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing admin ID")
		return
	}

	// Converts the "admin_id" string to an integer
	adminID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid admin ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// All data being selected for this GET command
	query := `
		SELECT
			s.student_id, pe.first_name, pe.preferred_name, pe.middle_name, pe.last_name,
			pe.email, pe.phone_number, pe.pronouns, pe.sex, pe.gender, pe.birthday,
			pe.address, pe.city, pe.state, pe.zip_code, pe.country,
			s.year, s.start_year, s.planned_grad_year, s.housing, s.dining
		FROM pinned p
		JOIN student s ON p.student_id = s.student_id
		JOIN person pe ON s.student_id = pe.person_id
		WHERE p.admin_id = ?
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query, adminID)

	// Error message if QueryContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to obtain students pinned by admin")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	// Creates an empty slice to obtain results
	students := make([]models.Student, 0)

	// Reads each row returned by the database
	for rows.Next() {
		var s models.Student
		// Parses the current row into student struct
		if err := rows.Scan(
			&s.StudentID, &s.FirstName, &s.PreferredName, &s.MiddleName, &s.LastName,
			&s.Email, &s.PhoneNumber, &s.Pronouns, &s.Sex, &s.Gender, &s.Birthday,
			&s.Address, &s.City, &s.State, &s.ZipCode, &s.Country,
			&s.Year, &s.StartYear, &s.PlannedGradYear, &s.Housing, &s.Dining,
		); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse student record")
			log.Println("Row scan error:", err)
			return
		}
		students = append(students, s)
	}

	// Checks for errors during iteration
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Writes the slice as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, students)
}

func CreatePinned(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Empty variable for request struct
	var req models.Pinned
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents extra unexpected fields
	if err := decoder.Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		log.Println("JSON decode error:", err)
		return
	}

	// Validates required fields
	if req.AdminID == 0 || req.StudentID == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Executes written SQL to insert pinned record
	_, err := db.ExecContext(r.Context(),
		"INSERT INTO pinned (admin_id, student_id) VALUES (?, ?)",
		req.AdminID, req.StudentID,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert pinned record")
		log.Println("DB insert error:", err)
		return
	}

	// Writes JSON response confirming creation & sends a HTTP 201 response code
	utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "Pinned created successfully",
	})
}

func DeletePinned(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Parse query params
	adminIDStr := r.URL.Query().Get("admin_id")
	studentIDStr := r.URL.Query().Get("student_id")

	// Prevent deleting all records
	if adminIDStr == "" && studentIDStr == "" {
		utils.WriteError(w, http.StatusBadRequest, "Must provide at least admin_id or student_id")
		return
	}

	// Build base query dynamically
	query := "DELETE FROM pinned WHERE 1=1"
	var args []interface{}

	if adminIDStr != "" {
		adminID, err := strconv.Atoi(adminIDStr)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "Invalid admin_id")
			return
		}
		query += " AND admin_id = ?"
		args = append(args, adminID)
	}

	if studentIDStr != "" {
		studentID, err := strconv.Atoi(studentIDStr)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "Invalid student_id")
			return
		}
		query += " AND student_id = ?"
		args = append(args, studentID)
	}

	// Executes written SQL to delete pinned record
	res, err := db.ExecContext(r.Context(), query, args...)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete pinned record(s)")
		log.Println("Delete error:", err)
		return
	}

	// Get number of rows affected
	rowsAffected, _ := res.RowsAffected()

	// Error message if no rows were deleted and it was a single delete
	if adminIDStr != "" && studentIDStr != "" && rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, "No pinned records found to delete")
		return
	}

	// Writes JSON response confirming deletion & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message":       "Pinned record(s) deleted successfully",
		"rows_affected": rowsAffected,
	})
}

func GetStuAccom(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// All data being selected for this GET command
	query := `
		SELECT student_id, accommodation_id
		FROM stu_accom
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query)

	// Error message if QueryContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to obtain student accommodations")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	// Creates an empty slice to obtain results
	stuAccomList := make([]models.StudentAccommodation, 0)

	// Reads each row returned by the database
	for rows.Next() {
		var sa models.StudentAccommodation
		// Parses the current data into fields of "sa" variable
		if err := rows.Scan(&sa.StudentID, &sa.AccommodationID); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse student accommodation")
			log.Println("Row scan error:", err)
			return
		}
		// Adds the obtained data to the slice
		stuAccomList = append(stuAccomList, sa)
	}

	// Checks for errors during iteration such as network interruptions and driver errors
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Writes the slice as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, stuAccomList)
}

func CreateStuAccom(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Empty variable for request struct
	var req models.StudentAccommodation
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents extra unexpected fields
	if err := decoder.Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		log.Println("JSON decode error:", err)
		return
	}

	// Validates required fields
	if req.StudentID == 0 || req.AccommodationID == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Executes written SQL to insert student accommodation record
	_, err := db.ExecContext(r.Context(),
		"INSERT INTO stu_accom (student_id, accommodation_id) VALUES (?, ?)",
		req.StudentID, req.AccommodationID,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert student accommodation")
		log.Println("DB insert error:", err)
		return
	}

	// Writes JSON response confirming creation & sends a HTTP 201 response code
	utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "Student accommodation created successfully",
	})
}

func DeleteStuAccom(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Parse query params
	studentIDStr := r.URL.Query().Get("student_id")
	accomIDStr := r.URL.Query().Get("accommodation_id")

	// Prevent deleting all records
	if studentIDStr == "" && accomIDStr == "" {
		utils.WriteError(w, http.StatusBadRequest, "Must provide at least student_id or accommodation_id")
		return
	}

	// Build base query dynamically
	query := "DELETE FROM stu_accom WHERE 1=1"
	var args []interface{}

	// Add student_id condition if provided
	if studentIDStr != "" {
		studentID, err := strconv.Atoi(studentIDStr)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "Invalid student_id")
			return
		}
		query += " AND student_id = ?"
		args = append(args, studentID)
	}

	// Add accommodation_id condition if provided
	if accomIDStr != "" {
		accomID, err := strconv.Atoi(accomIDStr)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "Invalid accommodation_id")
			return
		}
		query += " AND accommodation_id = ?"
		args = append(args, accomID)
	}

	// Executes written SQL to delete student-accommodation record(s)
	res, err := db.ExecContext(r.Context(), query, args...)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete student accommodation record(s)")
		log.Println("Delete error:", err)
		return
	}

	// Get number of rows affected
	rowsAffected, _ := res.RowsAffected()

	// Error message if no rows were deleted and it was a single delete
	if studentIDStr != "" && accomIDStr != "" && rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, "No student accommodation records found to delete")
		return
	}

	// Writes JSON response confirming deletion & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message":       "Student accommodation record(s) deleted successfully",
		"rows_affected": rowsAffected,
	})
}

func GetStuDis(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// All data being selected for this GET command
	query := `
		SELECT student_id, disability_id
		FROM stu_dis
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query)

	// Error message if QueryContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to obtain student disabilities")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	// Creates an empty slice to obtain results
	stuDisList := make([]models.StudentDisability, 0)

	// Reads each row returned by the database
	for rows.Next() {
		var sd models.StudentDisability
		// Parses the current data into fields of "sd" variable
		if err := rows.Scan(&sd.StudentID, &sd.DisabilityID); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse student disability")
			log.Println("Row scan error:", err)
			return
		}
		// Adds the obtained data to the slice
		stuDisList = append(stuDisList, sd)
	}

	// Checks for errors during iteration such as network interruptions and driver errors
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Writes the slice as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, stuDisList)
}

func CreateStuDis(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Empty variable for request struct
	var req models.StudentDisability
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents extra unexpected fields
	if err := decoder.Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		log.Println("JSON decode error:", err)
		return
	}

	// Validates required fields
	if req.StudentID == 0 || req.DisabilityID == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Executes written SQL to insert student disability record
	_, err := db.ExecContext(r.Context(),
		"INSERT INTO stu_dis (student_id, disability_id) VALUES (?, ?)",
		req.StudentID, req.DisabilityID,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert student disability")
		log.Println("DB insert error:", err)
		return
	}

	// Writes JSON response confirming creation & sends a HTTP 201 response code
	utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "Student disability created successfully",
	})
}

func DeleteStuDis(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Parse query params
	studentIDStr := r.URL.Query().Get("student_id")
	disabilityIDStr := r.URL.Query().Get("disability_id")

	// Prevent deleting all records
	if studentIDStr == "" && disabilityIDStr == "" {
		utils.WriteError(w, http.StatusBadRequest, "Must provide at least student_id or disability_id")
		return
	}

	// Build base query dynamically
	query := "DELETE FROM stu_dis WHERE 1=1"
	var args []interface{}

	if studentIDStr != "" {
		studentID, err := strconv.Atoi(studentIDStr)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "Invalid student_id")
			return
		}
		query += " AND student_id = ?"
		args = append(args, studentID)
	}

	if disabilityIDStr != "" {
		disabilityID, err := strconv.Atoi(disabilityIDStr)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "Invalid disability_id")
			return
		}
		query += " AND disability_id = ?"
		args = append(args, disabilityID)
	}

	// Executes written SQL to delete student_disability record
	res, err := db.ExecContext(r.Context(), query, args...)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete student_disability record(s)")
		log.Println("Delete error:", err)
		return
	}

	// Get number of rows affected
	rowsAffected, _ := res.RowsAffected()

	// Error message if no rows were deleted and it was a single delete
	if studentIDStr != "" && disabilityIDStr != "" && rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, "No stu_dis records found to delete")
		return
	}

	// Writes JSON response confirming deletion & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message":       "Student disability record(s) deleted successfully",
		"rows_affected": rowsAffected,
	})
}

func GetPocAdmin(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// All data being selected for this GET command
	query := `
		SELECT point_of_contact_id, admin_id
		FROM poc_admin
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query)

	// Error message if QueryContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to obtain POC admins")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	// Creates an empty slice to obtain results
	pocAdmins := make([]models.PocAdmin, 0)

	// Reads each row returned by the database
	for rows.Next() {
		var pa models.PocAdmin
		// Parses the current data into fields of "pa" variable
		if err := rows.Scan(&pa.PointOfContactID, &pa.AdminID); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse POC admin")
			log.Println("Row scan error:", err)
			return
		}
		// Adds the obtained data to the slice
		pocAdmins = append(pocAdmins, pa)
	}

	// Checks for errors during iteration such as network interruptions and driver errors
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Writes the slice as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, pocAdmins)
}

func CreatePocAdmin(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Empty variable for request struct
	var req models.PocAdmin
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents extra unexpected fields
	if err := decoder.Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		log.Println("JSON decode error:", err)
		return
	}

	// Validates required fields
	if req.PointOfContactID == 0 || req.AdminID == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Executes written SQL to insert POC admin record
	_, err := db.ExecContext(r.Context(),
		"INSERT INTO poc_admin (point_of_contact_id, admin_id) VALUES (?, ?)",
		req.PointOfContactID, req.AdminID,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert POC admin")
		log.Println("DB insert error:", err)
		return
	}

	// Writes JSON response confirming creation & sends a HTTP 201 response code
	utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "POC Admin created successfully",
	})
}

func DeletePocAdmin(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Parse query params
	pocIDStr := r.URL.Query().Get("point_of_contact_id")
	adminIDStr := r.URL.Query().Get("admin_id")

	// Prevent deleting all records
	if pocIDStr == "" && adminIDStr == "" {
		utils.WriteError(w, http.StatusBadRequest, "Must provide at least point_of_contact_id or admin_id")
		return
	}

	// Build base query dynamically
	query := "DELETE FROM poc_admin WHERE 1=1"
	var args []interface{}

	// Add point_of_contact_id condition if provided
	if pocIDStr != "" {
		pocID, err := strconv.Atoi(pocIDStr)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "Invalid point_of_contact_id")
			return
		}
		query += " AND point_of_contact_id = ?"
		args = append(args, pocID)
	}

	// Add admin_id condition if provided
	if adminIDStr != "" {
		adminID, err := strconv.Atoi(adminIDStr)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "Invalid admin_id")
			return
		}
		query += " AND admin_id = ?"
		args = append(args, adminID)
	}

	// Executes written SQL to delete poc-admin record(s)
	res, err := db.ExecContext(r.Context(), query, args...)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete poc-admin record(s)")
		log.Println("Delete error:", err)
		return
	}

	// Get number of rows affected
	rowsAffected, _ := res.RowsAffected()

	// Error message if no rows were deleted and it was a single delete
	if pocIDStr != "" && adminIDStr != "" && rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, "No pocadmin records found to delete")
		return
	}

	// Writes JSON response confirming deletion & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message":       "POC-admin record(s) deleted successfully",
		"rows_affected": rowsAffected,
	})
}
