package cmd

import (
	"fmt"
	"io"
	"syscall"

	"github.com/lob/aws-creds/config"
	"github.com/lob/aws-creds/input"
	"github.com/lob/aws-creds/okta"
)

func executeRefresh(cmd *Cmd) error {
	return executeRefreshWithPrompt(cmd, input.PromptPassword)
}

func executeRefreshWithPrompt(cmd *Cmd, prompt func(string, int, io.Writer) (string, error)) error {
	var profile *config.Profile
	for _, p := range cmd.Config.Profiles {
		if p.Name == cmd.Profile {
			profile = p
			break
		}
	}
	if profile == nil {
		return fmt.Errorf("profile %s not configured", cmd.Profile)
	}

	password, err := prompt(fmt.Sprintf("Enter password for %s: ", cmd.Config.Username), syscall.Stdin, cmd.Out)
	if err != nil {
		return err
	}

	return okta.Login(cmd.Config, password)
}
