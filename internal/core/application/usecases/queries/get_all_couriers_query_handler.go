package queries

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"
	inPorts "github.com/victor-tsykanov/delivery/internal/core/ports/in"
	"gorm.io/gorm"
)

type GetAllCouriersQueryHandler struct {
	db *gorm.DB
}

func NewGetAllCouriersQueryHandler(db *gorm.DB) (*GetAllCouriersQueryHandler, error) {
	return &GetAllCouriersQueryHandler{db: db}, nil
}

func (h *GetAllCouriersQueryHandler) Handle(ctx context.Context) ([]*inPorts.Courier, error) {
	rows, err := h.db.
		WithContext(ctx).
		Raw(
			"select id, name, location_x, location_y from couriers " +
				"where deleted_at is null " +
				"order by created_at",
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

	var couriers []*inPorts.Courier
	for rows.Next() {
		var (
			id        uuid.UUID
			name      string
			locationX int
			locationY int
		)

		err = rows.Scan(&id, &name, &locationX, &locationY)
		if err != nil {
			return nil, err
		}

		couriers = append(couriers, &inPorts.Courier{
			ID:   id,
			Name: name,
			Location: &inPorts.CourierLocation{
				X: locationX,
				Y: locationY,
			},
		})
	}

	return couriers, nil
}
