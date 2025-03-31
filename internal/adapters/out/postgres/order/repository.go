package order

import (
	"context"
	"fmt"

	"github.com/victor-tsykanov/delivery/internal/adapters/out/postgres"
	"github.com/victor-tsykanov/delivery/internal/common/errors"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/order"
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

func (r *Repository) Create(ctx context.Context, order *order.Order) error {
	orderRecord := toRecordFromDomainEntity(order)

	tx := postgres.GetTransactionFromContext(ctx, r.db)
	return tx.Create(&orderRecord).Error
}

func (r *Repository) Update(ctx context.Context, order *order.Order) error {
	orderRecord := toRecordFromDomainEntity(order)

	tx := postgres.GetTransactionFromContext(ctx, r.db)
	return tx.Save(&orderRecord).Error
}

func (r *Repository) Get(ctx context.Context, id order.ID) (*order.Order, error) {
	orderRecord := &Order{}

	tx := postgres.GetTransactionFromContext(ctx, r.db)
	result := tx.Find(&orderRecord, id)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	orderEntity, err := toDomainEntityFromRecord(orderRecord)
	if err != nil {
		return nil, fmt.Errorf("failed to map order record to domain entity: %w", err)
	}

	return orderEntity, nil
}

func (r *Repository) FindNew(ctx context.Context) ([]*order.Order, error) {
	return r.findByStatus(ctx, order.StatusCreated)
}

func (r *Repository) FindAssigned(ctx context.Context) ([]*order.Order, error) {
	return r.findByStatus(ctx, order.StatusAssigned)
}

func (r *Repository) findByStatus(ctx context.Context, status order.Status) ([]*order.Order, error) {
	orderRecords := make([]*Order, 0)

	tx := postgres.GetTransactionFromContext(ctx, r.db)
	result := tx.Find(&orderRecords, "status = ?", string(status))
	if result.Error != nil {
		return nil, result.Error
	}

	orders := make([]*order.Order, len(orderRecords))
	for i, orderRecord := range orderRecords {
		orderEntity, err := toDomainEntityFromRecord(orderRecord)
		if err != nil {
			return nil, fmt.Errorf("failed to map order record to domain entity: %w", err)
		}

		orders[i] = orderEntity
	}

	return orders, nil
}
