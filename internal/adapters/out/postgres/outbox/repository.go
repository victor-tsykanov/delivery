package outbox

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/victor-tsykanov/delivery/internal/adapters/out/postgres"
	"github.com/victor-tsykanov/delivery/internal/common/errors"
	"gorm.io/gorm"
)

type IRepository interface {
	Create(ctx context.Context, message *Message) error
	MarkProcessed(ctx context.Context, id uuid.UUID) error
	FindUnprocessed(limit int) ([]*Message, bool, error)
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) (*Repository, error) {
	if db == nil {
		return nil, errors.NewValueIsRequiredError("db")
	}

	return &Repository{db: db}, nil
}

func (r *Repository) Create(ctx context.Context, message *Message) error {
	tx := postgres.GetTransactionFromContext(ctx, r.db)

	return tx.Create(message).Error
}

func (r *Repository) MarkProcessed(ctx context.Context, id uuid.UUID) error {
	message, err := r.Get(ctx, id)
	if err != nil {
		return err
	}

	message.ProcessedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	err = r.update(ctx, message)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) FindUnprocessed(limit int) ([]*Message, bool, error) {
	var messages []*Message
	result := r.db.
		Order("created_at asc").
		Limit(limit + 1).
		Where("processed_at is null").Find(&messages)

	if result.Error != nil {
		return nil, false, result.Error
	}

	hasMore := len(messages) > limit
	itemsCount := min(len(messages), limit)

	return messages[:itemsCount], hasMore, nil
}

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*Message, error) {
	message := &Message{}

	tx := postgres.GetTransactionFromContext(ctx, r.db)
	result := tx.Find(&message, id)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.NewEntityNotfoundError("message", id)
	}

	return message, nil
}

func (r *Repository) update(ctx context.Context, message *Message) error {

	tx := postgres.GetTransactionFromContext(ctx, r.db)
	err := tx.Save(message).Error
	if err != nil {
		return err
	}

	return nil
}
