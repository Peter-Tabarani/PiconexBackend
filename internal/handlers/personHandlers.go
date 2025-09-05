package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"
	"github.com/gorilla/mux"
)

func GetPersons(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT
			id, first_name, preferred_name, middle_name, last_name,
			email, phone_number, pronouns, sex, gender,
			birthday, address, city, state, zip_code, country
		FROM person
	`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Failed to fetch persons", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var persons []models.Person

	for rows.Next() {
		var p models.Person
		if err := rows.Scan(
			&p.ID, &p.FirstName, &p.PreferredName, &p.MiddleName, &p.LastName,
			&p.Email, &p.PhoneNumber, &p.Pronouns, &p.Sex, &p.Gender,
			&p.Birthday, &p.Address, &p.City, &p.State, &p.ZipCode, &p.Country,
		); err != nil {
			http.Error(w, "Failed to scan person", http.StatusInternalServerError)
			return
		}
		persons = append(persons, p)
	}

	jsonBytes, err := json.MarshalIndent(persons, "", "    ")
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func GetPersonByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing person ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid person ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT
			id, first_name, preferred_name, middle_name, last_name,
			email, phone_number, pronouns, sex, gender,
			birthday, address, city, state, zip_code, country
		FROM person
		WHERE id = ?
	`

	row := db.QueryRow(query, id)
	var p models.Person
	if err := row.Scan(
		&p.ID, &p.FirstName, &p.PreferredName, &p.MiddleName, &p.LastName,
		&p.Email, &p.PhoneNumber, &p.Pronouns, &p.Sex, &p.Gender,
		&p.Birthday, &p.Address, &p.City, &p.State, &p.ZipCode, &p.Country,
	); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Person not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch person", http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
