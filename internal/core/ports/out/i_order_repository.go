package out

import (
	"context"

	"github.com/victor-tsykanov/delivery/internal/core/domain/model/order"
)

type IOrderRepository interface {
	Create(ctx context.Context, order *order.Order) error
	Update(ctx context.Context, order *order.Order) error
	Get(ctx context.Context, id order.ID) (*order.Order, error)
	FindNew(ctx context.Context) ([]*order.Order, error)
	FindAssigned(ctx context.Context) ([]*order.Order, error)
}
