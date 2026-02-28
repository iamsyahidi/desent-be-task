package middleware

import (
	"fmt"
	"personal/desent-be-task/models"

	"github.com/gofiber/fiber/v2"
)

// RecoveryMiddleware creates a middleware that recovers from panics
// Returns HTTP 500 with a generic error message when a panic occurs
// This prevents the server from crashing and provides consistent error responses
func RecoveryMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				// Log the panic for debugging (in production, use proper logging)
				fmt.Printf("Panic recovered: %v\n", r)

				// Return a generic error response
				// Don't expose internal panic details to clients
				c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
					Error: "internal server error",
				})
			}
		}()

		// Continue to next handler
		return c.Next()
	}
}
