package commands_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/victor-tsykanov/delivery/internal/core/application/usecases/commands"
	"github.com/victor-tsykanov/delivery/internal/core/domain/kernel"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/order"
	inPorts "github.com/victor-tsykanov/delivery/internal/core/ports/in"
	"github.com/victor-tsykanov/delivery/mocks/github.com/victor-tsykanov/delivery/internal_/common/persistence"
	outPorts "github.com/victor-tsykanov/delivery/mocks/github.com/victor-tsykanov/delivery/internal_/core/ports/out"
)

func TestCreateOrderCommandHandler_Handle(t *testing.T) {
	// Arrange
	ctx := context.Background()
	orderID := uuid.New()
	street := gofakeit.Address().Address
	location := kernel.RandomLocation()

	transactionManager := persistence.NewMockITransactionManager(t)
	orderRepository := outPorts.NewMockIOrderRepository(t)
	geoClient := outPorts.NewMockIGeoClient(t)

	transactionManager.
		EXPECT().
		Execute(ctx, mock.Anything).
		Run(func(ctx context.Context, fn func(context.Context) error) {
			err := fn(ctx)
			require.NoError(t, err)
		}).
		Return(nil).
		Once()

	orderRepository.
		EXPECT().
		Create(ctx, mock.MatchedBy(func(newOrder *order.Order) bool {
			return newOrder.ID() == orderID &&
				newOrder.Status() == order.StatusCreated &&
				location.Equals(newOrder.Location())
		})).
		Return(nil).
		Once()

	geoClient.
		EXPECT().
		GetLocation(ctx, street).
		Return(location, nil).
		Once()

	handler, err := commands.NewCreateOrderCommandHandler(transactionManager, orderRepository, geoClient)
	require.NoError(t, err)

	command, err := inPorts.NewCreateOrderCommand(orderID, street)
	require.NoError(t, err)

	// Act
	err = handler.Handle(ctx, *command)

	// Assert
	require.NoError(t, err)
}
