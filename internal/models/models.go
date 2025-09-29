package models

type Person struct {
	ID            int    `json:"id"`
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
	ID              int    `json:"id"`
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
	ID            int    `json:"id"`
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
	ActivityID int    `json:"activity_id"`
	Date       string `json:"date"`
	Time       string `json:"time"`
}

type Documentation struct {
	ActivityID int    `json:"activity_id"`
	Date       string `json:"date"`
	Time       string `json:"time"`
	File       []byte `json:"file"`
}

type PersonalDocumentation struct {
	ActivityID int    `json:"activity_id"`
	ID         int    `json:"id"`
	Date       string `json:"date"`
	Time       string `json:"time"`
	File       []byte `json:"file"`
}

type SpecificDocumentation struct {
	Activity_ID int    `json:"activity_id"`
	ID          int    `json:"id"`
	DocType     string `json:"doc_type"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	File        []byte `json:"file"`
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
	ActivityID int    `json:"activity_id"`
	Date       string `json:"date"`
	Time       string `json:"time"`
	EventDate  string `json:"event_date"`
	EventTime  string `json:"event_time"`
	//Duration   int    `json:"duration"` // NEED TO ADD. POC Handlers won't work until
	EventType string `json:"event_type"`
	ID        int    `json:"id"`
}

type Pinned struct {
	AdminID   int `json:"admin_id"`
	StudentID int `json:"student_id"`
}

type StudentDisability struct {
	ID           int `json:"id"`
	DisabilityID int `json:"disability_id"`
}

type StudentAccommodation struct {
	ID              int `json:"id"`
	AccommodationID int `json:"accommodation_id"`
}

type PocAdmin struct {
	ActivityID int `json:"activity_id"`
	ID         int `json:"id"`
}
