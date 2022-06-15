package generator

import (
	"bytes"
	"go/format"
	"text/template"
)

func (d Data) CreateTemplate(filename string) ([]byte, error) {
	buf := new(bytes.Buffer)

	// Template
	tmpl := template.Must(template.New("template.tmpl").ParseFiles("template.tmpl"))
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
