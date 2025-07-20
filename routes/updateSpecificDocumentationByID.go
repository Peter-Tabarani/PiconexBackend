package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UpdateSpecificDocumentationRequest struct {
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
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Update only specific_documentation table (doc_type, student_id)
	_, err = db.Exec(`
		UPDATE specific_documentation
		SET doc_type = ?, id = ?
		WHERE activity_id = ?`,
		req.DocType, req.StudentID, activityID,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update specific documentation: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Specific documentation updated successfully")
}
