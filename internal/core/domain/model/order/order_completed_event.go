package order

import (
	"github.com/google/uuid"
)

const EventTypeCompleted = "order.completed"

type CompletedEvent struct {
	id    uuid.UUID
	order Order
}

func NewCompletedEvent(order Order) *CompletedEvent {
	return &CompletedEvent{id: uuid.New(), order: order}
}

func (e *CompletedEvent) ID() uuid.UUID {
	return e.id
}

func (e *CompletedEvent) Type() string {
	return EventTypeCompleted
}

func (e *CompletedEvent) Order() Order {
	return e.order
}

func RestoreCompletedEvent(id uuid.UUID, order Order) *CompletedEvent {
	return &CompletedEvent{id: id, order: order}
}
