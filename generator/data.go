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

func CreateField(f *ast.Field) Field {
	field := Field{Name: f.Names[0].String()}
	switch f.Type.(type) {
	case *ast.StarExpr:
		field.Pointer = true
		field.Type = CreateFieldType(f.Type)
	case *ast.Ident, *ast.ArrayType, *ast.MapType:
		field.Type = CreateFieldType(f.Type)
	default:
		panic("unsupported ast type")
	}
	return field
}

type Field struct {
	Pointer bool
	Name    string
	Type    FieldType
}

// Used to generate template
func (f Field) Rule() string {
	return f.Type.Rule()
}

// Used to generate template
func (f Field) Convert() string {
	ident := "res." + f.Name + " = "
	switch typ := f.Type.(type) {
	case CustomType:
		return ident + typ.Convert(f.Name)
	case PointerType:
		switch typ := typ.Type.(type) {
		case PointerType:
			panic("double pointers")
		case ArrayType:
			if f.HasCustomBaseType() {
				return typ.CustomConvert(f.Name, 0, true, true)
			}
		case MapType:
			if f.HasCustomBaseType() {
				return typ.CustomConvert(f.Name, 0, true, false)
			}
		}
		return ident + Convert(f.Name, f.Type)
	case ArrayType:
		if f.HasCustomBaseType() {
			return typ.CustomConvert(f.Name, 0, false, true)
		}
	case MapType:
		if f.HasCustomBaseType() {
			return typ.CustomConvert(f.Name, 0, false, false)
		}
	}
	return ident + "*" + Convert(f.Name, f.Type)
}

// 	case ArrayType:
// 		switch in.BaseType().(type) {
// 		case CustomType:
// 			return typ.CustomConvert(name)
// 			// return in.Convert(name) + "// custom array"
// 		}
// 	case MapType:
// 		switch in.BaseType().(type) {
// 		case CustomType:
// 			return in.Convert(name) + "// custom map"
// 		}
// 	}
// 	return in.Convert(name)
// }

// Used to generate template
func (f Field) BaseType() FieldType {
	return f.Type.BaseType()
}

func (f Field) HasCustomBaseType() bool {
	switch f.BaseType().(type) {
	case CustomType:
		return true
	}
	return false
}

// Used to generate template
func (f Field) In() string {
	switch f.Type.(type) { // , * or not
	case PointerType:
		return inHelper(f.Type)
	default:
		return "*" + inHelper(f.Type)
	}
}

func inHelper(in FieldType) string {
	switch in.BaseType().(type) { // In or not
	case CustomType:
		return in.String() + "In"
	default: // primitive
		return in.String()
	}
}

// Used to generate template
func (f Field) IsPointer() bool {
	switch f.Type.(type) {
	case PointerType:
		return true
	default:
		return false
	}
}

type FieldType interface {
	String() string
	Rule() string
	// Convert(string) string
	BaseType() FieldType
	isFieldType()
}

func CreateFieldType(e ast.Expr) FieldType {
	switch x := e.(type) {
	case *ast.Ident:
		if x.Obj == nil { // Primitive
			return PrimitiveType{Value: x.Name}
		} else { // Custom
			return CustomType{Value: x.Name}
		}
	case *ast.ArrayType:
		return ArrayType{Type: CreateFieldType(x.Elt)}
	case *ast.MapType:
		return MapType{KeyType: CreateFieldType(x.Key), ValueType: CreateFieldType(x.Value)}
	case *ast.StarExpr:
		return PointerType{Type: CreateFieldType(x.X)}
	default:
		panic(fmt.Sprintf("unsupported type: %#v", x))
	}
}

// Data types

type PrimitiveType struct {
	Value string
}

type CustomType struct {
	Value string
}

type PointerType struct {
	Type FieldType
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
func (t PointerType) isFieldType()   {}
func (t ArrayType) isFieldType()     {}
func (t MapType) isFieldType()       {}

// String functions
func (t PrimitiveType) String() string { return t.Value }
func (t CustomType) String() string    { return t.Value }
func (t PointerType) String() string   { return "*" + t.Type.String() }
func (t ArrayType) String() string     { return "[]" + t.Type.String() }
func (t MapType) String() string       { return "map[" + t.KeyType.String() + "]" + t.ValueType.String() }

// Convert
func Convert(name string, t FieldType) string {
	switch t.BaseType().(type) {
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
func (t PointerType) Convert(name string) string {
	return "in." + name
}
func (t ArrayType) Convert(name string) string {
	return "in." + name
}
func (t MapType) Convert(name string) string {
	return "in." + name
}

// Rule functions
func (t PrimitiveType) Rule() string {
	return "validation.Rules[" + t.Value + "]"
}
func (t CustomType) Rule() string {
	return "validation.Rule[" + t.Value + "In]"
}
func (t PointerType) Rule() string {
	return t.Type.Rule()
}
func (t ArrayType) Rule() string {
	typWithoutBrackets := t.Type.String()
	switch t.BaseType().(type) {
	case CustomType:
		typWithoutBrackets += "In"
	}
	return "validation.ListRules[" + typWithoutBrackets + "]"
}
func (t MapType) Rule() string {
	valWithoutBrackets := t.ValueType.String()
	switch t.BaseType().(type) {
	case CustomType:
		valWithoutBrackets += "In"
	}
	return "validation.MapRules[" + t.KeyType.String() + "," + valWithoutBrackets + "]"
}

// Base type
func (t PrimitiveType) BaseType() FieldType { return t }
func (t CustomType) BaseType() FieldType    { return t }
func (t PointerType) BaseType() FieldType   { return t.Type.BaseType() }
func (t ArrayType) BaseType() FieldType     { return t.Type.BaseType() }
func (t MapType) BaseType() FieldType       { return t.ValueType.BaseType() }
