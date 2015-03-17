// package test contains a backend acceptance test suite that does not depend
// on the implementaiton and guarantees the same behavior for all backend implementations.
// see backend/membk and backend/etcdbk to see how the suite is used
package test

import (
	"testing"

	. "github.com/gravitational/hello/Godeps/_workspace/src/gopkg.in/check.v1"
	"github.com/gravitational/hello/backend"
)

func TestBackend(t *testing.T) { TestingT(t) }

type BackendSuite struct {
	B backend.GreetingBackend
}

// GreetingsCRUD tests simple CRUD cycle for greetings
func (s *BackendSuite) GreetingCRUD(c *C) {
	// Create
	c.Assert(s.B.UpsertGreeting("hello.us", "Hello"), IsNil)

	// Read
	g, err := s.B.GetGreeting("hello.us")
	c.Assert(err, IsNil)
	c.Assert(g, Equals, "Hello")

	// Update
	c.Assert(s.B.UpsertGreeting("hello.us", "Howdy"), IsNil)

	g, err = s.B.GetGreeting("hello.us")
	c.Assert(err, IsNil)
	c.Assert(g, Equals, "Howdy")

	// Delete
	c.Assert(s.B.DeleteGreeting("hello.us"), IsNil)

	g, err = s.B.GetGreeting("hello.us")
	c.Assert(err, FitsTypeOf, &backend.NotFoundError{})
}
