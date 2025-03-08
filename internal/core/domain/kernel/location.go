package kernel

import (
	"log"
	"math/rand"

	"github.com/victor-tsykanov/delivery/internal/common/errors"
	"github.com/victor-tsykanov/delivery/internal/common/math"
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

func (l *Location) X() int {
	return l.x
}

func (l *Location) Y() int {
	return l.y
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
	xDistance := math.Abs(l.x - other.x)
	yDistance := math.Abs(l.y - other.y)

	return int(xDistance + yDistance)
}
