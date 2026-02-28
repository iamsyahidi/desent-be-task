package main

import (
	"log"

	"personal/desent-be-task/config"
	"personal/desent-be-task/handlers"
	"personal/desent-be-task/middleware"
	"personal/desent-be-task/repository"
	"personal/desent-be-task/service"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration from environment variables
	cfg := config.Load()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
	})

	// Add recovery middleware to handle panics
	app.Use(middleware.RecoveryMiddleware())

	// Add request logging middleware
	app.Use(middleware.LoggerMiddleware())

	// Initialize repository
	bookRepo := repository.NewMemoryBookRepository()

	// Initialize services
	bookService := service.NewBookService(bookRepo)
	authService := service.NewAuthService(cfg.JWTSecret, cfg.JWTExpiration)

	// Initialize handlers
	bookHandler := handlers.NewBookHandler(bookService)
	authHandler := handlers.NewAuthHandler(authService)

	// Register routes
	setupRoutes(app, bookHandler, authHandler, authService)

	// Start HTTP server
	addr := ":" + cfg.Port
	log.Printf("Starting server on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupRoutes registers all application routes
func setupRoutes(app *fiber.App, bookHandler *handlers.BookHandler, authHandler *handlers.AuthHandler, authService *service.AuthService) {
	// Level 1: Ping endpoint (health check)
	app.Get("/ping", handlers.PingHandler)

	// Level 2: Echo endpoint
	app.Post("/echo", handlers.EchoHandler)

	// Level 5: Auth endpoint
	app.Post("/auth/token", authHandler.GenerateToken)

	// Level 3-4: Public book endpoints (CRUD operations)
	app.Post("/books", bookHandler.Create)
	// app.Get("/books", bookHandler.List)
	app.Get("/books/:id", bookHandler.GetByID)
	app.Put("/books/:id", bookHandler.Update)
	app.Delete("/books/:id", bookHandler.Delete)

	// Level 5-6: Protected book list endpoint (requires authentication)
	// This endpoint supports search by author and pagination
	app.Get("/books", middleware.AuthMiddleware(authService), bookHandler.List)
}

// customErrorHandler handles errors returned by handlers
func customErrorHandler(c *fiber.Ctx, err error) error {
	// Default to 500 Internal Server Error
	code := fiber.StatusInternalServerError

	// Check if it's a Fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Return error response
	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
