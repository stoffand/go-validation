package validation

import "fmt"

type BoolIsTrue struct{}

func (r BoolIsTrue) Validate(in bool) error {
	if !in {
		return fmt.Errorf("bool has to be true")
	}
	return nil
}

type BoolIsFalse struct{}

func (r BoolIsFalse) Validate(in bool) error {
	if in {
		return fmt.Errorf("bool has to be false")
	}
	return nil
}
