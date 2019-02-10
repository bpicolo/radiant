package storage

import (
	"fmt"
	"io/ioutil"

	"github.com/bpicolo/radiant/pkg/schema"
	yaml "gopkg.in/yaml.v2"
)

type config struct {
	Backends []*schema.Backend `yaml:"backends"`
	Queries  []*schema.Query   `yaml:"queries"`
	Aliases  []*schema.Alias   `yaml:"aliases"`
}

// Yaml is a static
type Yaml struct {
	config config
}

// NewYaml create a new yaml storage backend
func NewYaml(cfgPath string) (*Yaml, error) {
	dat, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	c := config{}
	err = yaml.Unmarshal(dat, c)
	if err != nil {
		return nil, err
	}

	return &Yaml{config: c}, nil
}

func (s *Yaml) GetAlias(name string) (*schema.Alias, error) {
	// TODO make this O(1). Not an actual perf issue until you have silly configs
	for _, alias := range s.config.Aliases {
		if alias.Name == name {
			return alias, nil
		}
	}
	return nil, fmt.Errorf("Alias `%s` not found", name)
}

func (s *Yaml) GetQuery(name string) (*schema.Query, error) {
	// TODO make this O(1). Not an actual perf issue until you have silly configs
	for _, query := range s.config.Queries {
		if query.Name == name {
			return query, nil
		}
	}
	return nil, fmt.Errorf("Query `%s` not found", name)
}

func (s *Yaml) Backends() []*schema.Backend {
	return s.config.Backends
}

// Writeable - yaml storage is not currently editable, it is statically-defined only
func (s *Yaml) Writeable() bool {
	return false
}
