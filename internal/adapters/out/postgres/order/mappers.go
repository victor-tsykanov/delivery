package order

import (
	"fmt"

	"github.com/victor-tsykanov/delivery/internal/core/domain/kernel"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/order"
)

func toRecordFromDomainEntity(entity *order.Order) *Order {
	location := entity.Location()

	return &Order{
		ID: entity.ID(),
		Location: Location{
			location.X(),
			location.Y(),
		},
		Status:    string(entity.Status()),
		CourierID: entity.CourierID(),
	}
}

func toDomainEntityFromRecord(record *Order) (*order.Order, error) {
	location, err := kernel.NewLocation(record.Location.X, record.Location.Y)
	if err != nil {
		return nil, fmt.Errorf("failed to create new kernel.Location from %v", record.Location)
	}

	return order.RestoreOrder(
		record.ID,
		location,
		order.Status(record.Status),
		record.CourierID,
	), nil
}
