package main

import (
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"

	"github.com/stoffand/go/validator/generator"
)

func main() {
	// Initialize
	fs := token.NewFileSet()

	// Get file
	arg := os.Args[1]
	f, err := parser.ParseFile(fs, arg, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("could not parse %s: %v", arg, err)
	}
	// Traverse fileData and create template
	fileData := generator.GetFileData(f)
	templateData, err := fileData.CreateTemplate(arg)
	if err != nil {
		panic(err)
	}
	// Get filename (without path)
	a := strings.Split(arg, "/")
	fileName := a[len(a)-1]
	// Write to file
	err = os.WriteFile("vgen_"+fileName, templateData, 0644)
	if err != nil {
		panic(err)
	}
}

// fmt.Printf("data: %+v\n", data)
// validation.Debug("data", data)
