package order

import (
	"github.com/google/uuid"
	"github.com/victor-tsykanov/delivery/internal/common/errors"
	"github.com/victor-tsykanov/delivery/internal/core/domain/kernel"
)

type Order struct {
	id        uuid.UUID
	location  kernel.Location
	status    Status
	courierID *uuid.UUID
}

func NewOrder(id uuid.UUID, location kernel.Location) (*Order, error) {
	if id == uuid.Nil {
		return nil, errors.NewValueIsRequiredError("id")
	}

	return &Order{
		id:       id,
		location: location,
		status:   StatusCreated,
	}, nil
}

func (o *Order) ID() uuid.UUID {
	return o.id
}

func (o *Order) Location() kernel.Location {
	return o.location
}

func (o *Order) Status() Status {
	return o.status
}

func (o *Order) CourierID() *uuid.UUID {
	return o.courierID
}

func (o *Order) Assign(courierID uuid.UUID) error {
	if o.status != StatusCreated {
		var message string

		switch o.status {
		case StatusAssigned:
			message = "order is already assigned"
		case StatusCompleted:
			message = "order is completed"
		default:
			message = "unsupported order status"
		}

		return errors.NewInvalidStateError(message)
	}

	o.courierID = &courierID
	o.status = StatusAssigned

	return nil
}

func (o *Order) Complete() error {
	if o.status != StatusAssigned {
		var message string

		switch o.status {
		case StatusCreated:
			message = "order needs to be assigned"
		case StatusCompleted:
			message = "order is already completed"
		default:
			message = "unsupported order status"
		}

		return errors.NewInvalidStateError(message)
	}

	o.status = StatusCompleted

	return nil
}
