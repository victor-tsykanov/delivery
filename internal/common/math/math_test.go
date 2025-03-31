package math_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victor-tsykanov/delivery/internal/common/math"
)

func TestAbs(t *testing.T) {
	assert.Equal(t, 42, math.Abs(42))
	assert.Equal(t, 0, math.Abs(0))
	assert.Equal(t, 42, math.Abs(-42))

	assert.Equal(t, int8(127), math.Abs(int8(127)))
	assert.Equal(t, int8(0), math.Abs(int8(0)))
	assert.Equal(t, int8(127), math.Abs(int8(-127)))

	assert.Equal(t, int16(32000), math.Abs(int16(32000)))
	assert.Equal(t, int16(0), math.Abs(int16(0)))
	assert.Equal(t, int16(32000), math.Abs(int16(-32000)))

	assert.Equal(t, int32(2147483647), math.Abs(int32(2147483647)))
	assert.Equal(t, int32(0), math.Abs(int32(0)))
	assert.Equal(t, int32(2147483647), math.Abs(int32(-2147483647)))

	assert.Equal(t, int64(9223372036854775807), math.Abs(int64(9223372036854775807)))
	assert.Equal(t, int64(0), math.Abs(int64(0)))
	assert.Equal(t, int64(9223372036854775807), math.Abs(int64(-9223372036854775807)))

	assert.Equal(t, float32(3.14), math.Abs(float32(3.14)))
	assert.Equal(t, float32(0), math.Abs(float32(0)))
	assert.Equal(t, float32(3.14), math.Abs(float32(-3.14)))

	assert.Equal(t, 3.14, math.Abs(3.14))
	assert.Equal(t, 0.0, math.Abs(0.0))
	assert.Equal(t, 3.14, math.Abs(-3.14))
}

func TestSign(t *testing.T) {
	assert.Equal(t, 1, math.Sign(42))
	assert.Equal(t, 1, math.Sign(0))
	assert.Equal(t, -1, math.Sign(-42))

	assert.Equal(t, int8(1), math.Sign(int8(127)))
	assert.Equal(t, int8(1), math.Sign(int8(0)))
	assert.Equal(t, int8(-1), math.Sign(int8(-127)))

	assert.Equal(t, int16(1), math.Sign(int16(32000)))
	assert.Equal(t, int16(1), math.Sign(int16(0)))
	assert.Equal(t, int16(-1), math.Sign(int16(-32000)))

	assert.Equal(t, int32(1), math.Sign(int32(2147483647)))
	assert.Equal(t, int32(1), math.Sign(int32(0)))
	assert.Equal(t, int32(-1), math.Sign(int32(-2147483647)))

	assert.Equal(t, int64(1), math.Sign(int64(9223372036854775807)))
	assert.Equal(t, int64(1), math.Sign(int64(0)))
	assert.Equal(t, int64(-1), math.Sign(int64(-9223372036854775807)))

	assert.Equal(t, float32(1), math.Sign(float32(3.14)))
	assert.Equal(t, float32(1), math.Sign(float32(0)))
	assert.Equal(t, float32(-1), math.Sign(float32(-3.14)))

	assert.Equal(t, 1.0, math.Sign(3.14))
	assert.Equal(t, 1.0, math.Sign(0.0))
	assert.Equal(t, -1.0, math.Sign(-3.14))
}
