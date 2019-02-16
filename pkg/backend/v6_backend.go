package backend

import (
	"context"
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/bpicolo/radiant/pkg/query"
	"github.com/bpicolo/radiant/pkg/schema"
	"github.com/olivere/elastic"
)

// V6Backend an elasticsearch backend for elasticsearch 6.x
type V6Backend struct {
	*backend
	client *elastic.Client
}

// NewV6Backend connect to an elasticsearch cluster currently using v6
func NewV6Backend(b *schema.Backend) (*V6Backend, error) {
	hosts := strings.Split(b.Host, ",")
	for i := range hosts {
		parsed, err := url.Parse(hosts[i])

		if err != nil {
			return nil, fmt.Errorf("Error parsing host for backend: %s", err)
		}
		if !parsed.IsAbs() {
			return nil, fmt.Errorf("Backend hosts must include protocol, got %s", hosts[i])
		}

		parsed.Path = ""
		parsed.RawQuery = ""
		parsed.Fragment = ""
		hosts[i] = parsed.String()
	}

	var es *elastic.Client
	var err error
	if b.EnableSniffing {
		es, err = elastic.NewClient(
			elastic.SetURL(hosts...),
		)
	} else {
		es, err = elastic.NewSimpleClient(
			elastic.SetURL(hosts...),
		)

	}

	if err != nil {
		return nil, err
	}

	proxyHost, _ := url.Parse(hosts[0])

	return &V6Backend{
		client: es,
		// TODO support a multi-master host reverse proxy
		backend: &backend{
			proxy: httputil.NewSingleHostReverseProxy(proxyHost),
		},
	}, nil
}

// Search performs an es query against a es 6 backend
func (b *V6Backend) Search(q *query.Query) (*elastic.SearchResult, error) {
	ctx := context.Background()

	log.Println("Performing search")
	log.Println(q.ESQuery)
	return b.client.Search().Query(elastic.NewRawStringQuery(q.ESQuery)).
		Index(q.Search.QueryDefinition.Index).
		From(q.Search.From).
		Size(q.Search.Size).Do(ctx)
}

// Stop shuts down the client
func (b *V6Backend) Stop() {
	b.client.Stop()
}

func (b *V6Backend) Backend() *backend {
	return b.backend
}
