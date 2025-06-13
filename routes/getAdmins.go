package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/models"

	_ "github.com/go-sql-driver/mysql"
)

func GetAdmins(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT
			p.id, p.first_name, p.preferred_name, p.middle_name, p.last_name,
			p.email, p.phone_number, p.pronouns, p.sex, p.gender,
			p.birthday, p.address, p.city, p.state, p.zip_code, p.country, a.title
		FROM admin a
		JOIN person p ON a.id = p.id
	`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var admins []models.Admin

	for rows.Next() {
		var a models.Admin
		err := rows.Scan(
			&a.ID, &a.FirstName, &a.PreferredName, &a.MiddleName, &a.LastName,
			&a.Email, &a.PhoneNumber, &a.Pronouns, &a.Sex, &a.Gender,
			&a.Birthday, &a.Address, &a.City, &a.State, &a.ZipCode, &a.Country,
			&a.Title,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		admins = append(admins, a)
	}

	jsonBytes, err := json.MarshalIndent(admins, "", "    ") // Pretty print with 4 spaces indent
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
