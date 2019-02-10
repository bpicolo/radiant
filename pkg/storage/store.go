package storage

import "github.com/bpicolo/radiant/pkg/schema"

type Store interface {
	Backends() []*schema.Backend
	GetAlias(name string) (*schema.Alias, error)
	GetQuery(name string) (*schema.Query, error)
	Writeable() bool
}
