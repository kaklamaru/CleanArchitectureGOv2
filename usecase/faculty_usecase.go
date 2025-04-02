package usecase

import (
	"go-clean-arch/repository"
	"go-clean-arch/structure/entity"
)

type FacultyBranchUsecase interface {
	CreateFaculty(faculty *entity.Faculty) error
	GetAllFaculties() ([]entity.Faculty, error)
	UpdateFacultyByID(faculty *entity.Faculty) error
	DeleteFacultyByID(facultyID uint) error
	// UpdateSuperUser(facultyID uint, superUserID uint) error

	// branch
	CreateBranch(branch *entity.Branch) error
	GetAllBranches() ([]entity.Branch, error)
	UpdateBranchByID(branch *entity.Branch) error
	DeleteBranchByID(branchID uint) error
}


type facultyBranchUsecase struct {
	facultyRepo repository.FacultyBranchRepository
}

func NewFacultyUsecase(facultyRepo repository.FacultyBranchRepository) FacultyBranchUsecase {
	return &facultyBranchUsecase{
		facultyRepo: facultyRepo,
	}
}

func (u *facultyBranchUsecase) CreateFaculty(faculty *entity.Faculty) error {
	return u.facultyRepo.CreateFaculty(faculty)
}

func (u *facultyBranchUsecase) GetAllFaculties() ([]entity.Faculty, error) {
	return u.facultyRepo.GetAllFaculties()
}

func (u *facultyBranchUsecase) UpdateFacultyByID(faculty *entity.Faculty) error {
	return u.facultyRepo.UpdateFacultyByID(faculty)
}

func (u *facultyBranchUsecase) DeleteFacultyByID(facultyID uint) error {
	return u.facultyRepo.DeleteFacultyByID(facultyID)
}

// func (u *facultyBranchUsecase) UpdateSuperUser(facultyID uint, superUserID uint) error {
// 	return u.facultyRepo.UpdateSuperUser(facultyID, &superUserID)
// }


// branch ------------------------------------------------------------------
func (u *facultyBranchUsecase) CreateBranch(branch *entity.Branch) error{
	return u.facultyRepo.CreateBranch(branch)
}

func (u *facultyBranchUsecase) GetAllBranches() ([]entity.Branch,error){
	return u.facultyRepo.GetAllBranches()
}

func (u *facultyBranchUsecase) UpdateBranchByID(branch *entity.Branch) error{
	return u.facultyRepo.UpdateBranchByID(branch)
}
func (u *facultyBranchUsecase) DeleteBranchByID(branchID uint) error{
	return u.facultyRepo.DeleteBranchByID(branchID)
}