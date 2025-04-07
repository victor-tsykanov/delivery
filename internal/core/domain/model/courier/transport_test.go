package courier_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/victor-tsykanov/delivery/internal/common/errors"
	"github.com/victor-tsykanov/delivery/internal/core/domain/kernel"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/courier"
)

func TestNewTransport(t *testing.T) {
	type args struct {
		id    courier.TransportID
		name  string
		speed int
	}

	carID := courier.NewTransportID()

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "valid",
			args: args{
				id:    carID,
				name:  "Car",
				speed: 2,
			},
			wantErr: nil,
		},
		{
			name: "empty name",
			args: args{
				id:    carID,
				name:  "",
				speed: 3,
			},
			wantErr: errors.NewValueIsRequiredError("name"),
		},
		{
			name: "too low speed",
			args: args{
				id:    carID,
				name:  "Car",
				speed: 0,
			},
			wantErr: errors.NewValueIsOutOfRangeError("speed", 0, 1, 3),
		},
		{
			name: "too high speed",
			args: args{
				id:    carID,
				name:  "Car",
				speed: 4,
			},
			wantErr: errors.NewValueIsOutOfRangeError("speed", 4, 1, 3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport, err := courier.NewTransport(tt.args.id, tt.args.name, tt.args.speed)

			assert.Equal(t, tt.wantErr, err)
			if tt.wantErr == nil {
				assert.Equal(t, tt.args.id, transport.ID())
				assert.Equal(t, tt.args.name, transport.Name())
				assert.Equal(t, tt.args.speed, transport.Speed())
			}
		})
	}
}

func TestTransport_Equals(t *testing.T) {
	transport := func(id courier.TransportID, name string, speed int) *courier.Transport {
		transport, err := courier.NewTransport(id, name, speed)
		require.NoError(t, err)

		return transport
	}

	id := courier.NewTransportID()

	tests := []struct {
		name   string
		first  *courier.Transport
		second *courier.Transport
		want   bool
	}{
		{
			name:   "equal",
			first:  transport(id, "First", 1),
			second: transport(id, "Second", 1),
			want:   true,
		},
		{
			name:   "not equal",
			first:  transport(id, "First", 1),
			second: transport(courier.NewTransportID(), "Second", 1),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isEqual := tt.first.Equals(*tt.second)

			assert.Equal(t, tt.want, isEqual)
		})
	}
}

func TestTransport_Move(t *testing.T) {
	location := func(x int, y int) *kernel.Location {
		location, err := kernel.NewLocation(x, y)
		require.NoError(t, err)

		return location
	}

	car, err := courier.NewTransport(courier.NewTransportID(), "Car", 3)
	require.NoError(t, err)

	tests := []struct {
		name string
		from *kernel.Location
		to   *kernel.Location
		want *kernel.Location
	}{
		{
			name: "move ↘️ within the range with remaining",
			from: location(1, 1),
			to:   location(2, 2),
			want: location(2, 2),
		},
		{
			name: "move ↖️ within the range with remaining",
			from: location(2, 2),
			to:   location(1, 1),
			want: location(1, 1),
		},
		{
			name: "move ↘️ within the range without remaining",
			from: location(1, 1),
			to:   location(2, 3),
			want: location(2, 3),
		},
		{
			name: "move ↖️ within the range without remaining",
			from: location(2, 3),
			to:   location(1, 1),
			want: location(1, 1),
		},
		{
			name: "move ↘️ beyond the range only along X",
			from: location(1, 1),
			to:   location(5, 5),
			want: location(4, 1),
		},
		{
			name: "move ↖️ beyond the range only along X",
			from: location(5, 5),
			to:   location(1, 1),
			want: location(2, 5),
		},
		{
			name: "move ↘️ beyond the range along X and Y",
			from: location(1, 1),
			to:   location(3, 7),
			want: location(3, 2),
		},
		{
			name: "move ↖️ beyond the range along X and Y",
			from: location(3, 7),
			to:   location(1, 1),
			want: location(1, 6),
		},
		{
			name: "move to the same location",
			from: location(1, 1),
			to:   location(1, 1),
			want: location(1, 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newLocation, err := car.Move(tt.from, tt.to)

			assert.Equal(t, tt.want, newLocation)
			assert.NoError(t, err)
		})
	}
}
