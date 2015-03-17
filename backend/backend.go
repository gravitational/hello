// Package backend represents an interface for accessing greetings storage
// and provides various implementations and testing examples.
package backend

import (
	"fmt"
)

// GreetingBackend is an interface to the backend (usually a database)
// that provides some storage functionality.
type GreetingBackend interface {

	// UpsertGreeting updates or inserts the greeting into the database
	UpsertGreeting(id, val string) error

	// GetGreeting returns a greeting stored in a database by it's id
	GetGreeting(id string) (string, error)

	// DeleteGreeting deletes greeting by ID
	DeleteGreeting(id string) error

	// Close closes all resources associated with this backend
	Close() error
}

// NotFoundError returns whenever the greeting requested is not found
type NotFoundError struct {
	ID string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("greeting with id '%v' not found", e.ID)
}
