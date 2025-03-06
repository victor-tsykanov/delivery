package errors

import "fmt"

type ValueIsRequiredError struct {
	paramName string
}

func NewValueIsRequiredError(paramName string) *ValueIsRequiredError {
	return &ValueIsRequiredError{
		paramName: paramName,
	}
}

func (e *ValueIsRequiredError) Error() string {
	return fmt.Sprintf("value of %v is required", e.paramName)
}
