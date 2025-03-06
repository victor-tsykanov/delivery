package math

type Numeric interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

func Abs[T Numeric](n T) T {
	if n < 0 {
		return -n
	}

	return n
}

func Sign[T Numeric](n T) T {
	if n < 0 {
		return T(-1)
	}

	return T(1)
}
