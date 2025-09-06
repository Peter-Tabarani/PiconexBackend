package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func DeleteAccommodationByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid accommodation ID", http.StatusBadRequest)
		return
	}

	// Step 1: Nullify any foreign keys referencing this accommodation (if needed)
	// For example, stu_accom references accommodation_id
	_, err = db.Exec(`UPDATE stu_accom SET accommodation_id = NULL WHERE accommodation_id = ?`, id)
	if err != nil {
		http.Error(w, "Failed to update stu_accom: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 2: Delete from accommodation
	_, err = db.Exec(`DELETE FROM accommodation WHERE accommodation_id = ?`, id)
	if err != nil {
		http.Error(w, "Failed to delete accommodation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Accommodation deleted successfully"})
}
