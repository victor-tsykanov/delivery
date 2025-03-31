package courier

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/victor-tsykanov/delivery/internal/common/errors"
	"github.com/victor-tsykanov/delivery/internal/core/domain/kernel"
)

type ID uuid.UUID

func NewID() ID {
	return ID(uuid.New())
}

type Courier struct {
	id        ID
	name      string
	transport *Transport
	location  *kernel.Location
	status    Status
}

func NewCourier(
	name string,
	transportName string,
	transportSpeed int,
	location *kernel.Location,
) (*Courier, error) {
	if name == "" {
		return nil, errors.NewValueIsRequiredError("name")
	}

	err := validateTransportName(transportName, "transportName")
	if err != nil {
		return nil, err
	}

	err = validateTransportSpeed(transportSpeed, "transportSpeed")
	if err != nil {
		return nil, err
	}

	if location == nil {
		return nil, errors.NewValueIsRequiredError("location")
	}

	transport, err := NewTransport(NewTransportID(), transportName, transportSpeed)
	if err != nil {
		return nil, err
	}

	return &Courier{
		id:        NewID(),
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

func (c *Courier) Move(targetLocation *kernel.Location) error {
	newLocation, err := c.transport.Move(c.location, targetLocation)
	if err != nil {
		return fmt.Errorf("transport failed to move: %w", err)
	}

	c.location = newLocation

	return nil
}

func (c *Courier) CalculateStepsToLocation(location *kernel.Location) (int, error) {
	var (
		steps           = 0
		err             error
		currentLocation = c.Location()
	)

	for {
		if currentLocation.Equals(*location) {
			break
		}

		currentLocation, err = c.transport.Move(currentLocation, location)
		if err != nil {
			return 0, fmt.Errorf("transport failed to move: %w", err)
		}

		steps++
	}

	return steps, nil
}

func (c *Courier) ID() ID {
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

func RestoreCourier(
	id ID,
	name string,
	transport *Transport,
	location *kernel.Location,
	status Status,
) *Courier {
	return &Courier{
		id:        id,
		name:      name,
		transport: transport,
		location:  location,
		status:    status,
	}
}
