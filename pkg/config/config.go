package config

import "github.com/bpicolo/radiant/pkg/schema"

type RadiantConfig struct {
	Backends []*schema.Backend `yaml:"backends"`
	Queries  []*schema.Query   `yaml:"queries"`
	Aliases  []*schema.Alias   `yaml:"aliases"`
}
