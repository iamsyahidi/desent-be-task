package handlers

import (
	"personal/desent-be-task/models"
	"personal/desent-be-task/service"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Credentials represents the login credentials
type Credentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// GenerateToken handles POST /auth/token requests
// Accepts username/password and returns a JWT token for valid credentials
// Returns 401 for invalid credentials
func (h *AuthHandler) GenerateToken(c *fiber.Ctx) error {
	var creds Credentials

	// Parse request body
	if err := c.BodyParser(&creds); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "invalid request body",
		})
	}

	// Simple validation for MVP - check that credentials are not empty
	if creds.Username == "" || creds.Password == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error: "invalid credentials",
		})
	}

	// For MVP, accept any non-empty username/password combination
	// In production, this would validate against a user database
	// Simple check: password must be at least 4 characters
	if len(creds.Password) < 4 {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error: "invalid credentials",
		})
	}

	// Generate JWT token
	token, err := h.authService.GenerateToken(creds.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "failed to generate token",
		})
	}

	// Return token response
	return c.JSON(models.TokenResponse{
		Token: token,
	})
}
