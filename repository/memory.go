package repository

import (
	"errors"
	"strings"
	"sync"

	"personal/desent-be-task/models"
)

var (
	// ErrBookNotFound is returned when a book is not found
	ErrBookNotFound = errors.New("book not found")
)

// BookRepository defines the interface for book data operations
type BookRepository interface {
	Create(book *models.Book) error
	FindAll() ([]*models.Book, error)
	FindByID(id string) (*models.Book, error)
	FindByAuthor(author string) ([]*models.Book, error)
	Update(id string, book *models.Book) error
	Delete(id string) error
}

// MemoryBookRepository implements BookRepository using in-memory storage
type MemoryBookRepository struct {
	mu    sync.RWMutex
	books map[string]*models.Book
}

// NewMemoryBookRepository creates a new in-memory book repository
func NewMemoryBookRepository() *MemoryBookRepository {
	return &MemoryBookRepository{
		books: make(map[string]*models.Book),
	}
}

// Create adds a new book to the repository
func (r *MemoryBookRepository) Create(book *models.Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.books[book.ID] = book
	return nil
}

// FindAll returns all books in the repository
func (r *MemoryBookRepository) FindAll() ([]*models.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	books := make([]*models.Book, 0, len(r.books))
	for _, book := range r.books {
		books = append(books, book)
	}

	return books, nil
}

// FindByID returns a book by its ID
func (r *MemoryBookRepository) FindByID(id string) (*models.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	book, exists := r.books[id]
	if !exists {
		return nil, ErrBookNotFound
	}

	return book, nil
}

// FindByAuthor returns all books by a specific author (case-insensitive)
func (r *MemoryBookRepository) FindByAuthor(author string) ([]*models.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	authorLower := strings.ToLower(author)
	books := make([]*models.Book, 0)

	for _, book := range r.books {
		if strings.ToLower(book.Author) == authorLower {
			books = append(books, book)
		}
	}

	return books, nil
}

// Update modifies an existing book
func (r *MemoryBookRepository) Update(id string, book *models.Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.books[id]; !exists {
		return ErrBookNotFound
	}

	book.ID = id
	r.books[id] = book
	return nil
}

// Delete removes a book from the repository
func (r *MemoryBookRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.books[id]; !exists {
		return ErrBookNotFound
	}

	delete(r.books, id)
	return nil
}
