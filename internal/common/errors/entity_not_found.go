package errors

import (
	"fmt"

	"github.com/google/uuid"
)

type EntityNotfoundError struct {
	kind string
	id   uuid.UUID
}

func NewEntityNotfoundError(kind string, id uuid.UUID) *EntityNotfoundError {
	return &EntityNotfoundError{
		kind: kind,
		id:   id,
	}
}

func (e *EntityNotfoundError) Error() string {
	return fmt.Sprintf(
		"%s with id=%s is not found",
		e.kind,
		e.id,
	)
}
