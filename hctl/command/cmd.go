package command

import (
	"fmt"
	"io"
	"os"

	"github.com/buger/goterm"
	"github.com/codegangsta/cli"
	"github.com/gravitational/hello/api"
	"strings"
)

type Command struct {
	url    string
	client *api.Client
	out    io.Writer
	in     io.Reader
}

func NewCommand() *Command {
	return &Command{
		out: os.Stdout,
		in:  os.Stdin,
	}
}

func (cmd *Command) Run(args []string) error {
	url, args, err := findURL(args)
	if err != nil {
		return err
	}
	cmd.url = url
	client, err := api.NewClient(cmd.url)
	if err != nil {
		return err
	}
	cmd.client = client

	app := cli.NewApp()
	app.Name = "hctl"
	app.Usage = "CLI for managing hello service"
	app.Flags = flags()

	app.Commands = []cli.Command{
		newGreetingCommand(cmd),
		newHelloCommand(cmd),
	}
	return app.Run(args)
}

func (cmd *Command) printError(err error) {
	fmt.Fprint(cmd.out, goterm.Color(fmt.Sprintf("ERROR: %s", err), goterm.RED)+"\n")
}

func (cmd *Command) printOK(message string, params ...interface{}) {
	fmt.Fprintf(cmd.out,
		goterm.Color(
			fmt.Sprintf("OK: %s\n", fmt.Sprintf(message, params...)), goterm.GREEN)+"\n")
}

func flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{Name: "hello", Value: DefaultHelloURL, Usage: "Hello URL"},
	}
}

const DefaultHelloURL = "localhost:8080"

// This function extracts url from the command line regardless of it's position
// this is a workaround, as cli libary does not support "superglobal" urls yet.
func findURL(args []string) (string, []string, error) {
	for i, arg := range args {
		if strings.HasPrefix(arg, "--hello=") || strings.HasPrefix(arg, "-hello=") {
			out := strings.Split(arg, "=")
			return out[1], cut(i, i+1, args), nil
		} else if strings.HasPrefix(arg, "-hello") || strings.HasPrefix(arg, "--hello") {
			// This argument should not be the last one
			if i > len(args)-2 {
				return "", nil, fmt.Errorf("provide a valid URL")
			}
			return args[i+1], cut(i, i+2, args), nil
		}
	}
	return "http://localhost:8080", args, nil
}

func cut(i, j int, args []string) []string {
	s := []string{}
	s = append(s, args[:i]...)
	return append(s, args[j:]...)
}
