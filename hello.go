// Package hello provides simple "hello, world!" output feature.
// Hello is a little bit over engineered, it is done to demonstrate concepts
// that can be used in  more complex projects
//
// Hello documentation is godoc friendly.
// Execute the following command to read the docs about any gravitational package:
//
//  godoc github.com/gravitational/
//
// Here's how you use the library:
//
//  import (
//      "fmt"
//
//      "github.com/gravitational/hello"
//      "github.com/gravitational/hello"
//  )
//
//  func main() {
//     b :=
//     h := hello.New()
//     fmt.Prinln(h.Hello())
//  }
package hello

import (
	"fmt"

	"github.com/gravitational/hello/Godeps/_workspace/src/github.com/mailgun/log" // Helloer interface represents "Hello, world!" functionality providers.
	"github.com/gravitational/hello/backend"
)

type Helloer interface {
	// Hello generates and returns "Hello, <username>!" message when called
	// with a string parameter
	Hello(prompt, username string) (string, error)
	// Close deallocates any resources that were allocated by instance of helloer
	Close() error
}

// New returns a new instance of Helloer
func New(b backend.GreetingBackend) Helloer {
	return &helloer{
		b: b,
	}
}

// helloer is an internal implementation of Helloer that uses fmt
type helloer struct {
	b backend.GreetingBackend
}

// Hello is Sprinf-based implementation and should not be used in high-perf
// environments as it generates a new string when called each time.
func (h *helloer) Hello(prompt, username string) (string, error) {
	log.Infof("Hello(%v, %v)", prompt, username)
	greeting, err := h.b.GetGreeting(prompt)
	if err != nil {
		return "", fmt.Errorf("error when retrieving backend: %v", err)
	}
	if username == "" {
		// try do be specific when returning errors, it makes
		// it easier for applications to provide better error
		// handling by testing the type of error returned.
		return "", &EmptyParamError{}
	}
	return fmt.Sprintf("%v, %v!", greeting, username), nil
}

func (h *helloer) Close() error {
	return nil
}

// EmptyParamError is retuned whenever someone omitted greeting name
// when calling Hello() method.
type EmptyParamError struct {
}

// Error provides human readable explanation of error
func (*EmptyParamError) Error() string {
	return "empty parameters are not allowed for Hello()"
}
