package handlers

import (
	"personal/desent-be-task/models"

	"github.com/gofiber/fiber/v2"
)

// EchoHandler handles POST /echo requests
// Accepts JSON body and returns the same JSON in the response
// Returns HTTP 400 for empty or invalid body
func EchoHandler(c *fiber.Ctx) error {
	var body map[string]any

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(body)
}
