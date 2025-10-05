package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func GetPersons(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// All data being selected for this GET command
	query := `
		SELECT
			person_id, first_name, preferred_name, middle_name, last_name,
			email, phone_number, pronouns, sex, gender,
			birthday, address, city, state, zip_code, country
		FROM person
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query)

	// Error message if QueryContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to obtain persons")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	// Creates an empty slice to obtain results
	persons := make([]models.Person, 0)

	// Reads each row returned by the database
	for rows.Next() {
		var p models.Person
		// Parses the current data into fields of "p" variable
		if err := rows.Scan(
			&p.PersonID, &p.FirstName, &p.PreferredName, &p.MiddleName, &p.LastName,
			&p.Email, &p.PhoneNumber, &p.Pronouns, &p.Sex, &p.Gender,
			&p.Birthday, &p.Address, &p.City, &p.State, &p.ZipCode, &p.Country,
		); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse persons")
			log.Println("Row scan error:", err)
			return
		}
		// Adds the obtained data to the slice
		persons = append(persons, p)
	}

	// Checks for errors during iteration such as network interruptions and driver errors
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Writes the slice as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, persons)
}

func GetPersonByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr, ok := vars["person_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing person ID")
		return
	}

	// Converts the "person_id" string to an integer
	personID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid person ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// All data being selected for this GET command
	query := `
		SELECT
			person_id, first_name, preferred_name, middle_name, last_name,
			email, phone_number, pronouns, sex, gender,
			birthday, address, city, state, zip_code, country
		FROM person
		WHERE id = ?
	`

	// Empty variable for person struct
	var p models.Person

	// Executes written SQL and retrieves only one row
	err = db.QueryRowContext(r.Context(), query, personID).Scan(
		&p.PersonID, &p.FirstName, &p.PreferredName, &p.MiddleName, &p.LastName,
		&p.Email, &p.PhoneNumber, &p.Pronouns, &p.Sex, &p.Gender,
		&p.Birthday, &p.Address, &p.City, &p.State, &p.ZipCode, &p.Country,
	)

	// Error message if no rows are found
	if err == sql.ErrNoRows {
		utils.WriteError(w, http.StatusNotFound, "Person not found")
		return
		// Error message if QueryRowContext or scan fails
	} else if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch person")
		log.Println("DB query error:", err)
		return
	}

	// Writes the struct as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, p)
}
