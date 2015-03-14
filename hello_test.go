package hello

import (
	"math/rand"
	"testing"
	"time"

	"github.com/mailgun/log"
	. "gopkg.in/check.v1"
	"github.com/gravitational/hello/backend/membk"
)

// We use gocheck: a rich test framework on top of vanilla Go standard test module
// read more about it here https://labix.org/gocheck
func TestHello(t *testing.T) { TestingT(t) }

type HelloSuite struct {
	h Helloer
}

var _ = Suite(&HelloSuite{})

// SetUpSuite can be used to set up anything that needs to be setup once
// per execution of the test suite.
func (s *HelloSuite) SetUpSuite(c *C) {
	rand.Seed(time.Now().Unix())
	// Sometimes it is helpful to turn on logging in tests too,
	// it helps to see what would the output be in the real application
	// and make it less verbose and more helpful.
	log.Init([]*log.LogConfig{&log.LogConfig{Name: "console"}})
}

// SetUpTest allocates necessary resources for each test.
// Make sure each test is independent from the others, use
// SetUpTest to recreate environment and resources needed for test case.
func (s *HelloSuite) SetUpTest(c *C) {
	b := membk.New()
	c.Assert(b.UpsertGreeting("hello.us", "Hello"), IsNil)
	c.Assert(b.UpsertGreeting("hello.sp", "Hola"), IsNil)
	s.h = New(b)
}

// TearDownTest deallocates resources allocated in the test.
// You can assert correct execution in TearDownTest too
func (s *HelloSuite) TearDownTest(c *C) {
	c.Assert(s.h.Close(), IsNil)
}

// Try using table tests whenever possible, they are great for testing various edgen conditions
// and set your mind to the right state (finding edgy conditions when something crashes).
func (s *HelloSuite) TestTableOK(c *C) {
	tcs := []struct {
		name     string
		prompt   string
		param    string
		expected string
	}{
		{name: "human name", prompt: "hello.us", param: "Sasha", expected: "Hello, Sasha!"},
		{name: "animal name", prompt: "hello.sp", param: "Dog", expected: "Hola, Dog!"},
	}
	for i, tc := range tcs {
		// make sure to provide comment parameter to assertions done in a test driven
		// test cases, otherwise it would be hard to find the condition that led to
		// the test failure
		comment := Commentf("test #%d (%v) prompt=%v, param=%v", i+1, tc.name, tc.prompt, tc.param)
		out, err := s.h.Hello(tc.prompt, tc.param)
		c.Assert(err, IsNil, comment)
		c.Assert(out, Equals, tc.expected, comment)
	}
}

// TestErrors demonstrates testing for edge conditions triggering errors
// Note that not only we are testing that error has been returned,
// we are also making sure that error type is the correct one.
func (s *HelloSuite) TestErrors(c *C) {
	tcs := []struct {
		name   string
		prompt string
		param  string
	}{
		{name: "unsupported prompt", prompt: "helo.!", param: "sasha"},
		{name: "empty prompt", prompt: "", param: "sasha"},
		{name: "empty name", prompt: "", param: "dog"},
	}
	for i, tc := range tcs {
		comment := Commentf("test #%d: name: %v, prompt: %v, param: %v", i+1, tc.name, tc.prompt, tc.param)
		_, err := s.h.Hello(tc.prompt, tc.param)
		c.Assert(err, NotNil, comment)
	}
}

// Benchmarks are useful to test performance of some well isolated components
func (s *HelloSuite) BenchmarkHello(c *C) {
	// make sure to turn off the tests before running benchmark
	log.Init([]*log.LogConfig{})
	// run the h.Hello() function c.N times with average string
	for n := 0; n < c.N; n++ {
		s.h.Hello("hello.us", "some name")
	}
}
