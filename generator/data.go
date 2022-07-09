package generator

import (
	"fmt"
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

type Alias struct {
	Name string
	Type string
}

type tags struct {
	include bool
}

func (d *Data) UseImport(imp string) {
	for i := 0; i < len(d.Imports); i++ {
		if d.Imports[i].Alias == imp {
			d.Imports[i].Used = true
		}
	}
}

func (d *Data) addAlias(a Alias) {
	d.Aliases = append(d.Aliases, a)
}

func (d *Data) addType(t Type) {
	d.Types = append(d.Types, t)
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
			switch tag {
			case "include", "i":
				tags.include = true
			case "":
				break
			default:
				panic(fmt.Sprintf("unknown tag: %s", tag)) // TODO add context to where this is
			}
		}
	}
	return tags
}
