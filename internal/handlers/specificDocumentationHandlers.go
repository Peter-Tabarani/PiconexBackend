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

func GetSpecificDocumentations(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT
		ac.activity_id, s.id, sd.doc_type, ac.date, ac.time, d.file
	FROM specific_documentation sd
	JOIN activity ac ON sd.activity_id = ac.activity_id
	JOIN student s ON sd.id = s.id
	JOIN documentation d ON d.activity_id = sd.activity_id
	`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var specific_documentations []models.Specific_Documentation

	for rows.Next() {
		var sd models.Specific_Documentation
		err := rows.Scan(
			&sd.Activity_ID, &sd.ID, &sd.DocType, &sd.Date, &sd.Time, &sd.File,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		specific_documentations = append(specific_documentations, sd)
	}

	jsonBytes, err := json.MarshalIndent(specific_documentations, "", "    ") // Pretty print with 4 spaces indent
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func GetSpecificDocumentationByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	activityIDStr := vars["activity_id"]
	activityID, err := strconv.Atoi(activityIDStr)
	if err != nil {
		http.Error(w, "Invalid activity ID", http.StatusBadRequest)
		return
	}

	query := `
        SELECT
            sd.activity_id,
            sd.id,
            sd.doc_type,
            a.date,
            a.time,
            d.file
        FROM specific_documentation sd
        JOIN activity a ON sd.activity_id = a.activity_id
        JOIN documentation d ON sd.activity_id = d.activity_id
        WHERE sd.activity_id = ?
    `

	row := db.QueryRow(query, activityID)

	var sd models.Specific_Documentation

	err = row.Scan(
		&sd.Activity_ID,
		&sd.ID,
		&sd.DocType,
		&sd.Date,
		&sd.Time,
		&sd.File,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No documentation found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.MarshalIndent(sd, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func GetSpecificDocumentationByStudentID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	studentID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT
			ac.activity_id, sd.id, sd.doc_type, ac.date, ac.time, d.file
		FROM specific_documentation sd
		JOIN documentation d ON sd.activity_id = d.activity_id
		JOIN activity ac ON ac.activity_id = sd.activity_id
		WHERE sd.id = ?
	`

	rows, err := db.Query(query, studentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var docs []models.Specific_Documentation

	for rows.Next() {
		var doc models.Specific_Documentation
		if err := rows.Scan(
			&doc.Activity_ID, &doc.ID, &doc.DocType, &doc.Date, &doc.Time, &doc.File,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		docs = append(docs, doc)
	}

	if len(docs) == 0 {
		http.Error(w, "No specific documentation found for the student", http.StatusNotFound)
		return
	}

	jsonBytes, err := json.MarshalIndent(docs, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

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

func DeleteSpecificDocumentation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Only allow DELETE method
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

	// Delete activity record (will cascade to documentation and specific_documentation)
	res, err := db.Exec("DELETE FROM activity WHERE activity_id = ?", activityID)
	if err != nil {
		http.Error(w, "Failed to delete activity and related documentation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking deletion result: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Activity not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Activity and related documentation deleted successfully"})
}

type UpdateSpecificDocumentationRequest struct {
	Date       string `json:"date"`
	Time       string `json:"time"`
	FileBase64 string `json:"file_base64"`
	DocType    string `json:"doc_type"`
	StudentID  int    `json:"student_id"`
}

func UpdateSpecificDocumentation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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

	// Update activity
	_, err = tx.Exec(`UPDATE activity SET date=?, time=? WHERE activity_id=?`, req.Date, req.Time, activityID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update activity: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update documentation
	_, err = tx.Exec(`UPDATE documentation SET file=? WHERE activity_id=?`, fileBytes, activityID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update documentation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update specific_documentation
	_, err = tx.Exec(`UPDATE specific_documentation SET doc_type=?, id=? WHERE activity_id=?`, req.DocType, req.StudentID, activityID)
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
