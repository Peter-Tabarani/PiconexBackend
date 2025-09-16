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
	// Error message for any request that is not GET
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// All data being selected for this GET command
	query := `
        SELECT accommodation_id, name, description
        FROM accommodation
    `
	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query)

	// Error message if QueryContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to obtain accommodations")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	// Creates an empty slice to obtain results
	accommodations := make([]models.Accommodation, 0)

	// Reads each row returned by the database
	for rows.Next() {
		var a models.Accommodation
		// Parses the current data into fields of "a" variable
		if err := rows.Scan(&a.Accommodation_ID, &a.Name, &a.Description); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse accommodations")
			log.Println("Row scan error:", err)
			return
		}

		// Adds the obtained data to the slice
		accommodations = append(accommodations, a)
	}

	// Checks for errors during iteration such as network interruptions and driver errors
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Writes the slice as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, accommodations)
}

func GetAccommodationByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not GET
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr, ok := vars["accommodation_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing accommodation ID")
		return
	}

	// Converts the "id" string to an integer
	accommodationID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid accommodation ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// All data being selected for this GET command
	query := `
        SELECT accommodation_id, name, description
        FROM accommodation
        WHERE accommodation_id = ?
    `

	// Empty variable for accommodation struct
	var a models.Accommodation

	// Executes query
	err = db.QueryRowContext(r.Context(), query, accommodationID).Scan(
		&a.Accommodation_ID, &a.Name, &a.Description,
	)

	// Error message if no rows are found
	if err == sql.ErrNoRows {
		utils.WriteError(w, http.StatusNotFound, "Accommodation not found")
		return
		// Error message if QueryRowContext or scan fails
	} else if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch accommodation")
		log.Println("DB query error:", err)
		return
	}

	// Writes the struct as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, a)
}

func GetAccommodationsByStudentID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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
	studentID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid student ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// All data being selected for this GET command
	query := `
		SELECT a.accommodation_id, a.name, a.description
		FROM stu_accom sa
		JOIN accommodation a ON sa.accommodation_id = a.accommodation_id
		WHERE sa.id = ?
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query, studentID)

	// Error message if QueryContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to obtain accommodations for student")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	// Creates an empty slice to obtain results
	accommodations := make([]models.Accommodation, 0)

	// Reads each row returned by the database
	for rows.Next() {
		var a models.Accommodation
		// Parses the current data into fields of "a" variable
		if err := rows.Scan(&a.Accommodation_ID, &a.Name, &a.Description); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse accommodations")
			log.Println("Row scan error:", err)
			return
		}

		// Adds the obtained data to the slice
		accommodations = append(accommodations, a)
	}

	// Checks for errors during iteration such as network interruptions and driver errors
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Writes the slice as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, accommodations)
}

func CreateAccommodation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not POST
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Empty variable for accommodation struct
	var a models.Accommodation

	// Decodes JSON body from the request into "a" variable
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents extra unexpected fields
	if err := decoder.Decode(&a); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		log.Println("JSON decode error:", err)
		return
	}

	// Validates required fields
	if a.Name == "" {
		utils.WriteError(w, http.StatusBadRequest, "Missing required fields: name")
		return
	}

	// Executes written SQL to insert a new accommodation
	res, err := db.ExecContext(r.Context(),
		"INSERT INTO accommodation (name, description) VALUES (?, ?)",
		a.Name, a.Description,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to insert accommodation")
		log.Println("DB insert error:", err)
		return
	}

	// Gets the ID of the newly inserted accommodation
	lastID, err := res.LastInsertId()

	// Error message if LastInsertId fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get last insert ID")
		log.Println("LastInsertId error:", err)
		return
	}

	// Writes JSON response including the new ID & sends a HTTP 201 response code
	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message":          "Accommodation created successfully",
		"accommodation_id": lastID,
	})
}

func DeleteAccommodation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not DELETE
	if r.Method != http.MethodDelete {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr, ok := vars["accommodation_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing accommodation ID")
		return
	}

	// Converts the "accommodation_id" string to an integer
	accommodationID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid accommodation ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// Executes written SQL to delete the accommodation
	res, err := db.ExecContext(r.Context(),
		"DELETE FROM accommodation WHERE accommodation_id = ?",
		accommodationID,
	)

	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete accommodation")
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
		utils.WriteError(w, http.StatusNotFound, "Accommodation not found")
		return
	}

	// Writes JSON response confirming deletion & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Accommodation deleted successfully",
	})
}

func UpdateAccommodation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Error message for any request that is not PUT
	if r.Method != http.MethodPut {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extracts path variables from the request
	vars := mux.Vars(r)

	// Reads the "accommodation_id" value from the path variables
	idStr := vars["accommodation_id"]

	// Converts the "accommodation_id" string to an integer
	accommodationID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid accommodation ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// Empty variable for accommodation struct
	var a models.Accommodation

	// Decodes JSON body from the request into "a"
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevents extra unexpected fields
	if err := decoder.Decode(&a); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		log.Println("JSON decode error:", err)
		return
	}

	// Validates required fields
	if a.Name == "" || a.Description == "" {
		utils.WriteError(w, http.StatusBadRequest, "Missing required fields: name or description")
		return
	}

	// Executes written SQL to update the accommodation
	res, err := db.ExecContext(
		r.Context(),
		`UPDATE accommodation SET name = ?, description = ? WHERE accommodation_id = ?`,
		a.Name, a.Description, accommodationID,
	)
	// Error message if ExecContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update accommodation")
		log.Println("DB update error:", err)
		return
	}

	// Gets the number of rows affected by the update
	rowsAffected, err := res.RowsAffected()

	// Error message if RowsAffected fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to check update result")
		log.Println("RowsAffected error:", err)
		return
	}

	// Error message if no rows were updated
	if rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, "Accommodation not found")
		return
	}

	// Writes JSON response confirming update & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Accommodation updated successfully",
	})
}
