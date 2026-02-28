package models

// Book represents a book entity in the system
type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
	Year   int    `json:"year" binding:"required,min=1000,max=9999"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// PaginatedResponse represents a paginated collection of books
type PaginatedResponse struct {
	Data       []Book `json:"data"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	TotalItems int    `json:"total_items"`
	TotalPages int    `json:"total_pages"`
}

// TokenResponse represents an authentication token response
type TokenResponse struct {
	Token string `json:"token"`
}
