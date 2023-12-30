package go_course_validation

import (
	"github.com/papey08/golang-fintech/validation/check"
	"github.com/papey08/golang-fintech/validation/parse"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

var ErrNotStruct = errors.New("wrong argument given, should be a struct")
var ErrInvalidValidatorSyntax = errors.New("invalid validator syntax")
var ErrValidateForUnexportedFields = errors.New("validation for unexported field is not allowed")

type ValidationError struct {
	Err error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	switch len(v) {
	case 0:
		return ""
	case 1:
		return v[0].Err.Error()
	default:
		var res strings.Builder
		for _, ve := range v {
			res.WriteString(ve.Err.Error() + "\n")
		}
		return res.String()
	}
}

func Validate(v any) error {
	vt := reflect.TypeOf(v)

	// check if v is not a struct
	if vt.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	fieldAmount := vt.NumField()
	validationErrors := make(ValidationErrors, 0, fieldAmount)

	for i := 0; i < fieldAmount; i++ {
		field := vt.Field(i)

		// checking if field is a struct for nested validation
		if validationTag, ok := field.Tag.Lookup("validate"); ok {

			// check if filed is not exported
			if !field.IsExported() {
				validationErrors = append(validationErrors, ValidationError{ErrValidateForUnexportedFields})
				continue
			}

			validateOperation, args := parse.ValidationParams(validationTag)

			// check if validate tag is invalid
			if validateOperation == parse.Wrong {
				validationErrors = append(validationErrors, ValidationError{ErrInvalidValidatorSyntax})
				continue
			}

			value := reflect.ValueOf(v).FieldByName(field.Name).Interface()

			// check if value or type of the field don't satisfy validate tag
			err := check.ValidField(value, validateOperation, args)
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{err})
			}

		} else if field.IsExported() { // check for nested struct
			value := reflect.ValueOf(v).FieldByName(field.Name).Interface()

			// check fields of exported nested struct
			if reflect.TypeOf(value).Kind() == reflect.Struct {
				err := Validate(value)
				if err != nil {
					validationErrors = append(validationErrors, ValidationError{err})
				}
			}
		}
	}

	if len(validationErrors) != 0 {
		return validationErrors
	} else {
		return nil
	}
}
