package main

import (
	"fmt"

	"net/http"
	"os"

	"github.com/gravitational/hello"
	"github.com/gravitational/hello/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/gravitational/hello/Godeps/_workspace/src/github.com/mailgun/log"
	"github.com/gravitational/hello/api"
	"github.com/gravitational/hello/backend"
	"github.com/gravitational/hello/backend/etcdbk"
)

func main() {
	app := cli.NewApp()
	app.Name = "hctl"
	app.Usage = "Clustering Hello World application"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "addr", Value: "localhost:8080", Usage: "hello listening host:port"},
		cli.StringFlag{Name: "shell", Value: "/bin/sh", Usage: "path to shell to launch for interactive sessions"},

		cli.StringFlag{Name: "backend", Value: "etcd", Usage: "backend type, currently only 'etcd'"},
		cli.StringFlag{Name: "backendConfig", Value: "", Usage: "backend-specific configuration string"},

		cli.StringFlag{Name: "log", Value: "console", Usage: "Log output, currently 'console' or 'syslog'"},
		cli.StringFlag{Name: "logSeverity", Value: "WARN", Usage: "Log severity, logs warning by default"},
	}
	app.Action = run
	app.Run(os.Args)
}

func run(c *cli.Context) {
	if err := start(c); err != nil {
		log.Errorf("service err %v", err)
		return
	}
	log.Infof("service exited gracefully")
}

func setupLogging(c *cli.Context) error {
	s, err := log.SeverityFromString(c.String("logSeverity"))
	if err != nil {
		return err
	}
	log.Init([]*log.LogConfig{&log.LogConfig{Name: c.String("log")}})
	log.SetSeverity(s)
	return nil
}

func start(c *cli.Context) error {
	if err := setupLogging(c); err != nil {
		return err
	}

	b, err := initBackend(c.String("backend"), c.String("backendConfig"))
	if err != nil {
		return err
	}

	return http.ListenAndServe(c.String("addr"), api.NewAPIServer(hello.New(b), b))
}

func initBackend(btype, bcfg string) (backend.GreetingBackend, error) {
	switch btype {
	case "etcd":
		return etcdbk.FromString(bcfg)
	}
	return nil, fmt.Errorf("unsupported backend type: %v", btype)
}
