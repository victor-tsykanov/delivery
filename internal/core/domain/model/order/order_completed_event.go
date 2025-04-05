package order

import "github.com/google/uuid"

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
	return "order.confirmed"
}

func (e *CompletedEvent) Order() Order {
	return e.order
}
