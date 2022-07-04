package generator

import (
	"go/ast"
	"regexp"
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
	Data *Data
}

func (v TypeVisitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	switch d := n.(type) {
	case *ast.File:
		v.Data.Pkg = d.Name.Name
		for _, i := range d.Imports {
			// Create import data
			imp := Import{Path: i.Path.Value}
			if i.Name != nil {
				imp.Alias = i.Name.Name
			} else {
				// Extract alias from path
				trimmed := strings.Trim(i.Path.Value, `"`)
				split := strings.Split(trimmed, "/") // TODO mac specific right now
				alias := split[len(split)-1]
				imp.Alias = alias
			}
			v.Data.Imports = append(v.Data.Imports, imp)
		}

	case *ast.TypeSpec:
		tName := d.Name.Name
		switch t := d.Type.(type) {
		case *ast.StructType:
			newType := Type{
				Name: tName,
			}
			for _, f := range t.Fields.List {
				args := createFieldArgs{typeName: tName, data: v.Data}
				if f.Tag != nil {
					args.parseTags(f.Tag.Value)
				}
				for _, n := range f.Names {
					args.fieldName = n.Name
					newType.addField(args, f.Type)
				}
			}
			v.Data.addType(newType)
		case *ast.Ident:
			v.Data.addAlias(Alias{Name: tName, Type: t.Name})
			// return v
			// panic("type aliases not supported")
		}
	}
	return v
}

func (i *createFieldArgs) parseTags(in string) {
	r := regexp.MustCompile(`vgen:"(.*)"`)
	match := r.FindStringSubmatch(in)
	if len(match) > 0 {
		args := match[1]
		args = strings.ReplaceAll(args, " ", "")
		split := strings.Split(args, ",")

		for _, tag := range split {
			if tag == "skip" {
				i.skip = true
			}
		}
	}
}
