package generator

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
