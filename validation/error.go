package validation

import (
	"encoding/json"
	"fmt"
)

// Error types
type RequiredErr struct{}

type ErrList struct {
	Errs []error `json:",omitempty"`
}

func (e *ErrList) Add(err error) {
	e.Errs = append(e.Errs, err)
}

type ErrMap[K comparable] struct {
	Errs map[K]error `json:",omitempty"`
}

func (e *ErrMap[K]) Add(index K, err error) {
	if e.Errs == nil {
		e.Errs = make(map[K]error)
	}
	e.Errs[index] = err
}

type StructErr struct {
	FailedFields map[string]error `json:",omitempty"`
}

func (e *StructErr) AddError(field string, err error) {
	if e.FailedFields == nil {
		e.FailedFields = make(map[string]error)
	}
	e.FailedFields[field] = err
}

type SliceError struct {
	FailedListRules    ErrList     `json:",omitempty"`
	FailedElementRules ErrMap[int] `json:",omitempty"`
}

func (e *SliceError) AddList(err error) {
	e.FailedListRules.Add(err)
}
func (e *SliceError) AddElem(index int, err error) {
	e.FailedElementRules.Add(index, err)
}

type MapError[K comparable, V any] struct {
	FailedMapRules   ErrList   `json:",omitempty"`
	FailedKeyRules   ErrMap[K] `json:",omitempty"`
	FailedValueRules ErrMap[K] `json:",omitempty"`
}

func (e *MapError[K, V]) AddMap(err error) {
	e.FailedMapRules.Add(err)
}
func (e *MapError[K, V]) AddKeys(index K, err error) {
	e.FailedKeyRules.Add(index, err)
}
func (e *MapError[K, V]) AddValues(index K, err error) {
	e.FailedValueRules.Add(index, err)
}

// Implement Error
func (e RequiredErr) Error() string {
	return "required"
}
func (e ErrList) Error() string {
	return fmt.Sprintf("%v", e.Errs)
}
func (e ErrMap[K]) Error() string {
	return fmt.Sprintf("%v", e.Errs)
}
func (e StructErr) Error() string {
	return fmt.Sprintf("%v", e.FailedFields)
}

func (e SliceError) Error() string {
	var list string
	var elems string
	if len(e.FailedListRules.Errs) > 0 {
		list = fmt.Sprintf("list: %v, ", e.FailedListRules.Errs)
	}
	if len(e.FailedElementRules.Errs) > 0 {
		elems = fmt.Sprintf("elements: %v", e.FailedElementRules.Errs)
	}
	return " " + list + elems
}

func (e MapError[K, V]) Error() string {
	var m string
	var keys string
	var values string
	if len(e.FailedMapRules.Errs) > 0 {
		m = fmt.Sprintf("map: %v, ", e.FailedMapRules.Errs)
	}
	if len(e.FailedKeyRules.Errs) > 0 {
		keys = fmt.Sprintf("keys: %v, ", e.FailedKeyRules.Errs)
	}
	if len(e.FailedValueRules.Errs) > 0 {
		values = fmt.Sprintf("values: %v", e.FailedValueRules.Errs)
	}
	return " " + m + keys + values
}

// Custom marshalling
func (e RequiredErr) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Error())
}
func (e ErrList) MarshalJSON() ([]byte, error) {
	str := make([]string, len(e.Errs))
	for i, v := range e.Errs {
		str[i] = v.Error()
	}
	return json.Marshal(str)
}

func (e SliceError) MarshalJSON() ([]byte, error) {
	res := struct {
		List  *ErrList     `json:",omitempty"`
		Elems *ErrMap[int] `json:",omitempty"`
	}{}
	if len(e.FailedListRules.Errs) > 0 {
		res.List = &e.FailedListRules
	}
	if len(e.FailedElementRules.Errs) > 0 {
		res.Elems = &e.FailedElementRules
	}
	return json.Marshal(res)
}

func (e MapError[K, V]) MarshalJSON() ([]byte, error) {
	res := struct {
		Map    *ErrList   `json:",omitempty"`
		Keys   *ErrMap[K] `json:",omitempty"`
		Values *ErrMap[K] `json:",omitempty"`
	}{}
	if len(e.FailedMapRules.Errs) > 0 {
		res.Map = &e.FailedMapRules
	}
	if len(e.FailedKeyRules.Errs) > 0 {
		res.Keys = &e.FailedKeyRules
	}
	if len(e.FailedValueRules.Errs) > 0 {
		res.Values = &e.FailedValueRules
	}
	return json.Marshal(res)
}

// TODO Using any now, better solution?
func (e ErrMap[K]) MarshalJSON() ([]byte, error) {
	m := make(map[string]any, len(e.Errs)) // Json only allows string keys
	for k, v := range e.Errs {
		switch v.(type) {
		case RequiredErr, ErrList, ErrMap[K], StructErr, SliceError, MapError[K, error]:
			m[fmt.Sprintf("%v", k)] = v
		default: // Errors not defined by package
			m[fmt.Sprintf("%v", k)] = v.Error()
		}
	}
	return json.Marshal(m)
}
func (e StructErr) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.FailedFields)
}
