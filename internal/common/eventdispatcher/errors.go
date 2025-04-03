package eventdispatcher

import "fmt"

type UnknownEventError struct {
	eventType string
}

func NewUnknownEventError(eventType string) *UnknownEventError {
	return &UnknownEventError{eventType: eventType}
}

func (e *UnknownEventError) Error() string {
	return fmt.Sprintf(
		"unknown event type: %s",
		e.eventType,
	)
}
