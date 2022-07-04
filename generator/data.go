package generator

import (
	"fmt"
	"go/ast"
	"regexp"
	"strings"
)

// TODO move data, field and dataType into individual files

// ImportPath string?
type Data struct {
	Pkg     string
	Types   []Type
	Aliases []Alias
	Imports []Import
}

type Import struct {
	Alias string
	Path  string
	Used  bool
}

func (d *Data) UseImport(imp string) {
	for i := 0; i < len(d.Imports); i++ {
		if d.Imports[i].Alias == imp {
			d.Imports[i].Used = true
		}
	}
}

type tags struct {
	skip bool
}

func parseTags(in string) tags {
	tags := tags{}

	r := regexp.MustCompile(`vgen:"(.*)"`)
	match := r.FindStringSubmatch(in)
	if len(match) > 0 {
		args := match[1]
		args = strings.ReplaceAll(args, " ", "")
		split := strings.Split(args, ",")

		for _, tag := range split {
			if tag == "skip" {
				tags.skip = true
			}
		}
	}
	return tags
}

func (d *Data) addAlias(a Alias) {
	d.Aliases = append(d.Aliases, a)
}

type Alias struct {
	Name string
	Type string
}

func (d *Data) addType(t Type) {
	d.Types = append(d.Types, t)
}

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
	if args.tags.skip {
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

	switch t.Type.(type) {
	case CustomType:
		return ident + "in." + t.Name + ".Convert()"
	case ImportedType:
		return ident + "in." + t.Name + ".Convert()"
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

// Rule functions
func (t PrimitiveType) rule() string {
	return "validation.Rule[" + t.String() + "]"
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
	return "validation.ListRule[" + typWithoutBrackets + "]"
}
func (t MapType) rule() string {
	valWithoutBrackets := t.ValueType.String()
	switch t.baseType().(type) {
	case CustomType:
		valWithoutBrackets += "In"
	}
	return "validation.MapRule[" + t.KeyType.String() + "," + valWithoutBrackets + "]"
}

// Base type
func (t PrimitiveType) baseType() FieldType { return t }
func (t CustomType) baseType() FieldType    { return t }
func (t ImportedType) baseType() FieldType  { return t.Type }
func (t ArrayType) baseType() FieldType     { return t.Type.baseType() }
func (t MapType) baseType() FieldType       { return t.ValueType.baseType() }
