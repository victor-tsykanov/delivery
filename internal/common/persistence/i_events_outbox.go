package persistence

import (
	"context"

	"github.com/victor-tsykanov/delivery/internal/common/ddd"
)

type IEventsOutbox interface {
	StoreAggregateEvents(ctx context.Context, aggregate ddd.IAggregateRoot) error
}
