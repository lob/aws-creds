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

	fmt.Println("Configuring profile settings...")
	err = configureProfiles(cmd)
	if err != nil {
		return err
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

	fmt.Print("\n")
	return nil
}

func configureProfiles(cmd *CMD) error {
	cont := true
	for cont {
		name, err := input.Prompt("Profile name: ", cmd.In, cmd.Out)
		if err != nil {
			return err
		}

		roleARN, err := input.Prompt("Role ARN (e.g. arn:aws:iam::123456789001:role/EngineeringRole): ", cmd.In, cmd.Out)
		if err != nil {
			return err
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

		more, err := input.Prompt("Do you want to configure more profiles? [y/N]: ", cmd.In, cmd.Out)
		if err != nil {
			return err
		}
		fmt.Print("\n")

		cont = strings.ToLower(more)[0] == 'y'
	}
	return nil
}
