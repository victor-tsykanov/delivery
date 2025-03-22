package in

import (
	"context"

	"github.com/google/uuid"
)

type IGetAllCouriersQueryHandler interface {
	Handle(ctx context.Context) ([]*Courier, error)
}

type Courier struct {
	ID       uuid.UUID
	Name     string
	Location *CourierLocation
}

type CourierLocation struct {
	X int
	Y int
}
