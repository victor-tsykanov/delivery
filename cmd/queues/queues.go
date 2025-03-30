package queues

import (
	"context"

	"github.com/victor-tsykanov/delivery/cmd/app"
)

func Consume(ctx context.Context, root *app.CompositionRoot) {
	go root.QueueConsumers.BasketConfirmedConsumer.Consume(ctx)
}
