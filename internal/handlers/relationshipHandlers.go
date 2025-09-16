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
	// Error message for any request that is not GET
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// All data being selected for this GET command
	query := `
		SELECT admin_id, student_id
		FROM pinned
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query)
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

	// Checks for errors during iteration
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Writes the slice as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, pinnedList)
}

func GetPinnedByAdminID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not GET
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing admin ID")
		return
	}

	// Converts the "id" string to an integer
	adminID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid admin ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// All data being selected for this GET command
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

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query, adminID)
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
		var discardAdminID int // used because SQL returns admin_id first

		// Parses the current row into student struct
		if err := rows.Scan(
			&discardAdminID,
			&s.ID, &s.FirstName, &s.PreferredName, &s.MiddleName, &s.LastName,
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
	// Error message for any request that is not POST
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

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
	// Error message for any request that is not DELETE
	if r.Method != http.MethodDelete {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extracts path variables from the request
	vars := mux.Vars(r)
	adminIDStr, ok1 := vars["admin_id"]
	studentIDStr, ok2 := vars["student_id"]
	if !ok1 || !ok2 {
		utils.WriteError(w, http.StatusBadRequest, "Missing admin or student ID")
		return
	}

	// Converts path variables to integers
	adminID, err := strconv.Atoi(adminIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid admin ID")
		log.Println("Invalid admin ID parse error:", err)
		return
	}
	studentID, err := strconv.Atoi(studentIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid student ID")
		log.Println("Invalid student ID parse error:", err)
		return
	}

	// Executes written SQL to delete pinned record
	res, err := db.ExecContext(r.Context(),
		"DELETE FROM pinned WHERE admin_id = ? AND student_id = ?",
		adminID, studentID,
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete pinned record")
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
	if rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, "Pinned record not found")
		return
	}

	// Writes JSON response confirming deletion & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Pinned deleted successfully",
	})
}

func GetStuAccom(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not GET
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// All data being selected for this GET command
	query := `
		SELECT id, accommodation_id
		FROM stu_accom
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query)
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
		if err := rows.Scan(&sa.ID, &sa.AccommodationID); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse student accommodation")
			log.Println("Row scan error:", err)
			return
		}
		// Adds the obtained data to the slice
		stuAccomList = append(stuAccomList, sa)
	}

	// Checks for errors during iteration
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Writes the slice as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, stuAccomList)
}

func CreateStuAccom(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not POST
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

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
	if req.ID == 0 || req.AccommodationID == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Executes written SQL to insert student accommodation record
	_, err := db.ExecContext(r.Context(),
		"INSERT INTO stu_accom (id, accommodation_id) VALUES (?, ?)",
		req.ID, req.AccommodationID,
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
	// Error message for any request that is not DELETE
	if r.Method != http.MethodDelete {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extracts path variables from the request
	vars := mux.Vars(r)
	studentIDStr, ok1 := vars["id"]
	accomIDStr, ok2 := vars["accommodation_id"]
	if !ok1 || !ok2 {
		utils.WriteError(w, http.StatusBadRequest, "Missing student or accommodation ID")
		return
	}

	// Converts path variables to integers
	studentID, err := strconv.Atoi(studentIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid student ID")
		log.Println("Invalid student ID parse error:", err)
		return
	}
	accomID, err := strconv.Atoi(accomIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid accommodation ID")
		log.Println("Invalid accommodation ID parse error:", err)
		return
	}

	// Executes written SQL to delete student accommodation record
	res, err := db.ExecContext(r.Context(),
		"DELETE FROM stu_accom WHERE id = ? AND accommodation_id = ?",
		studentID, accomID,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete student accommodation")
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
	if rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, "Student accommodation record not found")
		return
	}

	// Writes JSON response confirming deletion & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Student accommodation deleted successfully",
	})
}

func GetStuDis(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not GET
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// All data being selected for this GET command
	query := `
		SELECT id, disability_id
		FROM stu_dis
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query)
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
		if err := rows.Scan(&sd.ID, &sd.DisabilityID); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse student disability")
			log.Println("Row scan error:", err)
			return
		}
		// Adds the obtained data to the slice
		stuDisList = append(stuDisList, sd)
	}

	// Checks for errors during iteration
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Writes the slice as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, stuDisList)
}

func CreateStuDis(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not POST
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

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
	if req.ID == 0 || req.DisabilityID == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Executes written SQL to insert student disability record
	_, err := db.ExecContext(r.Context(),
		"INSERT INTO stu_dis (id, disability_id) VALUES (?, ?)",
		req.ID, req.DisabilityID,
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
	// Error message for any request that is not DELETE
	if r.Method != http.MethodDelete {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extracts path variables from the request
	vars := mux.Vars(r)
	studentIDStr, ok1 := vars["id"]
	disabilityIDStr, ok2 := vars["disability_id"]
	if !ok1 || !ok2 {
		utils.WriteError(w, http.StatusBadRequest, "Missing student or disability ID")
		return
	}

	// Converts path variables to integers
	studentID, err := strconv.Atoi(studentIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid student ID")
		log.Println("Invalid student ID parse error:", err)
		return
	}
	disabilityID, err := strconv.Atoi(disabilityIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid disability ID")
		log.Println("Invalid disability ID parse error:", err)
		return
	}

	// Executes written SQL to delete student disability record
	res, err := db.ExecContext(r.Context(),
		"DELETE FROM stu_dis WHERE id = ? AND disability_id = ?",
		studentID, disabilityID,
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete student disability")
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
	if rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, "Student disability record not found")
		return
	}

	// Writes JSON response confirming deletion & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Student disability deleted successfully",
	})
}

func GetPocAdmin(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not GET
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// All data being selected for this GET command
	query := `
		SELECT activity_id, admin_id
		FROM poc_adm
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query)
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
		if err := rows.Scan(&pa.ActivityID, &pa.AdminID); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse POC admin")
			log.Println("Row scan error:", err)
			return
		}
		// Adds the obtained data to the slice
		pocAdmins = append(pocAdmins, pa)
	}

	// Checks for errors during iteration
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Writes the slice as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, pocAdmins)
}

func CreatePocAdmin(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not POST
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

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
	if req.ActivityID == 0 || req.AdminID == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Executes written SQL to insert POC admin record
	_, err := db.ExecContext(r.Context(),
		"INSERT INTO poc_adm (activity_id, admin_id) VALUES (?, ?)",
		req.ActivityID, req.AdminID,
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
	// Error message for any request that is not DELETE
	if r.Method != http.MethodDelete {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extracts path variables from the request
	vars := mux.Vars(r)
	activityIDStr, ok1 := vars["activity_id"]
	adminIDStr, ok2 := vars["id"]
	if !ok1 || !ok2 {
		utils.WriteError(w, http.StatusBadRequest, "Missing activity or admin ID")
		return
	}

	// Converts path variables to integers
	activityID, err := strconv.Atoi(activityIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid activity ID")
		log.Println("Invalid activity ID parse error:", err)
		return
	}
	adminID, err := strconv.Atoi(adminIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid admin ID")
		log.Println("Invalid admin ID parse error:", err)
		return
	}

	// Executes written SQL to delete POC admin record
	res, err := db.ExecContext(r.Context(),
		"DELETE FROM poc_adm WHERE activity_id = ? AND admin_id = ?",
		activityID, adminID,
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete POC admin")
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
	if rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, "POC admin record not found")
		return
	}

	// Writes JSON response confirming deletion & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "POC Admin deleted successfully",
	})
}
