package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/themilar/plibrary/internal"
)

// Book represents a book in the system
// swagger:model Book
type Book struct {
	// The unique ID of the book
	// required: true
	// example: 13
	ID int64 `json:"id"`
	// The title of the book
	// required: true
	// example: Black Panther
	Title     string    `json:"title" validate:"required,max=56"`
	CreatedAt time.Time `json:"-"`
	// The year the book was published
	// required: true
	// example: 2018
	Published int `json:"published" validate:"required,publication_date"`
	// The number of pages in the book
	// example: 134
	Pages int `json:"pages,omitempty,string" validate:"required,gt=0"`
	// Genres of the book
	// example: ["sci-fi", "action", "adventure"]
	Genres []string `json:"genres,omitempty" validate:"required,unique,gt=0,lt=6"`
	// Version number of the book record
	// example: 1
	Version int `json:"version"`
}
type BookModel struct {
	DB *pgxpool.Pool
}

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Books BookModel
}

func NewModels(db *pgxpool.Pool) Models {
	return Models{
		Books: BookModel{DB: db},
	}
}
func (b BookModel) All(title string, genres []string, filters internal.Filters) ([]*Book, *internal.PaginationMetadata, error) {
	query := fmt.Sprintf(`SELECT COUNT(*) OVER(), id,created_at,title,published,pages,genres,version 
	FROM books 
	WHERE (LOWER(title)=LOWER($1) OR $1='')
	AND (genres@>$2 OR $2='{}') 
	ORDER BY %s %s, id ASC
	LIMIT $3 OFFSET $4`, filters.SortColumn(), filters.SortDirection())
	params := []any{title, genres, filters.Limit(), filters.Offset()}
	rows, err := b.DB.Query(context.Background(), query, params...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	totalRecords := 0
	books := []*Book{}
	for rows.Next() {
		var book Book
		err := rows.Scan(&totalRecords, &book.ID, &book.CreatedAt, &book.Title, &book.Published, &book.Pages, &book.Genres, &book.Version)
		if err != nil {
			return nil, nil, err
		}
		books = append(books, &book)
	}
	if err = rows.Err(); err != nil {
		return nil, nil, err
	}
	metadata := internal.CalculateMetadata(totalRecords, filters.Page, filters.Size)
	return books, metadata, nil
}
func (b BookModel) FullTextSearch(title string) ([]*Book, error) {
	query := `SELECT * FROM books WHERE (to_tsvector('simple',title) @@ plainto_tsquery('simple',$1) OR $1='') ORDER BY id`
	rows, err := b.DB.Query(context.Background(), query, title)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := []*Book{}
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.CreatedAt, &book.Title, &book.Published, &book.Pages, &book.Genres, &book.Version)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return books, nil
}

func (b BookModel) Insert(book *Book) error {
	query := `INSERT INTO books (title,published,pages,genres)
	VALUES ($1,$2,$3,$4) RETURNING id,created_at,version`
	params := []any{book.Title, book.Published, book.Pages, book.Genres}
	return b.DB.QueryRow(context.Background(), query, params...).Scan(&book.ID, &book.CreatedAt, &book.Version)
}
func (b BookModel) Get(id int64) (*Book, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `SELECT id,created_at,title,published,pages,genres,version FROM books WHERE id=$1`
	var book Book
	err := b.DB.QueryRow(context.Background(), query, id).Scan(&book.ID, &book.CreatedAt, &book.Title, &book.Published, &book.Pages, &book.Genres, &book.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &book, nil
}
func (b BookModel) Update(book *Book) error {
	query := `UPDATE books SET title=$1,published=$2,pages=$3,genres=$4,version=version+1 WHERE id=$5 AND version=$6 RETURNING version`
	params := []any{book.Title, book.Published, book.Pages, book.Genres, book.ID, book.Version}
	err := b.DB.QueryRow(context.Background(), query, params...).Scan(&book.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}
func (b BookModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `DELETE FROM books WHERE ID=$1`
	result, err := b.DB.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

type JsonValidationError struct {
	Errors map[string]string
}

func (jv *JsonValidationError) AddError(key, message string) {
	if _, ok := jv.Errors[key]; !ok {
		jv.Errors[key] = message
	}
}
func validateDate(fl validator.FieldLevel) bool {
	if fl.Field().Int() > int64(time.Now().Year()) {
		return false
	} else if fl.Field().Int() < 1430 {
		return false
	}
	return true
}
func (b *Book) Validate() map[string]string {

	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("publication_date", validateDate)
	jv := JsonValidationError{
		Errors: make(map[string]string),
	}
	err := validate.Struct(b)
	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			for _, e := range validateErrs {
				switch {
				case e.Tag() == "required":
					jv.AddError(strings.ToLower(e.Field()), "must be provided")
				case e.Tag() == "max":
					jv.AddError(strings.ToLower(e.Field()), fmt.Sprintf("above the character limit: %v", e.Param()))
				case e.Tag() == "gt":
					jv.AddError(strings.ToLower(e.Field()), "must be above 0")
				case e.Tag() == "lt":
					jv.AddError(strings.ToLower(e.Field()), "must not exceed 5 items")
				case e.Tag() == "publication_date":
					jv.AddError(strings.ToLower(e.Field()), fmt.Sprintf("publication date cannot exceed the range: 1430-%v", time.Now().Year()))
				case e.Tag() == "unique":
					jv.AddError(strings.ToLower(e.Field()), "cannot contain duplicate genres")
				}
			}
			return jv.Errors
		}
		// from here you can create your own error messages in whatever language you wish
		return nil
	}
	return nil
}
