package generator

import (
	"fmt"
	"go/ast"
)

// ImportPath string?
type Data struct {
	Pkg     string
	Types   []Type
	Imports []Import
}

func (d *Data) UseImport(imp string) {
	for i := 0; i < len(d.Imports); i++ {
		if d.Imports[i].Alias == imp {
			d.Imports[i].Used = true
		}
	}
}

type Import struct {
	Alias string
	Path  string
	Used  bool
}

func (d *Data) addType(t Type) {
	d.Types = append(d.Types, t)
}

type Type struct {
	Name   string
	Fields []Field
}

func (t *Type) addField(f Field) {
	t.Fields = append(t.Fields, f)
}

type Field struct {
	Pointer bool
	Name    string
	Type    FieldType
}

func createFields(tName string, f *ast.Field, data *Data) []Field {
	var fields []Field
	for _, v := range f.Names {
		field := Field{Name: v.Name}
		switch typ := f.Type.(type) {
		case *ast.StarExpr:
			field.Pointer = true
			field.Type = createFieldType(tName, v.Name, typ.X, data)
		default:
			field.Type = createFieldType(tName, v.Name, typ, data)
		}
		fields = append(fields, field)
	}
	return fields
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
	var hasCustomBaseType bool
	switch f.BaseType().(type) {
	case CustomType:
		hasCustomBaseType = true
	}
	ident := "res." + f.Name + " = "
	if f.Pointer {
		if hasCustomBaseType {
			switch typ := f.Type.(type) {
			case ArrayType:
				return typ.customConvert(f.Name, 0, true, true)
			case MapType:
				return typ.customConvert(f.Name, 0, true, false)
			}
		}
		return ident + Convert(f.Name, f.Type)
	} else {
		switch typ := f.Type.(type) {
		case CustomType, ImportedType:
			return ident + Convert(f.Name, f.Type)
		case ArrayType:
			if hasCustomBaseType {
				return typ.customConvert(f.Name, 0, false, true)
			}
		case MapType:
			if hasCustomBaseType {
				return typ.customConvert(f.Name, 0, false, false)
			}
		}
		return ident + "*" + Convert(f.Name, f.Type)
	}
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

func createFieldType(tName, fName string, e ast.Expr, data *Data) FieldType {
	switch x := e.(type) {
	case *ast.Ident:
		if isBaseType(x.Name) {
			return PrimitiveType{Value: x.Name}
		}
		return CustomType{Value: x.Name}
	case *ast.ArrayType:
		return ArrayType{Type: createFieldType(tName, fName, x.Elt, data)}
	case *ast.MapType:
		return MapType{KeyType: createFieldType(tName, fName, x.Key, data), ValueType: createFieldType(tName, fName, x.Value, data)}
	case *ast.StarExpr:
		panic(fmt.Sprintf("struct %s, field %s: more than one pointer are not allowed", tName, fName))
	case *ast.SelectorExpr: // imported type (always custom)
		imp, ok := x.X.(*ast.Ident)
		typ := x.Sel
		if !ok || typ == nil {
			panic("imported type syntax error")
		}
		data.UseImport(imp.Name)
		return ImportedType{
			Import: imp.Name,
			Type:   CustomType{Value: typ.Name},
		}
	default:
		panic(fmt.Sprintf("struct %s, field %s: unsupported type", tName, fName))
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

// Data types

type PrimitiveType struct {
	Value string
}

type CustomType struct {
	Value string
}

type ImportedType struct {
	Import string
	Type   CustomType
}

type ArrayType struct {
	Type FieldType
}

type MapType struct {
	KeyType   FieldType
	ValueType FieldType
}

// Interface functions
func (t PrimitiveType) isFieldType() {}
func (t CustomType) isFieldType()    {}
func (t ImportedType) isFieldType()  {}
func (t ArrayType) isFieldType()     {}
func (t MapType) isFieldType()       {}

// String functions
func (t PrimitiveType) String() string { return t.Value }
func (t CustomType) String() string    { return t.Value }
func (t ImportedType) String() string  { return t.Import + "." + t.Type.Value }
func (t ArrayType) String() string     { return "[]" + t.Type.String() }
func (t MapType) String() string       { return "map[" + t.KeyType.String() + "]" + t.ValueType.String() }

// Convert
func Convert(name string, t FieldType) string {
	switch x := t.(type) {
	case CustomType:
		return "in." + x.String() + ".Convert()"
	case ImportedType:
		return "in." + x.Type.String() + ".Convert()"
	}
	return "in." + name
}

// func (t PrimitiveType) Convert(name string) string {
// 	return "in." + name
// }
// func (t CustomType) Convert(name string) string {
// 	return "in." + t.String() + ".Convert()"
// }
// func (t ImportedType) Convert(name string) string {
// 	return t.Type.Convert(name)
// }
// func (t ArrayType) Convert(name string) string {
// 	return "in." + name
// }
// func (t MapType) Convert(name string) string {
// 	return "in." + name
// }

// Rule functions
func (t PrimitiveType) rule() string {
	return "validation.Rules[" + t.String() + "]"
}
func (t CustomType) rule() string {
	return "validation.Rule[" + t.String() + "In]"
}
func (t ImportedType) rule() string {
	return "validation.Rule[" + t.String() + "In]"
}
func (t ArrayType) rule() string {
	typWithoutBrackets := t.Type.String()
	switch t.baseType().(type) {
	case CustomType:
		typWithoutBrackets += "In"
	}
	return "validation.ListRules[" + typWithoutBrackets + "]"
}
func (t MapType) rule() string {
	valWithoutBrackets := t.ValueType.String()
	switch t.baseType().(type) {
	case CustomType:
		valWithoutBrackets += "In"
	}
	return "validation.MapRules[" + t.KeyType.String() + "," + valWithoutBrackets + "]"
}

// Base type
func (t PrimitiveType) baseType() FieldType { return t }
func (t CustomType) baseType() FieldType    { return t }
func (t ImportedType) baseType() FieldType  { return t.Type }
func (t ArrayType) baseType() FieldType     { return t.Type.baseType() }
func (t MapType) baseType() FieldType       { return t.ValueType.baseType() }
