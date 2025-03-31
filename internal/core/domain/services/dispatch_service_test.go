package services_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/victor-tsykanov/delivery/internal/common/errors"
	"github.com/victor-tsykanov/delivery/internal/core/domain/kernel"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/courier"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/order"
	"github.com/victor-tsykanov/delivery/internal/core/domain/services"
)

func TestDispatchService_Dispatch(t *testing.T) {
	orderLocation, err := kernel.NewLocation(1, 1)
	require.NoError(t, err)

	theOrder, err := order.NewOrder(uuid.New(), *orderLocation)
	require.NoError(t, err)

	tests := []struct {
		name     string
		couriers []*courier.Courier
		wantIdx  int
		wantErr  error
	}{
		{
			name: "dispatch to closest",
			couriers: []*courier.Courier{
				courier.Fixtures.FreeCourierAtLocationWithSpeed(7, 1, 2),
				courier.Fixtures.FreeCourierAtLocationWithSpeed(5, 1, 2),
			},
			wantIdx: 1,
			wantErr: nil,
		},
		{
			name: "dispatch to first of equally close",
			couriers: []*courier.Courier{
				courier.Fixtures.FreeCourierAtLocationWithSpeed(7, 1, 2),
				courier.Fixtures.FreeCourierAtLocationWithSpeed(1, 7, 2),
				courier.Fixtures.FreeCourierAtLocationWithSpeed(1, 9, 2),
			},
			wantIdx: 0,
			wantErr: nil,
		},
		{
			name: "dispatch to fastest",
			couriers: []*courier.Courier{
				courier.Fixtures.FreeCourierAtLocationWithSpeed(7, 1, 2),
				courier.Fixtures.FreeCourierAtLocationWithSpeed(7, 1, 3),
			},
			wantIdx: 1,
			wantErr: nil,
		},
		{
			name:     "fail to dispatch if no couriers are passed",
			couriers: []*courier.Courier{},
			wantIdx:  0,
			wantErr:  errors.NewValueIsRequiredError("couriers"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := services.NewDispatchService()
			require.NoError(t, err)

			chosenCourier, err := service.Dispatch(theOrder, tt.couriers)

			assert.Equal(t, tt.wantErr, err)
			if err == nil {
				assert.Equal(t, tt.couriers[tt.wantIdx], chosenCourier)
			}
		})
	}
}
