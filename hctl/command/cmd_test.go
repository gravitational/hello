package command

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gravitational/hello"
	"github.com/gravitational/hello/api"
	"github.com/gravitational/hello/backend/membk"

	. "github.com/gravitational/hello/Godeps/_workspace/src/gopkg.in/check.v1"
)

const OK = ".*OK.*"

func TestTeleportCLI(t *testing.T) { TestingT(t) }

type CmdSuite struct {
	srv *httptest.Server
	clt *api.Client
	cmd *Command
	out *bytes.Buffer
	bk  *membk.MemBackend
}

var _ = Suite(&CmdSuite{})

func (s *CmdSuite) SetUpSuite(c *C) {
}

func (s *CmdSuite) SetUpTest(c *C) {
	s.bk = membk.New()
	h := hello.New(s.bk)

	s.srv = httptest.NewServer(api.NewAPIServer(h, s.bk))
	clt, err := api.NewClient(s.srv.URL)
	c.Assert(err, IsNil)
	s.clt = clt

	s.out = &bytes.Buffer{}
	s.cmd = &Command{out: s.out, url: s.srv.URL}
}

func (s *CmdSuite) TearDownTest(c *C) {
	s.srv.Close()
}

func (s *CmdSuite) runString(in string) string {
	return s.run(strings.Split(in, " ")...)
}

func (s *CmdSuite) run(params ...string) string {
	args := []string{"hctl"}
	args = append(args, params...)
	args = append(args, fmt.Sprintf("--hello=%s", s.srv.URL))
	s.out = &bytes.Buffer{}
	s.cmd = &Command{out: s.out, url: s.srv.URL}
	s.cmd.Run(args)
	return strings.Replace(s.out.String(), "\n", " ", -1)
}

func (s *CmdSuite) TestGreetingCRUD(c *C) {
	c.Assert(
		s.run("greeting", "upsert", "-id", "hello.us", "-val", "Hello"),
		Matches, fmt.Sprintf(".*%v.*", "upserted"))
	c.Assert(s.bk.Greetings["hello.us"], Equals, "Hello")

	c.Assert(
		s.run("greeting", "get", "-id", "hello.us"),
		Matches, fmt.Sprintf(".*%v.*", "Hello"))

	c.Assert(
		s.run("greeting", "upsert", "-id", "hello.us", "-val", "Howdy"),
		Matches, fmt.Sprintf(".*%v.*", "upserted"))
	c.Assert(s.bk.Greetings["hello.us"], Equals, "Howdy")

	c.Assert(
		s.run("greeting", "delete", "-id", "hello.us"),
		Matches, fmt.Sprintf(".*%v.*", "deleted"))
	_, ok := s.bk.Greetings["hello.us"]
	c.Assert(ok, Equals, false)
}

func (s *CmdSuite) TestHello(c *C) {
	c.Assert(s.bk.UpsertGreeting("hello.us", "Hello"), IsNil)

	c.Assert(
		s.run("hello", "-id", "hello.us", "-name", "Dog"),
		Matches, fmt.Sprintf(".*%v.*", "Hello, Dog!"))
}
