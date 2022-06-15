package generator

import (
	"fmt"
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
		newType := Type{
			Name: d.Name.Name,
		}
		switch t := d.Type.(type) {
		case *ast.StructType:
			fmt.Printf("Struct: %v\n", d.Name)
			for _, f := range t.Fields.List {
				newType.AddField(CreateField(f))
				fmt.Printf("\t%v: %#v\n", f.Names[0], f.Type)
			}
		case *ast.Ident:
			panic("type aliases not implemented")
		}
		v.Data.AddType(newType)
		fmt.Printf("\n-------------\n\n")
	}
	return v
}
