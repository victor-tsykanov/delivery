package events

import (
	"context"
	"fmt"

	"github.com/victor-tsykanov/delivery/internal/core/domain/model/order"
	outPorts "github.com/victor-tsykanov/delivery/internal/core/ports/out"
)

type OrderCompletedEventHandler struct {
	producer outPorts.IOrderCompletedEventProducer
}

func NewOrderCompletedEventHandler(
	producer outPorts.IOrderCompletedEventProducer,
) (*OrderCompletedEventHandler, error) {
	return &OrderCompletedEventHandler{producer: producer}, nil
}

func (h *OrderCompletedEventHandler) Handle(ctx context.Context, event *order.CompletedEvent) error {
	err := h.producer.Publish(ctx, *event)
	if err != nil {
		return fmt.Errorf("failed to publish %s event: %w", event.Type(), err)
	}

	return nil
}
