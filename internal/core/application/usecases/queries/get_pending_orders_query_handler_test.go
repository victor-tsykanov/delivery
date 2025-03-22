package queries_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/victor-tsykanov/delivery/internal/common/testutils"
	"github.com/victor-tsykanov/delivery/internal/core/application/usecases/queries"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/order"
)

type GetPendingOrdersQueryHandlerTestSuite struct {
	testutils.DBTestSuite
}

func TestGetPendingOrdersQueryHandlerTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(GetPendingOrdersQueryHandlerTestSuite))
}

func (s *GetPendingOrdersQueryHandlerTestSuite) TestHandle() {
	// Arrange
	ctx := context.Background()

	assignedOrder := s.createOrder(order.StatusAssigned)
	_ = s.createOrder(order.StatusCompleted)
	createdOrder := s.createOrder(order.StatusCreated)

	handler, err := queries.NewGetPendingOrdersQueryHandler(s.DB())
	s.Require().NoError(err)

	// Act
	orders, err := handler.Handle(ctx)

	// Assert
	s.Require().NoError(err)

	s.Assert().Len(orders, 2)

	s.Assert().Equal(assignedOrder.id, orders[0].ID)
	s.Assert().Equal(assignedOrder.locationX, orders[0].Location.X)
	s.Assert().Equal(assignedOrder.locationY, orders[0].Location.Y)

	s.Assert().Equal(createdOrder.id, orders[1].ID)
	s.Assert().Equal(createdOrder.locationX, orders[1].Location.X)
	s.Assert().Equal(createdOrder.locationY, orders[1].Location.Y)
}

type orderRecord struct {
	id        uuid.UUID
	locationX int
	locationY int
}

func (s *GetPendingOrdersQueryHandlerTestSuite) createOrder(status order.Status) *orderRecord {
	record := &orderRecord{
		id:        uuid.New(),
		locationX: gofakeit.Number(0, 10),
		locationY: gofakeit.Number(0, 10),
	}

	result := s.DB().Exec(
		"insert into orders (id, location_x, location_y, status, created_at, updated_at) "+
			"values (?, ?, ?, ?, now(), now())",
		record.id,
		record.locationX,
		record.locationY,
		string(status),
	)
	s.Require().NoError(result.Error)

	return record
}
