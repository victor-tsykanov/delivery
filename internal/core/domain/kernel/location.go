package kernel

import (
	"log"
	"math"
	"math/rand"

	"github.com/victor-tsykanov/delivery/internal/core/domain/errors"
)

const (
	minX = 1
	maxX = 10
	minY = 1
	maxY = 10
)

type Location struct {
	x int
	y int
}

func NewLocation(x int, y int) (*Location, error) {
	if x < minX || x > maxX {
		return nil, errors.NewValueIsOutOfRangeError("x", x, minX, maxX)
	}

	if y < minY || y > maxY {
		return nil, errors.NewValueIsOutOfRangeError("y", y, minY, maxY)
	}

	return &Location{x: x, y: y}, nil
}

//nolint:gosec
func RandomLocation() *Location {
	location, err := NewLocation(
		minX+rand.Intn(maxX),
		minY+rand.Intn(maxY),
	)
	if err != nil {
		log.Fatalf("random location is invalid: %v", err)
	}

	return location
}

func (l *Location) Equals(other Location) bool {
	return *l == other
}

func (l *Location) DistanceTo(other Location) int {
	xDistance := math.Abs(float64(l.x - other.x))
	yDistance := math.Abs(float64(l.y - other.y))

	return int(xDistance + yDistance)
}
