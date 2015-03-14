// package membk implements in-memory backend, that is used in tests to mock database
package membk

import (
	"github.com/gravitational/hello/backend"
)

type MemBackend struct {
	Greetings map[string]string
}

func New() *MemBackend {
	return &MemBackend{
		Greetings: make(map[string]string),
	}
}

// UpsertGreeting updates or inserts the greeting into the database
func (b *MemBackend) UpsertGreeting(id, val string) error {
	b.Greetings[id] = val
	return nil
}

// GetGreeting returns a greeting stored in a database by it's id
func (b *MemBackend) GetGreeting(id string) (string, error) {
	g, ok := b.Greetings[id]
	if !ok {
		return "", &backend.NotFoundError{ID: id}
	}
	return g, nil
}

// DeleteGreeting deletes greeting by ID
func (b *MemBackend) DeleteGreeting(id string) error {
	_, ok := b.Greetings[id]
	if !ok {
		return &backend.NotFoundError{}
	}
	delete(b.Greetings, id)
	return nil
}

// Close closes all resources associated with this backend
func (b *MemBackend) Close() error {
	return nil
}
