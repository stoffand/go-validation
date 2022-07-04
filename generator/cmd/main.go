package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"

	"github.com/stoffand/go-validation/generator"
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

		// Create new file path wih vgen prefix
		base := filepath.Base(fName)
		dir := filepath.Dir(fName)
		path := filepath.Join(dir, "vgen_"+base)

		// if readonly file already exists, set writeable
		if _, err := os.Stat(path); err == nil {
			err := os.Chmod(path, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: could not make file writeable: %v\n", fName, err)
				continue
			}
		}

		// Write to file
		err = os.WriteFile(path, templateData, 0444) // normal 0644, read-only 0444
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: could write to files: %v\n", fName, err)
			continue
		}
	}
}
