package courier

import (
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/victor-tsykanov/delivery/internal/common/errors"
	"github.com/victor-tsykanov/delivery/internal/core/domain/kernel"
)

func TestNewCourier(t *testing.T) {
	location := kernel.RandomLocation()
	transport, err := NewTransport(uuid.New(), "Car", 3)
	require.NoError(t, err)

	tests := []struct {
		name         string
		courrierName string
		transport    *Transport
		location     *kernel.Location
		wantErr      error
	}{
		{
			name:         "valid",
			courrierName: "John Doe",
			transport:    transport,
			location:     location,
			wantErr:      nil,
		},
		{
			name:         "empty name",
			courrierName: "",
			transport:    transport,
			location:     location,
			wantErr:      errors.NewValueIsRequiredError("name"),
		},
		{
			name:         "nil transport",
			courrierName: "John Doe",
			transport:    nil,
			location:     location,
			wantErr:      errors.NewValueIsRequiredError("transport"),
		},
		{
			name:         "nil location",
			courrierName: "John Doe",
			transport:    transport,
			location:     nil,
			wantErr:      errors.NewValueIsRequiredError("location"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			courier, err := NewCourier(tt.courrierName, tt.transport, tt.location)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, courier)

				return
			}

			assert.NotNil(t, courier)
			assert.Equal(t, tt.courrierName, courier.Name())
			assert.Equal(t, tt.transport, courier.Transport())
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
	transport, err := NewTransport(uuid.New(), "Car", 3)
	require.NoError(t, err)

	currentLocation := newLocation(1, 1)
	targetLocation := newLocation(1, 9)
	expectedLocation := newLocation(1, 4)

	courier, err := NewCourier("John Doe", transport, currentLocation)
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
			transport, err := NewTransport(uuid.New(), "Car", 3)
			require.NoError(t, err)

			courier, err := NewCourier("John Doe", transport, tt.from)
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
