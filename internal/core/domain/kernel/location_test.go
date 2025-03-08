package kernel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victor-tsykanov/delivery/internal/common/errors"
)

func TestNewLocation(t *testing.T) {
	type args struct {
		x int
		y int
	}

	tests := []struct {
		name    string
		args    args
		want    *Location
		wantErr error
	}{
		{
			name: "valid value",
			args: args{
				x: 4,
				y: 7,
			},
			want: &Location{
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
			want: nil,
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
			location, err := NewLocation(tt.args.x, tt.args.y)

			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, location)
		})
	}
}

func TestLocation_Equals(t *testing.T) {
	tests := []struct {
		name   string
		first  Location
		second Location
		want   bool
	}{
		{
			name:   "equal",
			first:  Location{1, 2},
			second: Location{1, 2},
			want:   true,
		},
		{
			name:   "not equal",
			first:  Location{5, 2},
			second: Location{1, 2},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isEqual := tt.first.Equals(tt.second)

			assert.Equal(t, tt.want, isEqual)
		})
	}
}

func TestRandomLocation(t *testing.T) {
	location := RandomLocation()

	assert.NotEmpty(t, *location)
}

func TestLocation_DistanceTo(t *testing.T) {
	tests := []struct {
		name   string
		first  Location
		second Location
		want   int
	}{
		{
			name:   "different locations",
			second: Location{4, 9},
			first:  Location{2, 6},
			want:   5,
		},
		{
			name:   "same location",
			first:  Location{2, 4},
			second: Location{2, 4},
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			distance := tt.first.DistanceTo(tt.second)

			assert.Equal(t, tt.want, distance)
		})
	}
}
