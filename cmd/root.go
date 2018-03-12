package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/lob/aws-creds/config"
)

// CMD contains the necessary information for the CLI to function.
type CMD struct {
	Command string
	Config  *config.Config
	In      io.Reader
	Out     io.Writer
}

const defaultConfigFilepath = "/.aws-creds/config"

var configFilepath = flag.String("config", "", fmt.Sprintf("config file (default is $HOME%s)", defaultConfigFilepath))
var help = flag.Bool("help", false, "print this help text")

// Execute runs the CLI application.
func Execute() {
	flag.Parse()

	if err := execute(flag.Args(), os.Stdin, os.Stdout); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func execute(args []string, in io.Reader, out io.Writer) error {
	if *help {
		flag.Usage()
		return nil
	}

	if *configFilepath == "" {
		*configFilepath = os.Getenv("HOME") + defaultConfigFilepath
	}

	cmd := &CMD{In: in, Out: out}
	cmd.Command = ""
	if len(args) > 0 {
		cmd.Command = args[0]
	}

	cmd.Config = config.New(*configFilepath)
	err := cmd.Config.Load()
	if err != nil && cmd.Command != "configure" {
		return err
	}

	switch cmd.Command {
	case "configure":
		return executeConfigure(cmd)
	case "":
		return nil
	default:
		return errors.New("unknown command")
	}
}
