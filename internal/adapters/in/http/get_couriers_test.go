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

func TestServer_GetCouriers(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		// Arrange
		firstCourierID := uuid.New()
		firstCourierLocation := kernel.RandomLocation()
		secondCourierID := uuid.New()
		secondCourierLocation := kernel.RandomLocation()

		mockQueryHandler := inPortsMocks.NewMockIGetAllCouriersQueryHandler(t)
		mockQueryHandler.
			EXPECT().
			Handle(ctx).
			Return(
				[]*inPorts.Courier{
					{
						ID:   firstCourierID,
						Name: "first",
						Location: &inPorts.CourierLocation{
							X: firstCourierLocation.X(),
							Y: firstCourierLocation.Y(),
						},
					},
					{
						ID:   secondCourierID,
						Name: "second",
						Location: &inPorts.CourierLocation{
							X: secondCourierLocation.X(),
							Y: secondCourierLocation.Y(),
						},
					},
				},
				nil,
			).
			Once()

		server, err := http.NewServer(
			inPortsMocks.NewMockICreateOrderCommandHandler(t),
			mockQueryHandler,
			inPortsMocks.NewMockIGetPendingOrdersQueryHandler(t),
		)
		require.NoError(t, err)

		// Act
		response, err := server.GetCouriers(ctx, servers.GetCouriersRequestObject{})

		// Assert
		assert.NoError(t, err)
		assert.Equal(
			t,
			servers.GetCouriers200JSONResponse(
				[]servers.Courier{
					{
						Id: firstCourierID,
						Location: servers.Location{
							X: firstCourierLocation.X(),
							Y: firstCourierLocation.Y(),
						},
						Name: "first",
					},
					{
						Id: secondCourierID,
						Location: servers.Location{
							X: secondCourierLocation.X(),
							Y: secondCourierLocation.Y(),
						},
						Name: "second",
					},
				}),
			response,
		)
	})

	t.Run("error", func(t *testing.T) {
		// Arrange
		expectedErr := errors.New("boom")
		mockQueryHandler := inPortsMocks.NewMockIGetAllCouriersQueryHandler(t)
		mockQueryHandler.
			EXPECT().
			Handle(ctx).
			Return(nil, expectedErr).
			Once()

		server, err := http.NewServer(
			inPortsMocks.NewMockICreateOrderCommandHandler(t),
			mockQueryHandler,
			inPortsMocks.NewMockIGetPendingOrdersQueryHandler(t),
		)
		require.NoError(t, err)

		// Act
		response, err := server.GetCouriers(ctx, servers.GetCouriersRequestObject{})

		// Assert
		assert.Equal(t, echo.NewHTTPError(netHttp.StatusInternalServerError, expectedErr), err)
		assert.Nil(t, response)
	})
}
