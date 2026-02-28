package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// PingHandler handles GET /ping requests
// Returns JSON {"success": true} with HTTP 200
func PingHandler(c *fiber.Ctx) error {
	log.Printf("[PING] Health check from %s", c.IP())
	return c.JSON(fiber.Map{"success": true})
}
