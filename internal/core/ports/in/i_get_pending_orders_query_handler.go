package in

import (
	"context"

	"github.com/google/uuid"
)

type IGetPendingOrdersQueryHandler interface {
	Handle(ctx context.Context) ([]*Order, error)
}

type Order struct {
	ID       uuid.UUID
	Location *OrderLocation
}

type OrderLocation struct {
	X int
	Y int
}
