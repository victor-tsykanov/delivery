package commands_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/victor-tsykanov/delivery/internal/core/application/usecases/commands"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/courier"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/order"
	outPorts "github.com/victor-tsykanov/delivery/mocks/github.com/victor-tsykanov/delivery/internal_/core/ports/out"
)

func TestMoveCouriersCommandHandler_Handle(t *testing.T) {
	// Arrange
	ctx := context.Background()

	transactionManager := outPorts.NewMockITransactionManager(t)
	orderRepository := outPorts.NewMockIOrderRepository(t)
	courierRepository := outPorts.NewMockICourierRepository(t)

	courier1 := courier.Fixtures.BusyCourierAtLocationWithSpeed(1, 1, 1)
	courier2 := courier.Fixtures.BusyCourierAtLocationWithSpeed(3, 3, 1)

	order1 := order.Fixtures.OrderWithTargetLocationAssignedToCourier(1, 2, courier1.ID())
	order2 := order.Fixtures.OrderWithTargetLocationAssignedToCourier(3, 7, courier2.ID())

	transactionManager.
		EXPECT().
		Execute(ctx, mock.Anything).
		Run(func(ctx context.Context, fn func(context.Context) error) {
			err := fn(ctx)
			require.NoError(t, err)
		}).
		Return(nil).
		Twice()

	orderRepository.
		EXPECT().
		FindAssigned(ctx).
		Return([]*order.Order{order1, order2}, nil).
		Once()

	// Mock first courier move
	courierRepository.
		EXPECT().
		Get(ctx, *order1.CourierID()).
		Return(courier1, nil).
		Once()

	orderRepository.
		EXPECT().
		Update(ctx, mock.MatchedBy(func(updatedOrder *order.Order) bool {
			return updatedOrder.ID() == order1.ID() && updatedOrder.Status() == order.StatusCompleted
		})).
		Return(nil).
		Once()

	courierRepository.
		EXPECT().
		Update(ctx, mock.MatchedBy(func(updatedCourier *courier.Courier) bool {
			return updatedCourier.ID() == courier1.ID() &&
				updatedCourier.Location().X() == 1 &&
				updatedCourier.Location().Y() == 2 &&
				updatedCourier.Status() == courier.StatusFree
		})).
		Return(nil).
		Once()

	// Mock second courier move
	courierRepository.
		EXPECT().
		Get(ctx, *order2.CourierID()).
		Return(courier2, nil).
		Once()

	courierRepository.
		EXPECT().
		Update(ctx, mock.MatchedBy(func(updatedCourier *courier.Courier) bool {
			return updatedCourier.ID() == courier2.ID() &&
				updatedCourier.Location().X() == 3 &&
				updatedCourier.Location().Y() == 4 &&
				updatedCourier.Status() == courier.StatusBusy
		})).
		Return(nil).
		Once()

	handler, err := commands.NewMoveCouriersCommandHandler(
		transactionManager,
		courierRepository,
		orderRepository,
	)
	require.NoError(t, err)

	// Act
	err = handler.Handle(ctx)

	// Assert
	require.NoError(t, err)

	assert.Equal(t, order.StatusCompleted, order1.Status())
	assert.Equal(t, courier.StatusFree, courier1.Status())
	assert.Equal(t, 1, courier1.Location().X())
	assert.Equal(t, 2, courier1.Location().Y())

	assert.Equal(t, order.StatusAssigned, order2.Status())
	assert.Equal(t, courier.StatusBusy, courier2.Status())
	assert.Equal(t, 3, courier2.Location().X())
	assert.Equal(t, 4, courier2.Location().Y())
}
