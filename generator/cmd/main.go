package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"

	"github.com/stoffand/go-validator/generator"
)

// how to handle errors (with multiple files)
// TODO add flag with which types to parse
func main() {
	// Input files
	flag.Parse()
	fileNames := flag.Args()
	if len(fileNames) == 0 {
		log.Fatal("no files specified")
	}

	// Initialize
	fs := token.NewFileSet()

	// Generate for each file
	for _, fName := range fileNames {
		// Parse file
		f, err := parser.ParseFile(fs, fName, nil, parser.ParseComments)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: could not parse file\n", fName)
			continue
		}

		// Traverse fileData
		fileData := generator.GetFileData(f)
		if len(fileData.Types) == 0 {
			fmt.Fprintf(os.Stderr, "%s: did not contain any parseable types\n", fName)
			continue
		}

		// Create template
		templateData, err := fileData.CreateTemplate()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: could not create template\n", fName)
			continue
		}

		// Get filename (without path)
		a := strings.Split(fName, "/")
		fileName := a[len(a)-1]

		// Write to file
		err = os.WriteFile("vgen_"+fileName, templateData, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: could write to files\n", fName)
			continue
		}
	}
}
