package order

import (
	"github.com/google/uuid"
	"github.com/victor-tsykanov/delivery/internal/core/domain/kernel"
)

type fixtures struct{}

var Fixtures = &fixtures{}

func (f *fixtures) UnassignedOrder() *Order {
	return &Order{
		id:       uuid.New(),
		location: *kernel.RandomLocation(),
		status:   StatusCreated,
	}
}

func (f *fixtures) AssignedOrder() *Order {
	courierID := uuid.New()

	return &Order{
		id:        uuid.New(),
		location:  *kernel.RandomLocation(),
		status:    StatusAssigned,
		courierID: &courierID,
	}
}

func (f *fixtures) CompletedOrder() *Order {
	courierID := uuid.New()

	return &Order{
		id:        uuid.New(),
		location:  *kernel.RandomLocation(),
		status:    StatusCompleted,
		courierID: &courierID,
	}
}
