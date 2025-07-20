package routes

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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

	// Parse multipart form, limit max memory to 10MB (adjust as needed)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get form values (optional updates)
	docType := r.FormValue("doc_type")
	studentIDStr := r.FormValue("student_id")

	var studentID int
	if studentIDStr != "" {
		studentID, err = strconv.Atoi(studentIDStr)
		if err != nil {
			http.Error(w, "Invalid student ID", http.StatusBadRequest)
			return
		}
	}

	// Get file from form (optional)
	var fileBytes []byte
	file, _, err := r.FormFile("file")
	if err == nil {
		defer file.Close()
		fileBytes, err = io.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading file: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else if err != http.ErrMissingFile {
		// Only error if the error is something other than missing file (which is allowed)
		http.Error(w, "Error retrieving file: "+err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	// Update documentation table (file) only if file was provided
	if len(fileBytes) > 0 {
		docUpdate := `UPDATE documentation SET file=? WHERE activity_id=?`
		_, err = tx.Exec(docUpdate, fileBytes, activityID)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Failed to update documentation: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Update specific_documentation table only if docType or studentID was provided
	if docType != "" || studentIDStr != "" {
		// Build dynamic query based on provided fields
		query := "UPDATE specific_documentation SET "
		args := []interface{}{}
		setClauses := []string{}

		if docType != "" {
			setClauses = append(setClauses, "doc_type=?")
			args = append(args, docType)
		}
		if studentIDStr != "" {
			setClauses = append(setClauses, "id=?")
			args = append(args, studentID)
		}

		query += join(setClauses, ", ")
		query += " WHERE activity_id=?"
		args = append(args, activityID)

		_, err = tx.Exec(query, args...)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Failed to update specific documentation: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Specific documentation updated successfully"})
}

// helper function to join strings (could use strings.Join but for []string)
func join(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
