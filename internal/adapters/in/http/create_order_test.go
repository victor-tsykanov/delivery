package http_test

import (
	"context"
	"errors"
	netHttp "net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/victor-tsykanov/delivery/internal/adapters/in/http"
	inPortsMocks "github.com/victor-tsykanov/delivery/mocks/github.com/victor-tsykanov/delivery/internal_/core/ports/in"
	"github.com/victor-tsykanov/delivery/pkg/servers"
)

func TestServer_CreateOrder(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockCommandHandler := inPortsMocks.NewMockICreateOrderCommandHandler(t)
		mockCommandHandler.
			EXPECT().
			Handle(ctx, mock.Anything).
			Return(nil).
			Once()

		server, err := http.NewServer(
			mockCommandHandler,
			inPortsMocks.NewMockIGetAllCouriersQueryHandler(t),
			inPortsMocks.NewMockIGetPendingOrdersQueryHandler(t),
		)
		require.NoError(t, err)

		// Act
		response, err := server.CreateOrder(ctx, servers.CreateOrderRequestObject{})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, servers.CreateOrder201Response{}, response)
	})

	t.Run("error", func(t *testing.T) {
		// Arrange
		expectedErr := errors.New("boom")
		mockCommandHandler := inPortsMocks.NewMockICreateOrderCommandHandler(t)
		mockCommandHandler.
			EXPECT().
			Handle(ctx, mock.Anything).
			Return(expectedErr).
			Once()

		server, err := http.NewServer(
			mockCommandHandler,
			inPortsMocks.NewMockIGetAllCouriersQueryHandler(t),
			inPortsMocks.NewMockIGetPendingOrdersQueryHandler(t),
		)
		require.NoError(t, err)

		// Act
		response, err := server.CreateOrder(ctx, servers.CreateOrderRequestObject{})

		// Assert
		assert.Equal(t, echo.NewHTTPError(netHttp.StatusUnprocessableEntity, expectedErr), err)
		assert.Nil(t, response)
	})
}
