package order

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/victor-tsykanov/delivery/internal/common/ddd"
	"github.com/victor-tsykanov/delivery/internal/core/domain/kernel"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/courier"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/order"
)

type CompletedEventSerializer struct {
}

func NewCompletedEventSerializer() *CompletedEventSerializer {
	return &CompletedEventSerializer{}
}

func (s *CompletedEventSerializer) Serialize(event ddd.IDomainEvent) ([]byte, error) {
	completedEvent, ok := event.(*order.CompletedEvent)
	if !ok {
		return nil, fmt.Errorf("invalid event type: %s", event.Type())
	}

	orderEntity := completedEvent.Order()
	location := orderEntity.Location()
	dto := &completedEventDTO{
		ID: event.ID(),
		Order: orderDTO{
			ID: uuid.UUID(orderEntity.ID()),
			Location: locationDTO{
				X: location.X(),
				Y: location.Y(),
			},
			Status:    string(orderEntity.Status()),
			CourierID: (*uuid.UUID)(orderEntity.CourierID()),
		},
	}

	data, err := json.Marshal(dto)
	if err != nil {
		return nil, fmt.Errorf("failed to convert order.CompletedEvent to JSON: %w", err)
	}

	return data, nil
}

func (s *CompletedEventSerializer) Deserialize(data []byte) (ddd.IDomainEvent, error) {
	dto := &completedEventDTO{}
	err := json.Unmarshal(data, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to convert JSON to completed event DTO: %w", err)
	}

	location, err := kernel.NewLocation(dto.Order.Location.X, dto.Order.Location.Y)
	if err != nil {
		return nil, fmt.Errorf("failed create Location: %w", err)
	}

	orderEntity := order.RestoreOrder(
		order.ID(dto.Order.ID),
		location,
		order.Status(dto.Order.Status),
		(*courier.ID)(dto.Order.CourierID),
	)

	return order.RestoreCompletedEvent(dto.ID, *orderEntity), nil
}

type completedEventDTO struct {
	ID    uuid.UUID
	Order orderDTO
}
type orderDTO struct {
	ID        uuid.UUID
	Location  locationDTO
	Status    string
	CourierID *uuid.UUID
}

type locationDTO struct {
	X int
	Y int
}
