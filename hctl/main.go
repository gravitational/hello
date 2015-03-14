package main

import (
	"os"

	"github.com/mailgun/log"
	"github.com/gravitational/hello/hctl/command"
)

func main() {

	log.Init([]*log.LogConfig{&log.LogConfig{Name: "console"}})

	cmd := command.NewCommand()
	err := cmd.Run(os.Args)
	if err != nil {
		log.Infof("error: %s\n", err)
	}
}
