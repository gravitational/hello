// package api implements HTTP API server and a client wrapper
package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gravitational/hello"
	"github.com/gravitational/form"
	"github.com/julienschmidt/httprouter" // APIServer is a http.Handler server requests to Helo server
	"github.com/gravitational/hello/backend"
)

type APIServer struct {
	httprouter.Router
	h hello.Helloer
	b backend.GreetingBackend
}

// NewAPIServer returns http.Handler compatible HTTP server
//
//  srv := NewAPIServer(h, b)
//  http.ListenAndServe(srv)
//
func NewAPIServer(h hello.Helloer, b backend.GreetingBackend) *APIServer {
	srv := &APIServer{
		h: h,
		b: b,
	}
	srv.Router = *httprouter.New()

	// Greetings CRUD
	srv.POST("/v1/greetings", srv.upsertGreeting)
	srv.GET("/v1/greetings/:prompt", srv.getGreeting)
	srv.DELETE("/v1/greetings/:prompt", srv.deleteGreeting)

	// Say hello
	srv.POST("/v1/hello", srv.hello)

	return srv
}

// note that these functions are not exported, so godoc is not mentioning them, as they are really
// implementation detail that is not visible to users
func (s *APIServer) upsertGreeting(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var prompt, value string
	err := form.Parse(r,
		form.String("prompt", &prompt, form.Required()),
		form.String("value", &value, form.Required()))
	if err != nil {
		replyErr(w, err)
		return
	}
	if err := s.b.UpsertGreeting(prompt, value); err != nil {
		replyErr(w, err)
		return
	}
	reply(w, http.StatusOK, greetingResponse{Greeting: greeting{Prompt: prompt, Value: value}})
}

func (s *APIServer) getGreeting(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	prompt := p[0].Value
	val, err := s.b.GetGreeting(prompt)
	if err != nil {
		replyErr(w, err)
		return
	}
	reply(w, http.StatusOK, &greetingResponse{Greeting: greeting{Prompt: prompt, Value: val}})
}

func (s *APIServer) deleteGreeting(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	prompt := p[0].Value
	if err := s.b.DeleteGreeting(prompt); err != nil {
		replyErr(w, err)
		return
	}

	reply(w, http.StatusOK, message(fmt.Sprintf("greeting '%v' deleted", prompt)))
}

func (s *APIServer) hello(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var prompt, name string

	err := form.Parse(r,
		form.String("prompt", &prompt, form.Required()),
		form.String("name", &name, form.Required()))

	if err != nil {
		replyErr(w, err)
		return
	}
	hello, err := s.h.Hello(prompt, name)
	if err != nil {
		replyErr(w, err)
		return
	}

	reply(w, http.StatusOK, &helloResponse{Value: hello})
}

type greetingResponse struct {
	Greeting greeting `json:"greeting"`
}

type greeting struct {
	Prompt string `json:"prompt"`
	Value  string `json:"value"`
}

type helloResponse struct {
	Value string `json:"val"`
}

func message(msg string) map[string]interface{} {
	return map[string]interface{}{"message": msg}
}

func replyErr(w http.ResponseWriter, e error) {
	switch err := e.(type) {
	case *backend.NotFoundError:
		reply(w, http.StatusNotFound, message(err.Error()))
	default:
		reply(w, http.StatusBadRequest, message(err.Error()))
	}
}

func reply(w http.ResponseWriter, code int, message interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	out, err := json.Marshal(message)
	if err != nil {
		out = []byte(`{"msg": "internal marshal error"}`)
	}
	w.Write(out)
}
