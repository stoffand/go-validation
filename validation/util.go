package validation

import (
	"encoding/json"
	"fmt"
	"log"
)

func P[T any](in T) *T {
	return &in
}

func MapToError(in map[string]string) error {
	if in == nil {
		return nil
	}
	fmt.Println(in)
	j, err := json.Marshal(in)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Errorf(string(j))
}

func Debug(where string, a any) {
	j, err := json.MarshalIndent(a, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("DEBUG at %v:\n%v\n", where, string(j))
}
