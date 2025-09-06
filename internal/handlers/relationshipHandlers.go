package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"
	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func GetPinned(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT admin_id, student_id
		FROM pinned
	`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pinnedList []models.Pinned

	for rows.Next() {
		var p models.Pinned
		if err := rows.Scan(&p.AdminID, &p.StudentID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pinnedList = append(pinnedList, p)
	}

	jsonBytes, err := json.MarshalIndent(pinnedList, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

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

type CreatePinnedRequest struct {
	AdminID   int `json:"admin_id"`
	StudentID int `json:"student_id"`
}

func CreatePinned(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreatePinnedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(
		"INSERT INTO pinned (admin_id, student_id) VALUES (?, ?)",
		req.AdminID, req.StudentID,
	)
	if err != nil {
		http.Error(w, "Failed to insert pinned: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Pinned created successfully",
	})
}

func DeletePinned(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")

	vars := mux.Vars(r)

	adminIDStr := vars["admin_id"]
	studentIDStr := vars["student_id"]

	adminID, err := strconv.Atoi(adminIDStr)
	if err != nil {
		http.Error(w, "Invalid admin ID", http.StatusBadRequest)
		return
	}

	studentID, err := strconv.Atoi(studentIDStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	// Delete the pinned record
	_, err = db.Exec(`DELETE FROM pinned WHERE admin_id = ? AND student_id = ?`, adminID, studentID)
	if err != nil {
		http.Error(w, "Failed to delete pinned: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Pinned deleted successfully"})
}

func GetStuAccom(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, accommodation_id
		FROM stu_accom
	`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var stuAccomList []models.StudentAccommodation

	for rows.Next() {
		var sa models.StudentAccommodation
		if err := rows.Scan(&sa.ID, &sa.AccommodationID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		stuAccomList = append(stuAccomList, sa)
	}

	jsonBytes, err := json.MarshalIndent(stuAccomList, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

type CreateStuAccomRequest struct {
	ID              int `json:"id"`               // student_id
	AccommodationID int `json:"accommodation_id"` // accommodation_id
}

func CreateStuAccom(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateStuAccomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(
		"INSERT INTO stu_accom (id, accommodation_id) VALUES (?, ?)",
		req.ID, req.AccommodationID,
	)
	if err != nil {
		http.Error(w, "Failed to insert stu_accom: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Student accommodation created successfully",
	})
}

func DeleteStuAccom(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")

	vars := mux.Vars(r)

	studentIDStr := vars["id"]
	accomIDStr := vars["accommodation_id"]

	studentID, err := strconv.Atoi(studentIDStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	accomID, err := strconv.Atoi(accomIDStr)
	if err != nil {
		http.Error(w, "Invalid accommodation ID", http.StatusBadRequest)
		return
	}

	// Delete the student accommodation record
	_, err = db.Exec(`DELETE FROM stu_accom WHERE id = ? AND accommodation_id = ?`, studentID, accomID)
	if err != nil {
		http.Error(w, "Failed to delete stu_accom: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Student accommodation deleted successfully"})
}

func GetStuDis(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, disability_id
		FROM stu_dis
	`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var stuDisList []models.StudentDisability

	for rows.Next() {
		var sd models.StudentDisability
		if err := rows.Scan(&sd.ID, &sd.DisabilityID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		stuDisList = append(stuDisList, sd)
	}

	jsonBytes, err := json.MarshalIndent(stuDisList, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

type CreateStuDisRequest struct {
	ID           int `json:"id"`            // student id
	DisabilityID int `json:"disability_id"` // disability id
}

func CreateStuDis(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateStuDisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(
		"INSERT INTO stu_dis (id, disability_id) VALUES (?, ?)",
		req.ID, req.DisabilityID,
	)
	if err != nil {
		http.Error(w, "Failed to insert stu_dis: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "StuDis created successfully",
	})
}

func DeleteStuDis(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")

	vars := mux.Vars(r)

	idStr := vars["id"]
	disabilityIDStr := vars["disability_id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	disabilityID, err := strconv.Atoi(disabilityIDStr)
	if err != nil {
		http.Error(w, "Invalid disability ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`DELETE FROM stu_dis WHERE id = ? AND disability_id = ?`, id, disabilityID)
	if err != nil {
		http.Error(w, "Failed to delete stu_dis: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "StuDis deleted successfully"})
}

func GetPocAdmin(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	query := `
        SELECT activity_id, admin_id
        FROM poc_adm
    `

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pocAdmins []models.PocAdmin

	for rows.Next() {
		var pa models.PocAdmin
		if err := rows.Scan(&pa.ActivityID, &pa.AdminID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pocAdmins = append(pocAdmins, pa)
	}

	jsonBytes, err := json.MarshalIndent(pocAdmins, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

type CreatePOCAdminRequest struct {
	ActivityID int `json:"activity_id"`
	AdminID    int `json:"admin_id"`
}

func CreatePocAdmin(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreatePOCAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(
		"INSERT INTO poc_adm (activity_id, admin_id) VALUES (?, ?)",
		req.ActivityID, req.AdminID,
	)
	if err != nil {
		http.Error(w, "Failed to insert poc_adm: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "POC Admin created successfully",
	})
}

func DeletePocAdmin(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")

	vars := mux.Vars(r)

	activityIDStr := vars["activity_id"]
	adminIDStr := vars["admin_id"]

	activityID, err := strconv.Atoi(activityIDStr)
	if err != nil {
		http.Error(w, "Invalid activity ID", http.StatusBadRequest)
		return
	}

	adminID, err := strconv.Atoi(adminIDStr)
	if err != nil {
		http.Error(w, "Invalid admin ID", http.StatusBadRequest)
		return
	}

	// Delete the POC admin record
	_, err = db.Exec(`DELETE FROM poc_adm WHERE activity_id = ? AND admin_id = ?`, activityID, adminID)
	if err != nil {
		http.Error(w, "Failed to delete poc_adm: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "POC Admin deleted successfully"})
}
