package generator

import (
	"fmt"
	"go/ast"
)

// ImportPath string?
type Data struct {
	Pkg   string
	Types []Type
}

func (d *Data) AddType(t Type) {
	d.Types = append(d.Types, t)
}

type Type struct {
	Name   string
	Fields []Field
}

func (t *Type) AddField(f Field) {
	t.Fields = append(t.Fields, f)
}

type Field struct {
	Pointer bool
	Name    string
	Type    FieldType
}

func CreateField(tName string, f *ast.Field) Field {
	fName := f.Names[0].String()
	field := Field{Name: fName}
	switch typ := f.Type.(type) {
	case *ast.StarExpr:
		field.Pointer = true
		field.Type = CreateFieldType(tName, fName, typ.X)
	default:
		field.Type = CreateFieldType(tName, fName, typ)
	}
	return field
}
func CreateFields(tName string, f *ast.Field) []Field {
	var fields []Field
	for _, v := range f.Names {
		field := Field{Name: v.Name}
		switch typ := f.Type.(type) {
		case *ast.StarExpr:
			field.Pointer = true
			field.Type = CreateFieldType(tName, v.Name, typ.X)
		default:
			field.Type = CreateFieldType(tName, v.Name, typ)
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
				return typ.CustomConvert(f.Name, 0, true, true)
			case MapType:
				return typ.CustomConvert(f.Name, 0, true, false)
			}
		}
		return ident + Convert(f.Name, f.Type)
	} else {
		switch typ := f.Type.(type) {
		case CustomType:
			return ident + typ.Convert(f.Name)
		case ArrayType:
			if hasCustomBaseType {
				return typ.CustomConvert(f.Name, 0, false, true)
			}
		case MapType:
			if hasCustomBaseType {
				return typ.CustomConvert(f.Name, 0, false, false)
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
	default: // primitive
		return "*" + f.Type.String()
	}
}

type FieldType interface {
	String() string
	rule() string
	baseType() FieldType
	isFieldType()
}

func CreateFieldType(tName, fName string, e ast.Expr) FieldType {
	switch x := e.(type) {
	case *ast.Ident:
		if x.Obj == nil { // Primitive
			return PrimitiveType{Value: x.Name}
		} else { // Custom
			return CustomType{Value: x.Name}
		}
	case *ast.ArrayType:
		return ArrayType{Type: CreateFieldType(tName, fName, x.Elt)}
	case *ast.MapType:
		return MapType{KeyType: CreateFieldType(tName, fName, x.Key), ValueType: CreateFieldType(tName, fName, x.Value)}
	case *ast.StarExpr:
		panic(fmt.Sprintf("struct %s, field %s: more than one pointer are not allowed", tName, fName))
	default:
		panic(fmt.Sprintf("struct %s, field %s: unsupported type", tName, fName))
	}
}

// Data types

type PrimitiveType struct {
	Value string
}

type CustomType struct {
	Value string
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
func (t ArrayType) isFieldType()     {}
func (t MapType) isFieldType()       {}

// String functions
func (t PrimitiveType) String() string { return t.Value }
func (t CustomType) String() string    { return t.Value }
func (t ArrayType) String() string     { return "[]" + t.Type.String() }
func (t MapType) String() string       { return "map[" + t.KeyType.String() + "]" + t.ValueType.String() }

// Convert
func Convert(name string, t FieldType) string {
	switch t.baseType().(type) {
	case CustomType:
		return "in." + t.String() + ".Convert()"
	}
	return "in." + name
}

func (t PrimitiveType) Convert(name string) string {
	return "in." + name
}
func (t CustomType) Convert(name string) string {
	return "in." + t.String() + ".Convert()"
}
func (t ArrayType) Convert(name string) string {
	return "in." + name
}
func (t MapType) Convert(name string) string {
	return "in." + name
}

// Rule functions
func (t PrimitiveType) rule() string {
	return "validation.Rules[" + t.Value + "]"
}
func (t CustomType) rule() string {
	return "validation.Rule[" + t.Value + "In]"
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
func (t ArrayType) baseType() FieldType     { return t.Type.baseType() }
func (t MapType) baseType() FieldType       { return t.ValueType.baseType() }
