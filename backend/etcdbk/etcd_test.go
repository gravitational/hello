package etcdbk

import (
	"os"
	"strings"
	"testing"

	"github.com/gravitational/hello/Godeps/_workspace/src/github.com/coreos/go-etcd/etcd"
	"github.com/gravitational/hello/backend/test"

	. "github.com/gravitational/hello/Godeps/_workspace/src/gopkg.in/check.v1"
)

func TestEtcd(t *testing.T) { TestingT(t) }

type EtcdSuite struct {
	suite      test.BackendSuite
	nodes      []string
	etcdPrefix string
	client     *etcd.Client
	key        string
}

var _ = Suite(&EtcdSuite{
	etcdPrefix: "/hello_test",
})

func (s *EtcdSuite) SetUpSuite(c *C) {
	nodes_string := os.Getenv("TEST_ETCD_NODES")
	if nodes_string == "" {
		// Skips the entire suite
		c.Skip("This test requires etcd, provide comma separated nodes in VULCAND_TEST_ETCD_NODES environment variable")
		return
	}
	s.nodes = strings.Split(nodes_string, ",")
}

func (s *EtcdSuite) SetUpTest(c *C) {
	// Initiate a backend with a registry
	b, err := New(s.nodes, s.etcdPrefix)
	c.Assert(err, IsNil)
	s.client = b.(*bk).client

	// Delete all values under the given prefix
	_, err = s.client.Get(s.etcdPrefix, false, false)
	if err != nil {
		if !notFound(err) {
			c.Assert(err, IsNil)
		}
	} else {
		_, err = s.client.Delete(s.etcdPrefix, true)
		if !notFound(err) {
			c.Assert(err, IsNil)
		}
	}

	// Set up suite
	s.suite.B = b
}

func (s *EtcdSuite) TestGreetingCRUD(c *C) {
	s.suite.GreetingCRUD(c)
}
