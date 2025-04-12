package internal

import (
	"errors"
	"fmt"
	"maps"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Filters struct {
	Page int    `validate:"max=1000,min=1"`
	Size int    `validate:"max=20,min=1"`
	Sort string `validate:"oneofci=title published pages -title -published -pages"`
}
type FilterValidationErrors struct {
	Errors map[string]string
}

func (fve *FilterValidationErrors) AddError(key, message string) {
	if _, ok := fve.Errors[key]; !ok {
		fve.Errors[key] = message
	}
}
func (f Filters) SortColumn() string {
	sortField, ok := reflect.TypeOf(f).FieldByName("Sort")
	validateTag := sortField.Tag.Get("validate")
	vtSlice := strings.Split(validateTag, "=")
	safeSortValues := vtSlice[1]
	if !ok {
		panic("that field does not exist")
	}

	for _, safeValue := range strings.Split(safeSortValues, " ") {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	panic("unsafe sort param: " + f.Sort)
}
func (f Filters) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func ValidateFilters(f Filters, fte map[string]string) map[string]string {
	v := validator.New(validator.WithRequiredStructEnabled())
	var fve = FilterValidationErrors{
		Errors: make(map[string]string),
	}
	err := v.Struct(f)
	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			for _, e := range validateErrs {
				fmt.Println(e.Tag(), e.Param(), e.Field(), e.Value())
				switch {
				case e.Tag() == "max":
					fve.AddError(strings.ToLower(e.Field()), fmt.Sprintf("value must be less than: %v", e.Param()))
				case e.Tag() == "min":
					fve.AddError(strings.ToLower(e.Field()), fmt.Sprintf("value must be greater than: %v", e.Param()))
				case e.Tag() == "oneofci":
					fve.AddError(strings.ToLower(e.Field()), fmt.Sprintf("can only contain values: %v", e.Param()))
				}
			}
			maps.Copy(fve.Errors, fte)
			return fve.Errors
		}
		// clear(Fve)
		return nil
	}
	return nil
}

// how to concatenate two maps
