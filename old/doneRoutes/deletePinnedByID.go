package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func DeletePinnedByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")

	vars := mux.Vars(r)

	adminIDStr := vars["admin_id"]
	studentIDStr := vars["student_id"]

	adminID, err := strconv.Atoi(adminIDStr)
	if err != nil {
		http.Error(w, "Invalid admin ID", http.StatusBadRequest)
		return
	}

	studentID, err := strconv.Atoi(studentIDStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	// Delete the pinned record
	_, err = db.Exec(`DELETE FROM pinned WHERE admin_id = ? AND student_id = ?`, adminID, studentID)
	if err != nil {
		http.Error(w, "Failed to delete pinned: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Pinned deleted successfully"})
}
