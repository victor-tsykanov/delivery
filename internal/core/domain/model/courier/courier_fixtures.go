package courier

import (
	"log"

	"github.com/google/uuid"
	"github.com/victor-tsykanov/delivery/internal/core/domain/kernel"
)

type fixtures struct{}

var Fixtures = &fixtures{}

func (f *fixtures) FreeCourier() *Courier {
	return &Courier{
		id:        uuid.New(),
		location:  kernel.RandomLocation(),
		transport: newTransport(),
		status:    StatusFree,
	}
}

func (f *fixtures) BusyCourier() *Courier {
	return &Courier{
		id:        uuid.New(),
		location:  kernel.RandomLocation(),
		transport: newTransport(),
		status:    StatusBusy,
	}
}

func newTransport() *Transport {
	transport, err := NewTransport(uuid.New(), "Car", 3)
	if err != nil {
		log.Fatalf("failed to create transport: %v", err)
	}

	return transport
}
