package entity

type Faculty struct {
	FacultyID   uint   `gorm:"primaryKey;autoIncrement" json:"faculty_id"`
	FacultyCode string `gorm:"unique;not null" json:"faculty_code"`
	FacultyName string `gorm:"unique;not null" json:"faculty_name"`
	SuperUser   *uint   `gorm:"default:null" json:"super_user"`
	Teacher   Teacher `gorm:"foreignKey:SuperUser;references:UserID" json:"teacher"`
}

type Branch struct {
	BranchID   uint    `gorm:"primaryKey;autoIncrement" json:"branch_id"`
	BranchCode string  `gorm:"unique;not null" json:"branch_code"`
	BranchName string  `gorm:"unique;not null" json:"branch_name"`
	FacultyId  uint    `gorm:"not null" json:"faculty_id"`
	Faculty    Faculty `gorm:"foreignKey:FacultyId;references:FacultyID" json:"faculty"`
}
