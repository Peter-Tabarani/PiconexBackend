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

func GetAccommodations(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error Message For Any Request That Is Not GET
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	// All Data Being Selected For This GET Command
	query := `
        SELECT accommodation_id, name, description
        FROM accommodation
    `
	// Executes Written SQL
	rows, err := db.QueryContext(r.Context(), query)

	// Error Message If QueryContext Fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed To Fetch Accommodations")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	// Creates An Empty Slice To Obtain Results
	accommodations := make([]models.Accommodation, 0)

	// Reads Each Row Returned By The Database
	for rows.Next() {

		// Empty Variable For Accommodation Struct
		var am models.Accommodation

		// Reads The Current Data Into Fields Of (am) Variable
		if err := rows.Scan(&am.Accommodation_ID, &am.Name, &am.Description); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed To Read Accommodations")
			log.Println("Row scan error:", err)
			return
		}

		// Adds The Obtained Data To The Slice
		accommodations = append(accommodations, am)
	}

	// Checks For Errors During Iteration Such As Network Interruptions and Driver Errors
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Error During Iteration")
		log.Println("Rows error:", err)
		return
	}

	// Writes The Slice As JSON & Sends A HTTP 200 Response Code
	utils.WriteJSON(w, http.StatusOK, accommodations)
}

func GetAccommodationByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing accommodation ID")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid accommodation ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	query := `
        SELECT accommodation_id, name, description
        FROM accommodation
        WHERE accommodation_id = ?
    `

	var accom models.Accommodation
	err = db.QueryRowContext(r.Context(), query, id).Scan(
		&accom.Accommodation_ID, &accom.Name, &accom.Description,
	)
	if err == sql.ErrNoRows {
		utils.WriteError(w, http.StatusNotFound, "Accommodation not found")
		return
	} else if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch accommodation")
		log.Println("DB query error:", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, accom)
}

func GetAccommodationsByStudentID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]

	studentID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid student ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	query := `
		SELECT
			a.accommodation_id, a.name, a.description
		FROM stu_accom sa
		JOIN accommodation a ON sa.accommodation_id = a.accommodation_id
		WHERE sa.id = ?
	`

	rows, err := db.QueryContext(r.Context(), query, studentID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch accommodations for student")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	accommodations := make([]models.Accommodation, 0)
	for rows.Next() {
		var a models.Accommodation
		if err := rows.Scan(&a.Accommodation_ID, &a.Name, &a.Description); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse accommodations")
			log.Println("Row scan error:", err)
			return
		}
		accommodations = append(accommodations, a)
	}

	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Error reading accommodations")
		log.Println("Rows error:", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, accommodations)
}

func CreateAccommodation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var a models.Accommodation
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&a); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		log.Println("JSON decode error:", err)
		return
	}

	if a.Name == "" || a.Description == "" {
		utils.WriteError(w, http.StatusBadRequest, "Missing required fields: name or description")
		return
	}

	res, err := db.ExecContext(r.Context(),
		"INSERT INTO accommodation (name, description) VALUES (?, ?)",
		a.Name, a.Description,
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert accommodation")
		log.Println("DB insert error:", err)
		return
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get last insert ID")
		log.Println("LastInsertId error:", err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message":          "Accommodation created successfully",
		"accommodation_id": lastID,
	})
}

func DeleteAccommodation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	accomID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid accommodation ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	res, err := db.ExecContext(r.Context(), "DELETE FROM accommodation WHERE accommodation_id = ?", accomID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete accommodation")
		log.Println("DB delete error:", err)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get rows affected")
		log.Println("RowsAffected error:", err)
		return
	}

	if rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, "Accommodation not found")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Accommodation deleted successfully",
	})
}

func UpdateAccommodation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	accommodationID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid accommodation ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	var a models.Accommodation
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&a); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		log.Println("JSON decode error:", err)
		return
	}

	if a.Name == "" || a.Description == "" {
		utils.WriteError(w, http.StatusBadRequest, "Missing required fields: name or description")
		return
	}

	res, err := db.ExecContext(
		r.Context(),
		`UPDATE accommodation SET name = ?, description = ? WHERE accommodation_id = ?`,
		a.Name, a.Description, accommodationID,
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update accommodation")
		log.Println("DB update error:", err)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to check update result")
		log.Println("RowsAffected error:", err)
		return
	}

	if rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, "Accommodation not found")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Accommodation updated successfully",
	})
}
