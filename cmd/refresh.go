package cmd

import (
	"errors"
	"fmt"

	"github.com/lob/aws-creds/config"
	"github.com/lob/aws-creds/okta"
)

func executeRefresh(cmd *Cmd) error {
	if cmd.Profile == "" {
		return errors.New("profile is required")
	}
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

	msg := fmt.Sprintf("Enter password for %s: ", cmd.Config.Username)
	password, err := cmd.Input.PromptPassword(msg)
	if err != nil {
		return err
	}

	return okta.Login(cmd.Config, cmd.Input, password)
}
