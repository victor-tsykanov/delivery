package eventserializers

import (
	"fmt"

	orderSerializers "github.com/victor-tsykanov/delivery/internal/adapters/out/postgres/outbox/eventserializers/order"
	"github.com/victor-tsykanov/delivery/internal/common/ddd"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/order"
)

type Registry struct {
	serializers map[string]ISerializer
}

func NewRegistry() (*Registry, error) {
	return &Registry{
		serializers: map[string]ISerializer{
			order.EventTypeCompleted: orderSerializers.NewCompletedEventSerializer(),
		},
	}, nil
}

func (r *Registry) Serialize(event ddd.IDomainEvent) ([]byte, error) {
	serializer, err := r.getEventSerializer(event.Type())
	if err != nil {
		return nil, err
	}

	data, err := serializer.Serialize(event)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize event: %w", err)
	}

	return data, nil
}

func (r *Registry) Deserialize(eventType string, data []byte) (ddd.IDomainEvent, error) {
	serializer, err := r.getEventSerializer(eventType)
	if err != nil {
		return nil, err
	}

	event, err := serializer.Deserialize(data)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize event: %w", err)
	}

	return event, nil
}

func (r *Registry) getEventSerializer(eventType string) (ISerializer, error) {
	serializer, ok := r.serializers[eventType]
	if !ok {
		return nil, fmt.Errorf("unsupported event type: %s", eventType)
	}

	return serializer, nil
}
