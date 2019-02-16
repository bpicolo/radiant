package schema

import (
	"encoding/json"
	"fmt"

	"github.com/ghodss/yaml"
	yamlbase "gopkg.in/yaml.v2"

	"github.com/xeipuuv/gojsonschema"
)

// QueryDefinition represents an elasticsearch query. Supported types are the default,
// which is go-templated json, yaml, as well as lua-scripted
type QueryDefinition struct {
	Name         string      `yaml:"name"`
	Backend      string      `yaml:"backend"`
	Index        string      `yaml:"index"`
	Type         string      `yaml:"type"`
	QuerySource  string      `yaml:"source"`
	SchemaSource interface{} `yaml:"schema"`

	schema *gojsonschema.Schema
}

// UnmarshalYAML performs custom unmarshaling of query
func (q *QueryDefinition) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var inner struct {
		Name         string      `yaml:"name"`
		Backend      string      `yaml:"backend"`
		Index        string      `yaml:"index"`
		Type         string      `yaml:"type"`
		QuerySource  string      `yaml:"source"`
		SchemaSource interface{} `yaml:"schema"`
	}

	if err := unmarshal(&inner); err != nil {
		return err
	}

	if inner.SchemaSource != nil {
		sl := gojsonschema.NewSchemaLoader()

		// This is a bit hacky seeming. We have YAML, but want to convert
		// it into a JSON schema. go-yaml generates map[inferface{}]interface{}
		// because yaml can have non-string keys, unlike json.
		// The yaml-wrapping package helps do pieces of this
		schemaBytes, err := yamlbase.Marshal(inner.SchemaSource)
		if err != nil {
			return fmt.Errorf("Error parsing query schema: %s", err)
		}

		jsonBytes, err := yaml.YAMLToJSON(schemaBytes)
		if err != nil {
			return fmt.Errorf("Error parsing query schema: %s", err)
		}

		var schema map[string]interface{}
		err = json.Unmarshal(jsonBytes, &schema)
		if err != nil {
			return fmt.Errorf("Failed to parse schema as map[string]interface{}: %s", err)
		}

		loader := gojsonschema.NewGoLoader(schema)

		compiledSchema, err := sl.Compile(loader)
		if err != nil {
			return fmt.Errorf("Specified schema is not a valid jsonschema spec: %s", err)
		}
		q.schema = compiledSchema
	}

	q.Name = inner.Name
	q.Backend = inner.Backend
	q.Index = inner.Index
	q.Type = inner.Type
	q.QuerySource = inner.QuerySource
	q.SchemaSource = inner.SchemaSource

	return nil
}

// Schema returns the underlying schema for the search query
func (q *QueryDefinition) Schema() *gojsonschema.Schema {
	return q.schema
}
