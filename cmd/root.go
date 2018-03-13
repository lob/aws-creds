package cmd

import (
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

const (
	configureCommand = "configure"
	refreshCommand   = ""
)

var (
	defaultConfigFilepath = os.Getenv("HOME") + "/.aws-creds/config"
	configFilepath        = flag.String("config", defaultConfigFilepath, "config file")
	help                  = flag.Bool("help", false, "print this help text")
)

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

	cmd := &CMD{In: in, Out: out}
	cmd.Command = ""
	if len(args) > 0 {
		cmd.Command = args[0]
	}

	cmd.Config = config.New(*configFilepath)
	err := cmd.Config.Load()

	switch cmd.Command {
	case configureCommand:
		return executeConfigure(cmd)
	case refreshCommand:
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unknown command: %s", cmd.Command)
	}
}
