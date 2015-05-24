package backend

import (
	"errors"
	"fmt"
)

type Backend interface {
	List(string) map[string]string
	Close(url string)
}

var (
	backends    map[string]Backend
	ErrNotFound = errors.New("Could not find the backend.")
)

func init() {
	backends = make(map[string]Backend)
}

func Register(name string, backend Backend) error {
	if _, exists := backends[name]; exists {
		return fmt.Errorf("Scheme already registered %s", name)
	}
	backends[name] = backend

	return nil
}

func New(name string) (Backend, error) {
	if p, exists := backends[name]; exists {
		return p, nil
	}

	return nil, ErrNotFound
}

func List() []string {
	keys := make([]string, 0, len(backends))
	for k := range backends {
		keys = append(keys, k)
	}
	return keys
}
