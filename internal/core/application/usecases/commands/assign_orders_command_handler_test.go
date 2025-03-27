package commands_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/victor-tsykanov/delivery/internal/common/errors"
	"github.com/victor-tsykanov/delivery/internal/core/application/usecases/commands"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/courier"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/order"
	"github.com/victor-tsykanov/delivery/mocks/github.com/victor-tsykanov/delivery/internal_/common/persistence"
	"github.com/victor-tsykanov/delivery/mocks/github.com/victor-tsykanov/delivery/internal_/core/domain/services"
	outPorts "github.com/victor-tsykanov/delivery/mocks/github.com/victor-tsykanov/delivery/internal_/core/ports/out"
)

func TestAssignOrdersCommandHandler_Handle(t *testing.T) {
	// Arrange
	ctx := context.Background()

	order1 := order.Fixtures.UnassignedOrder()
	order2 := order.Fixtures.UnassignedOrder()

	courier1 := courier.Fixtures.FreeCourier()
	courier2 := courier.Fixtures.FreeCourier()
	courier3 := courier.Fixtures.FreeCourier()

	transactionManager := persistence.NewMockITransactionManager(t)
	orderRepository := outPorts.NewMockIOrderRepository(t)
	courierRepository := outPorts.NewMockICourierRepository(t)
	dispatchService := services.NewMockIDispatchService(t)

	orderRepository.
		EXPECT().
		FindNew(ctx).
		Return([]*order.Order{order1, order2}, nil).
		Once()

	findFreeCalls := 0
	courierRepository.
		EXPECT().
		FindFree(ctx).
		RunAndReturn(func(_ context.Context) ([]*courier.Courier, error) {
			if findFreeCalls == 0 {
				findFreeCalls++
				return []*courier.Courier{courier1, courier2, courier3}, nil
			}

			return []*courier.Courier{courier2, courier3}, nil
		}).
		Twice()

	transactionManager.
		EXPECT().
		Execute(ctx, mock.Anything).
		Run(func(ctx context.Context, fn func(context.Context) error) {
			err := fn(ctx)
			require.NoError(t, err)
		}).
		Return(nil).
		Twice()

	// First order
	dispatchService.
		EXPECT().
		Dispatch(order1, []*courier.Courier{courier1, courier2, courier3}).
		Return(courier1, nil).
		Once()

	courierRepository.
		EXPECT().
		Update(ctx, courier1).
		Return(nil).
		Once()

	orderRepository.
		EXPECT().
		Update(ctx, order1).
		Return(nil).
		Once()

	// Second order
	dispatchService.
		EXPECT().
		Dispatch(order2, []*courier.Courier{courier2, courier3}).
		Return(courier2, nil)

	courierRepository.
		EXPECT().
		Update(ctx, courier2).
		Return(nil).
		Once()

	orderRepository.
		EXPECT().
		Update(ctx, order2).
		Return(nil).
		Once()

	handler, err := commands.NewAssignOrdersCommandHandler(
		transactionManager,
		dispatchService,
		courierRepository,
		orderRepository,
	)
	require.NoError(t, err)

	// Act
	err = handler.Handle(ctx)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, courier.StatusBusy, courier1.Status())
	assert.Equal(t, courier.StatusBusy, courier2.Status())
	assert.Equal(t, courier.StatusFree, courier3.Status())
	assert.Equal(t, order.StatusAssigned, order1.Status())
	assert.Equal(t, order.StatusAssigned, order2.Status())
}

func TestAssignOrdersCommandHandler_Handle_FailsWhenNotEnoughFreeCouriersAvailable(t *testing.T) {
	// Arrange
	ctx := context.Background()

	order1 := order.Fixtures.UnassignedOrder()
	order2 := order.Fixtures.UnassignedOrder()

	courier1 := courier.Fixtures.FreeCourier()

	transactionManager := persistence.NewMockITransactionManager(t)
	orderRepository := outPorts.NewMockIOrderRepository(t)
	courierRepository := outPorts.NewMockICourierRepository(t)
	dispatchService := services.NewMockIDispatchService(t)

	orderRepository.
		EXPECT().
		FindNew(ctx).
		Return([]*order.Order{order1, order2}, nil).
		Once()

	findFreeCalls := 0
	courierRepository.
		EXPECT().
		FindFree(ctx).
		RunAndReturn(func(_ context.Context) ([]*courier.Courier, error) {
			if findFreeCalls == 0 {
				findFreeCalls++
				return []*courier.Courier{courier1}, nil
			}

			return []*courier.Courier{}, nil
		}).
		Twice()

	dispatchService.
		EXPECT().
		Dispatch(order1, []*courier.Courier{courier1}).
		Return(courier1, nil).
		Once()

	transactionManager.
		EXPECT().
		Execute(ctx, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Twice()

	courierRepository.
		EXPECT().
		Update(ctx, courier1).
		Return(nil).
		Once()

	orderRepository.
		EXPECT().
		Update(ctx, order1).
		Return(nil).
		Once()

	handler, err := commands.NewAssignOrdersCommandHandler(
		transactionManager,
		dispatchService,
		courierRepository,
		orderRepository,
	)
	require.NoError(t, err)

	// Act
	err = handler.Handle(ctx)

	// Assert
	assert.ErrorIs(t, err, errors.NoAvailableCouriersError)
	assert.Equal(t, courier.StatusBusy, courier1.Status())
	assert.Equal(t, order.StatusAssigned, order1.Status())
	assert.Equal(t, order.StatusCreated, order2.Status())
}
