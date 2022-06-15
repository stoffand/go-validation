package validation

import "fmt"

type MapRule[K comparable, V any] interface {
	Validate(map[K]V) error
}

type MapRules[K comparable, V any] struct {
	MapRules   Rule[map[K]V]
	KeyRules   Rule[K]
	ValueRules Rule[V]
}

func (r MapRules[K, V]) Validate(in map[K]V) error {
	var MapErrs MapError[K, V]

	// Map rules
	switch rr := r.MapRules.(type) {
	case Rules[map[K]V]:
		for _, rule := range rr {
			if err := rule.Validate(in); err != nil {
				MapErrs.AddMap(err)
			}
		}
	case Rule[map[K]V]:
		if err := rr.Validate(in); err != nil {
			MapErrs.AddMap(err)
		}
	}

	for k, v := range in { // loop over key value pairs
		// Key rules
		switch rr := r.KeyRules.(type) {
		case Rules[K]: // list of rules
			var errs ErrList
			for _, rule := range rr {
				if err := rule.Validate(k); err != nil {
					errs.Add(err)
				}
			}
			if len(errs.Errs) > 0 {
				MapErrs.AddKeys(k, errs)
			}
		case Rule[K]: // single rule
			if err := rr.Validate(k); err != nil {
				MapErrs.AddKeys(k, err)
			}
		}
		// Value rules
		switch rr := r.ValueRules.(type) {
		case Rules[V]: // list of rules
			var errs ErrList
			for _, rule := range rr {
				if err := rule.Validate(v); err != nil {
					errs.Add(err)
				}
			}
			if len(errs.Errs) > 0 {
				MapErrs.AddValues(k, errs)
			}
		case Rule[V]: // single rule
			if err := rr.Validate(v); err != nil {
				MapErrs.AddValues(k, err)
			}
		}
	}
	// fmt.Println(MapErrs.FailedMapRules.Errs, MapErrs.FailedKeyRules.Errs, MapErrs.FailedValueRules.Errs)
	if len(MapErrs.FailedMapRules.Errs) == 0 && len(MapErrs.FailedKeyRules.Errs) == 0 && len(MapErrs.FailedValueRules.Errs) == 0 {
		return nil
	}
	return MapErrs
}

// Map rules

type MapMinLenghtRule[K comparable, V any] struct {
	Len       int
	Inclusive bool
}

func (r MapMinLenghtRule[K, V]) Validate(input map[K]V) error {
	if r.Inclusive {
		if len(input) < int(r.Len) {
			return fmt.Errorf("map length needs to be longer than or equal to %d", r.Len)
		}
		return nil
	} else {
		if len(input) <= int(r.Len) {
			return fmt.Errorf("map length needs to be longer than %d", r.Len)
		}
		return nil
	}
}

type MapMaxLenghtRule[K comparable, V any] struct {
	Len       int
	Inclusive bool
}

func (r MapMaxLenghtRule[K, V]) Validate(input map[K]V) error {
	if r.Inclusive {
		if len(input) > int(r.Len) {
			return fmt.Errorf("map length needs to be shorter than or equal to %d", r.Len)
		}
		return nil
	} else {
		if len(input) >= int(r.Len) {
			return fmt.Errorf("map length needs to be shorter than %d", r.Len)
		}
		return nil
	}
}

type MapExactLenghtRule[K comparable, V any] struct {
	Len int
}

func (r MapExactLenghtRule[K, V]) Validate(input map[K]V) error {
	if len(input) != int(r.Len) {
		return fmt.Errorf("map length needs to be exactly %d", r.Len)
	}
	return nil
}
