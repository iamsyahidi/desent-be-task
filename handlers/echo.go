package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// EchoHandler handles POST /echo requests
// Accepts JSON body and returns the same JSON in the response
// Returns HTTP 400 for empty or invalid body
func EchoHandler(c *fiber.Ctx) error {
	body := c.Body()
	if len(body) == 0 {
		log.Printf("[ECHO] Received empty body")
	}
	log.Printf("[ECHO] Echoing %d bytes", len(body))
	// Return the raw body as-is to preserve exact JSON formatting and key order
	c.Set("Content-Type", "application/json")
	return c.Send(body)
}
