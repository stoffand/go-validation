package validation

import (
	"fmt"
	"log"
	"regexp"
)

type StringMinLengthRule struct {
	Len       int32
	Inclusive bool
}

func (r StringMinLengthRule) Validate(input string) error {
	if r.Inclusive {
		if len(input) < int(r.Len) {
			return fmt.Errorf("string length needs to be longer than or equal to %d", r.Len)
		}
		return nil

	} else {
		if len(input) <= int(r.Len) {
			return fmt.Errorf("string length needs to be longer than %d", r.Len)
		}
		return nil

	}
}

type StringMaxLengthRule struct {
	Len       int32
	Inclusive bool
}

func (r StringMaxLengthRule) Validate(input string) error {
	if r.Inclusive {
		if len(input) > int(r.Len) {
			return fmt.Errorf("string length needs to be shorter than or equal to %d", r.Len)
		}
		return nil
	} else {
		if len(input) >= int(r.Len) {
			return fmt.Errorf("string length needs to be shorter than %d", r.Len)
		}
		return nil
	}
}

type StringExactLengthRule struct {
	Len int32
}

func (r StringExactLengthRule) Validate(input string) error {
	if len(input) != int(r.Len) {
		return fmt.Errorf("string length needs to be exactly %d", r.Len)
	}
	return nil
}

type StringRegexRule struct {
	Regex string
	// Description string
}

func (r StringRegexRule) Validate(input string) error {
	v, err := regexp.MatchString(r.Regex, input)
	if err != nil {
		log.Printf("Error matching regex: %v", err)
	}
	if !v {
		return fmt.Errorf("string did not match regex `%s`", r.Regex)
	}
	return nil
}

type StringOneOf struct {
	List []string
}

func (r StringOneOf) Validate(input string) error {
	for _, v := range r.List {
		if input == v {
			return nil
		}
	}
	return fmt.Errorf("string was not one of %v", r.List)
}

type StringNotOneOf struct {
	List []string
}

func (r StringNotOneOf) Validate(input string) error {
	for _, v := range r.List {
		if input == v {
			return fmt.Errorf("string was one of illegal %v", r.List)
		}
	}
	return nil
}
