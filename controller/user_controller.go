package controller

import (
	"fmt"
	"go-clean-arch/pkg/utility"
	"go-clean-arch/structure/request"
	"go-clean-arch/usecase"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userUsecase usecase.UserUsecase
}

func NewUserController(userUsecase usecase.UserUsecase) *UserController {
	return &UserController{userUsecase: userUsecase}
}

func (c *UserController) RegisterTeacher(ctx *fiber.Ctx) error {
	var req request.RegisterTeacher

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	if err := c.userUsecase.CreateTeacher(&req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "registered successfully",
	})
}

func (c *UserController) RegisterStudent(ctx *fiber.Ctx) error {
	var req request.RegisterStudent

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	if err := c.userUsecase.CreateStudent(&req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "registered successfully",
	})
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request",
		})
	}

	token, role, err := c.userUsecase.GetUserByEmail(req.Email, req.Password)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "token",                        // ชื่อคุกกี้
		Value:    token,                          // ค่า JWT
		Expires:  time.Now().Add(24 * time.Hour), // วันหมดอายุ (24 ชั่วโมง)
		HTTPOnly: false,                          // ป้องกันการเข้าถึงผ่าน JavaScript
		Secure:   false,                          // ใช้งานเฉพาะ HTTPS (แนะนำสำหรับ Production)
		SameSite: "Lax",                          // นโยบาย SameSite
	})

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"role":    role,
	})
}

func (c *UserController) GetUserByClaims(ctx *fiber.Ctx) error {
	claims, err := utility.GetClaimsFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "failed to retrieve claims",
		})
	}
	userData, err := c.userUsecase.GetUserByClaims(claims)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(userData)
}

func (c *UserController) GetAllTeacher(ctx *fiber.Ctx) error {
	allTeacher, err := c.userUsecase.GetAllTeacher()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to get teachers: %v", err),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(allTeacher)
}

func (c *UserController) GetAllStudent(ctx *fiber.Ctx) error {
	allStudent, err := c.userUsecase.GetAllStudent()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to get students: %v", err),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(allStudent)
}

func (c *UserController) UpdateTeacher(ctx *fiber.Ctx) error {
	var req request.RegisterTeacher
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}
	claims, err := utility.GetClaimsFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "failed to retrieve claims",
		})
	}
	if err := c.userUsecase.UpdateTeacherByID(&req, claims); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update teacher",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "teacher updated successfully",
	})

}

func (c *UserController) UpdateStudent(ctx *fiber.Ctx) error {
	var req request.RegisterStudent
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}
	claims, err := utility.GetClaimsFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "failed to retrieve claims",
		})
	}
	if err := c.userUsecase.UpdateStudentByID(&req, claims); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "fialed to update student",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "student updated successfully",
	})
}

func (c *UserController) UpdateRoleByID(ctx *fiber.Ctx) error {
	var req struct {
		UserID uint   `json:"user_id"`
		Role   string `json:"role"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"bad request",
		})
	}
	if err:= c.userUsecase.UpdateRoleByID(req.UserID,req.Role);err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":"Failed to update role",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"massage":"role updated successfully",
	})
}

func (c *UserController) SendEvent(ctx *fiber.Ctx) error{
	claims, err := utility.GetClaimsFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "failed to retrieve claims",
		})
	}
	idStr := ctx.Params("year")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id format",
		})
	}
	year := uint(id)

	if err := c.userUsecase.SendEvent(year,claims); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Send all event successfully",
	})
}

func (c *UserController) GetStudentsAndYearsByCertifier(ctx *fiber.Ctx) error{
	claims, err := utility.GetClaimsFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "failed to retrieve claims",
		})
	}
	allStudent, err := c.userUsecase.GetStudentsAndYearsByCertifier(claims)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to get: %v", err),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(allStudent)

}


func (c *UserController) UpdateStatusDones(ctx *fiber.Ctx) error {
	var req struct {
		Status  bool   `json:"status"`
		Comment string `json:"comment"`
	}

	idStr := ctx.Params("userid")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid UserID",
		})
	}
	userID := uint(idInt)
	claims, err := utility.GetClaimsFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "failed to retrieve claims",
		})
	}
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return fmt.Errorf("invalid user_id in claims")
	}
	certifierID := uint(userIDFloat)


	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if err := c.userUsecase.UpdateStatusDones(certifierID, userID, req.Status, req.Comment); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Checking all event successfully",
	})
}