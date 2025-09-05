package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type CreateStuAccomRequest struct {
	ID              int `json:"id"`               // student_id
	AccommodationID int `json:"accommodation_id"` // accommodation_id
}

func CreateStuAccom(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateStuAccomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(
		"INSERT INTO stu_accom (id, accommodation_id) VALUES (?, ?)",
		req.ID, req.AccommodationID,
	)
	if err != nil {
		http.Error(w, "Failed to insert stu_accom: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Student accommodation created successfully",
	})
}
