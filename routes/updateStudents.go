package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Peter-Tabarani/PiconexBackend/models"
)

func UpdateStudent(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var req struct {
		Person  models.Person  `json:"person"`
		Student models.Student `json:"student"`
	}

	// Decode the request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validation (example)
	sex := req.Person.Sex
	if sex != "male" && sex != "female" {
		http.Error(w, "Invalid sex value, must be 'male' or 'female'", http.StatusBadRequest)
		return
	}

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
		return
	}

	// Update person
	personQuery := `
		UPDATE person SET first_name=?, preferred_name=?, middle_name=?, last_name=?, 
			email=?, phone_number=?, pronouns=?, sex=?, gender=?, birthday=?, 
			address=?, city=?, state=?, zip_code=?, country=?
		WHERE id=?`

	_, err = tx.Exec(personQuery,
		req.Person.FirstName, req.Person.PreferredName, req.Person.MiddleName, req.Person.LastName,
		req.Person.Email, req.Person.PhoneNumber, req.Person.Pronouns, req.Person.Sex,
		req.Person.Gender, req.Person.Birthday, req.Person.Address, req.Person.City,
		req.Person.State, req.Person.ZipCode, req.Person.Country, req.Person.ID)

	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update person: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update student
	studentQuery := `
		UPDATE student SET year=?, start_year=?, planned_grad_year=?, housing=?, dining=?
		WHERE id=?`

	_, err = tx.Exec(studentQuery,
		req.Student.Year, req.Student.StartYear, req.Student.PlannedGradYear,
		req.Student.Housing, req.Student.Dining, req.Person.ID)

	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update student: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Student updated successfully"))
}
