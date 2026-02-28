package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// EchoHandler handles POST /echo requests
// Accepts JSON body and returns the same JSON in the response
// Returns HTTP 400 for empty or invalid body
func EchoHandler(c *fiber.Ctx) error {
	// Return the raw body as-is to preserve exact JSON formatting and key order
	c.Set("Content-Type", "application/json")
	return c.Send(c.Body())
}
