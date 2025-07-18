package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/models"
	"github.com/gorilla/mux"
)

func UpdateStudentByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var s models.Student
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Update only fields that were sent (optional: dynamic SQL builder)
	_, err = db.Exec(`
		UPDATE person 
		SET gender = ? 
		WHERE id = ?`, s.Gender, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Student updated successfully")
}
