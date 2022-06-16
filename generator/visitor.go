package generator

import (
	"go/ast"
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
			imp := Import{Path: i.Path.Value}
			if i.Name != nil {
				imp.Alias = i.Name.Name
			} else { // TODO mac specific right now
				trimmed := strings.Trim(i.Path.Value, `"`)
				split := strings.Split(trimmed, "/")
				alias := split[len(split)-1]
				imp.Alias = alias
			}
			v.Data.Imports = append(v.Data.Imports, imp)
		}
	case *ast.TypeSpec:
		tName := d.Name.Name
		newType := Type{
			Name: tName,
		}
		switch t := d.Type.(type) {
		case *ast.StructType:
			for _, f := range t.Fields.List {
				fields := createFields(tName, f, v.Data)
				for _, v := range fields {
					newType.addField(v)
				}
			}
		case *ast.Ident:
			panic("type aliases not implemented")
		}
		v.Data.addType(newType)
	}
	return v
}
