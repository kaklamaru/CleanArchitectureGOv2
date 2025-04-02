package entity

import "time"

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
	EventId uint `gorm:"primaryKey" json:"event_id"`
	User    uint `gorm:"primaryKey" json:"user_id"`
	// Event     Event   `gorm:"foreignKey:EventId;references:EventID" json:"event"`
	Event     Event   `gorm:"foreignKey:EventId;references:EventID;constraint:OnDelete:CASCADE;" json:"event"`
	Student   Student `gorm:"foreignKey:User;references:UserID" json:"student"`
	Certifier uint    `gorm:"default:null" json:"certifier"`
	Teacher   Teacher `gorm:"foreignKey:Certifier;references:UserID" json:"teacher"`
	Status    bool    `json:"status"`
	Comment   string  `json:"comment"`
	File      string  `gorm:"size:255" json:"file"`
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
	// Certifier   uint      `gorm:"default:null" json:"certifier"`
	// Teacher     Teacher   `gorm:"foreignKey:Certifier;references:UserID" json:"teacher"`
	// Status      bool      `json:"status"`
	// Comment     string    `json:"comment"`
}

type Done struct {
	User      uint    `gorm:"primaryKey" json:"user_id"`
	Student   Student `gorm:"foreignKey:User;references:UserID" json:"student"`
	Certifier uint    `gorm:"default:null" json:"certifier"`
	Teacher   Teacher `gorm:"foreignKey:Certifier;references:UserID" json:"teacher"`
	Status    bool    `json:"status"`
	Comment   string  `json:"comment"`
}

type News struct {
    NewsID    uint      `gorm:"primaryKey;autoIncrement" json:"news_id"`
    Title     string    `json:"title"`
    UserID    uint      `json:"user_id"`
    User      User      `gorm:"foreignKey:UserID;references:UserID" json:"user"`
    Message   string    `json:"message"`
    IsRead    bool      `gorm:"default:false" json:"is_read"`
    CreatedAt time.Time `json:"created_at"`
}


type Permission struct {
	BranchIDs      string `json:"branches"`
	Years          string `json:"years"`
	AllowAllBranch bool   `json:"allow_all_branch"`
	AllowAllYear   bool   `json:"allow_all_year"`
}
