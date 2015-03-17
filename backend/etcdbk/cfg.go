package etcdbk

import (
	"encoding/json"
	"fmt"

	"github.com/gravitational/hello/backend"
)

// cfg represents JSON config for etcd backlend
type cfg struct {
	Nodes []string `json:"nodes"`
	Key   string   `json:"key"`
}

// FromString initializes the backend from backend-specific configuration string
//
//   backend.FromString(`{"nodes": ["http://localhost:4001], "key": "/hello"}`)
//
func FromString(v string) (backend.GreetingBackend, error) {
	if len(v) == 0 {
		return nil, fmt.Errorf(`please supply a valid dictionary, e.g. {"nodes": ["http://localhost:4001]}`)
	}
	var c *cfg
	if err := json.Unmarshal([]byte(v), &c); err != nil {
		return nil, fmt.Errorf("invalid backend configuration format, err: %v", err)
	}
	return New(c.Nodes, c.Key)
}
