package storage

import (
	"fmt"

	"github.com/bpicolo/radiant/pkg/config"
	"github.com/bpicolo/radiant/pkg/schema"
)

type Static struct {
	config *config.RadiantConfig
}

func NewStatic(cfg *config.RadiantConfig) (*Static, error) {
	return &Static{config: cfg}, nil
}

func (s *Static) GetAlias(name string) (*schema.Alias, error) {
	// TODO make this O(1). Not an actual perf issue until you have silly configs
	for _, alias := range s.config.Aliases {
		if alias.Name == name {
			return alias, nil
		}
	}
	return nil, fmt.Errorf("Alias `%s` not found", name)
}

func (s *Static) GetQueryDefinition(name string) (*schema.QueryDefinition, error) {
	// TODO make this O(1). Not an actual perf issue until you have silly configs
	for _, query := range s.config.Queries {
		if query.Name == name {
			return query, nil
		}
	}
	return nil, fmt.Errorf("Query `%s` not found", name)
}

func (s *Static) Backends() []*schema.Backend {
	return s.config.Backends
}

// Writeable - Static storage is not currently editable, it is statically-defined only
func (s *Static) Writeable() bool {
	return false
}
