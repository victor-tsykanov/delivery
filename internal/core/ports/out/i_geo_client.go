package out

import (
	"context"

	"github.com/victor-tsykanov/delivery/internal/core/domain/kernel"
)

type IGeoClient interface {
	GetLocation(ctx context.Context, street string) (*kernel.Location, error)
}
