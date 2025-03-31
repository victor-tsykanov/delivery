package out

import (
	"context"

	"github.com/victor-tsykanov/delivery/internal/core/domain/model/courier"
)

type ICourierRepository interface {
	Create(ctx context.Context, courier *courier.Courier) error
	Update(ctx context.Context, courier *courier.Courier) error
	Get(ctx context.Context, id courier.ID) (*courier.Courier, error)
	FindFree(ctx context.Context) ([]*courier.Courier, error)
}
