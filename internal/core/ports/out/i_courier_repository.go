package out

import (
	"context"

	"github.com/google/uuid"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/courier"
)

type ICourierRepository interface {
	Create(ctx context.Context, courier *courier.Courier) error
	Update(ctx context.Context, courier *courier.Courier) error
	Get(ctx context.Context, id uuid.UUID) (*courier.Courier, error)
	FindFree(ctx context.Context) ([]*courier.Courier, error)
}
