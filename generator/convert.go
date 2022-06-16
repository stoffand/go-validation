package generator

import (
	"bytes"
	"fmt"
	"text/template"
)

// TODO find better way to know if you came from arr or map (using boolean right now)

var (
	arrMain = `  var res1 {{ .FullType }} 
		for _, v1 := range *in.{{ .Name }} {
			{{ .Inner }}
		}
		res.{{ .Name }} = {{ if .Pointer }}&{{ end }}res1 `
	arrOuterArr = ` var res{{ .Depth2 }} {{ .FullType }} 
		for _, v{{ .Depth2 }} := range v{{ .Depth1 }} { 
			{{ .Inner }}
		} 
		res{{ .Depth1 }} = append(res{{ .Depth1 }}, res{{ .Depth2 }})`
	arrOuterMap = ` var res{{ .Depth2 }} {{ .FullType }} 
		for _, v{{ .Depth2 }} := range v{{ .Depth1 }} { 
			{{ .Inner }}
		} 
		res{{ .Depth1 }}[k{{ .Depth1 }}] = res{{ .Depth2 }}`
	mapMain = ` res1 := make({{ .FullType }})
 		for k1, v1 := range *in.{{ .Name }} {
			 {{ .Inner }}
 		}
 		res.{{ .Name }} = {{ if .Pointer }}&{{ end }}res1 `
	mapOuterMap = ` res{{ .Depth2 }} := make({{ .FullType }})
 		for k{{ .Depth2 }}, v{{ .Depth2 }} := range v{{ .Depth1 }} {
			 {{ .Inner }}
 		}
		res{{ .Depth1 }}[k{{ .Depth1 }}] = res{{ .Depth2 }}`
	mapOuterArr = ` res{{ .Depth2 }} := make({{ .FullType }})
 		for k{{ .Depth2 }}, v{{ .Depth2 }} := range v{{ .Depth1 }} {
			 {{ .Inner }}
 		}
		res{{ .Depth1 }} = append(res{{ .Depth1 }}, res{{ .Depth2 }})`
)

func (t ArrayType) CustomConvert(name string, depth int, pointer bool, lastArr bool) string {
	// Get inner recursively
	var inner string
	switch typ := t.Type.(type) {
	case CustomType:
		inner = fmt.Sprintf("res%d = append(res%d, v%d.Convert())", depth+1, depth+1, depth+1)
	case ArrayType:
		inner = typ.CustomConvert(name, depth+1, pointer, true)
	case MapType:
		inner = typ.CustomConvert(name, depth+1, pointer, true)
	default:
		// return fmt.Sprintf(" %v ", depth+1) + typ.CustomConvert(name, depth+1)
		panic("custom convert was not one of custom, arr or map")
	}
	fmt.Printf("inner: %v\n", inner)

	// Template
	var tmpl *template.Template
	var err error
	if depth == 0 {
		tmpl, err = template.New("tmpl").Parse(arrMain)
	} else if lastArr {
		tmpl, err = template.New("tmpl").Parse(arrOuterArr)
	} else {
		tmpl, err = template.New("tmpl").Parse(arrOuterMap)
	}
	if err != nil {
		panic(err)
	}

	data := struct {
		FullType string
		Type     string
		Name     string
		Inner    string
		Depth1   int
		Depth2   int
		Pointer  bool
	}{
		FullType: t.String(),
		Type:     t.baseType().String(),
		Name:     name,
		Inner:    inner,
		Depth1:   depth,
		Depth2:   depth + 1,
		Pointer:  pointer,
	}

	// Execute
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, data)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func (t MapType) CustomConvert(name string, depth int, pointer bool, lastArr bool) string {
	// Get inner recursively
	var inner string
	switch typ := t.ValueType.(type) {
	case CustomType:
		inner = fmt.Sprintf("res%d[k%d] = v%d.Convert()", depth+1, depth+1, depth+1)
	case ArrayType:
		inner = typ.CustomConvert(name, depth+1, pointer, false)
	case MapType:
		inner = typ.CustomConvert(name, depth+1, pointer, false)
	default:
		// return fmt.Sprintf(" %v ", depth+1) + typ.CustomConvert(name, depth+1)
		panic("custom convert was not one of custom, arr or map")
	}
	fmt.Printf("inner: %v\n", inner)

	// Template
	var tmpl *template.Template
	var err error
	if depth == 0 {
		tmpl, err = template.New("tmpl").Parse(mapMain)
	} else if !lastArr {
		tmpl, err = template.New("tmpl").Parse(mapOuterMap)
	} else {
		tmpl, err = template.New("tmpl").Parse(mapOuterArr)
	}
	if err != nil {
		panic(err)
	}

	data := struct {
		FullType string
		Type     string
		Name     string
		Inner    string
		Depth1   int
		Depth2   int
		Pointer  bool
	}{
		FullType: t.String(),
		Type:     t.baseType().String(),
		Name:     name,
		Inner:    inner,
		Depth1:   depth,
		Depth2:   depth + 1,
		Pointer:  pointer,
	}

	// Execute
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, data)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

// func (t MapType) CustomConvert(name string, depth int) string {
// 	depth++
// 	switch typ := t.ValueType.(type) {
// 	case CustomType:
// 		return typ.String()
// 	case ArrayType:
// 		return typ.CustomConvert(name, depth, false)
// 	case MapType:
// 		return typ.CustomConvert(name, depth)
// 	}
// 	panic("custom convert was not one of custom, arr or map")
// }
