package middleware

import (
	"log"
	"strings"

	"personal/desent-be-task/models"
	"personal/desent-be-task/service"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware validates JWT tokens in the Authorization header
func AuthMiddleware(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			log.Printf("[AUTH] Missing Authorization header")
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error: "missing authorization header",
			})
		}

		// Check Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Printf("[AUTH] Invalid Authorization header format")
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error: "invalid authorization header format",
			})
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			log.Printf("[AUTH] Empty token")
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error: "empty token",
			})
		}

		// Validate token
		claims, err := authService.ValidateToken(token)
		if err != nil {
			log.Printf("[AUTH] Invalid token: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error: "invalid token",
			})
		}

		// Check if token is valid
		if !claims.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error: "invalid or expired token",
			})
		}

		// Token is valid, proceed to next handler
		return c.Next()
	}
}
