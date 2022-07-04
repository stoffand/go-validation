package validation

import (
	"encoding/json"
	"fmt"
	"log"
)

func P[T any](in T) *T {
	return &in
}

func Debug(where string, a any) {
	j, err := json.MarshalIndent(a, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("DEBUG at %v:\n%v\n", where, string(j))
}
