// package etcdbk implements Etcd powered greetings backend
package etcdbk

import (
	"fmt"
	"strings"

	"github.com/gravitational/hello/Godeps/_workspace/src/github.com/coreos/go-etcd/etcd"
	"github.com/gravitational/hello/backend"
)

type BackendOption func(b *bk) error

func Consistency(v string) BackendOption {
	return func(b *bk) error {
		b.etcdConsistency = v
		return nil
	}
}

type bk struct {
	nodes []string

	etcdConsistency string
	etcdKey         string
	client          *etcd.Client
	cancelC         chan bool
	stopC           chan bool
}

func New(nodes []string, etcdKey string, options ...BackendOption) (backend.GreetingBackend, error) {
	if len(nodes) == 0 {
		return nil, fmt.Errorf("empty list of etcd nodes, supply at least one node")
	}

	if len(etcdKey) == 0 {
		return nil, fmt.Errorf("supply a valid etcd key")
	}

	b := &bk{
		nodes:   nodes,
		etcdKey: etcdKey,
		cancelC: make(chan bool, 1),
		stopC:   make(chan bool, 1),
	}
	b.etcdConsistency = etcd.WEAK_CONSISTENCY
	for _, o := range options {
		if err := o(b); err != nil {
			return nil, err
		}
	}
	if err := b.reconnect(); err != nil {
		return nil, err
	}
	return b, nil
}

func (b *bk) Close() error {
	return nil
}

func (b *bk) key(keys ...string) string {
	return strings.Join(append([]string{b.etcdKey}, keys...), "/")
}

func (b *bk) reconnect() error {
	b.client = etcd.NewClient(b.nodes)
	return nil
}

func (b *bk) UpsertGreeting(id, greeting string) error {
	_, err := b.client.Set(b.key("greetings", id), greeting, 0)
	return convertErr(err)
}

func (b *bk) GetGreeting(id string) (string, error) {
	re, err := b.client.Get(b.key("greetings", id), false, false)
	if err != nil {
		return "", convertErr(err)
	}
	return re.Node.Value, nil
}

// DeleteUser deletes a user with all the keys from the backend
func (b *bk) DeleteGreeting(id string) error {
	_, err := b.client.Delete(b.key("greetings", id), true)
	return convertErr(err)
}

func notFound(e error) bool {
	err, ok := e.(*etcd.EtcdError)
	return ok && err.ErrorCode == 100
}

func convertErr(e error) error {
	if e == nil {
		return nil
	}
	switch err := e.(type) {
	case *etcd.EtcdError:
		if err.ErrorCode == 100 {
			return &backend.NotFoundError{ID: err.Cause}
		}
	}
	return e
}
