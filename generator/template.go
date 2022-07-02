package generator

import (
	"bytes"
	"go/format"
	"text/template"
)

func (d Data) CreateTemplate() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Template
	tmpl := template.Must(template.New("template.tmpl").Parse(tmplStr))
	err := tmpl.Execute(buf, d)
	if err != nil {
		panic(err)
	}

	// Format
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	return formatted, nil
}

const (
	tmplStr = ` // Generated file DO NOT EDIT

	package {{ .Pkg }}
	
	import (
		"github.com/stoffand/go-validation/validation"
		{{- range .Imports }}
			{{ if .Used }} {{ .Alias }} {{ .Path }}	{{ end }}
		{{- end }}
	)
	
	{{ range .Types }} 
	
	// Input struct
	type {{ .Name }}In struct {
		{{- range .Fields }}
			{{ .Name }} {{ .In }}
		{{- end }}
	}
	
	// Convert input struct to original
	func (in {{ .Name }}In) Convert() {{ .Name }} {
		res := {{ .Name }}{}
		{{- range .Fields }}
			if in.{{ .Name }} != nil {
				{{ .Convert }}
			}
		{{- end }}
		return res
	}

	// Rules
	type {{ .Name }}Rules struct {
		{{- range .Fields }}
			{{ .Name }} {{ .Rule }}            
		{{- end }}
	}
	
	// Validate required and validation rules
	func (r {{ .Name }}Rules) Validate(in {{ .Name }}In) error {
		var errs validation.StructErr
		{{- range .Fields }}
			if in.{{ .Name }} != nil {
				if err := r.{{ .Name }}.Validate(*in.{{ .Name }}); err != nil {
					errs.AddError("{{ .Name }}", err)
				}
			}  {{ if not .Pointer }} else {
				errs.AddError("{{ .Name }}", validation.RequiredErr{})
			} {{ end }}
		{{- end }} 
		if len(errs.FailedFields) == 0 {
			return nil
		}
		return errs
	}
	
	func (r {{ .Name }}Rules) ValidatedConvert(in {{ .Name }}In) ({{ .Name }}, error) {
		err := r.Validate(in)
		if err != nil {
			return {{ .Name }}{}, err
		}
		return in.Convert(), nil
	}
	
	{{ end }} `
)
