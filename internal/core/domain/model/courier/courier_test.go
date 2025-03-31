package courier_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/victor-tsykanov/delivery/internal/common/errors"
	"github.com/victor-tsykanov/delivery/internal/core/domain/kernel"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/courier"
)

func TestNewCourier(t *testing.T) {
	location := kernel.RandomLocation()

	tests := []struct {
		name           string
		courierName    string
		transportName  string
		transportSpeed int
		location       *kernel.Location
		wantErr        error
	}{
		{
			name:           "valid",
			courierName:    "John Doe",
			transportName:  "Car",
			transportSpeed: 3,
			location:       location,
			wantErr:        nil,
		},
		{
			name:           "empty name",
			courierName:    "",
			transportName:  "Car",
			transportSpeed: 3,
			location:       location,
			wantErr:        errors.NewValueIsRequiredError("name"),
		},
		{
			name:           "empty transport name",
			courierName:    "John Doe",
			transportName:  "",
			transportSpeed: 3,
			location:       location,
			wantErr:        errors.NewValueIsRequiredError("transportName"),
		},
		{
			name:           "invalid transport speed",
			courierName:    "John Doe",
			transportName:  "Car",
			transportSpeed: 3000,
			location:       location,
			wantErr:        errors.NewValueIsOutOfRangeError("transportSpeed", 3000, 1, 3),
		},
		{
			name:           "nil location",
			courierName:    "John Doe",
			transportName:  "Car",
			transportSpeed: 3,
			location:       nil,
			wantErr:        errors.NewValueIsRequiredError("location"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theCourier, err := courier.NewCourier(
				tt.courierName,
				tt.transportName,
				tt.transportSpeed,
				tt.location,
			)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, theCourier)

				return
			}

			assert.NotNil(t, theCourier)
			assert.Equal(t, tt.courierName, theCourier.Name())
			assert.Equal(t, tt.transportName, theCourier.Transport().Name())
			assert.Equal(t, tt.transportSpeed, theCourier.Transport().Speed())
			assert.Equal(t, tt.location, theCourier.Location())
			assert.Equal(t, courier.StatusFree, theCourier.Status())
		})
	}
}

func TestCourier_SetBusy(t *testing.T) {
	tests := []struct {
		name       string
		courier    *courier.Courier
		wantStatus courier.Status
		wantErr    error
	}{
		{
			name:       "free courier",
			courier:    courier.Fixtures.FreeCourier(),
			wantStatus: courier.StatusBusy,
			wantErr:    nil,
		},
		{
			name:       "busy courier",
			courier:    courier.Fixtures.BusyCourier(),
			wantStatus: courier.StatusBusy,
			wantErr:    errors.NewInvalidStateError("courier must be free"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.courier.SetBusy()

			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantStatus, tt.courier.Status())
		})
	}
}

func TestCourier_SetFree(t *testing.T) {
	tests := []struct {
		name       string
		courier    *courier.Courier
		wantStatus courier.Status
		wantErr    error
	}{
		{
			name:       "free courier",
			courier:    courier.Fixtures.FreeCourier(),
			wantStatus: courier.StatusFree,
			wantErr:    errors.NewInvalidStateError("courier must be busy"),
		},
		{
			name:       "busy courier",
			courier:    courier.Fixtures.BusyCourier(),
			wantStatus: courier.StatusFree,
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.courier.SetFree()

			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantStatus, tt.courier.Status())
		})
	}
}

func TestCourier_Move(t *testing.T) {
	currentLocation := kernel.MustNewLocation(1, 1)
	targetLocation := kernel.MustNewLocation(1, 9)
	expectedLocation := kernel.MustNewLocation(1, 4)

	john, err := courier.NewCourier("John Doe", "Car", 3, currentLocation)
	require.NoError(t, err)

	err = john.Move(targetLocation)

	assert.NoError(t, err)
	assert.Equal(t, expectedLocation, john.Location())
}

func TestCourier_CalculateStepsToLocation(t *testing.T) {
	tests := []struct {
		name string
		from *kernel.Location
		to   *kernel.Location
		want int
	}{
		{
			name: "same location",
			from: kernel.MustNewLocation(1, 1),
			to:   kernel.MustNewLocation(1, 1),
			want: 0,
		},
		{
			name: "steps without remaining",
			from: kernel.MustNewLocation(1, 1),
			to:   kernel.MustNewLocation(7, 1),
			want: 2,
		},
		{
			name: "steps with remaining",
			from: kernel.MustNewLocation(1, 1),
			to:   kernel.MustNewLocation(9, 1),
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			john, err := courier.NewCourier("John Doe", "Car", 3, tt.from)
			require.NoError(t, err)

			steps, err := john.CalculateStepsToLocation(tt.to)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, steps)
		})
	}
}
