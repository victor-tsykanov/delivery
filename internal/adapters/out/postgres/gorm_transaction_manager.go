package postgres

import (
	"context"

	"gorm.io/gorm"
)

type contextKey string

const transactionKey contextKey = "gormTransaction"

type GormTransactionManager struct {
	db *gorm.DB
}

func NewGormTransactionManager(db *gorm.DB) (*GormTransactionManager, error) {
	return &GormTransactionManager{db: db}, nil
}

func (u *GormTransactionManager) Execute(ctx context.Context, fn func(context.Context) error) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, transactionKey, tx)
		return fn(txCtx)
	})
}

func GetTransactionFromContext(ctx context.Context, db *gorm.DB) *gorm.DB {
	tx, ok := ctx.Value(transactionKey).(*gorm.DB)
	if ok {
		return tx
	}

	return db
}
