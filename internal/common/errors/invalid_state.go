package errors

import "fmt"

type InvalidStateError struct {
	message string
}

func NewInvalidStateError(message string) *InvalidStateError {
	return &InvalidStateError{
		message: message,
	}
}

func (e *InvalidStateError) Error() string {
	return fmt.Sprintf(
		"invalid state: %s",
		e.message,
	)
}
