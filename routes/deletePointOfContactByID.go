package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func DeletePointOfContactByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")

	vars := mux.Vars(r)
	activityIDStr := vars["activity_id"]

	activityID, err := strconv.Atoi(activityIDStr)
	if err != nil {
		http.Error(w, "Invalid activity ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`DELETE FROM point_of_contact WHERE activity_id = ?`, activityID)
	if err != nil {
		http.Error(w, "Failed to delete point_of_contact: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Point of Contact deleted successfully"})
}
