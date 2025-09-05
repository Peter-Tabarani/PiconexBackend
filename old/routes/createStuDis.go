package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type CreateStuDisRequest struct {
	ID           int `json:"id"`            // student id
	DisabilityID int `json:"disability_id"` // disability id
}

func CreateStuDis(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateStuDisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(
		"INSERT INTO stu_dis (id, disability_id) VALUES (?, ?)",
		req.ID, req.DisabilityID,
	)
	if err != nil {
		http.Error(w, "Failed to insert stu_dis: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "StuDis created successfully",
	})
}
