package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func DeleteDisabilityByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")

	vars := mux.Vars(r)
	idStr := vars["disability_id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid disability ID", http.StatusBadRequest)
		return
	}

	// Step 1: Nullify or remove references in stu_dis
	_, err = db.Exec(`DELETE FROM stu_dis WHERE disability_id = ?`, id)
	if err != nil {
		http.Error(w, "Failed to update stu_dis: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 2: Delete from disability
	_, err = db.Exec(`DELETE FROM disability WHERE disability_id = ?`, id)
	if err != nil {
		http.Error(w, "Failed to delete disability: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Disability deleted successfully"})
}
