package query

import "github.com/bpicolo/radiant/pkg/schema"

type Query struct {
	ESQuery string
	Search  *schema.Search
}
