package order

import (
	"github.com/google/uuid"
	"github.com/victor-tsykanov/delivery/internal/common/errors"
	"github.com/victor-tsykanov/delivery/internal/core/domain/kernel"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/courier"
)

type ID uuid.UUID

func (i ID) IsNil() bool {
	return uuid.UUID(i) == uuid.Nil
}

func NewID() ID {
	return ID(uuid.New())
}

type Order struct {
	id        ID
	location  kernel.Location
	status    Status
	courierID *courier.ID
}

func NewOrder(id ID, location kernel.Location) (*Order, error) {
	if id.IsNil() {
		return nil, errors.NewValueIsRequiredError("id")
	}

	return &Order{
		id:       id,
		location: location,
		status:   StatusCreated,
	}, nil
}

func (o *Order) ID() ID {
	return o.id
}

func (o *Order) Location() kernel.Location {
	return o.location
}

func (o *Order) Status() Status {
	return o.status
}

func (o *Order) CourierID() *courier.ID {
	return o.courierID
}

func (o *Order) Assign(courier *courier.Courier) error {
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

	courierID := courier.ID()
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

func RestoreOrder(id ID, location *kernel.Location, status Status, courierID *courier.ID) *Order {
	return &Order{
		id:        id,
		location:  *location,
		status:    status,
		courierID: courierID,
	}
}
