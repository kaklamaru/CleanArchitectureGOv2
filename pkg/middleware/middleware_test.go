package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go-clean-arch/pkg/jwt"
)

// TestJWTMiddleware tests the JWTMiddleware
func TestJWTMiddleware(t *testing.T) {
	// Mock JWTService with a simple secret key
	jwtService := &jwt.JWTService{
		SecretKey: "test-secret",
	}

	// Create a Fiber app
	app := fiber.New()
	app.Use(JWTMiddleware(jwtService))

	// Create a protected route
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Authorized",
		})
	})

	t.Run("No token provided", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		resp, _ := app.Test(req, -1)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Invalid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Cookie", "token=invalidtoken")
		resp, _ := app.Test(req, -1)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Valid token", func(t *testing.T) {
		token, _ := jwtService.GenerateJWT(1, "admin")

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Cookie", "token="+token)

		resp, _ := app.Test(req, -1)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestRoleMiddleware tests the RoleMiddleware
func TestRoleMiddleware(t *testing.T) {
	t.Run("Authorized role", func(t *testing.T) {
		// Create a Fiber app
		app := fiber.New()

		// Mock middleware to set role
		app.Use(func(c *fiber.Ctx) error {
			c.Locals("role", "admin") // Mock role as "admin"
			return c.Next()
		})

		// Use RoleMiddleware to allow admin and superadmin
		app.Use(RoleMiddleware("superadmin", "admin"))

		// Create a protected route
		app.Get("/admin", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"message": "Welcome, Admin or Superadmin!",
			})
		})

		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		resp, _ := app.Test(req, -1)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Unauthorized role", func(t *testing.T) {
		// Create a Fiber app
		app := fiber.New()

		// Mock middleware to set role
		app.Use(func(c *fiber.Ctx) error {
			c.Locals("role", "student") // Mock role as "student"
			return c.Next()
		})

		// Use RoleMiddleware to allow admin and superadmin
		app.Use(RoleMiddleware("superadmin", "admin"))

		// Create a protected route
		app.Get("/admin", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"message": "Welcome, Admin or Superadmin!",
			})
		})

		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		resp, _ := app.Test(req, -1)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}
