package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

func GetActivitiesByDate(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extracts path variables from the request
	vars := mux.Vars(r)
	date := vars["date"]
	if date == "" {
		utils.WriteError(w, http.StatusBadRequest, "Date is required")
		return
	}

	// All data being selected for this GET command
	query := `
		SELECT activity_id, activity_datetime
		FROM activity
		WHERE DATE(activity_datetime) = ?
	`

	// Executes written SQL
	rows, err := db.QueryContext(r.Context(), query, date)

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

type ActivitySummary struct {
	ActivityID       int       `json:"activity_id"`
	Summary          string    `json:"summary"`
	ActivityDateTime time.Time `json:"datetime"`
}

func GetActivitiesSummary(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// --- Step 1: Fetch all activities that occurred on this date ---
	query := `
		SELECT activity_id, activity_datetime
		FROM activity
	`
	rows, err := db.QueryContext(r.Context(), query)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to query activities")
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	log.Println(rows)

	activities := make([]models.Activity, 0)
	for rows.Next() {
		var a models.Activity
		if err := rows.Scan(&a.ActivityID, &a.ActivityDateTime); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse activity rows")
			log.Println("Row scan error:", err)
			return
		}
		activities = append(activities, a)
	}
	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Operational error")
		log.Println("Rows error:", err)
		return
	}

	// --- Step 2: For each activity, check its type and build a summary ---
	summaries := make([]ActivitySummary, 0)

	for _, a := range activities {
		// --- CASE 1: Point of Contact (meeting) ---
		var poc models.PointOfContact
		err := db.QueryRowContext(r.Context(), `
			SELECT point_of_contact_id, student_id, event_datetime
			FROM point_of_contact
			WHERE point_of_contact_id = ?
		`, a.ActivityID).Scan(&poc.PointOfContactID, &poc.StudentID, &poc.EventDateTime)

		if err == nil {
			// Find student
			var student models.Student
			err = db.QueryRowContext(r.Context(), `
				SELECT first_name, preferred_name
				FROM person
				WHERE person_id = ?
			`, poc.StudentID).Scan(&student.FirstName, &student.PreferredName)
			if err != nil && err != sql.ErrNoRows {
				log.Println("Student lookup error:", err)
				continue
			}
			studentName := student.PreferredName
			if studentName == "" {
				studentName = student.FirstName
			}
			if studentName == "" {
				studentName = "Student"
			}

			// Find linked admins
			adminRows, err := db.QueryContext(r.Context(), `
				SELECT p.first_name, p.preferred_name
				FROM admin a
				INNER JOIN person p ON a.admin_id = p.person_id
				INNER JOIN poc_admin pa ON pa.admin_id = a.admin_id
				WHERE pa.point_of_contact_id = ?
			`, poc.PointOfContactID)
			if err != nil {
				log.Println("Admin join error:", err)
				continue
			}

			adminNames := ""
			for adminRows.Next() {
				var firstName, preferredName sql.NullString
				adminRows.Scan(&firstName, &preferredName)
				name := preferredName.String
				if name == "" {
					name = firstName.String
				}
				if name != "" {
					if adminNames != "" {
						adminNames += ", "
					}
					adminNames += name
				}
			}
			adminRows.Close()
			if adminNames == "" {
				adminNames = "an administrator"
			}

			// Format readable summary
			summary := fmt.Sprintf(
				"%s scheduled a meeting with %s on %s at %s",
				studentName,
				adminNames,
				poc.EventDateTime.Format("1/2/2006"),
				poc.EventDateTime.Format("3:04 PM"),
			)
			summaries = append(summaries, ActivitySummary{
				ActivityID:       a.ActivityID,
				Summary:          summary,
				ActivityDateTime: a.ActivityDateTime,
			})
			continue
		}

		// --- CASE 2: Specific Documentation (upload) ---
		var doc models.SpecificDocumentation
		err = db.QueryRowContext(r.Context(), `
			SELECT specific_documentation_id, student_id, doc_type
			FROM specific_documentation
			WHERE specific_documentation_id = ?
		`, a.ActivityID).Scan(&doc.SpecificDocumentationID, &doc.StudentID, &doc.DocType)
		if err == nil {
			var student models.Student
			err = db.QueryRowContext(r.Context(), `
				SELECT first_name, preferred_name
				FROM person
				WHERE person_id = ?
			`, doc.StudentID).Scan(&student.FirstName, &student.PreferredName)
			if err != nil && err != sql.ErrNoRows {
				log.Println("Student lookup error:", err)
				continue
			}
			studentName := student.PreferredName
			if studentName == "" {
				studentName = student.FirstName
			}
			if studentName == "" {
				studentName = "Student"
			}

			summary := fmt.Sprintf("%s uploaded documentation", studentName)
			if doc.DocType != "" {
				summary += fmt.Sprintf(": %s", doc.DocType)
			}

			summaries = append(summaries, ActivitySummary{
				ActivityID:       a.ActivityID,
				Summary:          summary,
				ActivityDateTime: a.ActivityDateTime,
			})
			continue
		}
	}

	// --- Step 3: Sort summaries by time descending (optional, usually done in SQL) ---
	// In this handler, we’ll just return as-is for simplicity

	utils.WriteJSON(w, http.StatusOK, summaries)
}
