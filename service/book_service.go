package service

import (
	"personal/desent-be-task/models"
	"personal/desent-be-task/repository"

	"github.com/google/uuid"
)

// BookService handles business logic for book operations
type BookService struct {
	repo repository.BookRepository
}

// NewBookService creates a new BookService
func NewBookService(repo repository.BookRepository) *BookService {
	return &BookService{
		repo: repo,
	}
}

// CreateBook creates a new book with a generated UUID
func (s *BookService) CreateBook(book *models.Book) error {
	book.ID = uuid.New().String()
	return s.repo.Create(book)
}

// GetBooks retrieves books with optional author filtering and pagination
func (s *BookService) GetBooks(author string, page, limit int) (*models.PaginatedResponse, error) {
	var books []*models.Book
	var err error

	if author != "" {
		books, err = s.repo.FindByAuthor(author)
	} else {
		books, err = s.repo.FindAll()
	}

	if err != nil {
		return nil, err
	}

	return paginate(books, page, limit), nil
}

// GetBookByID retrieves a book by its ID
func (s *BookService) GetBookByID(id string) (*models.Book, error) {
	return s.repo.FindByID(id)
}

// UpdateBook updates an existing book
func (s *BookService) UpdateBook(id string, book *models.Book) error {
	return s.repo.Update(id, book)
}

// DeleteBook deletes a book by its ID
func (s *BookService) DeleteBook(id string) error {
	return s.repo.Delete(id)
}

// GetAllBooks retrieves all books with optional author filtering (no pagination)
func (s *BookService) GetAllBooks(author string) ([]*models.Book, error) {
	if author != "" {
		return s.repo.FindByAuthor(author)
	}
	return s.repo.FindAll()
}

// paginate applies pagination to a slice of books
func paginate(books []*models.Book, page, limit int) *models.PaginatedResponse {
	totalItems := len(books)
	totalPages := (totalItems + limit - 1) / limit

	// Calculate start and end indices
	start := (page - 1) * limit
	end := start + limit

	// Handle out of bounds
	if start >= totalItems {
		return &models.PaginatedResponse{
			Data:       []models.Book{},
			Page:       page,
			Limit:      limit,
			TotalItems: totalItems,
			TotalPages: totalPages,
		}
	}

	if end > totalItems {
		end = totalItems
	}

	// Extract the page of books
	pageBooks := books[start:end]

	// Convert []*Book to []Book for response
	data := make([]models.Book, len(pageBooks))
	for i, book := range pageBooks {
		data[i] = *book
	}

	return &models.PaginatedResponse{
		Data:       data,
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
