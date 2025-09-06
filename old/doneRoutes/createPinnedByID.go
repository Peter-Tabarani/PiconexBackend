package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type CreatePinnedRequest struct {
	AdminID   int `json:"admin_id"`
	StudentID int `json:"student_id"`
}

func CreatePinnedByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreatePinnedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(
		"INSERT INTO pinned (admin_id, student_id) VALUES (?, ?)",
		req.AdminID, req.StudentID,
	)
	if err != nil {
		http.Error(w, "Failed to insert pinned: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Pinned created successfully",
	})
}
