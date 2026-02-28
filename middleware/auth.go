package middleware

import (
	"strings"

	"personal/desent-be-task/models"
	"personal/desent-be-task/service"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware creates a middleware that validates JWT tokens
// Extracts Bearer token from Authorization header and validates it
// Returns 401 if token is missing, invalid, or expired
func AuthMiddleware(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error: "missing authorization header",
			})
		}

		// Extract token from "Bearer <token>" format
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Check if Bearer prefix was present
		if tokenString == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error: "invalid authorization header format",
			})
		}

		// Validate token using AuthService
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error: "invalid or expired token",
			})
		}

		// Check if token is valid
		if !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error: "invalid or expired token",
			})
		}

		// Token is valid, proceed to next handler
		return c.Next()
	}
}
