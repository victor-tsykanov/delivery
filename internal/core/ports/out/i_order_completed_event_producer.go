package out

import (
	"context"

	"github.com/victor-tsykanov/delivery/internal/core/domain/model/order"
)

type IOrderCompletedEventProducer interface {
	Publish(ctx context.Context, event order.CompletedEvent) error
}
