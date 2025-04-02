package repository

import (
	"errors"
	"fmt"
	"go-clean-arch/structure/entity"

	"gorm.io/gorm"
)

type FacultyBranchRepository interface {
	CreateFaculty(faculty *entity.Faculty) error
	GetAllFaculties() ([]entity.Faculty, error)
	UpdateFacultyByID(faculty *entity.Faculty) error
	DeleteFacultyByID(facultyID uint) error
	// UpdateSuperUser(facultyID uint, superUserID *uint) error

	// branch
	CreateBranch(branch *entity.Branch) error
	GetAllBranches() ([]entity.Branch, error)
	UpdateBranchByID(branch *entity.Branch) error
	DeleteBranchByID(branchID uint) error
	BranchExists(branchID uint) (bool, error)
}

type facultyBranchRepository struct {
	db *gorm.DB
}

func NewFacultyRepositiry(db *gorm.DB) FacultyBranchRepository {
	return &facultyBranchRepository{
		db: db,
	}
}

func (r *facultyBranchRepository) CreateFaculty(faculty *entity.Faculty) error {
	if faculty.SuperUser != nil && *faculty.SuperUser == 0 {
		faculty.SuperUser = nil
	}
	return r.db.Create(faculty).Error
}

func (r *facultyBranchRepository) GetAllFaculties() ([]entity.Faculty, error) {
	var faculties []entity.Faculty
	if err := r.db.Preload("Teacher").Find(&faculties).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch faculties: %w", err)
	}
	return faculties, nil
}

func (r *facultyBranchRepository) UpdateFacultyByID(faculty *entity.Faculty) error {
	var existing entity.Faculty
	if err := r.db.First(&existing, "faculty_id = ?", faculty.FacultyID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("faculty with ID %d not found", faculty.FacultyID)
		}
		return err
	}

	existing.FacultyCode = faculty.FacultyCode
	existing.FacultyName = faculty.FacultyName
	if faculty.SuperUser != nil && *faculty.SuperUser == 0 {
		existing.SuperUser = nil
	} else {
		existing.SuperUser = faculty.SuperUser
	}

	return r.db.Save(&existing).Error
}

func (r *facultyBranchRepository) DeleteFacultyByID(facultyID uint) error {
	var existingFaculty entity.Faculty
	if err := r.db.First(&existingFaculty, "faculty_id = ?", facultyID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("faculty with ID %d not found", facultyID)
		}
		return err
	}

	if err := r.db.Delete(existingFaculty).Error; err != nil {
		return fmt.Errorf("failed to delete faculty with ID %d: %w", existingFaculty.FacultyID, err)
	}

	return nil
}

// func (r *facultyBranchRepository) UpdateSuperUser(facultyID uint, superUserID *uint) error {
// 	var existingFaculty entity.Faculty
// 	if err := r.db.First(&existingFaculty, "faculty_id = ?", facultyID).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return fmt.Errorf("faculty with ID %d not found", facultyID)
// 		}
// 		return fmt.Errorf("error fetching faculty: %w", err)
// 	}

// 	if superUserID == nil || *superUserID == 0 {
// 		existingFaculty.SuperUser = nil
// 	} else {
// 		var user entity.Teacher
// 		if err := r.db.First(&user, "user_id = ?", *superUserID).Error; err != nil {
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				return fmt.Errorf("user with ID %d not found", *superUserID)
// 			}
// 			return fmt.Errorf("error fetching user: %w", err)
// 		}
// 		existingFaculty.SuperUser = superUserID
// 	}

// 	if err := r.db.Save(&existingFaculty).Error; err != nil {
// 		return fmt.Errorf("error updating superuser: %w", err)
// 	}

// 	return nil
// }

// Branch----------------------------------------------------------------------
func (r *facultyBranchRepository) CreateBranch(branch *entity.Branch) error {
	return r.db.Create(branch).Error
}

func (r *facultyBranchRepository) GetAllBranches() ([]entity.Branch, error) {
	var branches []entity.Branch
	if err := r.db.Preload("Faculty.Teacher").Find(&branches).Error; err != nil {
		return nil, err
	}
	return branches, nil
}

func (r *facultyBranchRepository) UpdateBranchByID(branch *entity.Branch) error {
	var existing entity.Branch
	if err := r.db.First(&existing, "branch_id = ?", branch.BranchID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("branch with ID %d not found", branch.BranchID)
		}
		return err
	}

	existing.BranchCode = branch.BranchCode
	existing.BranchName = branch.BranchName
	existing.FacultyId = branch.FacultyId

	return r.db.Save(&existing).Error
}

func (r *facultyBranchRepository) DeleteBranchByID(branchID uint) error {
	var existingBranch entity.Branch
	if err := r.db.First(&existingBranch, "branch_id = ?", branchID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("branch with ID %d not found", branchID)
		}
		return err
	}

	if err := r.db.Delete(existingBranch).Error; err != nil {
		return fmt.Errorf("failed to delete branch with ID %d: %w", existingBranch.BranchID, err)
	}

	return nil
}

func (r *facultyBranchRepository) BranchExists(branchID uint) (bool, error) {
	var branch entity.Branch
	if err := r.db.Where("branch_id = ?", branchID).First(&branch).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if branch with ID %d exists: %w", branchID, err)
	}
	return true, nil
}

func (r *facultyBranchRepository) GetAllBranchesByFaculty(facultyId int) ([]entity.Branch, error) {
	var branches []entity.Branch
	result := r.db.Where("faculty_id = ?", facultyId).Find(&branches)
	if result.Error != nil {
		return nil, result.Error
	}
	return branches, nil
}

func (r *facultyBranchRepository) GetBranch(id uint) (*entity.Branch, error) {
	var branch entity.Branch
	if err := r.db.Preload("Faculty").First(&branch, id).Error; err != nil {
		return nil, err
	}
	return &branch, nil
}
