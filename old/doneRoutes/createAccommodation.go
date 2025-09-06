package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type CreateAccommodationRequest struct {
	Accommodation models.Accommodation `json:"accommodation"`
}

func CreateAccommodation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateAccommodationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	res, err := db.Exec(
		"INSERT INTO accommodation (name, description) VALUES (?, ?)",
		req.Accommodation.Name, req.Accommodation.Description,
	)
	if err != nil {
		http.Error(w, "Failed to insert accommodation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	lastID, _ := res.LastInsertId()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":          "Accommodation created successfully",
		"accommodation_id": lastID,
	})
}
