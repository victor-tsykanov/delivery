package order

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/victor-tsykanov/delivery/internal/common/errors"
	"github.com/victor-tsykanov/delivery/internal/core/domain/kernel"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/courier"
)

func TestNewOrder(t *testing.T) {
	testCases := []struct {
		name      string
		id        uuid.UUID
		location  kernel.Location
		wantErr   error
		wantOrder bool
	}{
		{
			name:      "valid order",
			id:        uuid.New(),
			location:  *kernel.RandomLocation(),
			wantErr:   nil,
			wantOrder: true,
		},
		{
			name:      "empty id",
			id:        uuid.Nil,
			location:  *kernel.RandomLocation(),
			wantErr:   errors.NewValueIsRequiredError("id"),
			wantOrder: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			order, err := NewOrder(tt.id, tt.location)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, order)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, order)
			assert.Equal(t, tt.id, order.ID())
			assert.Equal(t, tt.location, order.Location())
			assert.Equal(t, StatusCreated, order.Status())
			assert.Nil(t, order.CourierID())
		})
	}
}

func TestOrder_Assign(t *testing.T) {
	tests := []struct {
		name    string
		order   *Order
		wantErr error
	}{
		{
			name:    "unassigned order",
			order:   Fixtures.UnassignedOrder(),
			wantErr: nil,
		},
		{
			name:    "assigned order",
			order:   Fixtures.AssignedOrder(),
			wantErr: errors.NewInvalidStateError("order is already assigned"),
		},
		{
			name:    "completed order",
			order:   Fixtures.CompletedOrder(),
			wantErr: errors.NewInvalidStateError("order is completed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bob := courier.Fixtures.FreeCourier()
			originalStatus := tt.order.Status()
			originalCourierID := tt.order.CourierID()

			err := tt.order.Assign(bob)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
				assert.Equal(t, originalStatus, tt.order.Status())
				assert.Equal(t, originalCourierID, tt.order.CourierID())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, StatusAssigned, tt.order.Status())
			assert.Equal(t, bob.ID(), *(tt.order.CourierID()))
		})
	}
}

func TestOrder_Complete(t *testing.T) {
	tests := []struct {
		name    string
		order   *Order
		wantErr error
	}{
		{
			name:    "assigned order",
			order:   Fixtures.AssignedOrder(),
			wantErr: nil,
		},
		{
			name:    "unassigned order",
			order:   Fixtures.UnassignedOrder(),
			wantErr: errors.NewInvalidStateError("order needs to be assigned"),
		},
		{
			name:    "completed order",
			order:   Fixtures.CompletedOrder(),
			wantErr: errors.NewInvalidStateError("order is already completed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalStatus := tt.order.Status()

			err := tt.order.Complete()

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
				assert.Equal(t, originalStatus, tt.order.Status())

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, StatusCompleted, tt.order.Status())
		})
	}
}
