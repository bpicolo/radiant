package backend

import (
	"github.com/bpicolo/radiant/pkg/schema"
)

type backendMeta struct {
	backend schema.Backend
}

// Manager manages a set of elasticsearch backends
type Manager struct {
	backends map[string]Backend
}

// NewManager initializes a default backend manager
func NewManager() *Manager {
	return &Manager{
		backends: make(map[string]Backend),
	}
}

// AddBackend adds a supported backend to the backend manager
func (m *Manager) AddBackend(backend schema.Backend) error {
	b, err := discover(backend)
	if err != nil {
		return err
	}

	m.backends[backend.Name] = b
	return nil
}

func (m *Manager) Stop() {
	for _, b := range m.backends {
		b.Stop()
	}
}
