package courier_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/suite"
	"github.com/victor-tsykanov/delivery/internal/adapters/out/postgres/courier"
	"github.com/victor-tsykanov/delivery/internal/common/testutils"
	"github.com/victor-tsykanov/delivery/internal/core/domain/kernel"
	domainCourier "github.com/victor-tsykanov/delivery/internal/core/domain/model/courier"
)

type CourierRepositoryTestSuite struct {
	testutils.DBTestSuite
	repository *courier.Repository
}

func (s *CourierRepositoryTestSuite) SetupTest() {
	s.DBTestSuite.SetupTest()

	repository, err := courier.NewRepository(s.DB())
	s.Require().NoError(err)

	s.repository = repository
}

func TestCourierRepositoryTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(CourierRepositoryTestSuite))
}

func (s *CourierRepositoryTestSuite) TestCreate() {
	// Arrange
	courierEntity := s.createCourierEntity()

	// Act
	err := s.repository.Create(context.Background(), courierEntity)

	// Assert
	s.Require().NoError(err)

	courierRecord := &courier.Courier{}
	err = s.DB().
		Preload("Transport").
		First(courierRecord, "id = ?", courierEntity.ID()).
		Error
	s.Require().NoError(err)

	s.Assert().Equal(courierEntity.ID(), courierRecord.ID)
	s.Assert().Equal(courierEntity.Name(), courierRecord.Name)
	s.Assert().Equal(courierEntity.Transport().ID(), courierRecord.Transport.ID)
	s.Assert().Equal(courierEntity.Transport().Name(), courierRecord.Transport.Name)
	s.Assert().Equal(courierEntity.Transport().Speed(), courierRecord.Transport.Speed)
	s.Assert().Equal(courierEntity.Location().X(), courierRecord.Location.X)
	s.Assert().Equal(courierEntity.Location().Y(), courierRecord.Location.Y)
	s.Assert().Equal(string(courierEntity.Status()), courierRecord.Status)
}

func (s *CourierRepositoryTestSuite) TestUpdate() {
	// Arrange
	ctx := context.Background()

	courierEntity := s.createCourierEntity()
	err := s.repository.Create(ctx, courierEntity)
	s.Require().NoError(err)

	err = courierEntity.Move(kernel.RandomLocation())
	s.Require().NoError(err)
	err = courierEntity.SetBusy()
	s.Require().NoError(err)

	// Act
	err = s.repository.Update(ctx, courierEntity)

	// Assert
	s.Require().NoError(err)

	courierRecord := &courier.Courier{}
	err = s.DB().
		Preload("Transport").
		First(courierRecord, "id = ?", courierEntity.ID()).
		Error
	s.Require().NoError(err)

	s.Assert().Equal(courierEntity.ID(), courierRecord.ID)
	s.Assert().Equal(courierEntity.Name(), courierRecord.Name)
	s.Assert().Equal(courierEntity.Transport().ID(), courierRecord.Transport.ID)
	s.Assert().Equal(courierEntity.Transport().Name(), courierRecord.Transport.Name)
	s.Assert().Equal(courierEntity.Transport().Speed(), courierRecord.Transport.Speed)
	s.Assert().Equal(courierEntity.Location().X(), courierRecord.Location.X)
	s.Assert().Equal(courierEntity.Location().Y(), courierRecord.Location.Y)
	s.Assert().Equal(string(courierEntity.Status()), courierRecord.Status)
}

func (s *CourierRepositoryTestSuite) TestGet() {
	// Arrange
	ctx := context.Background()

	courierEntity := s.createCourierEntity()
	err := s.repository.Create(ctx, courierEntity)
	s.Require().NoError(err)

	// Act
	courierEntityFromDB, err := s.repository.Get(ctx, courierEntity.ID())

	// Assert
	s.Require().NoError(err)

	s.Assert().Equal(courierEntity.ID(), courierEntityFromDB.ID())
	s.Assert().Equal(courierEntity.Name(), courierEntityFromDB.Name())
	s.Assert().Equal(courierEntity.Transport().ID(), courierEntityFromDB.Transport().ID())
	s.Assert().Equal(courierEntity.Transport().Name(), courierEntityFromDB.Transport().Name())
	s.Assert().Equal(courierEntity.Transport().Speed(), courierEntityFromDB.Transport().Speed())
	s.Assert().Equal(courierEntity.Location().X(), courierEntityFromDB.Location().X())
	s.Assert().Equal(courierEntity.Location().Y(), courierEntityFromDB.Location().Y())
	s.Assert().Equal(courierEntity.Status(), courierEntityFromDB.Status())
}

func (s *CourierRepositoryTestSuite) TestFindFree() {
	// Arrange
	ctx := context.Background()

	freeCourier := s.createCourierEntity()
	err := s.repository.Create(ctx, freeCourier)
	s.Require().NoError(err)

	busyCourier := s.createCourierEntity()
	err = busyCourier.SetBusy()
	s.Require().NoError(err)
	err = s.repository.Create(ctx, busyCourier)
	s.Require().NoError(err)

	// Act
	freeCouriers, err := s.repository.FindFree(ctx)

	// Assert
	s.Require().NoError(err)

	s.Assert().Len(freeCouriers, 1)
	s.Assert().Equal(freeCourier.ID(), freeCouriers[0].ID())
}

func (s *CourierRepositoryTestSuite) createCourierEntity() *domainCourier.Courier {
	location := kernel.RandomLocation()
	courierEntity, err := domainCourier.NewCourier(
		gofakeit.Name(),
		gofakeit.CarType(),
		3,
		location,
	)
	s.Require().NoError(err)

	return courierEntity
}
