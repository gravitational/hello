package command

import (
	"github.com/codegangsta/cli"
)

func newGreetingCommand(c *Command) cli.Command {
	return cli.Command{
		Name:  "greeting",
		Usage: "Operations with stored greetings",
		Subcommands: []cli.Command{
			{
				Name:   "upsert",
				Usage:  "Upsert greeting",
				Action: c.upsertGreeting,
				Flags: []cli.Flag{
					cli.StringFlag{Name: "id", Usage: "Greeting id"},
					cli.StringFlag{Name: "val, v", Usage: "Greeting value"},
				},
			},
			{
				Name:   "get",
				Usage:  "Get greeting by id",
				Action: c.getGreeting,
				Flags: []cli.Flag{
					cli.StringFlag{Name: "id", Usage: "Greeting id to get"},
				},
			},
			{
				Name:   "delete",
				Usage:  "Delete greeting by ID",
				Action: c.deleteGreeting,
				Flags: []cli.Flag{
					cli.StringFlag{Name: "id", Usage: "Greeting id to delete"},
				},
			},
		},
	}
}

func (cmd *Command) upsertGreeting(c *cli.Context) {
	err := cmd.client.UpsertGreeting(c.String("id"), c.String("val"))
	if err != nil {
		cmd.printError(err)
		return
	}
	cmd.printOK("greeting %v upserted", c.String("id"))
}

func (cmd *Command) deleteGreeting(c *cli.Context) {
	if err := cmd.client.DeleteGreeting(c.String("id")); err != nil {
		cmd.printError(err)
		return
	}
	cmd.printOK("greeting %v deleted", c.String("user"))
}

func (cmd *Command) getGreeting(c *cli.Context) {
	val, err := cmd.client.GetGreeting(c.String("id"))
	if err != nil {
		cmd.printError(err)
		return
	}
	cmd.printOK("Greeting: %v %v", c.String("id"), val)
}
