package kernel_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victor-tsykanov/delivery/internal/common/errors"
	"github.com/victor-tsykanov/delivery/internal/core/domain/kernel"
)

func TestNewLocation(t *testing.T) {
	type args struct {
		x int
		y int
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "valid value",
			args: args{
				x: 4,
				y: 7,
			},
			wantErr: nil,
		},
		{
			name: "y coordinate out of range",
			args: args{
				x: 4,
				y: 0,
			},
			wantErr: errors.NewValueIsOutOfRangeError(
				"y",
				0,
				1,
				10,
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location, err := kernel.NewLocation(tt.args.x, tt.args.y)

			assert.Equal(t, tt.wantErr, err)
			if tt.wantErr == nil {
				assert.Equal(t, tt.args.x, location.X())
				assert.Equal(t, tt.args.y, location.Y())
			}
		})
	}
}

func TestLocation_Equals(t *testing.T) {
	tests := []struct {
		name   string
		first  *kernel.Location
		second *kernel.Location
		want   bool
	}{
		{
			name:   "equal",
			first:  kernel.MustNewLocation(1, 2),
			second: kernel.MustNewLocation(1, 2),
			want:   true,
		},
		{
			name:   "not equal",
			first:  kernel.MustNewLocation(5, 2),
			second: kernel.MustNewLocation(1, 2),
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

func TestRandomLocation(t *testing.T) {
	location := kernel.RandomLocation()

	assert.NotEmpty(t, *location)
}

func TestLocation_DistanceTo(t *testing.T) {
	tests := []struct {
		name   string
		first  *kernel.Location
		second *kernel.Location
		want   int
	}{
		{
			name:   "different locations",
			second: kernel.MustNewLocation(4, 9),
			first:  kernel.MustNewLocation(2, 6),
			want:   5,
		},
		{
			name:   "same location",
			first:  kernel.MustNewLocation(2, 4),
			second: kernel.MustNewLocation(2, 4),
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			distance := tt.first.DistanceTo(*tt.second)

			assert.Equal(t, tt.want, distance)
		})
	}
}
