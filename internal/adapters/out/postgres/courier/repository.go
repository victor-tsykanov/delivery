package courier

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/victor-tsykanov/delivery/internal/adapters/out/postgres"
	"github.com/victor-tsykanov/delivery/internal/common/errors"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/courier"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) (*Repository, error) {
	if db == nil {
		return nil, errors.NewValueIsRequiredError("db")
	}

	return &Repository{db: db}, nil
}

func (r *Repository) Create(ctx context.Context, courier *courier.Courier) error {
	courierRecord := toRecordFromDomainEntity(courier)

	tx := postgres.GetTransactionFromContext(ctx, r.db)
	return tx.Create(courierRecord).Error
}

func (r *Repository) Update(ctx context.Context, courier *courier.Courier) error {
	courierRecord := toRecordFromDomainEntity(courier)

	tx := postgres.GetTransactionFromContext(ctx, r.db)
	return tx.Save(courierRecord).Error
}

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*courier.Courier, error) {
	courierRecord := &Courier{}

	tx := postgres.GetTransactionFromContext(ctx, r.db)
	result := tx.
		Preload("Transport").
		First(courierRecord, "id = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	courierEntity, err := toDomainEntityFromRecord(courierRecord)
	if err != nil {
		return nil, fmt.Errorf("failed to map courier record to domain entity: %w", err)
	}

	return courierEntity, nil
}

func (r *Repository) FindFree(ctx context.Context) ([]*courier.Courier, error) {
	courierRecords := make([]*Courier, 0)

	tx := postgres.GetTransactionFromContext(ctx, r.db)
	result := tx.
		Preload("Transport").
		Find(&courierRecords, "status = ?", string(courier.StatusFree))
	if result.Error != nil {
		return nil, result.Error
	}

	couriers := make([]*courier.Courier, len(courierRecords))
	for i, courierRecord := range courierRecords {
		courierEntity, err := toDomainEntityFromRecord(courierRecord)
		if err != nil {
			return nil, fmt.Errorf("failed to map courier record to domain entity: %w", err)
		}

		couriers[i] = courierEntity
	}

	return couriers, nil
}
