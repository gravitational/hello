package api

import (
	"net/http/httptest"
	"testing"

	"github.com/gravitational/hello"
	"github.com/gravitational/hello/backend"
	"github.com/gravitational/hello/backend/membk"

	. "gopkg.in/check.v1"
)

func TestAPI(t *testing.T) { TestingT(t) }

type APISuite struct {
	srv *httptest.Server
	clt *Client
	bk  *membk.MemBackend
}

var _ = Suite(&APISuite{})

func (s *APISuite) SetUpSuite(c *C) {
}

func (s *APISuite) SetUpTest(c *C) {
	s.bk = membk.New()

	h := hello.New(s.bk)
	s.srv = httptest.NewServer(
		NewAPIServer(h, s.bk))
	clt, err := NewClient(s.srv.URL)
	c.Assert(err, IsNil)
	c.Assert(clt, NotNil)
	s.clt = clt
}

func (s *APISuite) TearDownTest(c *C) {
	s.srv.Close()
}

func (s *APISuite) TestGreetingsCRUD(c *C) {
	// Create
	c.Assert(s.clt.UpsertGreeting("hello.us", "Hello"), IsNil)
	c.Assert(s.bk.Greetings["hello.us"], Equals, "Hello")

	// Read
	g, err := s.clt.GetGreeting("hello.us")
	c.Assert(err, IsNil)
	c.Assert(g, Equals, "Hello")

	// Update
	c.Assert(s.clt.UpsertGreeting("hello.us", "Howdy"), IsNil)
	c.Assert(s.bk.Greetings["hello.us"], Equals, "Howdy")

	g, err = s.clt.GetGreeting("hello.us")
	c.Assert(err, IsNil)
	c.Assert(g, Equals, "Howdy")

	// Delete
	c.Assert(s.clt.DeleteGreeting("hello.us"), IsNil)
	_, ok := s.bk.Greetings["hello.us"]
	c.Assert(ok, Equals, false)

	// Make sure it's not found
	_, err = s.clt.GetGreeting("hello.us")
	c.Assert(err, FitsTypeOf, &backend.NotFoundError{})
}

func (s *APISuite) TestHello(c *C) {
	c.Assert(s.clt.UpsertGreeting("hello.us", "Hello"), IsNil)
	c.Assert(s.clt.UpsertGreeting("hello.sp", "Hola"), IsNil)

	hello, err := s.clt.Hello("hello.us", "John")
	c.Assert(err, IsNil)
	c.Assert(hello, Equals, "Hello, John!")

	hello, err = s.clt.Hello("hello.sp", "John")
	c.Assert(err, IsNil)
	c.Assert(hello, Equals, "Hola, John!")
}
