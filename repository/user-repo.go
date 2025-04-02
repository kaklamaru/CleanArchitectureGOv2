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

	GetTeacherByID(userID uint) (*entity.Teacher, error)
	GetStudentByID(userID uint) (*entity.Student, error)
	GetAllTeacher() ([]response.TeacherResponse, error)
	GetAllStudent() ([]entity.Student, error)
	GetAllStudentID() ([]uint,error)

	UpdateTeacherByID(teacher *entity.Teacher) error
	UpdateStudentByID(student *entity.Student) error
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
func (r *userRepository) GetTeacherByID(userID uint) (*entity.Teacher, error) {
	var teacher entity.Teacher
	if err := r.db.Where("user_id = ?", userID).First(&teacher).Error; err != nil {
		return nil, err
	}
	return &teacher, nil
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
