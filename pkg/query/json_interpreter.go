package query

import (
	"bytes"
	"html/template"

	"github.com/bpicolo/radiant/pkg/schema"
)

type JsonInterpreter struct {
}

func (i *JsonInterpreter) Interpret(search schema.Search) (*Query, error) {
	tmpl, err := template.New(search.Query.Name).Parse(search.Query.Source)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer

	err = tmpl.Execute(&buf, search.Context)
	if err != nil {
		return nil, err
	}

	return &Query{
		Search:  search,
		ESQuery: buf.String(),
	}, nil
}
