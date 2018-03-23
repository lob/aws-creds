package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"

	"github.com/lob/aws-creds/config"
	"github.com/lob/aws-creds/input"
)

// Cmd contains the necessary information for the CLI to function.
type Cmd struct {
	Command string
	Config  *config.Config
	Profile string
	Input   input.Prompter
	STS     stsiface.STSAPI
}

const (
	configureCommand = "configure"
	refreshCommand   = ""
)

var (
	version = "master"

	defaultConfigFilepath = os.Getenv("HOME") + "/.aws-creds/config"
	configFilepath        = flag.String("c", defaultConfigFilepath, "config file")
	profile               = flag.String("p", "", "AWS profile to retrieve credentials for (required)")
	printVersion          = flag.Bool("v", false, "print the version")
	printHelp             = flag.Bool("h", false, "print this help text")
)

// Execute runs the CLI application.
func Execute(p input.Prompter) {
	flag.Parse()

	if err := execute(flag.Args(), p); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func execute(args []string, p input.Prompter) error {
	if *printVersion {
		fmt.Println(version)
		return nil
	}

	if *printHelp {
		flag.Usage()
		return nil
	}

	sess := session.Must(session.NewSession())
	conf, err := config.New(*configFilepath)
	if err != nil {
		return err
	}

	cmd := &Cmd{
		Command: "",
		Config:  conf,
		Profile: *profile,
		Input:   p,
		STS:     sts.New(sess),
	}
	if len(args) > 0 {
		cmd.Command = args[0]
	}

	err = cmd.Config.Load()

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
