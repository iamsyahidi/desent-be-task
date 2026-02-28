package handlers

import (
	"errors"
	"strconv"

	"personal/desent-be-task/models"
	"personal/desent-be-task/repository"
	"personal/desent-be-task/service"

	"github.com/gofiber/fiber/v2"
)

// BookHandler handles HTTP requests for book operations
type BookHandler struct {
	service *service.BookService
}

// NewBookHandler creates a new BookHandler
func NewBookHandler(service *service.BookService) *BookHandler {
	return &BookHandler{
		service: service,
	}
}

// Create handles POST /books requests
// Creates a new book and returns HTTP 201 with the created book
// Returns HTTP 400 for validation errors
func (h *BookHandler) Create(c *fiber.Ctx) error {
	var book models.Book

	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
	}

	// Validate required fields
	if book.Title == "" || book.Author == "" || book.Year == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "missing required fields"})
	}

	// Validate year range
	if book.Year < 1000 || book.Year > 9999 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "year must be between 1000 and 9999"})
	}

	if err := h.service.CreateBook(&book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(book)
}

// GetByID handles GET /books/:id requests
// Returns the book with the specified ID or HTTP 404 if not found
func (h *BookHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	book, err := h.service.GetBookByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrBookNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "book not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(book)
}

// List handles GET /books requests
// Supports author query parameter for filtering
// Supports page and limit query parameters for pagination
// Returns paginated list of books
func (h *BookHandler) List(c *fiber.Ctx) error {
	author := c.Query("author")

	// Parse pagination parameters with defaults
	page := parseIntQuery(c, "page", 1)
	limit := parseIntQuery(c, "limit", 10)

	// Validate pagination parameters
	if page < 1 || limit < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid pagination parameters"})
	}

	result, err := h.service.GetBooks(author, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(result)
}

// Update handles PUT /books/:id requests
// Updates an existing book and returns HTTP 200 with the updated book
// Returns HTTP 404 if book not found
// Returns HTTP 400 for validation errors
func (h *BookHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var book models.Book
	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
	}

	// Validate required fields
	if book.Title == "" || book.Author == "" || book.Year == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "missing required fields"})
	}

	// Validate year range
	if book.Year < 1000 || book.Year > 9999 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "year must be between 1000 and 9999"})
	}

	if err := h.service.UpdateBook(id, &book); err != nil {
		if errors.Is(err, repository.ErrBookNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "book not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(book)
}

// Delete handles DELETE /books/:id requests
// Deletes a book and returns HTTP 204 on success
// Returns HTTP 404 if book not found
func (h *BookHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.service.DeleteBook(id); err != nil {
		if errors.Is(err, repository.ErrBookNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "book not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// parseIntQuery parses an integer query parameter with a default value
func parseIntQuery(c *fiber.Ctx, key string, defaultValue int) int {
	valueStr := c.Query(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}
