package backend

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/bpicolo/radiant/pkg/schema"
	version "github.com/hashicorp/go-version"
)

type esVersionInfo struct {
	Number string `json:"number"`
}

type esInfo struct {
	Version esVersionInfo `json:"version"`
}

type backend struct {
	proxy *httputil.ReverseProxy
}

func (b *backend) Proxy() *httputil.ReverseProxy {
	return b.proxy
}

// Backend represents an elasticsearch backend
type Backend interface {
	Backend() *backend
	Stop()
}

func discover(b *schema.Backend) (Backend, error) {
	hosts := strings.Split(b.Host, ",")
	if len(hosts) < 1 {
		return nil, errors.New("Empty hosts: nothing to discover")
	}

	info, err := getInfo(hosts[0])
	if err != nil {
		return nil, err
	}

	v6Constraint, _ := version.NewConstraint(">= 6.0.0")

	version, err := version.NewVersion(info.Version.Number)
	if v6Constraint.Check(version) {
		return NewV6Backend(b)
	}

	return nil, errors.New("Unsupported elasticsearch version")
}

func getInfo(host string) (*esInfo, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	r, err := client.Get(host)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	info := esInfo{}
	err = json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}
