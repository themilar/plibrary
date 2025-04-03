package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
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
