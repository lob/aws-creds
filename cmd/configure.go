package cmd

import (
	"fmt"
	"strings"

	"github.com/lob/aws-creds/config"
	"github.com/lob/aws-creds/input"
)

func executeConfigure(cmd *CMD) error {
	fmt.Println("Configuring global settings...")
	err := configureGlobal(cmd)
	if err != nil {
		return err
	}
	fmt.Print("\n")

	fmt.Println("Configuring profile settings...")
	cont := true
	for cont {
		cont, err = configureProfile(cmd)
		if err != nil {
			return err
		}
	}

	err = cmd.Config.Save()
	if err == nil {
		fmt.Println("Configuration saved!")
	}
	return err
}

func configureGlobal(cmd *CMD) error {
	username, err := input.Prompt("Okta username: ", cmd.In, cmd.Out)
	if err != nil {
		return err
	}

	org, err := input.Prompt("Okta org (e.g. for https://example.okta.com, the org is example): ", cmd.In, cmd.Out)
	if err != nil {
		return err
	}

	cmd.Config.Username = username
	cmd.Config.OktaOrgURL = fmt.Sprintf("https://%s.okta.com", org)
	return nil
}

func configureProfile(cmd *CMD) (bool, error) {
	name, err := input.Prompt("Profile name: ", cmd.In, cmd.Out)
	if err != nil {
		return false, err
	}

	roleARN, err := input.Prompt("Role ARN (e.g. arn:aws:iam::123456789001:role/EngineeringRole): ", cmd.In, cmd.Out)
	if err != nil {
		return false, err
	}

	found := false
	for _, p := range cmd.Config.Profiles {
		if p.Name == name {
			found = true
			p.RoleARN = roleARN
			break
		}
	}
	if !found {
		cmd.Config.Profiles = append(cmd.Config.Profiles, &config.Profile{
			Name:    name,
			RoleARN: roleARN,
		})
	}

	cont, err := input.Prompt("Do you want to configure more profiles? [y/N]: ", cmd.In, cmd.Out)
	if err != nil {
		return false, err
	}
	fmt.Print("\n")

	return strings.ToLower(cont)[0] == 'y', nil
}
