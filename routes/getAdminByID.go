package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Peter-Tabarani/PiconexBackend/models"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func GetAdminByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT 
			a.id,
			p.first_name, p.preferred_name, p.middle_name, p.last_name,
			p.email, p.phone_number, p.pronouns, p.sex, p.gender,
			p.birthday, p.address, p.city, p.state, p.zip_code, p.country,
			a.title
		FROM admin a
		JOIN person p ON a.id = p.id
		WHERE a.id = ?
	`

	var admin models.Admin
	err = db.QueryRow(query, id).Scan(
		&admin.ID,
		&admin.FirstName,
		&admin.PreferredName,
		&admin.MiddleName,
		&admin.LastName,
		&admin.Email,
		&admin.PhoneNumber,
		&admin.Pronouns,
		&admin.Sex,
		&admin.Gender,
		&admin.Birthday,
		&admin.Address,
		&admin.City,
		&admin.State,
		&admin.ZipCode,
		&admin.Country,
		&admin.Title, // extra field
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Admin not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	jsonBytes, err := json.MarshalIndent(admin, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
