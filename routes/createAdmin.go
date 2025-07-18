package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/models"

	_ "github.com/go-sql-driver/mysql"
)

type CreateAdminRequest struct {
	Person models.Person `json:"person"`
	Title  string        `json:"title"`
}

func CreateAdmin(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	personQuery := `
		INSERT INTO person (
			first_name, preferred_name, middle_name, last_name, email,
			phone_number, pronouns, sex, gender, birthday,
			address, city, state, zip_code, country
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := tx.Exec(personQuery,
		req.Person.FirstName, req.Person.PreferredName, req.Person.MiddleName, req.Person.LastName,
		req.Person.Email, req.Person.PhoneNumber, req.Person.Pronouns, req.Person.Sex,
		req.Person.Gender, req.Person.Birthday, req.Person.Address, req.Person.City,
		req.Person.State, req.Person.ZipCode, req.Person.Country,
	)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to insert into person: "+err.Error(), http.StatusInternalServerError)
		return
	}

	personID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to retrieve person ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	adminQuery := `
		INSERT INTO admin (
			id, title
		) VALUES (?, ?)
	`

	_, err = tx.Exec(adminQuery, personID, req.Title)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to insert into admin: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Admin created successfully",
		"adminId": personID,
	})
}
