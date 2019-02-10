package query

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/bpicolo/radiant/pkg/schema"
	"github.com/ghodss/yaml"
)

type YamlInterpreter struct {
}

func (i *YamlInterpreter) Interpret(search *schema.Search) (*Query, error) {
	tmpl := template.
		New(search.Query.Name).
		Funcs(sprig.TxtFuncMap())

	tmpl, err := tmpl.Parse(search.Query.Source)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer

	err = tmpl.Execute(&buf, search.Context)
	if err != nil {
		return nil, err
	}

	jsonQuery, err := yaml.YAMLToJSON(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("Failed to convert search query to JSON: %s", err)
	}

	return &Query{
		Search:  search,
		ESQuery: string(jsonQuery),
	}, nil
}
