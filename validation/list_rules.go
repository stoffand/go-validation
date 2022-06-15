package validation

import (
	"fmt"
)

type ListRule[T any] interface {
	Validate([]T) error
}

type ListRules[T any] struct {
	ListRules    Rule[[]T]
	ElementRules Rule[T]
}

func (r ListRules[T]) Validate(in []T) error {
	var SliceErrs SliceError
	switch rr := r.ListRules.(type) {
	case Rules[[]T]:
		for _, rule := range rr {
			if err := rule.Validate(in); err != nil {
				SliceErrs.AddList(err)
			}
		}
	case Rule[[]T]:
		if err := rr.Validate(in); err != nil {
			SliceErrs.AddList(err)
		}
	}
	for i, v := range in { // Loop over elements
		switch rr := r.ElementRules.(type) {
		case Rules[T]: // if rules list
			for _, rule := range rr {
				if err := rule.Validate(v); err != nil {
					SliceErrs.AddElem(i, err)
				}
			}
		case Rule[T]: // if single rule
			if err := rr.Validate(v); err != nil {
				SliceErrs.AddElem(i, err)
			}
		}
	}
	if len(SliceErrs.FailedListRules.Errs) == 0 && len(SliceErrs.FailedElementRules.Errs) == 0 {
		return nil
	}
	return SliceErrs
}

type ListMinLenghtRule[T any] struct {
	Len       int
	Inclusive bool
}

func (r ListMinLenghtRule[T]) Validate(input []T) error {
	if r.Inclusive {
		if len(input) < int(r.Len) {
			return fmt.Errorf("list length needs to be longer than or equal to %d", r.Len)
		}
		return nil
	} else {
		if len(input) <= int(r.Len) {
			return fmt.Errorf("list length needs to be longer than %d", r.Len)
		}
		return nil
	}
}

type ListMaxLenghtRule[T any] struct {
	Len       int
	Inclusive bool
}

func (r ListMaxLenghtRule[T]) Validate(input []T) error {
	if r.Inclusive {
		if len(input) > int(r.Len) {
			return fmt.Errorf("list length needs to be shorter than or equal to %d", r.Len)
		}
		return nil
	} else {
		if len(input) >= int(r.Len) {
			return fmt.Errorf("list length needs to be shorter than %d", r.Len)
		}
		return nil
	}
}

type ListExactLenghtRule[T any] struct {
	Len int
}

func (r ListExactLenghtRule[T]) Validate(input []T) error {
	if len(input) != int(r.Len) {
		return fmt.Errorf("list length needs to be exactly %d", r.Len)
	}
	return nil
}
