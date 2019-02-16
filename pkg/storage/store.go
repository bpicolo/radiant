package storage

import "github.com/bpicolo/radiant/pkg/schema"

type Store interface {
	Backends() []*schema.Backend
	GetAlias(name string) (*schema.Alias, error)
	GetQueryDefinition(name string) (*schema.QueryDefinition, error)
	Writeable() bool
}
