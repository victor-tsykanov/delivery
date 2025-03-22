package in

import (
	"context"

	"github.com/google/uuid"
	"github.com/victor-tsykanov/delivery/internal/common/errors"
)

type ICreateOrderCommandHandler interface {
	Handle(context.Context, CreateOrderCommand) error
}

type CreateOrderCommand struct {
	id     uuid.UUID
	street string
}

func (c *CreateOrderCommand) ID() uuid.UUID {
	return c.id
}

func (c *CreateOrderCommand) Street() string {
	return c.street
}

func NewCreateOrderCommand(id uuid.UUID, street string) (*CreateOrderCommand, error) {
	if id == uuid.Nil {
		return nil, errors.NewValueIsRequiredError("id")
	}

	if street == "" {
		return nil, errors.NewValueIsRequiredError("street")
	}

	return &CreateOrderCommand{id: id, street: street}, nil
}
