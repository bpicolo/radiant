package query

import "github.com/bpicolo/radiant/pkg/schema"

type Interpreter interface {
	Interpret(query schema.Search) (*Query, error)
}
