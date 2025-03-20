package courier

import (
	"fmt"

	"github.com/victor-tsykanov/delivery/internal/core/domain/kernel"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/courier"
)

func toRecordFromDomainEntity(entity *courier.Courier) *Courier {
	return &Courier{
		ID:   entity.ID(),
		Name: entity.Name(),
		Transport: Transport{
			ID:    entity.Transport().ID(),
			Name:  entity.Transport().Name(),
			Speed: entity.Transport().Speed(),
		},
		Location: Location{
			X: entity.Location().X(),
			Y: entity.Location().Y(),
		},
		Status: string(entity.Status()),
	}
}

func toDomainEntityFromRecord(record *Courier) (*courier.Courier, error) {
	location, err := kernel.NewLocation(record.Location.X, record.Location.Y)
	if err != nil {
		return nil, fmt.Errorf("failed to create new kernel.Location from %v", record.Location)
	}

	transport := courier.RestoreTransport(
		record.Transport.ID,
		record.Transport.Name,
		record.Transport.Speed,
	)

	return courier.RestoreCourier(
		record.ID,
		record.Name,
		transport,
		location,
		courier.Status(record.Status),
	), nil
}
