package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Filters struct {
	// change from number to validate integer
	Page         string `validate:"integer,number,gt=0,max=10000000"`
	Size         string `validate:"integer,number,gt=0,lt=101"`
	Sort         string `validate:"oneofci=id title published pages"`
	SafeSortList []string
}
type FilterValidationErrors struct {
	Errors map[string]string
}

func (fve *FilterValidationErrors) AddError(key, message string) {
	if _, ok := fve.Errors[key]; !ok {
		fve.Errors[key] = message
	}
}
func validateInteger(fl validator.FieldLevel) bool {
	fmt.Println(fl.Field(), fl.Field().Kind())
	return reflect.TypeOf(1).Kind().String() == "int"
}
func validateInteger10(fl any) bool {
	return reflect.TypeOf(fl) != reflect.TypeOf(1)
}

func ValidateFilters(f Filters) map[string]string {
	v := validator.New(validator.WithRequiredStructEnabled())
	var fve = FilterValidationErrors{
		Errors: make(map[string]string),
	}
	v.RegisterValidation("integer", validateInteger)
	err := v.Struct(f)
	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			for _, e := range validateErrs {
				fmt.Println(e.Tag(), validateInteger10(10), validateInteger10("hgfdg"))
				switch {
				case e.Tag() == "integer":
					fve.AddError(strings.ToLower(e.Field()), "must be an integer")
				case e.Tag() == "max":
					fve.AddError(strings.ToLower(e.Field()), fmt.Sprintf("above the character limit: %v", e.Param()))
				case e.Tag() == "gt":
					fve.AddError(strings.ToLower(e.Field()), fmt.Sprintf("must be greater than: %v", e.Param()))
				case e.Tag() == "lt":
					fve.AddError(strings.ToLower(e.Field()), "must not exceed 5 items")
				case e.Tag() == "oneofci":
					fve.AddError(strings.ToLower(e.Field()), fmt.Sprintf("can only contain values: %v", e.Param()))
				}
			}
			return fve.Errors
		}
		return nil
	}
	return nil
}
