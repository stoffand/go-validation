package generator

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

// how to handle errors (with multiple files)
// TODO add flag with which types to parse
func Generate(paths []string, verbose bool) {
	// Initialize
	fs := token.NewFileSet()

	var warnings []string
	var successful []string
	// Generate for each file
	for _, path := range paths {
		parseFileOrDir(fs, path, &successful, &warnings)
	}

	fmt.Printf("%d sucessfully parsed files\n", len(successful))
	if verbose {
		for _, v := range successful {
			fmt.Println(v)
		}
		fmt.Println()
	}
	fmt.Printf("%d warnigns parsed files\n", len(warnings))
	if verbose {
		for _, v := range warnings {
			fmt.Println(v)
		}
	}
}

// Walk through file tree and try to parse each file
func parseFileOrDir(fs *token.FileSet, root string, s, w *[]string) {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		parseFile(fs, path, s, w)
		return nil
	})
	if err != nil {
		panic(fmt.Sprintf("could not walk %s: %v", root, err))
	}
}

const PREFIX string = "vgen"

func parseFile(fs *token.FileSet, filePath string, s, w *[]string) {
	// Skip generated files
	file := filepath.Base(filePath)
	if file[:len(PREFIX)] == PREFIX {
		return
	}

	// Parse file
	f, err := parser.ParseFile(fs, filePath, nil, parser.ParseComments)
	if err != nil {
		*w = append(*w, fmt.Sprintf("%s: could not parse file", filePath))
		return
	}

	// Traverse fileData
	fileData := GetFileData(f)
	if len(fileData.Types) == 0 {
		*w = append(*w, fmt.Sprintf("%s: did not contain any parseable types", filePath))
		return
	}

	// Create template
	templateData, err := fileData.CreateTemplate()
	if err != nil {
		*w = append(*w, fmt.Sprintf("%s: could not create template", filePath))
		return
	}

	// Create new file path wih vgen prefix
	base := filepath.Base(filePath)
	dir := filepath.Dir(filePath)
	path := filepath.Join(dir, PREFIX+"_"+base)

	// if readonly file already exists, set writeable
	if _, err := os.Stat(path); err == nil {
		err := os.Chmod(path, 0644)
		if err != nil {
			*w = append(*w, fmt.Sprintf("%s: could not make file writeable: %v", filePath, err))
			return
		}
	}

	// Write to file
	err = os.WriteFile(path, templateData, 0444) // normal 0644, read-only 0444
	if err != nil {
		*w = append(*w, fmt.Sprintf("%s: could write to files: %v", filePath, err))
		return
	}

	// Add successful files
	*s = append(*s, filePath)
}
