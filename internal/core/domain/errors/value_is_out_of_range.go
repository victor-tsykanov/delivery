package errors

import "fmt"

type ValueIsOutOfRangeError[T any] struct {
	paramName string
	value     T
	min       T
	max       T
}

func NewValueIsOutOfRangeError[T any](paramName string, value T, minValue T, maxValue T) *ValueIsOutOfRangeError[T] {
	return &ValueIsOutOfRangeError[T]{
		paramName: paramName,
		value:     value,
		min:       minValue,
		max:       maxValue,
	}
}

func (v *ValueIsOutOfRangeError[T]) Error() string {
	return fmt.Sprintf(
		"value %v of %v is out of range, must be between %v and %v",
		v.value,
		v.paramName,
		v.min,
		v.max,
	)
}
