package main

import "github.com/themilar/plibrary/internal/models"

// bookSearchSuccessResponse is the response with a list of books that match the query parameter
// swagger:response successResponse
type bookSearchSuccessResponse struct {
	// A list of books
	// in: body
	Body struct {
		// The list of books
		// required: true
		Books []models.Book `json:"books"`
	}
}

// bookDetailSuccessResponse contains the response for a single book object
type bookDetailSuccessResponse struct {
	Body struct {
		book models.Book
	}
}

// PaginationMetadata contains pagination information
// swagger:model
type PaginationMetadata struct {
	// The current page number
	// required: true
	// example: 2
	CurrentPage int `json:"current_page"`

	// Number of items per page
	// required: true
	// example: 2
	PageSize int `json:"page_size"`

	// The first page number
	// required: true
	// example: 1
	FirstPage int `json:"first_page"`

	// The last page number
	// required: true
	// example: 7
	LastPage int `json:"last_page"`
}

// PaginatedBooksResponse is the response with a paginated list of books
// swagger:response paginatedBooksResponse
type PaginatedBooksResponse struct {
	// A paginated list of books
	// in: body
	Body struct {
		// The list of books
		// required: true
		Books []models.Book `json:"books"`

		// Pagination metadata
		// required: true
		Metadata PaginationMetadata `json:"metadata"`
	}
}

// Error response
// swagger:response errorResponse
type errorResponse struct {
	// The error message
	// in: body
	Body struct {
		// The error message
		// example: Internal server error
		Message string `json:"message"`
	}
}
