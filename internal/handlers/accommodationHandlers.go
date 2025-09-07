package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func GetAccommodations(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	query := `
        SELECT accommodation_id, name, description
        FROM accommodation
    `

	rows, err := db.QueryContext(r.Context(), query)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch accommodations")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	var accommodations []models.Accommodation
	for rows.Next() {
		var am models.Accommodation
		if err := rows.Scan(&am.Accommodation_ID, &am.Name, &am.Description); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse accommodations")
			log.Println("Row scan error:", err)
			return
		}
		accommodations = append(accommodations, am)
	}

	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Error reading accommodations")
		log.Println("Rows error:", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, accommodations)
}

func GetAccommodationByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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

	var accommodations []models.Accommodation
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

type CreateAccommodationRequest struct {
	Accommodation models.Accommodation `json:"accommodation"`
}

func CreateAccommodation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateAccommodationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	res, err := db.Exec(
		"INSERT INTO accommodation (name, description) VALUES (?, ?)",
		req.Accommodation.Name, req.Accommodation.Description,
	)
	if err != nil {
		http.Error(w, "Failed to insert accommodation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	lastID, _ := res.LastInsertId()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":          "Accommodation created successfully",
		"accommodation_id": lastID,
	})
}

// PROBLEM: When deleting an accommodation that isn't there, it submits a success message
func DeleteAccommodation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid accommodation ID", http.StatusBadRequest)
		return
	}

	// Step 1: Nullify any foreign keys referencing this accommodation (if needed)
	// For example, stu_accom references accommodation_id
	_, err = db.Exec(`UPDATE stu_accom SET accommodation_id = NULL WHERE accommodation_id = ?`, id)
	if err != nil {
		http.Error(w, "Failed to update stu_accom: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 2: Delete from accommodation
	_, err = db.Exec(`DELETE FROM accommodation WHERE accommodation_id = ?`, id)
	if err != nil {
		http.Error(w, "Failed to delete accommodation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Accommodation deleted successfully"})
}

// FAILING
func UpdateAccommodation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var a models.Accommodation
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Update only fields that were sent
	_, err = db.Exec(`
		UPDATE accommodation
		SET name = ?, description = ?
		WHERE accommodation_id = ?`, a.Name, a.Description, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Accommodation updated successfully")
}
