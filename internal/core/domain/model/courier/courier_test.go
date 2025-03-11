package courier

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/victor-tsykanov/delivery/internal/common/errors"
	"github.com/victor-tsykanov/delivery/internal/core/domain/kernel"
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
			courier, err := NewCourier(
				tt.courierName,
				tt.transportName,
				tt.transportSpeed,
				tt.location,
			)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, courier)

				return
			}

			assert.NotNil(t, courier)
			assert.Equal(t, tt.courierName, courier.Name())
			assert.Equal(t, tt.transportName, courier.Transport().Name())
			assert.Equal(t, tt.transportSpeed, courier.Transport().Speed())
			assert.Equal(t, tt.location, courier.Location())
			assert.Equal(t, StatusFree, courier.Status())
		})
	}
}

func TestCourier_SetBusy(t *testing.T) {
	tests := []struct {
		name       string
		courier    *Courier
		wantStatus Status
		wantErr    error
	}{
		{
			name:       "free courier",
			courier:    Fixtures.FreeCourier(),
			wantStatus: StatusBusy,
			wantErr:    nil,
		},
		{
			name:       "busy courier",
			courier:    Fixtures.BusyCourier(),
			wantStatus: StatusBusy,
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
		courier    *Courier
		wantStatus Status
		wantErr    error
	}{
		{
			name:       "free courier",
			courier:    Fixtures.FreeCourier(),
			wantStatus: StatusFree,
			wantErr:    errors.NewInvalidStateError("courier must be busy"),
		},
		{
			name:       "busy courier",
			courier:    Fixtures.BusyCourier(),
			wantStatus: StatusFree,
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
	currentLocation := newLocation(1, 1)
	targetLocation := newLocation(1, 9)
	expectedLocation := newLocation(1, 4)

	courier, err := NewCourier("John Doe", "Car", 3, currentLocation)
	require.NoError(t, err)

	err = courier.Move(targetLocation)

	assert.NoError(t, err)
	assert.Equal(t, expectedLocation, courier.Location())
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
			from: newLocation(1, 1),
			to:   newLocation(1, 1),
			want: 0,
		},
		{
			name: "steps without remaining",
			from: newLocation(1, 1),
			to:   newLocation(7, 1),
			want: 2,
		},
		{
			name: "steps with remaining",
			from: newLocation(1, 1),
			to:   newLocation(9, 1),
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			courier, err := NewCourier("John Doe", "Car", 3, tt.from)
			require.NoError(t, err)

			steps, err := courier.CalculateStepsToLocation(tt.to)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, steps)
		})
	}
}

func newLocation(x int, y int) *kernel.Location {
	location, err := kernel.NewLocation(x, y)
	if err != nil {
		log.Fatalf("failed to create location: %v", err)
	}

	return location
}
