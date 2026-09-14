package Middleware

import (
	"strings"

	"github.com/Lordlance-lanre/Go-Restaurant-Mysql-DB.git/Helpers"
	"github.com/gofiber/fiber/v3"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func AuthGuard(c fiber.Ctx) error {
	// Prefer Authorization header with Bearer token, fallback to cookie
	authHeader := c.Get("Authorization")
	var token string
	if authHeader != "" {
		// Expect format: "Bearer <token>"
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid Authorization header format",
			})
		}
	} else {
		token = c.Cookies("jwt")
	}

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized Access: missing token",
		})
	}

	if err := Helpers.ValidateToken(token); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized Access: invalid or expired token",
		})
	}

	// Enforce JSON headers: for GET require Accept header, for others require Content-Type
	method := c.Method()
	if method == "GET" {
		accept := c.Get("Accept")
		if accept == "" || !strings.Contains(accept, "application/json") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Accept header must include application/json",
			})
		}
	} else {
		contentType := c.Get("Content-Type")
		if contentType == "" || !strings.Contains(contentType, "application/json") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Content-Type must be application/json",
			})
		}
	}

	return c.Next()

}
