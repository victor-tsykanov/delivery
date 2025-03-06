package math

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbs(t *testing.T) {
	assert.Equal(t, 42, Abs(42))
	assert.Equal(t, 0, Abs(0))
	assert.Equal(t, 42, Abs(-42))

	assert.Equal(t, int8(127), Abs(int8(127)))
	assert.Equal(t, int8(0), Abs(int8(0)))
	assert.Equal(t, int8(127), Abs(int8(-127)))

	assert.Equal(t, int16(32000), Abs(int16(32000)))
	assert.Equal(t, int16(0), Abs(int16(0)))
	assert.Equal(t, int16(32000), Abs(int16(-32000)))

	assert.Equal(t, int32(2147483647), Abs(int32(2147483647)))
	assert.Equal(t, int32(0), Abs(int32(0)))
	assert.Equal(t, int32(2147483647), Abs(int32(-2147483647)))

	assert.Equal(t, int64(9223372036854775807), Abs(int64(9223372036854775807)))
	assert.Equal(t, int64(0), Abs(int64(0)))
	assert.Equal(t, int64(9223372036854775807), Abs(int64(-9223372036854775807)))

	assert.Equal(t, float32(3.14), Abs(float32(3.14)))
	assert.Equal(t, float32(0), Abs(float32(0)))
	assert.Equal(t, float32(3.14), Abs(float32(-3.14)))

	assert.Equal(t, 3.14, Abs(3.14))
	assert.Equal(t, 0.0, Abs(0.0))
	assert.Equal(t, 3.14, Abs(-3.14))
}

func TestSign(t *testing.T) {
	assert.Equal(t, 1, Sign(42))
	assert.Equal(t, 1, Sign(0))
	assert.Equal(t, -1, Sign(-42))

	assert.Equal(t, int8(1), Sign(int8(127)))
	assert.Equal(t, int8(1), Sign(int8(0)))
	assert.Equal(t, int8(-1), Sign(int8(-127)))

	assert.Equal(t, int16(1), Sign(int16(32000)))
	assert.Equal(t, int16(1), Sign(int16(0)))
	assert.Equal(t, int16(-1), Sign(int16(-32000)))

	assert.Equal(t, int32(1), Sign(int32(2147483647)))
	assert.Equal(t, int32(1), Sign(int32(0)))
	assert.Equal(t, int32(-1), Sign(int32(-2147483647)))

	assert.Equal(t, int64(1), Sign(int64(9223372036854775807)))
	assert.Equal(t, int64(1), Sign(int64(0)))
	assert.Equal(t, int64(-1), Sign(int64(-9223372036854775807)))

	assert.Equal(t, float32(1), Sign(float32(3.14)))
	assert.Equal(t, float32(1), Sign(float32(0)))
	assert.Equal(t, float32(-1), Sign(float32(-3.14)))

	assert.Equal(t, 1.0, Sign(3.14))
	assert.Equal(t, 1.0, Sign(0.0))
	assert.Equal(t, -1.0, Sign(-3.14))
}
