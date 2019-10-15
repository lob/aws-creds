package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/lob/aws-creds/pkg/config"
	"github.com/lob/aws-creds/pkg/input"
)

type awsProfiles []string

func (i *awsProfiles) String() string {
	return strings.Join(*i, ", ")
}

func (i *awsProfiles) Set(value string) error {
	*i = append(*i, value)
	return nil
}

// Cmd contains the necessary information for the CLI to function.
type Cmd struct {
	Command  string
	Config   *config.Config
	Profiles awsProfiles
	Input    input.Prompter
	STS      stsiface.STSAPI
}

const (
	configureCommand = "configure"
	refreshCommand   = ""
)

var profiles awsProfiles

var (
	version = "development build"

	defaultConfigFilepath = os.Getenv("HOME") + "/.aws-creds/config"
	configFilepath        = flag.String("c", defaultConfigFilepath, fmt.Sprintf("config file (default: %q)", defaultConfigFilepath))
	allProfiles           = flag.Bool("a", false, "retrieve credentials for all configured AWS profiles")
	printVersion          = flag.Bool("v", false, "print the version")
	printHelp             = flag.Bool("h", false, "print this help text")
)

func init() {
	flag.Var(&profiles, "p", "AWS profiles to retrieve credentials for (required)")
}

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
		fmt.Printf("aws-creds %s\n", version)
		return nil
	}

	if *printHelp {
		printUsage()
		return nil
	}

	sess := session.Must(session.NewSession())
	conf, err := config.New(*configFilepath)
	if err != nil {
		return err
	}

	cmd := &Cmd{
		Command:  "",
		Config:   conf,
		Profiles: profiles,
		Input:    p,
		STS:      sts.New(sess),
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

		if *allProfiles {
			if len(profiles) > 0 {
				return errors.New("you can provide profiles with '-p' or select all with '-a' but not both")
			}

			var allProfiles awsProfiles
			for _, p := range conf.Profiles {
				allProfiles = append(allProfiles, p.Name)
			}

			cmd.Profiles = allProfiles
		}

		return executeRefresh(cmd)
	default:
		return fmt.Errorf("unknown command: %s", cmd.Command)
	}
}

func printUsage() {
	usage := `aws-creds %s
CLI tool to authenticate with Okta as the IdP to fetch AWS credentials

Usage:
  aws-creds [options] [commands]

Available Commands:
  configure	configure aws-creds with Okta configs and AWS profiles

Flags:
`

	flag.VisitAll(func(f *flag.Flag) {
		usage += fmt.Sprintf("  -%s\t%s\n", f.Name, f.Usage)
	})
	fmt.Printf(usage, version)
}
