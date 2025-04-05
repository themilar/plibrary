package models

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Book struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title" validate:"required,max=56"`
	CreatedAt time.Time `json:"-"`
	Published int       `json:"published" validate:"required,publication_date"`
	Pages     int       `json:"pages,omitempty,string" validate:"required,gt=0"`
	Genres    []string  `json:"genres,omitempty" validate:"required,unique,gt=0,lt=6"`
	Version   int       `json:"version"`
}
type BookModel struct {
	DB *pgxpool.Pool
}

var ErrRecordNotFound = errors.New("record not found")

type Models struct {
	Books BookModel
}

func NewModels(db *pgxpool.Pool) Models {
	return Models{
		Books: BookModel{DB: db},
	}
}

func (b BookModel) Insert(book *Book) error {
	query := `INSERT INTO books (title,published,pages,genres)
	VALUES ($1,$2,$3,$4) RETURNING id,created_at,version`
	params := []interface{}{book.Title, book.Published, book.Pages, book.Genres}
	return b.DB.QueryRow(context.Background(), query, params...).Scan(&book.ID, &book.CreatedAt, &book.Version)
}
func (b BookModel) Get(id int64) (*Book, error) {
	return nil, nil
}
func (b BookModel) Update(book *Book) error {
	return nil
}
func (b BookModel) Delete(id int64) error {
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
					jv.AddError(strings.ToLower(e.Field()), "must be provided")
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
