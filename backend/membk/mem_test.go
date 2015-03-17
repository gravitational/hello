package membk

import (
	"testing"

	. "github.com/gravitational/hello/Godeps/_workspace/src/gopkg.in/check.v1"
	"github.com/gravitational/hello/backend/test"
)

// Test suite for memory backend is very simple, it fully relies on acceptance test
// to make sure the behavior is expected
func TestMem(t *testing.T) { TestingT(t) }

type MemSuite struct {
	bk    *MemBackend
	suite test.BackendSuite
}

var _ = Suite(&MemSuite{})

func (s *MemSuite) SetUpTest(c *C) {
	s.bk = New()
	s.suite.B = s.bk
}

func (s *MemSuite) TearDownTest(c *C) {
	c.Assert(s.bk.Close(), IsNil)
}

func (s *MemSuite) TestGretingsCRUD(c *C) {
	s.suite.GreetingCRUD(c)
}
