package routes

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type CreateSpecificDocumentationRequest struct {
	StudentID  int    `json:"student_id"`
	DocType    string `json:"doc_type"`    // enum('trad','al','in')
	Date       string `json:"date"`        // format: "YYYY-MM-DD"
	Time       string `json:"time"`        // format: "HH:MM:SS"
	FileBase64 string `json:"file_base64"` // base64-encoded file blob
}

func CreateSpecificDocumentation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateSpecificDocumentationRequest
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

	// Insert into specific_documentation
	specificDocQuery := `INSERT INTO specific_documentation (activity_id, id, doc_type) VALUES (?, ?, ?)`
	_, err = tx.Exec(specificDocQuery, activityID, req.StudentID, req.DocType)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to insert into specific_documentation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "Specific documentation created successfully",
		"activity_id": activityID,
		"student_id":  req.StudentID,
		"doc_type":    req.DocType,
	})
}
