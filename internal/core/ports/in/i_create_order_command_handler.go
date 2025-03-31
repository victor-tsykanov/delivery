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
	basketID uuid.UUID
	street   string
}

func (c *CreateOrderCommand) BasketID() uuid.UUID {
	return c.basketID
}

func (c *CreateOrderCommand) Street() string {
	return c.street
}

func NewCreateOrderCommand(basketID uuid.UUID, street string) (*CreateOrderCommand, error) {
	if basketID == uuid.Nil {
		return nil, errors.NewValueIsRequiredError("basketID")
	}

	if street == "" {
		return nil, errors.NewValueIsRequiredError("street")
	}

	return &CreateOrderCommand{basketID: basketID, street: street}, nil
}
