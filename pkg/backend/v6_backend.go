package backend

import (
	"strings"

	"github.com/bpicolo/radiant/pkg/schema"
	"github.com/olivere/elastic"
)

// V6Backend an elasticsearch backend for elasticsearch 6.x
type V6Backend struct {
	client *elastic.Client
}

// NewV6Backend connect to an elasticsearch cluster currently using v6
func NewV6Backend(b schema.Backend) (*V6Backend, error) {
	hosts := strings.Split(b.Host, ",")

	var es *elastic.Client
	var err error
	if b.Simple {
		es, err = elastic.NewSimpleClient(
			elastic.SetURL(hosts...),
		)
	} else {
		es, err = elastic.NewClient(
			elastic.SetURL(hosts...),
		)
	}

	if err != nil {
		return nil, err
	}

	return &V6Backend{client: es}, nil
}

// Stop shuts down the client
func (b *V6Backend) Stop() {
	b.client.Stop()
}
