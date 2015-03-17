package api

import (
	"encoding/json"
	"fmt"

	"net/http"
	"net/url"

	"github.com/gravitational/hello/Godeps/_workspace/src/github.com/gravitational/roundtrip" // Client is an HTTP RPC client to the running Hello server
	"github.com/gravitational/hello/backend"
)

type Client struct {
	roundtrip.Client
}

// NewClient returns a new instance of the client connected to the Hello server
// that is reachable by address addr
func NewClient(addr string) (*Client, error) {
	c, err := roundtrip.NewClient(addr, CurrentVersion)
	if err != nil {
		return nil, err
	}
	return &Client{*c}, nil
}

// UpsertGreeting updates or inserts the greeting into the database backend
//
//     c.UpsertGreeting("hello.us", "Hello")
//
func (c *Client) UpsertGreeting(prompt, value string) error {
	_, err := convert(
		c.PostForm(
			c.Endpoint("greetings"),
			url.Values{"prompt": []string{prompt}, "value": []string{value}}))
	return err
}

// GetGreeting returns the value of the greeting by it's prompt id
//
//     val, err := c.GetGreeting("hello.us", "Hello")
//
func (c *Client) GetGreeting(prompt string) (string, error) {
	body, err := convert(
		c.Get(c.Endpoint("greetings", prompt), url.Values{}))
	if err != nil {
		return "", err
	}
	var g *greetingResponse
	if err := json.Unmarshal(body, &g); err != nil {
		return "", err
	}
	return g.Greeting.Value, nil
}

// DeleteGreeting deletes the greeting from DB by it's ID
//
//     err := c.DeleteGreeting("hello.us")
//
func (c *Client) DeleteGreeting(prompt string) error {
	_, err := convert(c.Delete(c.Endpoint("greetings", prompt)))
	return err
}

// Hello generates the Hello sentence by prompt id and a name
//
//     h, err := c.Hello("hello.us", "Dog") // Hello, Dog!
//     h, err := c.Hello("hello.sp", "Dog") // Hola, Dog!
//
func (c *Client) Hello(prompt, name string) (string, error) {
	body, err := convert(
		c.PostForm(
			c.Endpoint("hello"),
			url.Values{"prompt": []string{prompt}, "name": []string{name}}))
	if err != nil {
		return "", err
	}
	var h *helloResponse
	if err := json.Unmarshal(body, &h); err != nil {
		return "", err
	}
	return h.Value, nil
}

func (c *Client) Close() error {
	return nil
}

// convert converts generic HTTP response codes to hello-specific errors
func convert(re *roundtrip.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	if re.Code() == http.StatusNotFound {
		return nil, &backend.NotFoundError{ID: string(re.Bytes())}
	}
	if re.Code() >= 200 && re.Code() < 300 {
		return re.Bytes(), nil
	}
	return nil, fmt.Errorf(string(re.Bytes()))
}

// CurrentVersion is a current API version prefix
const CurrentVersion = "v1"
