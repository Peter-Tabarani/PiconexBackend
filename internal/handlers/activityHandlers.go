package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Peter-Tabarani/PiconexBackend/internal/models"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"
	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func GetActivities(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// All data being selected for this GET command
	query := `
		SELECT
    	    activity_id, activity_datetime
    	FROM activity
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query)

	// Error message if QueryContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to obtain activities")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	// Creates an empty slice to obtain results
	activities := make([]models.Activity, 0)

	// Reads each row returned by the database
	for rows.Next() {
		var a models.Activity
		// Parses the current data into fields of "a" variable
		if err := rows.Scan(&a.ActivityID, &a.ActivityDateTime); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse activities")
			log.Println("Row scan error:", err)
			return
		}

		// Adds the obtained data to the slice
		activities = append(activities, a)
	}

	// Checks for errors during iteration such as network interruptions and driver errors
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational Error")
		log.Println("Rows error:", err)
		return
	}

	// Writes the slice as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, activities)
}

func GetActivityByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extracts path variables from the request
	vars := mux.Vars(r)
	idStr, ok := vars["activity_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Missing activity ID")
		return
	}

	// Converts the "activity_id" string to an integer
	activityID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid activity ID")
		log.Println("Invalid ID parse error:", err)
		return
	}

	// All data being selected for this GET command
	query := `
        SELECT activity_id, activity_datetime
        FROM activity
        WHERE activity_id = ?
    `

	// Empty variable for activity struct
	var a models.Activity

	// Executes written SQL and retrieves only one row
	err = db.QueryRowContext(r.Context(), query, activityID).Scan(
		&a.ActivityID, &a.ActivityDateTime,
	)

	// Error message if no rows are found
	if err == sql.ErrNoRows {
		utils.WriteError(w, http.StatusNotFound, "Activity not found")
		return
		// Error message if QueryRowContext or scan fails
	} else if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch activity")
		log.Println("DB query error:", err)
		return
	}

	// Writes the struct as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, a)
}

func GetActivitiesSummary(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extracts query parameters from the request URL
	dateStr := r.URL.Query().Get("date")
	tzStr := r.URL.Query().Get("tz")
	studentIDStr := r.URL.Query().Get("student_id")
	adminIDStr := r.URL.Query().Get("admin_id")

	// Sets the default timezone to UTC or loads the provided timezone
	loc := time.UTC
	if tzStr != "" {
		var err error
		loc, err = time.LoadLocation(tzStr)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "Invalid timezone")
			log.Println("Timezone parse error:", err)
			return
		}
	}

	// Base SQL query to retrieve all activities
	query := `
		SELECT activity_id, activity_datetime
		FROM activity
	`
	args := []any{}
	where := []string{}

	// Optional date filter — converts date string to time range
	if dateStr != "" {
		targetDate, err := time.ParseInLocation("2006-01-02", dateStr, loc)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "Invalid date format (expected YYYY-MM-DD)")
			log.Println("Date parse error:", err)
			return
		}
		start := targetDate
		end := targetDate.Add(24 * time.Hour)
		where = append(where, "activity_datetime >= ? AND activity_datetime < ?")
		args = append(args, start.UTC(), end.UTC())
	}

	// Optional student filter — restricts to activities linked to a student
	if studentIDStr != "" {
		where = append(where, `
			activity_id IN (
				SELECT point_of_contact_id FROM point_of_contact WHERE student_id = ?
				UNION
				SELECT specific_documentation_id FROM specific_documentation WHERE student_id = ?
			)
		`)
		args = append(args, studentIDStr, studentIDStr)
	}

	// Optional admin filter — restricts to activities linked to an admin
	if adminIDStr != "" {
		where = append(where, `
			activity_id IN (
				SELECT poc.point_of_contact_id
				FROM point_of_contact poc
				INNER JOIN poc_admin pa ON pa.point_of_contact_id = poc.point_of_contact_id
				WHERE pa.admin_id = ?
			)
		`)
		args = append(args, adminIDStr)
	}

	// Adds all conditions to the final query
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	query += " ORDER BY activity_datetime DESC"

	// Executes written SQL query
	rows, err := db.QueryContext(r.Context(), query, args...)

	// Error message if QueryContext fails
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to obtain activities")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	// Defines supporting structs for output
	type Person struct {
		ID            int    `json:"id"`
		FirstName     string `json:"first_name"`
		PreferredName string `json:"preferred_name"`
	}

	type ActivityData struct {
		ActivityID       int        `json:"activity_id"`
		Type             string     `json:"type"`
		Student          Person     `json:"student"`
		Admins           []Person   `json:"admins,omitempty"`
		DocType          *string    `json:"doc_type,omitempty"`
		FileName         *string    `json:"file_name,omitempty"`
		ActivityDateTime time.Time  `json:"activity_datetime"`
		EventDateTime    *time.Time `json:"event_datetime,omitempty"`
	}

	// Creates an empty slice to store results
	activities := make([]ActivityData, 0)

	// Reads each row returned by the database
	for rows.Next() {
		var a ActivityData
		// Parses current row data into ActivityData fields
		if err := rows.Scan(&a.ActivityID, &a.ActivityDateTime); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse activities")
			log.Println("Row scan error:", err)
			return
		}

		// --- CASE 1: Point of Contact ---
		var poc struct {
			PointOfContactID int
			StudentID        int
			EventDateTime    time.Time
		}
		err = db.QueryRowContext(r.Context(), `
			SELECT point_of_contact_id, student_id, event_datetime
			FROM point_of_contact
			WHERE point_of_contact_id = ?
		`, a.ActivityID).Scan(&poc.PointOfContactID, &poc.StudentID, &poc.EventDateTime)

		// Found: Point of Contact type
		if err == nil {
			var student Person
			if err := db.QueryRowContext(r.Context(), `
				SELECT person_id, first_name, preferred_name
				FROM person
				WHERE person_id = ?
			`, poc.StudentID).Scan(&student.ID, &student.FirstName, &student.PreferredName); err != nil {
				log.Println("Student query error:", err)
			}

			adminRows, _ := db.QueryContext(r.Context(), `
				SELECT p.person_id, p.first_name, p.preferred_name
				FROM admin a
				INNER JOIN person p ON a.admin_id = p.person_id
				INNER JOIN poc_admin pa ON pa.admin_id = a.admin_id
				WHERE pa.point_of_contact_id = ?
			`, poc.PointOfContactID)

			admins := []Person{}
			for adminRows.Next() {
				var adm Person
				if err := adminRows.Scan(&adm.ID, &adm.FirstName, &adm.PreferredName); err == nil {
					admins = append(admins, adm)
				}
			}
			adminRows.Close()

			a.Type = "point_of_contact"
			a.Student = student
			a.Admins = admins
			a.EventDateTime = &poc.EventDateTime
			activities = append(activities, a)
			continue
		}

		// --- CASE 2: Specific Documentation ---
		var doc struct {
			ID        int
			StudentID int
			DocType   string
		}
		err = db.QueryRowContext(r.Context(), `
			SELECT specific_documentation_id, student_id, doc_type
			FROM specific_documentation
			WHERE specific_documentation_id = ?
		`, a.ActivityID).Scan(&doc.ID, &doc.StudentID, &doc.DocType)

		// Found: Specific Documentation type
		if err == nil {
			var student Person
			if err := db.QueryRowContext(r.Context(), `
				SELECT person_id, first_name, preferred_name
				FROM person
				WHERE person_id = ?
			`, doc.StudentID).Scan(&student.ID, &student.FirstName, &student.PreferredName); err != nil {
				log.Println("Student query error:", err)
			}

			a.Type = "specific_documentation"
			a.Student = student
			a.DocType = &doc.DocType
			activities = append(activities, a)
		}
	}

	// Checks for iteration errors
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational error")
		log.Println("Rows error:", err)
		return
	}

	// Writes the slice as JSON & sends a HTTP 200 response code
	utils.WriteJSON(w, http.StatusOK, activities)
}
