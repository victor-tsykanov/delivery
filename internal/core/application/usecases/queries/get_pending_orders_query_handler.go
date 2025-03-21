package queries

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/order"
	inPorts "github.com/victor-tsykanov/delivery/internal/core/ports/in"
	"gorm.io/gorm"
)

type GetPendingOrdersQueryHandler struct {
	db *gorm.DB
}

func NewGetPendingOrdersQueryHandler(db *gorm.DB) (*GetPendingOrdersQueryHandler, error) {
	return &GetPendingOrdersQueryHandler{db: db}, nil
}

func (h *GetPendingOrdersQueryHandler) Handle(ctx context.Context) ([]*inPorts.Order, error) {
	rows, err := h.db.
		WithContext(ctx).
		Raw(
			"select id, location_x, location_y from orders "+
				"where deleted_at is null "+
				"and status in (?, ?) "+
				"order by created_at",
			string(order.StatusCreated),
			string(order.StatusAssigned),
		).
		Rows()
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("failed to close cursor: %v", err)
		}
	}(rows)

	var orders []*inPorts.Order
	for rows.Next() {
		var (
			id        uuid.UUID
			locationX int
			locationY int
		)

		err = rows.Scan(&id, &locationX, &locationY)
		if err != nil {
			return nil, err
		}

		orders = append(orders, &inPorts.Order{
			ID: id,
			Location: &inPorts.OrderLocation{
				X: locationX,
				Y: locationY,
			},
		})
	}

	return orders, nil
}
