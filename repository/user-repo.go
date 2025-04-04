package repository

import (
	"errors"
	"fmt"
	"go-clean-arch/structure/entity"
	"go-clean-arch/structure/response"
	"strings"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateTeacher(user *entity.User, teacher *entity.Teacher) error
	CreateStudent(user *entity.User, student *entity.Student) error
	GetUserByEmail(email string) (*entity.User, error)
	CreateDones(userID uint,year uint,superUserID uint)error
	GetTotalWorkingHours(userID uint, year uint) (uint, uint, error) 

	// GetTeacherByID(userID uint) (*entity.Teacher, error)
	GetTeacherByID(userID uint) (*entity.Teacher, bool, error)
	GetStudentByID(userID uint) (*entity.Student, error)
	GetAllTeacher() ([]response.TeacherResponse, error)
	GetAllStudent() ([]entity.Student, error)
	GetAllStudentID() ([]uint,error)
	GetSuperUserForStudent(userID uint) (*uint, error) 
	GetStudentsAndYearsByCertifier(certifierID uint) ([]response.StudentYear, error)
	GetDone(userID uint,year uint) (*entity.Done,error) 

	UpdateTeacherByID(teacher *entity.Teacher) error
	UpdateStudentByID(student *entity.Student) error
	UpdateStatusDones(certifierID uint, userID uint, status bool, comment string) error 

	UpdateRoleByID(userID uint, role string) error
	
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateTeacher(user *entity.User, teacher *entity.Teacher) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	teacher.UserID = user.UserID

	if err := tx.Create(teacher).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepository) CreateStudent(user *entity.User, student *entity.Student) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	student.UserID = user.UserID

	if err := tx.Create(student).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// information
func (r *userRepository) GetTeacherByID(userID uint) (*entity.Teacher, bool, error) {
	var teacher entity.Teacher

	// ดึง teacher ก่อน
	if err := r.db.Where("user_id = ?", userID).First(&teacher).Error; err != nil {
		return nil, false, err
	}

	// ตรวจสอบว่ามี faculty ไหนที่ super_user = userID
	var count int64
	if err := r.db.Model(&entity.Faculty{}).
		Where("super_user = ?", userID).
		Count(&count).Error; err != nil {
		return &teacher, false, err
	}

	isSuperUser := count > 0

	return &teacher, isSuperUser, nil
}


func (r *userRepository) GetStudentByID(userID uint) (*entity.Student, error) {
	var student entity.Student
	if err := r.db.Preload("Branch.Faculty").Where("user_id=?", userID).First(&student).Error; err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *userRepository) GetAllTeacher() ([]response.TeacherResponse, error) {
	var teachers []response.TeacherResponse
	if err := r.db.Table("teachers").
		Select("teachers.user_id, teachers.title_name, teachers.first_name, teachers.last_name, teachers.phone, teachers.code, users.role").
		Joins("JOIN users ON users.user_id = teachers.user_id").
		Scan(&teachers).Error; err != nil {
		return nil, fmt.Errorf("repository: failed to retrieve teachers: %w", err)
	}
	return teachers, nil
}

func (r *userRepository) GetAllStudent() ([]entity.Student, error) {
	var students []entity.Student
	if err := r.db.Preload("Branch.Faculty").Find(&students).Error; err != nil {
		return nil, fmt.Errorf("repository: failed to retrieve students: %w", err)
	}
	return students, nil
}

func (r *userRepository) GetAllStudentID() ([]uint,error){
	var student []entity.Student
    err := r.db.Find(&student).Error
    if err != nil {
        return nil, err
    }
    var userIDs []uint
    for _, stdid := range student {
        userIDs = append(userIDs, stdid.UserID)
    }
    return userIDs, nil
}



func (r *userRepository) UpdateTeacherByID(teacher *entity.Teacher) error {
	var existing entity.Teacher
	if err := r.db.First(&existing, "user_id = ?", teacher.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("teacher with ID %d not found", teacher.UserID)
		}
		return err
	}
	existing.TitleName = teacher.TitleName
	existing.FirstName = teacher.FirstName
	existing.LastName = teacher.LastName
	existing.Phone = teacher.Phone
	existing.Code = teacher.Code

	if err := r.db.Save(&existing).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate") || strings.Contains(err.Error(), "duplicate") {
			return fmt.Errorf("phone number %s already exists", teacher.Phone)
		}
		return err
	}
	return nil
}

func (r *userRepository) UpdateStudentByID(student *entity.Student) error {
	var existing entity.Student
	if err := r.db.First(&existing, "user_id =? ", student.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("student with ID %d not found", student.UserID)
		}
		return err
	}
	existing.TitleName = student.TitleName
	existing.FirstName = student.FirstName
	existing.LastName = student.LastName
	existing.Phone = student.Phone
	existing.Code = student.Code
	existing.Year = student.Year
	existing.BranchId = student.BranchId

	if err := r.db.Save(&existing).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate") || strings.Contains(err.Error(), "duplicate") {
			return fmt.Errorf("phone or code already exists: %v", err)
		}
		return err
	}
	return nil
}

func (r *userRepository) UpdateRoleByID(userID uint, role string) error {
	var existing entity.User
	if err := r.db.First(&existing, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user with ID %d not found", userID)
		}
		return err
	}
	existing.Role = role
	if err := r.db.Save(&existing).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetSuperUserForStudent(userID uint) (*uint, error) {
	var superUserID *uint
	err := r.db.Model(&entity.Student{}).
		Select("faculties.super_user").
		Joins("JOIN branches ON students.branch_id = branches.branch_id").
		Joins("JOIN faculties ON branches.faculty_id = faculties.faculty_id").
		Where("students.user_id = ?", userID).
		Scan(&superUserID).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find super user for student: %w", err)
	}
	return superUserID, nil
}

func (r *userRepository) CreateDones(userID uint, year uint, superUserID uint) error {
	var done entity.Done

	err := r.db.Where("user = ? AND year = ?", userID, year).First(&done).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newDone := entity.Done{
				User:    userID,
				Certifier: superUserID,
				Year:      year,
				Status:    false,
			}
			return r.db.Create(&newDone).Error
		}
		return err
	}

	if !done.Status {
		done.Comment = ""
		done.Certifier = superUserID 
		return r.db.Save(&done).Error
	}

	return fmt.Errorf("document for year %d already approved", year)
}

func (r *userRepository) GetTotalWorkingHours(userID uint, year uint) (uint, uint, error) {
	var eventOutsideHours uint
	var eventInsideHours uint

	// รวมชั่วโมงจาก EventOutside
	err := r.db.Model(&entity.EventOutside{}).
		Select("COALESCE(SUM(working_hour), 0)").
		Where("user = ?", userID).
		Where("school_year = ?", year).
		Scan(&eventOutsideHours).Error
	if err != nil {
		return 0, 0, err
	}

	// รวมชั่วโมงจาก EventInside
	err = r.db.Model(&entity.EventInside{}).
		Joins("JOIN events ON event_insides.event_id = events.event_id").
		Select("COALESCE(SUM(events.working_hour), 0)").
		Where("event_insides.user = ?", userID).
		Where("events.school_year = ?", year).
		Where("event_insides.status = ?", true).
		Scan(&eventInsideHours).Error
	if err != nil {
		return 0, 0, err
	}

	return eventOutsideHours, eventInsideHours, nil
}

func (r *userRepository) GetStudentsAndYearsByCertifier(certifierID uint) ([]response.StudentYear, error) {
	var result []response.StudentYear

	err := r.db.Model(&entity.Done{}).
		Joins("JOIN students ON dones.user = students.user_id").
		Joins("JOIN branches ON students.branch_id = branches.branch_id").
		Joins("JOIN faculties ON branches.faculty_id = faculties.faculty_id").
		Where("dones.certifier = ?", certifierID).
		Where("dones.status = ?", false).
		Where("dones.comment = ?","").
		Select("students.user_id, students.title_name, students.first_name, students.last_name, students.phone, students.code, branches.branch_id, branches.branch_name, faculties.faculty_id, faculties.faculty_name, dones.year").
		Scan(&result).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get students and years: %v", err)
	}

	return result, nil
}

// func (r *userRepository) GetDone(userID uint,year uint) (*entity.Done,error){
// 	var dones entity.Done
// 	if err := r.db.Where("user=? AND year=?",userID,year).Find(&dones).Error; err != nil {
// 		return nil, fmt.Errorf("repository: failed to retrieve dones: %w", err)
// 	}
// 	return &dones, nil
// }

func (r *userRepository) GetDone(userID uint, year uint) (*entity.Done, error) {
	var done entity.Done
	err := r.db.Where("user = ? AND year = ?", userID, year).First(&done).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // ยังไม่ส่งข้อมูล
	}
	if err != nil {
		return nil, err
	}
	return &done, nil
}


func (r *userRepository) UpdateStatusDones(certifierID uint, userID uint, status bool, comment string) error {
	updates := map[string]interface{}{
		"status":  status,
		"comment": comment,
	}
	if err := r.db.Model(&entity.Done{}).
		Where("certifier = ? AND user = ?", certifierID, userID).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}
	return nil
}