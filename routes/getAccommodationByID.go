package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"
	"github.com/gorilla/mux"
)

func GetAccommodationByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"] // match the route

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid accommodation ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT accommodation_id, name, description
		FROM accommodation
		WHERE accommodation_id = ?
	`

	var accom models.Accommodation
	err = db.QueryRow(query, id).Scan(&accom.Accommodation_ID, &accom.Name, &accom.Description)
	if err == sql.ErrNoRows {
		http.Error(w, "Accommodation not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accom)
}
