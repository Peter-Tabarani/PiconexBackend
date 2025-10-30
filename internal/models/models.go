package models

import "time"

type Person struct {
	PersonID      int    `json:"person_id"`
	FirstName     string `json:"first_name"`
	PreferredName string `json:"preferred_name"`
	MiddleName    string `json:"middle_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	PhoneNumber   string `json:"phone_number"`
	Pronouns      string `json:"pronouns"`
	Sex           string `json:"sex"`
	Gender        string `json:"gender"`
	Birthday      string `json:"birthday"`
	Address       string `json:"address"`
	City          string `json:"city"`
	State         string `json:"state"`
	ZipCode       string `json:"zip_code"`
	Country       string `json:"country"`
}

type Student struct {
	StudentID       int    `json:"student_id"`
	FirstName       string `json:"first_name"`
	PreferredName   string `json:"preferred_name"`
	MiddleName      string `json:"middle_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	PhoneNumber     string `json:"phone_number"`
	Pronouns        string `json:"pronouns"`
	Sex             string `json:"sex"`
	Gender          string `json:"gender"`
	Birthday        string `json:"birthday"`
	Address         string `json:"address"`
	City            string `json:"city"`
	State           string `json:"state"`
	ZipCode         string `json:"zip_code"`
	Country         string `json:"country"`
	Year            string `json:"year"`
	StartYear       int    `json:"start_year"`
	PlannedGradYear int    `json:"planned_grad_year"`
	Housing         string `json:"housing"`
	Dining          string `json:"dining"`
}

type Admin struct {
	AdminID       int    `json:"admin_id"`
	FirstName     string `json:"first_name"`
	PreferredName string `json:"preferred_name"`
	MiddleName    string `json:"middle_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	PhoneNumber   string `json:"phone_number"`
	Pronouns      string `json:"pronouns"`
	Sex           string `json:"sex"`
	Gender        string `json:"gender"`
	Birthday      string `json:"birthday"`
	Address       string `json:"address"`
	City          string `json:"city"`
	State         string `json:"state"`
	ZipCode       string `json:"zip_code"`
	Country       string `json:"country"`
	Title         string `json:"title"`
}

type Activity struct {
	ActivityID       int       `json:"activity_id"`
	ActivityDateTime time.Time `json:"activity_datetime"`
}

type Documentation struct {
	DocumentationID  int       `json:"documentation_id"`
	ActivityDateTime time.Time `json:"activity_datetime"`
	FileName         string    `json:"file_name"`
	FilePath         string    `json:"file_path"`
	MimeType         string    `json:"mime_type"`
	SizeBytes        int64     `json:"size_bytes"`
	UploadedBy       *int      `json:"uploaded_by,omitempty"`
}

type PersonalDocumentation struct {
	PersonalDocumentationID int       `json:"personal_documentation_id"`
	ActivityDateTime        time.Time `json:"activity_datetime"`
	FileName                string    `json:"file_name"`
	FilePath                string    `json:"file_path"`
	MimeType                string    `json:"mime_type"`
	SizeBytes               int64     `json:"size_bytes"`
	UploadedBy              *int      `json:"uploaded_by,omitempty"`
	AdminID                 int       `json:"admin_id"`
}

type SpecificDocumentation struct {
	SpecificDocumentationID int       `json:"specific_documentation_id"`
	ActivityDateTime        time.Time `json:"activity_datetime"`
	FileName                string    `json:"file_name"`
	FilePath                string    `json:"file_path"`
	MimeType                string    `json:"mime_type"`
	SizeBytes               int64     `json:"size_bytes"`
	UploadedBy              *int      `json:"uploaded_by,omitempty"`
	DocType                 string    `json:"doc_type"`
	StudentID               int       `json:"student_id"`
}

type Disability struct {
	DisabilityID int    `json:"disability_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
}

type Accommodation struct {
	AccommodationID int    `json:"accommodation_id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
}

type PointOfContact struct {
	PointOfContactID int       `json:"point_of_contact_id"`
	ActivityDateTime time.Time `json:"activity_datetime"`
	EventDateTime    time.Time `json:"event_datetime"`
	Duration         int       `json:"duration"`
	EventType        string    `json:"event_type"`
	StudentID        int       `json:"student_id"`
}

type Pinned struct {
	AdminID   int `json:"admin_id"`
	StudentID int `json:"student_id"`
}

type StudentDisability struct {
	StudentID    int `json:"student_id"`
	DisabilityID int `json:"disability_id"`
}

type StudentAccommodation struct {
	StudentID       int `json:"student_id"`
	AccommodationID int `json:"accommodation_id"`
}

type PocAdmin struct {
	PointOfContactID int `json:"point_of_contact_id"`
	AdminID          int `json:"admin_id"`
}
