# Radiant

Radiant is a service abstraction that enables you to facilitate communication with many potential
Elasticsearch servers. In addition to providing transparent proxies to those clusters, it allows you
to create simple, shareable query endpoints to simplify querying from a variety of sources.

## Installation

`go get github.com/bpicolo/radiant`
`go install github.com/bpicolo/radiant`

## Getting Started

To start up, you'll have to define the available Elasticsearch clusters.

```yaml
backends:
  - name: main
    host: http://my-es-host:9200
  - name: users
    host: http://secondary-es-host:9200
```

For a multi-host ES cluster, you may provide a comma-separated list of hosts, though the reverse proxy
functionality currently targets only the first host for a given backend.

With backends set up, you can run Radiant and query your hosts.

```bash
radiant serve radiant.yaml &

curl -H"Radiant-Proxy-Backend: main" localhost:5000
```

## Defining a search API

Radiant allows you to define new search endpoints with simple YAML configurations. The default search directory for search definitions
is `./searches`.

Here's an example search, against a cluster pre-loaded with [Kibana's sample Shakespeare data](https://www.elastic.co/guide/en/kibana/current/tutorial-load-dataset.html):

```yaml
backend: main
index: shakespeare
name: shakespeare/LinesBySpeaker
source: |
  bool:
    filter:
      term:
        speaker: "{{ .speaker }}"
  {{ if .query}}
    must:
      terms:
        text_entry: {{ .query | toJson}}
  {{ end }}
```

Place this anywhere you desire in your `./searches` directory, e.g.

```
searches/
    shakespeare/
        LinesBySpeaker.yaml
```

With this defined, you can start Radiant and use your defined query:

```bash
radiant serve radiant.yaml

curl -H"Content-Type: application/json" "localhost:5000/search/shakespeare/LinesBySpeaker?from=0&size=1" -d '{"speaker": "DEMETRIUS", "query": ["thee", "no"]}' | jq .
{
  "took": 18,
  "hits": {
    "total": 20,
    "max_score": 1,
    "hits": [
      {
        "_score": 1,
        "_index": "shakespeare",
        "_type": "doc",
        "_id": "67569",
        "_source": {
          "type": "line",
          "line_id": 67570,
          "play_name": "A Midsummer nights dream",
          "speech_number": 24,
          "line_number": "2.1.198",
          "speaker": "DEMETRIUS",
          "text_entry": "Hence, get thee gone, and follow me no more."
        }
      }
    ]
  },
  "_shards": {
    "total": 5,
    "successful": 5,
    "failed": 0
  }
}
```

A command is also provided to try out your query templates while developing.

```bash
radiant query-check searches/shakespeare/speaker.yaml '{"query": ["bob"], "speaker": "Dave"}' | jq .
{
  "bool": {
    "filter": {
      "term": {
        "speaker": "Dave"
      }
    },
    "must": {
      "terms": {
        "text_entry": [
          "bob"
        ]
      }
    }
  }
}
```