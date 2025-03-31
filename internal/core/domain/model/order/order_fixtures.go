package order

import (
	"log"

	"github.com/victor-tsykanov/delivery/internal/core/domain/kernel"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/courier"
)

type fixtures struct{}

var Fixtures = &fixtures{}

func (f *fixtures) UnassignedOrder() *Order {
	return &Order{
		id:       NewID(),
		location: *kernel.RandomLocation(),
		status:   StatusCreated,
	}
}

func (f *fixtures) AssignedOrder() *Order {
	courierID := courier.NewID()

	return &Order{
		id:        NewID(),
		location:  *kernel.RandomLocation(),
		status:    StatusAssigned,
		courierID: &courierID,
	}
}

func (f *fixtures) OrderWithTargetLocationAssignedToCourier(x int, y int, courierID courier.ID) *Order {
	location, err := kernel.NewLocation(x, y)
	if err != nil {
		log.Fatalf("failed to create location: %v", err)
	}

	return &Order{
		id:        NewID(),
		location:  *location,
		status:    StatusAssigned,
		courierID: &courierID,
	}
}

func (f *fixtures) CompletedOrder() *Order {
	courierID := courier.NewID()

	return &Order{
		id:        NewID(),
		location:  *kernel.RandomLocation(),
		status:    StatusCompleted,
		courierID: &courierID,
	}
}
