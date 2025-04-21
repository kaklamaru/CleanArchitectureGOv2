// package entity

// type User struct {
// 	UserID   uint     `gorm:"primaryKey;autoIncrement" json:"user_id"`
// 	Email    string   `gorm:"unique;not null" json:"email"`
// 	Password string   `gorm:"not null" json:"password"`
// 	Role     string   `gorm:"default:'user'" json:"role"`
// 	Student  *Student `gorm:"foreignKey:UserID"`
// 	Teacher  *Teacher `gorm:"foreignKey:UserID"`
// 	News     []News   `gorm:"foreignKey:UserID" json:"news"`
// }
// type Teacher struct {
// 	UserID    uint   `gorm:"primaryKey" json:"user_id"`
// 	TitleName string `gorm:"not null" json:"title_name"`
// 	FirstName string `gorm:"not null" json:"first_name"`
// 	LastName  string `gorm:"not null" json:"last_name"`
// 	Phone     string `gorm:"unique" json:"phone"`
// 	Code      string `gorm:"unique" json:"code"`
// }

// type Student struct {
// 	UserID    uint   `gorm:"primaryKey" json:"user_id"`
// 	TitleName string `gorm:"not null" json:"title_name"`
// 	FirstName string `gorm:"not null" json:"first_name"`
// 	LastName  string `gorm:"not null" json:"last_name"`
// 	Phone     string `gorm:"unique" json:"phone"`
// 	Code      string `gorm:"unique" json:"code"`
// 	Year      uint   `gorm:"not null" json:"year"`
// 	BranchId  uint   `gorm:"not null" json:"branch_id"`
// 	Branch    Branch `gorm:"foreignKey:BranchId;references:BranchID" json:"branch"`
// }

package entity

import "time"

type User struct {
	UserID   uint     `gorm:"primaryKey;autoIncrement" json:"user_id"`
	Email    string   `gorm:"unique;not null" json:"email"`
	Password string   `gorm:"not null" json:"password"`
	Role     string   `gorm:"default:'user'" json:"role"`
	Student  *Student `gorm:"foreignKey:UserID"` // ใช้ UserID ในการเชื่อมโยง
	Teacher  *Teacher `gorm:"foreignKey:UserID"` // ใช้ UserID ในการเชื่อมโยง
	News     []News   `gorm:"foreignKey:UserID"`
	// News     []News   `gorm:"foreignKey:UserID" json:"news"`  // เชื่อมโยงกับ News
}

type News struct {
	NewsID    uint      `gorm:"primaryKey;autoIncrement" json:"news_id"`
	Title     string    `json:"title"`
	UserID    uint      `json:"user_id"`
	// User      User      `gorm:"foreignKey:UserID;references:UserID" json:"user"`
	Message   string    `json:"message"`
	IsRead    bool      `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

type Teacher struct {
	UserID    uint   `gorm:"primaryKey" json:"user_id"`
	TitleName string `gorm:"not null" json:"title_name"`
	FirstName string `gorm:"not null" json:"first_name"`
	LastName  string `gorm:"not null" json:"last_name"`
	Phone     string `gorm:"unique" json:"phone"`
	Code      string `gorm:"unique" json:"code"`
}

type Student struct {
	UserID    uint   `gorm:"primaryKey" json:"user_id"`
	TitleName string `gorm:"not null" json:"title_name"`
	FirstName string `gorm:"not null" json:"first_name"`
	LastName  string `gorm:"not null" json:"last_name"`
	Phone     string `gorm:"unique" json:"phone"`
	Code      string `gorm:"unique" json:"code"`
	Year      uint   `gorm:"not null" json:"year"`
	BranchId  uint   `gorm:"not null" json:"branch_id"`
	Branch    Branch `gorm:"foreignKey:BranchId;references:BranchID" json:"branch"`
}

type Faculty struct {
	FacultyID   uint    `gorm:"primaryKey;autoIncrement" json:"faculty_id"`
	FacultyCode string  `gorm:"unique;not null" json:"faculty_code"`
	FacultyName string  `gorm:"unique;not null" json:"faculty_name"`
	SuperUser   *uint   `gorm:"default:null" json:"super_user"`
	Teacher     Teacher `gorm:"foreignKey:SuperUser;references:UserID" json:"teacher"`
}

type Branch struct {
	BranchID   uint    `gorm:"primaryKey;autoIncrement" json:"branch_id"`
	BranchCode string  `gorm:"unique;not null" json:"branch_code"`
	BranchName string  `gorm:"unique;not null" json:"branch_name"`
	FacultyId  uint    `gorm:"not null" json:"faculty_id"`
	Faculty    Faculty `gorm:"foreignKey:FacultyId;references:FacultyID" json:"faculty"`
}

type Event struct {
	EventID        uint      `gorm:"primaryKey;autoIncrement" json:"event_id"`
	EventName      string    `gorm:"not null" json:"event_name"`
	Creator        uint      `gorm:"not null" json:"creator"`
	StartDate      time.Time `gorm:"not null" json:"start_date"`
	SchoolYear     uint      `gorm:"not null" json:"school_year" `
	WorkingHour    uint      `gorm:"not null" json:"working_hour"`
	FreeSpace      uint      `gorm:"not null" json:"free_space"`
	Location       string    `gorm:"not null" json:"location"`
	Detail         string    `json:"detail"`
	BranchIDs      string    `gorm:"type:json" json:"branches"`
	Years          string    `gorm:"type:json" json:"years"`
	AllowAllBranch bool      `json:"allow_all_branch"`
	AllowAllYear   bool      `json:"allow_all_year"`
	Status         bool      `gorm:"default:true" json:"status"`
	Teacher        Teacher   `gorm:"foreignKey:Creator;references:UserID" json:"teacher"`
}

type EventInside struct {
	EventId   uint    `gorm:"primaryKey" json:"event_id"`
	User      uint    `gorm:"primaryKey" json:"user_id"`
	Event     Event   `gorm:"foreignKey:EventId;references:EventID;constraint:OnDelete:CASCADE;" json:"event"`
	Student   Student `gorm:"foreignKey:User;references:UserID" json:"student"`
	Certifier uint    `gorm:"default:null" json:"certifier"`
	Teacher   Teacher `gorm:"foreignKey:Certifier;references:UserID" json:"teacher"`
	Status    bool    `json:"status"`
	Comment   string  `json:"comment"`
	File      string  `gorm:"size:255" json:"file"`
}

type Done struct {
	User      uint    `gorm:"primaryKey" json:"user_id"`
	Student   Student `gorm:"foreignKey:User;references:UserID" json:"student"`
	Certifier uint    `gorm:"default:null" json:"certifier"`
	Teacher   Teacher `gorm:"foreignKey:Certifier;references:UserID" json:"teacher"`
	Year      uint    `gorm:"not null" json:"year"`
	Status    bool    `json:"status"`
	Comment   string  `json:"comment"`
}

type EventOutside struct {
	EventID     uint      `gorm:"primaryKey;autoIncrement" json:"event_id"`
	User        uint      `gorm:"primaryKey" json:"user_id"`
	Student     Student   `gorm:"foreignKey:User;references:UserID" json:"student"`
	EventName   string    `gorm:"not null" json:"event_name"`
	SchoolYear  uint      `gorm:"not null" json:"school_year" `
	StartDate   time.Time `gorm:"not null" json:"start_date"`
	Intendant   string    `gorm:"not null" json:"intendent"`
	WorkingHour uint      `json:"working_hour"`
	Location    string    `gorm:"not null" json:"location"`
	File        string    `gorm:"size:255" json:"file"`
}
