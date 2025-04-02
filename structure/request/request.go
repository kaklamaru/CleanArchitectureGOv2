package request

type RegisterTeacher struct {
	// UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	TitleName string `json:"title_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Code      string `json:"code"`
}
type RegisterStudent struct {
	// UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	TitleName string `json:"title_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Code      string `json:"code"`
	Year      uint   `json:"year"`
	BranchId  uint   `json:"branch_id"`
}

type EventRequest struct {
	EventName   string `json:"event_name"`
	StartDate   string `json:"start_date"`
	WorkingHour uint   `json:"working_hour"`
	SchoolYear  uint   `json:"school_year"`
	Location    string `json:"location"`
	FreeSpace   uint   `json:"free_space"`
	Detail      string `json:"detail"`
	Branches    []uint `json:"branches"`
	Years       []uint `json:"years"`
}