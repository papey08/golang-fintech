package check

import (
	"github.com/pkg/errors"
	"homework/parse"
	"strconv"
)

type validatable interface {
	int | string
}

var ErrInvalidFieldValue = errors.New("value of field is not validate")
var ErrInvalidFieldType = errors.New("invalid type of field")

// ValidField checks field if it is complies with validation parameters
func ValidField(field any, vOp parse.ValidationOperation, args any) error {

	if values, ok := field.([]int); ok {
		return validSlice(values, vOp, args)
	}

	if values, ok := field.([]string); ok {
		return validSlice(values, vOp, args)
	}

	if value, ok := field.(int); ok {
		return validValue(value, vOp, args)
	}

	if value, ok := field.(string); ok {
		return validValue(value, vOp, args)
	}

	return ErrInvalidFieldType
}

// validSlice checks if every element of slice complies with validation parameters
func validSlice[T validatable](values []T, vOp parse.ValidationOperation, args any) error {
	for _, v := range values {
		if err := validValue(v, vOp, args); err != nil {
			return err
		}
	}
	return nil
}

// validValue checks if value complies with validation parameters
func validValue(value any, validateOperation parse.ValidationOperation, args any) error {

	if s, isString := value.(string); isString { // checking if string is valid
		switch validateOperation {

		case parse.Length:
			if !validLen(s, args.(int)) {
				return ErrInvalidFieldValue
			}

		case parse.In:
			if !validIn(s, args.([]string)) {
				return ErrInvalidFieldValue
			}

		case parse.Min:
			if !validMin(len(s), args.(int)) {
				return ErrInvalidFieldValue
			}

		case parse.Max:
			if !validMax(len(s), args.(int)) {
				return ErrInvalidFieldValue
			}

		}

	} else if n, isInt := value.(int); isInt { // checking if int is valid
		switch validateOperation {

		case parse.In:
			strArgs := args.([]string)
			intArgs := make([]int, len(strArgs))
			for i := range strArgs {
				temp, err := strconv.Atoi(strArgs[i])
				if err != nil {
					return ErrInvalidFieldType
				}
				intArgs[i] = temp
			}

			if !validIn(n, intArgs) {
				return ErrInvalidFieldValue
			}

		case parse.Min:
			if !validMin(n, args.(int)) {
				return ErrInvalidFieldValue
			}

		case parse.Max:
			if !validMax(n, args.(int)) {
				return ErrInvalidFieldValue
			}
		}
	}
	return nil
}
