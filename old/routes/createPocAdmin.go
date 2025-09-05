package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type CreatePOCAdminRequest struct {
	ActivityID int `json:"activity_id"`
	AdminID    int `json:"admin_id"`
}

func CreatePocAdmin(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreatePOCAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(
		"INSERT INTO poc_adm (activity_id, admin_id) VALUES (?, ?)",
		req.ActivityID, req.AdminID,
	)
	if err != nil {
		http.Error(w, "Failed to insert poc_adm: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "POC Admin created successfully",
	})
}
