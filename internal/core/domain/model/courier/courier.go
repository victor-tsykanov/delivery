package courier

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/victor-tsykanov/delivery/internal/common/errors"
	"github.com/victor-tsykanov/delivery/internal/core/domain/kernel"
)

type Courier struct {
	id        uuid.UUID
	name      string
	transport *Transport
	location  *kernel.Location
	status    Status
}

func NewCourier(name string, transport *Transport, location *kernel.Location) (*Courier, error) {
	if name == "" {
		return nil, errors.NewValueIsRequiredError("name")
	}

	if transport == nil {
		return nil, errors.NewValueIsRequiredError("transport")
	}

	if location == nil {
		return nil, errors.NewValueIsRequiredError("location")
	}

	return &Courier{
		id:        uuid.New(),
		name:      name,
		transport: transport,
		location:  location,
		status:    StatusFree,
	}, nil
}

func (c *Courier) SetBusy() error {
	if c.status != StatusFree {
		return errors.NewInvalidStateError("courier must be free")
	}

	c.status = StatusBusy

	return nil
}

func (c *Courier) SetFree() error {
	if c.status != StatusBusy {
		return errors.NewInvalidStateError("courier must be busy")
	}

	c.status = StatusFree

	return nil
}

func (c *Courier) Move(targetLocation kernel.Location) error {
	newLocation, err := c.transport.Move(c.location, &targetLocation)
	if err != nil {
		return fmt.Errorf("transport failed to move: %w", err)
	}

	c.location = newLocation

	return nil
}

func (c *Courier) CalculateStepsToLocation(location kernel.Location) (int, error) {
	var (
		steps           = 0
		err             error
		currentLocation = c.Location()
	)

	for {
		if currentLocation.Equals(location) {
			break
		}

		currentLocation, err = c.transport.Move(currentLocation, &location)
		if err != nil {
			return 0, fmt.Errorf("transport failed to move: %w", err)
		}

		steps++
	}

	return steps, nil
}

func (c *Courier) ID() uuid.UUID {
	return c.id
}

func (c *Courier) Name() string {
	return c.name
}

func (c *Courier) Transport() *Transport {
	return c.transport
}

func (c *Courier) Location() *kernel.Location {
	return c.location
}

func (c *Courier) Status() Status {
	return c.status
}
