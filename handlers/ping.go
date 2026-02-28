package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// PingHandler handles GET /ping requests
// Returns JSON {"status": "ok"} with HTTP 200
func PingHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}
