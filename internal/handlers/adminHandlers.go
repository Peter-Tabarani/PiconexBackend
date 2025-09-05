package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"
	"github.com/gorilla/mux"
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

func DeleteAdmin(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Only allow DELETE method
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	adminID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid admin ID", http.StatusBadRequest)
		return
	}

	// Delete person record (will cascade to admin)
	res, err := db.Exec("DELETE FROM person WHERE id = ?", adminID)
	if err != nil {
		http.Error(w, "Failed to delete admin: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking deletion result: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Admin not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Admin deleted successfully"})
}

type UpdateAdminRequest struct {
	Person models.Person `json:"person"`
	Title  string        `json:"title"`
}

func UpdateAdmin(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Only allow PUT method
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get admin ID from URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	adminID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid admin ID", http.StatusBadRequest)
		return
	}

	// Decode request body
	var req UpdateAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	// Update person table
	personUpdate := `
		UPDATE person SET
			first_name=?, preferred_name=?, middle_name=?, last_name=?,
			email=?, phone_number=?, pronouns=?, sex=?, gender=?,
			birthday=?, address=?, city=?, state=?, zip_code=?, country=?
		WHERE id=?
	`
	_, err = tx.Exec(personUpdate,
		req.Person.FirstName, req.Person.PreferredName, req.Person.MiddleName, req.Person.LastName,
		req.Person.Email, req.Person.PhoneNumber, req.Person.Pronouns, req.Person.Sex, req.Person.Gender,
		req.Person.Birthday, req.Person.Address, req.Person.City, req.Person.State, req.Person.ZipCode, req.Person.Country,
		adminID,
	)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update person: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update admin title
	adminUpdate := `UPDATE admin SET title=? WHERE id=?`
	_, err = tx.Exec(adminUpdate, req.Title, adminID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update admin: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Admin updated successfully"})
}
