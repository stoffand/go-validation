package test

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/stoffand/go-validator/test/different_package"
	nType "github.com/stoffand/go-validator/test/normal_type"
)

type Person struct {
	Name string
	Pet  Pet
	Home Home
	UID  uuid.UUID `vgen:"skip"`
	Car  different_package.Car
	Boat nType.Boat `vgen:"skip"`
}

type Pet struct {
	Name  string
	Breed string
}

func Test() {
	fmt.Printf("")
}
