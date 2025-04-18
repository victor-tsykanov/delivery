package order_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/victor-tsykanov/delivery/internal/adapters/out/postgres/order"
	"github.com/victor-tsykanov/delivery/internal/adapters/out/postgres/outbox"
	"github.com/victor-tsykanov/delivery/internal/adapters/out/postgres/outbox/eventserializers"
	"github.com/victor-tsykanov/delivery/internal/common/testutils"
	"github.com/victor-tsykanov/delivery/internal/core/domain/kernel"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/courier"
	domainOrder "github.com/victor-tsykanov/delivery/internal/core/domain/model/order"
	"github.com/victor-tsykanov/delivery/mocks/github.com/victor-tsykanov/delivery/internal_/common/eventdispatcher"
)

type OrderRepositoryTestSuite struct {
	testutils.DBTestSuite
	repository *order.Repository
}

func (s *OrderRepositoryTestSuite) SetupTest() {
	s.DBTestSuite.SetupTest()

	outboxRepository, err := outbox.NewRepository(s.DB())
	s.Require().NoError(err)

	eventSerializersRegistry, err := eventserializers.NewRegistry()
	s.Require().NoError(err)

	eventDispatcher := eventdispatcher.NewMockIEventDispatcher(s.T())

	eventsOutbox, err := outbox.NewEventsOutbox(outboxRepository, eventSerializersRegistry, eventDispatcher)
	s.Require().NoError(err)

	repository, err := order.NewRepository(s.DB(), eventsOutbox)
	s.Require().NoError(err)

	s.repository = repository
}

func TestOrderRepositoryTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(OrderRepositoryTestSuite))
}

func (s *OrderRepositoryTestSuite) TestCreate() {
	// Arrange
	orderLocation := kernel.RandomLocation()
	orderEntity, err := domainOrder.NewOrder(domainOrder.NewID(), *orderLocation)
	s.Require().NoError(err)

	// Act
	err = s.repository.Create(context.Background(), orderEntity)

	// Assert
	s.Require().NoError(err)

	orderRecord := &order.Order{}
	err = s.DB().First(orderRecord, "id = ?", orderEntity.ID()).Error
	s.Require().NoError(err)

	s.Assert().Equal(orderLocation.X(), orderRecord.Location.X)
	s.Assert().Equal(orderLocation.Y(), orderRecord.Location.Y)
	s.Assert().Equal(string(orderEntity.Status()), orderRecord.Status)
	s.Assert().Equal(orderEntity.CourierID(), (*courier.ID)(orderRecord.CourierID))
}

func (s *OrderRepositoryTestSuite) TestUpdate() {
	// Arrange
	orderLocation := kernel.RandomLocation()
	orderEntity, err := domainOrder.NewOrder(domainOrder.NewID(), *orderLocation)
	s.Require().NoError(err)

	err = s.repository.Create(context.Background(), orderEntity)
	s.Require().NoError(err)

	courierEntity := courier.Fixtures.FreeCourier()
	err = orderEntity.Assign(courierEntity)
	s.Require().NoError(err)

	err = orderEntity.Complete()
	s.Require().NoError(err)

	// Act
	err = s.repository.Update(context.Background(), orderEntity)

	// Assert
	s.Assert().NoError(err)

	orderRecord := &order.Order{}
	err = s.DB().First(orderRecord, "id = ?", orderEntity.ID()).Error
	s.Assert().NoError(err)

	s.Assert().Equal(orderLocation.X(), orderRecord.Location.X)
	s.Assert().Equal(orderLocation.Y(), orderRecord.Location.Y)
	s.Assert().Equal(string(orderEntity.Status()), orderRecord.Status)
	s.Assert().Equal(orderEntity.CourierID(), (*courier.ID)(orderRecord.CourierID))
}

func (s *OrderRepositoryTestSuite) TesGet() {
	// Arrange
	ctx := context.Background()

	orderEntity := s.createNewOrder()
	err := s.repository.Create(ctx, orderEntity)
	s.Require().NoError(err)

	// Act
	orderEntityFromDB, err := s.repository.Get(ctx, orderEntity.ID())

	// Assert
	s.Assert().NoError(err)

	s.Assert().Equal(orderEntity.Location(), orderEntityFromDB.Location())
	s.Assert().Equal(orderEntity.Status(), orderEntityFromDB.Status())
	s.Assert().Equal(orderEntity.CourierID(), orderEntityFromDB.CourierID())
}

func (s *OrderRepositoryTestSuite) TestFindNew() {
	// Arrange
	ctx := context.Background()

	newOrder := s.createNewOrder()
	err := s.repository.Create(ctx, newOrder)
	s.Require().NoError(err)

	assignedOrder := s.createAssignedOrder()
	err = s.repository.Create(ctx, assignedOrder)
	s.Require().NoError(err)

	// Act
	newOrders, err := s.repository.FindNew(ctx)

	// Assert
	s.Assert().NoError(err)
	s.Assert().Len(newOrders, 1)
	s.Assert().Equal(newOrder.ID(), newOrders[0].ID())
}

func (s *OrderRepositoryTestSuite) TestFindAssigned() {
	// Arrange
	ctx := context.Background()

	newOrder := s.createNewOrder()
	err := s.repository.Create(ctx, newOrder)
	s.Require().NoError(err)

	assignedOrder := s.createAssignedOrder()
	err = s.repository.Create(ctx, assignedOrder)
	s.Require().NoError(err)

	// Act
	assignedOrders, err := s.repository.FindAssigned(ctx)

	// Assert
	s.Assert().NoError(err)
	s.Assert().Len(assignedOrders, 1)
	s.Assert().Equal(assignedOrder.ID(), assignedOrders[0].ID())
}

func (s *OrderRepositoryTestSuite) createNewOrder() *domainOrder.Order {
	orderLocation := kernel.RandomLocation()
	orderEntity, err := domainOrder.NewOrder(domainOrder.NewID(), *orderLocation)
	s.Require().NoError(err)

	return orderEntity
}

func (s *OrderRepositoryTestSuite) createAssignedOrder() *domainOrder.Order {
	orderEntity := s.createNewOrder()
	courierEntity := courier.Fixtures.FreeCourier()

	err := orderEntity.Assign(courierEntity)
	s.Require().NoError(err)

	return orderEntity
}
