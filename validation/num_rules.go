package validation

import "fmt"

type Num interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

type NumLteRule[T Num] struct {
	Value int32
}

func (r NumLteRule[T]) Validate(input T) error {
	if input > T(r.Value) {
		return fmt.Errorf("num needs to be less than or equal to %d", r.Value)
	}
	return nil
}

type NumLtRule[T Num] struct {
	Value int32
}

func (r NumLtRule[T]) Validate(input T) error {
	if input >= T(r.Value) {
		return fmt.Errorf("num needs to be less than %d", r.Value)
	}
	return nil
}

type NumGteRule[T Num] struct {
	Value int32
}

func (r NumGteRule[T]) Validate(input T) error {
	if input < T(r.Value) {
		return fmt.Errorf("num needs to be greater than or equal to %d", r.Value)
	}
	return nil
}

type NumGtRule[T Num] struct {
	Value int32
}

func (r NumGtRule[T]) Validate(input T) error {
	if input <= T(r.Value) {
		return fmt.Errorf("num needs to be greater than %d", r.Value)
	}
	return nil
}

type NumOneOf[T Num] struct {
	List []T
}

func (r NumOneOf[T]) Validate(input T) error {
	for _, v := range r.List {
		if input == v {
			return nil
		}
	}
	return fmt.Errorf("Num was not one of %v", r.List)
}

type NumNotOneOf[T Num] struct {
	List []T
}

func (r NumNotOneOf[T]) Validate(input T) error {
	for _, v := range r.List {
		if input == v {
			return fmt.Errorf("Num was one of illegal %v", r.List)
		}
	}
	return nil
}

type NumPositive[T Num] struct{}

func (r NumPositive[T]) Validate(input T) error {
	if input < 0 {
		return fmt.Errorf("Num cant be negative")
	}
	return nil
}

type NumNegative[T Num] struct{}

func (r NumNegative[T]) Validate(input T) error {
	if input >= 0 {
		return fmt.Errorf("Num cant be positive")
	}
	return nil
}
