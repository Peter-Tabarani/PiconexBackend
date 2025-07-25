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
	Activity_ID int    `json:"activity_id"`
	Date        string `json:"date"`
	Time        string `json:"time"`
}

type Documentation struct {
	Activity_ID int    `json:"activity_id"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	File        []byte `json:"file"`
}

type Personal_Documentation struct {
	Activity_ID int    `json:"activity_id"`
	ID          int    `json:"id"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	File        []byte `json:"file"`
}

type Specific_Documentation struct {
	Activity_ID int    `json:"activity_id"`
	ID          int    `json:"id"`
	DocType     string `json:"doc_type"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	File        []byte `json:"file"`
}

type Disability struct {
	Disability_ID int    `json:"disability_id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
}

type Accommodation struct {
	Accommodation_ID int    `json:"accommodation_id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
}

type PointOfContact struct {
	Activity_ID int    `json:"activity_id"`
	Event_Date  string `json:"event_date"`
	Event_Time  string `json:"event_time"`
	Event_Type  string `json:"event_type"`
	Student_ID  *int   `json:"student_id,omitempty"`
	Admin_ID    *int   `json:"admin_id,omitempty"`
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
	AdminID    int `json:"admin_id"`
}
