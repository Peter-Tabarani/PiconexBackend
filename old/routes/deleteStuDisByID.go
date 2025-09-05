package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func DeleteStuDisByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")

	vars := mux.Vars(r)

	idStr := vars["id"]
	disabilityIDStr := vars["disability_id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	disabilityID, err := strconv.Atoi(disabilityIDStr)
	if err != nil {
		http.Error(w, "Invalid disability ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`DELETE FROM stu_dis WHERE id = ? AND disability_id = ?`, id, disabilityID)
	if err != nil {
		http.Error(w, "Failed to delete stu_dis: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "StuDis deleted successfully"})
}
