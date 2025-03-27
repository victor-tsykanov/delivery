package http_test

import (
	"context"
	"errors"
	netHttp "net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/victor-tsykanov/delivery/internal/adapters/in/http"
	"github.com/victor-tsykanov/delivery/internal/core/domain/kernel"
	inPorts "github.com/victor-tsykanov/delivery/internal/core/ports/in"
	inPortsMocks "github.com/victor-tsykanov/delivery/mocks/github.com/victor-tsykanov/delivery/internal_/core/ports/in"
	"github.com/victor-tsykanov/delivery/pkg/servers"
)

func TestServer_GetOrders(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		// Arrange
		firstOrderID := uuid.New()
		firstOrderLocation := kernel.RandomLocation()
		secondOrderID := uuid.New()
		secondOrderLocation := kernel.RandomLocation()

		mockQueryHandler := inPortsMocks.NewMockIGetPendingOrdersQueryHandler(t)
		mockQueryHandler.
			EXPECT().
			Handle(ctx).
			Return(
				[]*inPorts.Order{
					{
						ID: firstOrderID,
						Location: &inPorts.OrderLocation{
							X: firstOrderLocation.X(),
							Y: firstOrderLocation.Y(),
						},
					},
					{
						ID: secondOrderID,
						Location: &inPorts.OrderLocation{
							X: secondOrderLocation.X(),
							Y: secondOrderLocation.Y(),
						},
					},
				},
				nil,
			).
			Once()

		server, err := http.NewServer(
			inPortsMocks.NewMockICreateOrderCommandHandler(t),
			inPortsMocks.NewMockIGetAllCouriersQueryHandler(t),
			mockQueryHandler,
		)
		require.NoError(t, err)

		// Act
		response, err := server.GetOrders(ctx, servers.GetOrdersRequestObject{})

		// Assert
		assert.NoError(t, err)
		assert.Equal(
			t,
			servers.GetOrders200JSONResponse(
				[]servers.Order{
					{
						Id: firstOrderID,
						Location: servers.Location{
							X: firstOrderLocation.X(),
							Y: firstOrderLocation.Y(),
						},
					},
					{
						Id: secondOrderID,
						Location: servers.Location{
							X: secondOrderLocation.X(),
							Y: secondOrderLocation.Y(),
						},
					},
				}),
			response,
		)
	})

	t.Run("error", func(t *testing.T) {
		// Arrange
		expectedErr := errors.New("boom")
		mockQueryHandler := inPortsMocks.NewMockIGetPendingOrdersQueryHandler(t)
		mockQueryHandler.
			EXPECT().
			Handle(ctx).
			Return(nil, expectedErr).
			Once()

		server, err := http.NewServer(
			inPortsMocks.NewMockICreateOrderCommandHandler(t),
			inPortsMocks.NewMockIGetAllCouriersQueryHandler(t),
			mockQueryHandler,
		)
		require.NoError(t, err)

		// Act
		response, err := server.GetOrders(ctx, servers.GetOrdersRequestObject{})

		// Assert
		assert.Equal(t, echo.NewHTTPError(netHttp.StatusInternalServerError, expectedErr), err)
		assert.Nil(t, response)
	})
}
