package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/models"
	"github.com/gorilla/mux"
)

func UpdateDisabilityByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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
