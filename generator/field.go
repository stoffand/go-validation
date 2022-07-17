package generator

import (
	"fmt"
	"go/ast"
)

type Type struct {
	Name   string
	Fields []Field
}

func (t *Type) addField(args createFieldArgs, expr ast.Expr) {
	var err error
	newField := Field{Name: args.fieldName}

	// If not required
	if e, ok := expr.(*ast.StarExpr); ok {
		newField.Pointer = true
		expr = e.X
	}

	// Convert to fieldType
	newField.Type, err = createFieldType(args.data, expr)
	if err != nil {
		panic(fmt.Sprintf("type %v: field %v: %v", args.typeName, args.fieldName, err))
	}

	// Skip tag, custom type acts as primitive
	if !args.tags.include {
		switch newField.BaseType().(type) {
		case CustomType, ImportedType:
			newField.Type = PrimitiveType{Value: newField.Type.String()}
		}
	}

	t.Fields = append(t.Fields, newField)
}

type Field struct {
	Pointer bool
	Name    string
	Type    FieldType
}

type createFieldArgs struct {
	typeName  string
	fieldName string
	tags      tags
	data      *Data
}

func (f Field) String() string {
	if f.Pointer {
		return f.Name + ": *" + f.Type.String()
	}
	return f.Name + ": " + f.Type.String()
}

// Used to generate template
func (f Field) Rule() string {
	return f.Type.rule()
}

// Used to generate template
func (f Field) Convert() string {
	// Custom type
	if _, ok := f.BaseType().(CustomType); ok {
		switch typ := f.Type.(type) {
		case ArrayType:
			return typ.customConvert(f.Name, 0, f.Pointer, true)
		case MapType:
			return typ.customConvert(f.Name, 0, f.Pointer, false)
		}
	}
	return convertHelper(f)
}

func convertHelper(t Field) string {
	ident := "res." + t.Name + " = "
	customConvert := "in." + t.Name + ".Convert()"

	switch t.Type.(type) {
	case CustomType, ImportedType:
		if t.Pointer {
			return "tmp := " + customConvert + "\n" + ident + "&tmp"
		}
		return ident + customConvert
	}
	if t.Pointer {
		return ident + "in." + t.Name
	}
	return ident + "*in." + t.Name
}

// Used to generate template
func (f Field) BaseType() FieldType {
	return f.Type.baseType()
}

// Used to generate template
func (f Field) In() string {
	switch f.BaseType().(type) { // In or not
	case CustomType:
		return "*" + f.Type.String() + "In"
	case ImportedType:
		return "*" + f.Type.String() + "In"
	case PrimitiveType: // primitive
		return "*" + f.Type.String()
	}
	panic(fmt.Sprintf("unsupported basetype: %v", f))
}

type FieldType interface {
	String() string
	rule() string
	baseType() FieldType
	isFieldType()
}

// TODO clean up this shit
func createFieldType(data *Data, e ast.Expr) (FieldType, error) {
	switch x := e.(type) {

	case *ast.Ident:
		if isBaseType(x.Name) {
			return PrimitiveType{Value: x.Name}, nil
		}
		return CustomType{Value: x.Name}, nil

	case *ast.SelectorExpr: // imported type (always custom)
		imp, ok := x.X.(*ast.Ident)
		typ := x.Sel
		if !ok || typ == nil {
			return nil, fmt.Errorf("imported type syntax error")
		}

		// Add import
		data.UseImport(imp.Name)

		return ImportedType{
			Import: imp.Name,
			Type:   CustomType{Value: typ.Name},
		}, nil

	case *ast.ArrayType:
		t, err := createFieldType(data, x.Elt)
		if err != nil {
			return nil, err
		}
		return ArrayType{Type: t}, nil

	case *ast.MapType:
		k, err := createFieldType(data, x.Key)
		if err != nil {
			return nil, err
		}
		v, err := createFieldType(data, x.Value)
		if err != nil {
			return nil, err
		}
		return MapType{KeyType: k, ValueType: v}, nil

	case *ast.StarExpr:
		return nil, fmt.Errorf("more than one pointer is not allowed")
	default:
		return nil, fmt.Errorf("unsupported type")
	}
}

func isBaseType(in string) bool {
	switch in {
	case "int8", "int16", "int32", "int64",
		"uint8", "uint16", "uint32", "uint64",
		"int", "uint", "rune", "byte", "uintptr",
		"float32", "float64",
		"complex64", "complex128",
		"string", "bool":
		return true
	}
	return false
}
