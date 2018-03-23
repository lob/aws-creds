package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lob/aws-creds/aws"
	"github.com/lob/aws-creds/config"
	"github.com/lob/aws-creds/okta"
	"github.com/zalando/go-keyring"
)

const (
	keyringService = "aws-creds"
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

	password, err := getPassword(cmd)
	if err != nil {
		return err
	}

	saml, err := okta.Login(cmd.Config, cmd.Input, password)
	if err != nil {
		return err
	}

	creds, err := aws.GetCreds(cmd.STS, saml, profile)
	if err != nil {
		return err
	}

	return aws.WriteCreds(creds, profile, cmd.Config.CredentialsFilepath)
}

func getPassword(cmd *Cmd) (string, error) {
	var password string

	password, err := keyring.Get(keyringService, cmd.Config.Username)
	if err != nil && err != keyring.ErrNotFound {
		return "", err
	}
	if err == nil {
		fmt.Println("Password fetched from keyring.")
		return password, nil
	}
	fmt.Println("Password not found in keyring.")

	return promptPassword(cmd)
}

func promptPassword(cmd *Cmd) (string, error) {
	msg := fmt.Sprintf("Enter password for %s: ", cmd.Config.Username)
	password, err := cmd.Input.PromptPassword(msg)
	if err != nil {
		return "", err
	}
	save, err := cmd.Input.Prompt("Do you want to securely save your password in your system keyring? [y/N]: ")
	if err != nil {
		return "", err
	}
	if strings.ToLower(save)[0] == 'y' {
		err := keyring.Set(keyringService, cmd.Config.Username, password)
		if err != nil {
			return "", err
		}
		fmt.Println("Password saved!")
	}
	return password, nil
}
