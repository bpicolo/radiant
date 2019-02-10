package query

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"github.com/bpicolo/radiant/pkg/schema"
	"github.com/ghodss/yaml"
)

type YamlInterpreter struct {
}

func (i *YamlInterpreter) Interpret(search *schema.Search) (*Query, error) {
	tmpl, err := template.New(search.Query.Name).Parse(search.Query.Source)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer

	log.Printf("Source: %#v", search.Query.Source)
	log.Printf("Context: %#v", search.Context)
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
