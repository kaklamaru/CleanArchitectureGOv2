package controller

import (
	"fmt"
	"go-clean-arch/structure/entity"
	"go-clean-arch/usecase"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type FacultyBranchController struct {
	facultyUsecase usecase.FacultyBranchUsecase
}

func NewFacultyController(facultyUsecase usecase.FacultyBranchUsecase) *FacultyBranchController {
	return &FacultyBranchController{facultyUsecase: facultyUsecase}
}

func (c *FacultyBranchController) CreateFaculty(ctx *fiber.Ctx) error {
	var req entity.Faculty

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	if err := c.facultyUsecase.CreateFaculty(&req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "faculty created successfully",
	})

}

func (c *FacultyBranchController) GetAllFaculties(ctx *fiber.Ctx) error {
	faculties, err := c.facultyUsecase.GetAllFaculties()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to retrieve faculties",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(faculties)
}

func (c *FacultyBranchController) UpdateFacultyByID(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id format",
		})
	}
	facultyID := uint(id)
	var req entity.Faculty
	req.FacultyID = facultyID

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.FacultyCode == "" || req.FacultyName == "" {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "faculty code and name are required",
		})
	}

	if err := c.facultyUsecase.UpdateFacultyByID(&req); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fmt.Sprintf("faculty with ID %d not found", facultyID),
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to update faculty with ID %d : %s", facultyID, err),
		})
	}

	// สำเร็จ
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("faculty with ID %d updated successfully", facultyID),
	})
}

func (c *FacultyBranchController) DeleteFacultyByID(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id format",
		})
	}
	facultyID := uint(id)
	if err := c.facultyUsecase.DeleteFacultyByID(facultyID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fmt.Sprintf("faculty with ID %d not found", facultyID),
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to delete faculty with ID %d : %s", facultyID, err),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("faculty with ID %d deleted successfully", facultyID),
	})
}

// func (c *FacultyBranchController) UpdateSuperUser(ctx *fiber.Ctx) error {
// 	idStr := ctx.Params("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "invalid id format",
// 		})
// 	}
// 	facultyID := uint(id)

// 	idStr = ctx.Params("userid")
// 	id, err = strconv.Atoi(idStr)
// 	if err != nil {
// 		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "invalid id format",
// 		})
// 	}
// 	userID := uint(id)

// 	if err := c.facultyUsecase.UpdateSuperUser(facultyID, userID); err != nil {
// 		if strings.Contains(err.Error(), "not found") {
// 			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 				"error": err.Error(),
// 			})
// 		}
// 		if strings.Contains(err.Error(), "not a teacher") {
// 			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 				"error": err.Error(),
// 			})
// 		}
// 		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "failed to update super_user",
// 		})
// 	}

// 	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": fmt.Sprintf("super_user of faculty with ID %d updated successfully", facultyID),
// 	})

// }

// branch ------------------------------------------------------------------
func (c *FacultyBranchController) CreateBranch(ctx *fiber.Ctx) error {
	var req entity.Branch

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	if err := c.facultyUsecase.CreateBranch(&req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "branch created successfully",
	})
}

func (c *FacultyBranchController) GetAllBranches(ctx *fiber.Ctx) error {
	branches, err := c.facultyUsecase.GetAllBranches()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve branches",
		})
	}
	if len(branches) == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No branches found",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(branches)
}

func (c *FacultyBranchController) UpdateBranchByID(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id format",
		})
	}
	branchID := uint(id)
	var req entity.Branch
	req.BranchID = branchID

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.BranchCode == "" || req.BranchName == "" || req.FacultyId == 0 {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "branch code, branch name, and faculty id are required",
		})
	}

	if err := c.facultyUsecase.UpdateBranchByID(&req); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fmt.Sprintf("branch with ID %d not found", branchID),
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to update branch with ID %d : %s", branchID, err),
		})
	}

	// สำเร็จ
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("branch with ID %d updated successfully", branchID),
	})
}

func (c *FacultyBranchController) DeleteBranchByID(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id format",
		})
	}
	branchID := uint(id)

	if err := c.facultyUsecase.DeleteBranchByID(branchID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fmt.Sprintf("branch with ID %d not found", branchID),
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to delete branch with ID %d : %s", branchID, err),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("branch with ID %d deleted successfully", branchID),
	})
}
