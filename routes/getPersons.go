package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/models"

	_ "github.com/go-sql-driver/mysql"
)

func GetPersons(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT
			p.id, p.first_name, p.preferred_name, p.middle_name, p.last_name,
			p.email, p.phone_number, p.pronouns, p.sex, p.gender,
			p.birthday, p.address, p.city, p.state, p.zip_code, p.country
			FROM person p
	`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var persons []models.Person

	for rows.Next() {
		var p models.Person
		err := rows.Scan(
			&p.ID, &p.FirstName, &p.PreferredName, &p.MiddleName, &p.LastName,
			&p.Email, &p.PhoneNumber, &p.Pronouns, &p.Sex, &p.Gender,
			&p.Birthday, &p.Address, &p.City, &p.State, &p.ZipCode, &p.Country,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		persons = append(persons, p)
	}

	jsonBytes, err := json.MarshalIndent(persons, "", "    ") // Pretty print with 4 spaces indent
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
