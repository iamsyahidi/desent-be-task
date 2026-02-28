package middleware

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// LoggerMiddleware logs incoming API requests with method, path, status, duration, headers, and body
func LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Capture request details before processing
		method := c.Method()
		path := c.Path()
		ip := c.IP()
		reqBody := string(c.Body())
		reqHeaders := formatRequestHeaders(c)

		// Log incoming request
		log.Printf("[REQUEST] --> %s %s | IP: %s", method, path, ip)
		log.Printf("[REQUEST] Headers: %s", reqHeaders)
		if len(reqBody) > 0 {
			// Truncate body if too long (max 1000 chars)
			if len(reqBody) > 1000 {
				reqBody = reqBody[:1000] + "...(truncated)"
			}
			log.Printf("[REQUEST] Body: %s", reqBody)
		}

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get status code
		status := c.Response().StatusCode()

		// Get response headers and body
		respHeaders := formatResponseHeaders(c)
		respBody := string(c.Response().Body())

		// Log response details
		log.Printf("[RESPONSE] <-- %s %s | Status: %d | Duration: %v",
			method,
			path,
			status,
			duration,
		)
		log.Printf("[RESPONSE] Headers: %s", respHeaders)
		if len(respBody) > 0 {
			// Truncate body if too long (max 1000 chars)
			if len(respBody) > 1000 {
				respBody = respBody[:1000] + "...(truncated)"
			}
			log.Printf("[RESPONSE] Body: %s", respBody)
		}

		return err
	}
}

// formatRequestHeaders formats request headers as a string for logging
func formatRequestHeaders(c *fiber.Ctx) string {
	var headers []string
	c.Request().Header.VisitAll(func(key, value []byte) {
		keyStr := string(key)
		// Skip sensitive headers or truncate authorization
		if strings.EqualFold(keyStr, "Authorization") {
			valStr := string(value)
			if len(valStr) > 20 {
				valStr = valStr[:20] + "..."
			}
			headers = append(headers, keyStr+"="+valStr)
		} else {
			headers = append(headers, keyStr+"="+string(value))
		}
	})
	return strings.Join(headers, ", ")
}

// formatResponseHeaders formats response headers as a string for logging
func formatResponseHeaders(c *fiber.Ctx) string {
	var headers []string
	c.Response().Header.VisitAll(func(key, value []byte) {
		keyStr := string(key)
		// Skip sensitive headers
		if strings.EqualFold(keyStr, "Set-Cookie") {
			valStr := string(value)
			if len(valStr) > 20 {
				valStr = valStr[:20] + "..."
			}
			headers = append(headers, keyStr+"="+valStr)
		} else {
			headers = append(headers, keyStr+"="+string(value))
		}
	})
	return strings.Join(headers, ", ")
}
