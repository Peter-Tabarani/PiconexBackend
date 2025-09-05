package routes

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UpdatePersonalDocumentationRequest struct {
	Date       string `json:"date"`        // "YYYY-MM-DD"
	Time       string `json:"time"`        // "HH:MM:SS"
	FileBase64 string `json:"file_base64"` // base64-encoded file blob
	AdminID    int    `json:"admin_id"`    // personal_documentation.id
}

func UpdatePersonalDocumentationByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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
