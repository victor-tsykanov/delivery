package queries_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/victor-tsykanov/delivery/internal/common/testutils"
	"github.com/victor-tsykanov/delivery/internal/core/application/usecases/queries"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/courier"
)

type GetAllCouriersQueryHandlerTestSuite struct {
	testutils.DBTestSuite
}

func TestGetAllCouriersQueryHandlerTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(GetAllCouriersQueryHandlerTestSuite))
}

func (s *GetAllCouriersQueryHandlerTestSuite) TestHandle() {
	// Arrange
	ctx := context.Background()

	freeCourier := s.createCourier(courier.StatusFree)
	busyCourier := s.createCourier(courier.StatusBusy)

	handler, err := queries.NewGetAllCouriersQueryHandler(s.DB())
	s.Require().NoError(err)

	// Act
	couriers, err := handler.Handle(ctx)

	// Assert
	s.Require().NoError(err)

	s.Assert().Len(couriers, 2)

	s.Assert().Equal(freeCourier.id, couriers[0].ID)
	s.Assert().Equal(freeCourier.name, couriers[0].Name)
	s.Assert().Equal(freeCourier.locationX, couriers[0].Location.X)
	s.Assert().Equal(freeCourier.locationY, couriers[0].Location.Y)

	s.Assert().Equal(busyCourier.id, couriers[1].ID)
	s.Assert().Equal(busyCourier.name, couriers[1].Name)
	s.Assert().Equal(busyCourier.locationX, couriers[1].Location.X)
	s.Assert().Equal(busyCourier.locationY, couriers[1].Location.Y)
}

type courierRecord struct {
	id        uuid.UUID
	name      string
	locationX int
	locationY int
}

func (s *GetAllCouriersQueryHandlerTestSuite) createCourier(status courier.Status) *courierRecord {
	record := &courierRecord{
		id:        uuid.New(),
		name:      gofakeit.Name(),
		locationX: gofakeit.Number(0, 10),
		locationY: gofakeit.Number(0, 10),
	}

	result := s.DB().Exec(
		"insert into couriers (id, name, location_x, location_y, status, created_at, updated_at) "+
			"values (?, ?, ?, ?, ?, now(), now())",
		record.id,
		record.name,
		record.locationX,
		record.locationY,
		string(status),
	)
	s.Require().NoError(result.Error)

	return record
}
