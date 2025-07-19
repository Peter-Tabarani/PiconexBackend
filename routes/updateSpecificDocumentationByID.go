package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UpdateSpecificDocumentationRequest struct {
	File      []byte `json:"file"`
	DocType   string `json:"doc_type"` // enum('trad','al','in')
	StudentID int    `json:"student_id"`
}

func UpdateSpecificDocumentationByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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

	var req UpdateSpecificDocumentationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	// First update documentation table (file)
	docUpdate := `UPDATE documentation SET file=? WHERE activity_id=?`
	_, err = tx.Exec(docUpdate, req.File, activityID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update documentation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Then update specific_documentation table (doc_type, student_id)
	specDocUpdate := `
		UPDATE specific_documentation SET
			doc_type=?, id=?
		WHERE activity_id=?
	`
	_, err = tx.Exec(specDocUpdate, req.DocType, req.StudentID, activityID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update specific documentation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Specific documentation updated successfully"})
}
