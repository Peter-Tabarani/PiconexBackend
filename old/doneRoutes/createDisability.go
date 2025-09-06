package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type CreateDisabilityRequest struct {
	Disability models.Disability `json:"disability"`
}

func CreateDisability(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateDisabilityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	res, err := db.Exec(
		"INSERT INTO disability (name, description) VALUES (?, ?)",
		req.Disability.Name, req.Disability.Description,
	)
	if err != nil {
		http.Error(w, "Failed to insert disability: "+err.Error(), http.StatusInternalServerError)
		return
	}

	lastID, _ := res.LastInsertId()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "Disability created successfully",
		"disability_id": lastID,
	})
}
