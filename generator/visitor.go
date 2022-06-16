package generator

import (
	"go/ast"
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

	case *ast.TypeSpec:
		tName := d.Name.Name
		newType := Type{
			Name: tName,
		}
		switch t := d.Type.(type) {
		case *ast.StructType:
			for _, f := range t.Fields.List {
				fields := CreateFields(tName, f)
				for _, v := range fields {
					newType.AddField(v)
				}
			}
		case *ast.Ident:
			panic("type aliases not implemented")
		}
		v.Data.AddType(newType)
	}
	return v
}
