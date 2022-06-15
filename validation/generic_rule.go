package validation

// All rules follow this interface
type Rule[T any] interface {
	Validate(T) error
}

// Implementation of rules interface for a list of rules
// So this is a also a Rule[T]
type Rules[T any] []Rule[T]

func (r Rules[T]) Validate(in T) error {
	var err error
	var e []error
	for _, v := range r {
		if err = v.Validate(in); err != nil {
			e = append(e, err)
		}
	}
	if len(e) == 0 {
		return nil
	}
	return ErrList{Errs: e}
}
