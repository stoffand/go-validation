package generator

import (
	"go/ast"
	"path/filepath"
	"strings"
)

func GetFileData(file *ast.File) Data {
	data := Data{}
	visitor := TypeVisitor{Data: &data}
	ast.Walk(visitor, file)
	return data
}

// Use name visitor?
type TypeVisitor struct {
	Data     *Data
	SkipNext bool
}

// Visit searches for types and sends data to generate files in data.go
func (v TypeVisitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	switch d := n.(type) {
	case *ast.File:
		v.Data.Pkg = d.Name.Name
		for _, i := range d.Imports {
			imp := Import{Path: i.Path.Value} // Create import data
			if i.Name != nil {
				// If alias exists
				imp.Alias = i.Name.Name
			} else {
				// Extract alias from path
				trimmed := strings.Trim(i.Path.Value, `"`)
				imp.Alias = filepath.Base(trimmed)
			}
			v.Data.Imports = append(v.Data.Imports, imp)
		}
	case *ast.GenDecl:
		if d.Doc != nil {
			last := d.Doc.List[len(d.Doc.List)-1]
			tags := parseTags(last.Text)
			if tags.skip {
				v.SkipNext = true
			}
		}
	case *ast.TypeSpec:
		// Skip type
		if v.SkipNext {
			v.SkipNext = false
			return v
		}

		tName := d.Name.Name
		switch t := d.Type.(type) {
		// Create full type for structs
		case *ast.StructType:
			newType := Type{
				Name: tName,
			}
			for _, f := range t.Fields.List {

				args := createFieldArgs{typeName: tName, data: v.Data}
				if f.Tag != nil {
					args.tags = parseTags(f.Tag.Value)
				}
				for _, n := range f.Names {
					args.fieldName = n.Name
					newType.addField(args, f.Type)
				}
			}
			v.Data.addType(newType)
		// Create simpler type for aliases
		case *ast.Ident:
			v.Data.addAlias(Alias{Name: tName, Type: t.Name})
		}
	}
	return v
}
