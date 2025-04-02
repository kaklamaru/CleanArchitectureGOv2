package server

import (
	"fmt"
	"go-clean-arch/config"
	"go-clean-arch/database"
	"go-clean-arch/pkg/jwt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// โครงสร้างสำหรับเซิร์ฟเวอร์ Fiber
type fiberServer struct {
	app  *fiber.App
	port int
}

func NewServer(cfg *config.Config, db database.Database, jwt *jwt.JWTService) (Server, error) {
	if cfg.ServerPort == 0 {
		return nil, fmt.Errorf("Server port not specified in config")
	}
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // 10MB
	})
	setupCors(app)

	// กำหนด static files
	app.Static("/uploads", "./uploads")

	// กำหนด middleware สำหรับการกู้คืนจาก panic
	app.Use(recover.New())

	// กำหนด middleware สำหรับการบันทึก log ของการร้องขอ
	app.Use(logger.New())

	SetupRoutes(app,jwt,db)

	return &fiberServer{
		app:  app,
		port: cfg.ServerPort,
	}, nil
}

func setupCors(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000/, http://127.0.0.1:8080",
		// AllowOrigins:     cfg.CorsAllowOrigins,  // ใช้ค่าจาก config หรือ environment
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))
}

func (s *fiberServer) StartServer() error {
	serverUrl := fmt.Sprintf(":%d", s.port)
	return s.app.Listen(serverUrl)
}
