package middleware

import (
	"go-clean-arch/pkg/jwt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware(jwt *jwt.JWTService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// ดึง token จาก Cookie
		tokenString := ctx.Cookies("token")
		if tokenString == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing token",
			})
		}

		// ตรวจสอบและ decode token
		claims, err := jwt.ValidateJWT(tokenString)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// ตรวจสอบว่ามีข้อมูล role ใน claims หรือไม่
		role, ok := claims["role"].(string)
		if !ok || role == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token: role not found",
			})
		}

		// เก็บ claims และ role ลงใน Locals
		ctx.Locals("claims", claims)
		ctx.Locals("role", role) // ใช้ใน RoleMiddleware
		return ctx.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ดึง role จาก Locals
		userRole := c.Locals("role")
		if userRole == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: role not found",
			})
		}

		// แปลงเป็น string
		role, ok := userRole.(string)
		if !ok || role == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: invalid role format",
			})
		}

		// ตรวจสอบว่า role อยู่ใน allowedRoles หรือไม่
		for _, allowedRole := range allowedRoles {
			if strings.EqualFold(role, allowedRole) { // เช็คแบบ case insensitive
				return c.Next()
			}
		}

		// ถ้าไม่ตรงกับ allowedRoles ใดเลย
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden: you do not have permission to access this resource",
		})
	}
}
