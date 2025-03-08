package courier

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/victor-tsykanov/delivery/internal/common/errors"
	"github.com/victor-tsykanov/delivery/internal/common/math"
	"github.com/victor-tsykanov/delivery/internal/core/domain/kernel"
)

type Transport struct {
	id    uuid.UUID
	name  string
	speed int
}

func NewTransport(id uuid.UUID, name string, speed int) (*Transport, error) {
	if name == "" {
		return nil, errors.NewValueIsRequiredError("name")
	}

	const minSpeed = 1
	const maxSpeed = 3
	if speed < minSpeed || speed > maxSpeed {
		return nil, errors.NewValueIsOutOfRangeError("speed", speed, minSpeed, maxSpeed)
	}

	return &Transport{
		id:    id,
		name:  name,
		speed: speed,
	}, nil
}

func (t *Transport) ID() uuid.UUID {
	return t.id
}

func (t *Transport) Name() string {
	return t.name
}

func (t *Transport) Speed() int {
	return t.speed
}

func (t *Transport) Equals(other Transport) bool {
	return t.id == other.id
}

func (t *Transport) Move(from *kernel.Location, to *kernel.Location) (*kernel.Location, error) {
	xDistance := math.Abs(to.X() - from.X())
	yDistance := math.Abs(to.Y() - from.Y())

	xSteps := math.Sign(to.X()-from.X()) * min(xDistance, t.speed)
	remainingSteps := max(0, t.speed-math.Abs(xSteps))
	ySteps := math.Sign(to.Y()-from.Y()) * min(yDistance, remainingSteps)

	nextLocation, err := kernel.NewLocation(from.X()+xSteps, from.Y()+ySteps)
	if err != nil {
		return nil, fmt.Errorf("failed to create next location: %w", err)
	}

	return nextLocation, nil
}
