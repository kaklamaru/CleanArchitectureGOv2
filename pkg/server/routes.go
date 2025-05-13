package server

import (
	"go-clean-arch/controller"
	"go-clean-arch/database"
	"go-clean-arch/pkg/jwt"
	"go-clean-arch/pkg/middleware"
	"go-clean-arch/repository"
	"go-clean-arch/usecase"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, jwt *jwt.JWTService, db database.Database) {
	// repository
	userRepo := repository.NewUserRepository(db.GetDB())
	facBranRepo := repository.NewFacultyRepositiry(db.GetDB())
	eventRepo := repository.NewEventRepository(db.GetDB())

	// usecase
	userUsecase := usecase.NewUserUsecase(userRepo, *jwt)
	facBranUsecase := usecase.NewFacultyUsecase(facBranRepo)
	eventUsecase := usecase.NewEventUsecase(userRepo, facBranRepo, eventRepo)

	// controller
	userContro := controller.NewUserController(userUsecase)
	facBranContro := controller.NewFacultyController(facBranUsecase)
	eventContro := controller.NewEventController(eventUsecase)

	// login&register
	app.Post("/register/teacher", userContro.RegisterTeacher)
	app.Post("/register/student", userContro.RegisterStudent)
	app.Post("/login", userContro.Login)

	app.Post("/logout", func(c *fiber.Ctx) error {
		c.ClearCookie("token", "/")

		c.Cookie(&fiber.Cookie{
			Name:     "token",
			Value:    "",
			Path:     "/",
			Domain:   ".activitiesmanagement.online",
			Expires:  time.Now().Add(-1 * time.Hour),
			MaxAge:   -1,
			Secure:   true,
			HTTPOnly: true,
		})
		return c.SendStatus(fiber.StatusOK)
	})

	// middleware
	protected := app.Group("/protected", middleware.JWTMiddleware(jwt))
	admin := protected.Group("/admin", middleware.RoleMiddleware("superadmin", "admin"))
	admin.Get("/admin", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome, Admin or Superadmin!",
		})
	})
	teacher := protected.Group("/teacher", middleware.RoleMiddleware("superadmin", "admin", "teacher","officer"))
	teacher.Get("/teacher", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome, Teacher",
		})
	})
	student := protected.Group("student", middleware.RoleMiddleware("student"))
	student.Get("/student", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome, Student",
		})
	})

	// faculty&branch
	admin.Post("/faculty", facBranContro.CreateFaculty)
	app.Get("/faculties", facBranContro.GetAllFaculties)
	admin.Put("/faculty/:id", facBranContro.UpdateFacultyByID)
	admin.Delete("/faculty/:id", facBranContro.DeleteFacultyByID)
	// admin.Put("/faculty-staff/:id/:userid", facBranContro.UpdateSuperUser)

	admin.Post("/branch", facBranContro.CreateBranch)
	app.Get("/branches", facBranContro.GetAllBranches)
	admin.Put("/branch/:id", facBranContro.UpdateBranchByID)
	admin.Delete("/branch/:id", facBranContro.DeleteBranchByID)

	// user
	protected.Get("/userbyclaim", userContro.GetUserByClaims)
	teacher.Get("/allteacher", userContro.GetAllTeacher)
	teacher.Get("/allstudent", userContro.GetAllStudent)
	teacher.Put("/personalinfo", userContro.UpdateTeacher)
	student.Put("/personalinfo", userContro.UpdateStudent)
	admin.Put("/studentinfo", userContro.UpdateStudent)
	admin.Put("/role", userContro.UpdateRoleByID)

	// events
	teacher.Post("/event", eventContro.CreateEvent)
	app.Get("/events", eventContro.GetAllEvent)
	app.Get("/event/:id", eventContro.GetEventByID)
	teacher.Get("/myevents", eventContro.MyEvent)
	teacher.Get("/allowedevents", eventContro.AllAllowedEvent)
	teacher.Get("/currentevents", eventContro.AllCurrentEvent)
	teacher.Put("/status/:id", eventContro.ToggleEventStatus)
	teacher.Delete("/event/:id", eventContro.DeleteEventByID)
	teacher.Put("/event/:id", eventContro.UpdateEventByID)
	student.Get("myevents/:year", eventContro.MyEventThisYear)
	student.Get("myevents", eventContro.MyEventInside)

	// inside
	student.Post("/joinevent/:id", eventContro.JoinEvent)
	student.Delete("/unjoinevent/:id", eventContro.UnJoinEvent)
	student.Put("/upload/:id", eventContro.UploadFile)
	protected.Get("/file/:eventid/:userid", eventContro.GetFile)
	teacher.Get("/checklist/:id", eventContro.MyChecklist)
	teacher.Put("/check/:eventid/:userid", eventContro.UpdateEventStatusAndComment) //

	// outside
	student.Post("/outside", eventContro.CreateEventOutside)
	student.Delete("/outside/:id", eventContro.DeleteEventOutsideByID)
	student.Get("/download/:id", eventContro.CreateFile)
	student.Put("/upload-outside/:id", eventContro.UploadFileOutside)
	protected.Get("/file-outside/:eventid/:userid", eventContro.GetFileOutside)

	student.Post("/send-event/:year", userContro.SendEvent)

	teacher.Get("/superuser-check", userContro.GetStudentsAndYearsByCertifier)

	teacher.Get("/all-event/:userid/:year", eventContro.AllSendEventThisYear)
	teacher.Put("/check-all-event/:userid", userContro.UpdateStatusDones) //

	protected.Get("/news", userContro.GetNews)
	protected.Put("/news/:newsid", userContro.MarkNewsAsRead)

	admin.Get("/done", userContro.GetDoneFiltered)

	teacher.Get("/dashboard/:year", eventContro.DashboardData)
	teacher.Get("/event-7day", eventContro.GetUpcomingEventsWithin7Days)

	app.Get("/devfile", func(ctx *fiber.Ctx) error {
		// ดึง path มาจาก query param
		filePath := ctx.Query("path")
		if filePath == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "missing 'path' parameter",
			})
		}

		// แปลง path dev จาก "/app/uploads/..." เป็น "./uploads/..."
		// (ถ้าระบบจริงเก็บใน local folder แบบ relative)
		if strings.HasPrefix(filePath, "/app") {
			filePath = "." + strings.TrimPrefix(filePath, "/app")
		}

		// ตรวจสอบว่าไฟล์มีอยู่จริงไหม
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "file not found",
			})
		}

		return ctx.SendFile(filePath)
	})
}
