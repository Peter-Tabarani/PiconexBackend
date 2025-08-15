package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/models"

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

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	accomQuery := `
		INSERT INTO accommodation (
			name, description
		) VALUES (?, ?)
	`

	res, err := tx.Exec(accomQuery,
		req.Accommodation.Name,
		req.Accommodation.Description,
	)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to insert into accommodation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	accomID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to retrieve accommodation ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":         "Accommodation created successfully",
		"accommodationId": accomID,
	})
}
