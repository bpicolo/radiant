package config

import "github.com/bpicolo/radiant/pkg/schema"

type RadiantConfig struct {
	Backends []*schema.Backend         `yaml:"backends"`
	Queries  []*schema.QueryDefinition `yaml:"queries"`
	Aliases  []*schema.Alias           `yaml:"aliases"`
}
