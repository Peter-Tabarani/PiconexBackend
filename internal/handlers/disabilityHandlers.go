package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"
	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func GetDisabilities(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT
			ds.disability_id, ds.name, ds.description
		FROM disability ds
	`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var disabilities []models.Disability

	for rows.Next() {
		var ds models.Disability
		err := rows.Scan(
			&ds.Disability_ID, &ds.Name, &ds.Description,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		disabilities = append(disabilities, ds)
	}

	jsonBytes, err := json.MarshalIndent(disabilities, "", "    ") // Pretty print with 4 spaces indent
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func GetDisabilityByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["disability_id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid disability ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT disability_id, name, description
		FROM disability
		WHERE disability_id = ?
	`

	var dis models.Disability
	err = db.QueryRow(query, id).Scan(&dis.Disability_ID, &dis.Name, &dis.Description)
	if err == sql.ErrNoRows {
		http.Error(w, "Disability not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.MarshalIndent(dis, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func GetDisabilitiesByStudentID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	studentID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT
			d.disability_id, d.name, d.description
		FROM stu_dis sd
		JOIN disability d ON sd.disability_id = d.disability_id
		WHERE sd.id = ?
	`

	rows, err := db.Query(query, studentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var disabilities []models.Disability

	for rows.Next() {
		var d models.Disability
		if err := rows.Scan(&d.Disability_ID, &d.Name, &d.Description); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		disabilities = append(disabilities, d)
	}

	if len(disabilities) == 0 {
		http.Error(w, "No disabilities found for the student", http.StatusNotFound)
		return
	}

	jsonBytes, err := json.MarshalIndent(disabilities, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

type CreateDisabilityRequest struct {
	Disability models.Disability `json:"disability"`
}

func CreateDisability(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateDisabilityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	res, err := db.Exec(
		"INSERT INTO disability (name, description) VALUES (?, ?)",
		req.Disability.Name, req.Disability.Description,
	)
	if err != nil {
		http.Error(w, "Failed to insert disability: "+err.Error(), http.StatusInternalServerError)
		return
	}

	lastID, _ := res.LastInsertId()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "Disability created successfully",
		"disability_id": lastID,
	})
}

func DeleteDisability(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid disability ID", http.StatusBadRequest)
		return
	}

	// Step 1: Nullify or remove references in stu_dis
	_, err = db.Exec(`DELETE FROM stu_dis WHERE disability_id = ?`, id)
	if err != nil {
		http.Error(w, "Failed to update stu_dis: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 2: Delete from disability
	_, err = db.Exec(`DELETE FROM disability WHERE disability_id = ?`, id)
	if err != nil {
		http.Error(w, "Failed to delete disability: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Disability deleted successfully"})
}

func UpdateDisability(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var d models.Disability
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Update only fields that were sent
	_, err = db.Exec(`
		UPDATE disability
		SET name = ?, description = ?
		WHERE disability_id = ?`, d.Name, d.Description, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Disability updated successfully")
}
