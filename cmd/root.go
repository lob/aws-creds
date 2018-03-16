package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/lob/aws-creds/config"
)

// Cmd contains the necessary information for the CLI to function.
type Cmd struct {
	Command string
	Config  *config.Config
	Profile string
	In      io.Reader
	Out     io.Writer
}

const (
	configureCommand = "configure"
	refreshCommand   = ""
)

var (
	defaultConfigFilepath = os.Getenv("HOME") + "/.aws-creds/config"
	configFilepath        = flag.String("c", defaultConfigFilepath, "config file")
	profile               = flag.String("p", "", "AWS profile to retrieve credentials for (required)")
	help                  = flag.Bool("h", false, "print this help text")
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

	cmd := &Cmd{
		Command: "",
		Config:  config.New(*configFilepath),
		Profile: *profile,
		In:      in,
		Out:     out,
	}
	if len(args) > 0 {
		cmd.Command = args[0]
	}

	err := cmd.Config.Load()

	switch cmd.Command {
	case configureCommand:
		// TODO(robin): add debug log if err != nil
		return executeConfigure(cmd)
	case refreshCommand:
		if err != nil {
			return err
		}
		return executeRefresh(cmd)
	default:
		return fmt.Errorf("unknown command: %s", cmd.Command)
	}
}
