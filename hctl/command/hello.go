package command

import (
	"github.com/codegangsta/cli"
)

func newHelloCommand(c *Command) cli.Command {
	return cli.Command{
		Name:   "hello",
		Usage:  "Say hello",
		Action: c.sayHello,
		Flags: []cli.Flag{
			cli.StringFlag{Name: "id", Usage: "Greeting id"},
			cli.StringFlag{Name: "name", Usage: "Name to greet"},
		},
	}
}

func (cmd *Command) sayHello(c *cli.Context) {
	hello, err := cmd.client.Hello(c.String("id"), c.String("name"))
	if err != nil {
		cmd.printError(err)
		return
	}
	cmd.printOK(hello)
}
