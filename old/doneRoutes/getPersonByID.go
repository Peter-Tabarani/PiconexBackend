package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func GetPersonByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := `
        SELECT
            id, first_name, preferred_name, middle_name, last_name, email,
            phone_number, pronouns, sex, gender, birthday, address,
            city, state, zip_code, country
        FROM person
        WHERE id = ?
    `

	var person models.Person
	err = db.QueryRow(query, id).Scan(
		&person.ID,
		&person.FirstName,
		&person.PreferredName,
		&person.MiddleName,
		&person.LastName,
		&person.Email,
		&person.PhoneNumber,
		&person.Pronouns,
		&person.Sex,
		&person.Gender,
		&person.Birthday,
		&person.Address,
		&person.City,
		&person.State,
		&person.ZipCode,
		&person.Country,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Person not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	jsonBytes, err := json.MarshalIndent(person, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
