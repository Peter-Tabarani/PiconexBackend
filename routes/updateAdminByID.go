package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Peter-Tabarani/PiconexBackend/models"
	"github.com/gorilla/mux"
)

type UpdateAdminRequest struct {
	Person models.Person `json:"person"`
	Title  string        `json:"title"`
}

func UpdateAdminByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Only allow PUT method
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get admin ID from URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	adminID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid admin ID", http.StatusBadRequest)
		return
	}

	// Decode request body
	var req UpdateAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	// Update person table
	personUpdate := `
		UPDATE person SET
			first_name=?, preferred_name=?, middle_name=?, last_name=?,
			email=?, phone_number=?, pronouns=?, sex=?, gender=?,
			birthday=?, address=?, city=?, state=?, zip_code=?, country=?
		WHERE id=?
	`
	_, err = tx.Exec(personUpdate,
		req.Person.FirstName, req.Person.PreferredName, req.Person.MiddleName, req.Person.LastName,
		req.Person.Email, req.Person.PhoneNumber, req.Person.Pronouns, req.Person.Sex, req.Person.Gender,
		req.Person.Birthday, req.Person.Address, req.Person.City, req.Person.State, req.Person.ZipCode, req.Person.Country,
		adminID,
	)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update person: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update admin title
	adminUpdate := `UPDATE admin SET title=? WHERE id=?`
	_, err = tx.Exec(adminUpdate, req.Title, adminID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update admin: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Admin updated successfully"})
}
