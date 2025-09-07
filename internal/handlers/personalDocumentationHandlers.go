package handlers

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"
	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func GetPersonalDocumentations(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT
		ac.activity_id, a.id, ac.date, ac.time, d.file
	FROM personal_documentation pd
	JOIN activity ac ON pd.activity_id = ac.activity_id
	JOIN admin a ON pd.id = a.id
	JOIN documentation d ON d.activity_id = pd.activity_id
	`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var personal_documentations []models.Personal_Documentation

	for rows.Next() {
		var pd models.Personal_Documentation
		err := rows.Scan(
			&pd.Activity_ID, &pd.ID, &pd.Date, &pd.Time, &pd.File,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		personal_documentations = append(personal_documentations, pd)
	}

	jsonBytes, err := json.MarshalIndent(personal_documentations, "", "    ") // Pretty print with 4 spaces indent
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// FAILING
func GetPersonalDocumentationByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["activity_id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid activity ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT
			pd.activity_id,
			pd.id,
			a.date,
			a.time,
			d.file
		FROM personal_documentation pd
		JOIN documentation d ON pd.activity_id = d.activity_id
		JOIN activity a ON pd.activity_id = a.activity_id
		WHERE pd.activity_id = ?
	`

	var pd models.Personal_Documentation
	err = db.QueryRow(query, id).Scan(&pd.Activity_ID, &pd.ID, &pd.Date, &pd.Time, &pd.File)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Personal documentation not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	jsonBytes, err := json.MarshalIndent(pd, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// PROBLEM: Inconsistent parameter names
type CreatePersonalDocumentationRequest struct {
	AdminID    int    `json:"admin_id"`    // maps to personal_documentation.id
	Date       string `json:"date"`        // format: "YYYY-MM-DD"
	Time       string `json:"time"`        // format: "HH:MM:SS"
	FileBase64 string `json:"file_base64"` // base64-encoded file blob
}

// FAILING: Gave an success message, added to activity and documentation, but didn't add to personal-documentation
func CreatePersonalDocumentation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreatePersonalDocumentationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	fileBytes, err := base64.StdEncoding.DecodeString(req.FileBase64)
	if err != nil {
		http.Error(w, "Failed to decode file_base64: "+err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert into activity
	activityQuery := `INSERT INTO activity (date, time) VALUES (?, ?)`
	res, err := tx.Exec(activityQuery, req.Date, req.Time)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to insert into activity: "+err.Error(), http.StatusInternalServerError)
		return
	}

	activityID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to get activity ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert into documentation
	documentationQuery := `INSERT INTO documentation (activity_id, file) VALUES (?, ?)`
	_, err = tx.Exec(documentationQuery, activityID, fileBytes)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to insert into documentation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert into personal_documentation
	personalDocQuery := `INSERT INTO personal_documentation (activity_id, id) VALUES (?, ?)`
	_, err = tx.Exec(personalDocQuery, activityID, req.AdminID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to insert into personal_documentation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "Personal documentation created successfully",
		"activity_id": activityID,
		"admin_id":    req.AdminID,
	})
}

// FAILING
func DeletePersonalDocumentation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	activityIDStr := vars["activity_id"]
	activityID, err := strconv.Atoi(activityIDStr)
	if err != nil {
		http.Error(w, "Invalid activity ID", http.StatusBadRequest)
		return
	}

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete from activity (this will cascade to documentation & personal_documentation)
	deleteQuery := `DELETE FROM activity WHERE activity_id=?`
	_, err = tx.Exec(deleteQuery, activityID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to delete personal documentation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Personal documentation deleted successfully"}`))
}

type UpdatePersonalDocumentationRequest struct {
	Date       string `json:"date"`        // "YYYY-MM-DD"
	Time       string `json:"time"`        // "HH:MM:SS"
	FileBase64 string `json:"file_base64"` // base64-encoded file blob
	AdminID    int    `json:"admin_id"`    // personal_documentation.id
}

// FAILING
func UpdatePersonalDocumentation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	activityIDStr := vars["activity_id"]
	activityID, err := strconv.Atoi(activityIDStr)
	if err != nil {
		http.Error(w, "Invalid activity ID", http.StatusBadRequest)
		return
	}

	var req UpdatePersonalDocumentationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Decode base64 file data
	fileBytes, err := base64.StdEncoding.DecodeString(req.FileBase64)
	if err != nil {
		http.Error(w, "Invalid base64 file data: "+err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update activity (date, time)
	activityUpdate := `UPDATE activity SET date=?, time=? WHERE activity_id=?`
	_, err = tx.Exec(activityUpdate, req.Date, req.Time, activityID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update activity: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update documentation (file)
	docUpdate := `UPDATE documentation SET file=? WHERE activity_id=?`
	_, err = tx.Exec(docUpdate, fileBytes, activityID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update documentation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update personal_documentation (admin_id)
	personalDocUpdate := `UPDATE personal_documentation SET id=? WHERE activity_id=?`
	_, err = tx.Exec(personalDocUpdate, req.AdminID, activityID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update personal documentation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Personal documentation updated successfully"})
}
