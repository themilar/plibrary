package main

import "github.com/themilar/plibrary/internal/models"

// bookSearchSuccessResponse is the response with a list of books that match the query parameter
// swagger:response bookSearch
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
// swagger: response bookDetail
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
// swagger:response internalServerError
type errorResponse500 struct {
	// The error message
	// in: body
	Body struct {
		// The error message
		// example: Internal server error
		Message string `json:"message"`
	}
}

// swagger:response badRequest
type errorResponse400andotherclienterrors struct {
	// The error message
	// in: body
	Body struct {
		// The error message
		// example: badRequest
		Message string `json:"message"`
	}
}

// swagger:response noContent
type noContentResponse204 struct {
	// The error message
	// in: body
	Body struct {
	}
}

// swagger:response notFound
type errorResponse404 struct {
	// The error message
	// in: body
	Body struct {
		// The error message
		// example: requested resource not found
		Message string `json:"message"`
	}
}

// BookQueryParams contains all possible query parameters for book filter and pagination
// swagger:parameters listBooks
type BookQueryParams struct {
	// The page number for pagination
	// in: query
	// minimum: 1
	// maximum: 1000
	// default: 1
	// example: 2
	Page int `json:"page"`

	// Number of items per page
	// in: query
	// minimum: 1
	// maximum: 20
	// default: 12
	// example: 20
	Size int `json:"size"`

	// Filter by genre
	// in: query
	// example: sci-fi
	Genre string `json:"genre"`

	// Sort results by field
	// in: query
	// enum: ["title", "published", "pages","-title", "-published", "-pages"]
	// default: title
	// example: published
	Sort string `json:"sort"`
}

// SearchQueryParams contains the search query parameter
// swagger:parameters searchBooks
type SearchQueryParams struct {
	// Search query string to filter books by title or content
	// in: query
	// pattern: ^[a-zA-Z0-9 ]{1,50}$
	// example: fantasy adventure
	Q string `json:"q"`
}
