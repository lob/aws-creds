package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/lob/aws-creds/config"
)

const exampleEmbedLink = "https://example.okta.com/home/amazon_aws/0oa54k1gk2ukOJ9nGDt7/252"

var oktaRegex = regexp.MustCompile(`(https://.*\.okta\.com)(/home/[^/]*/[^/]*/[^/]*)`)

func executeConfigure(cmd *Cmd) error {
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

	return cmd.Config.Save()
}

func configureGlobal(cmd *Cmd) error {
	username, err := cmd.Input.Prompt("Okta username: ")
	if err != nil {
		return err
	}

	prompt := fmt.Sprintf("Okta AWS Embed Link (e.g. %s): ", exampleEmbedLink)
	link, err := cmd.Input.Prompt(prompt)
	if err != nil {
		return err
	}
	matches := oktaRegex.FindStringSubmatch(link)
	if len(matches) != 3 {
		return fmt.Errorf("%s doesn't look like an Embed Link", link)
	}

	cmd.Config.Username = username
	cmd.Config.OktaHost = matches[1]
	cmd.Config.OktaAppPath = matches[2]

	fmt.Print("\n")
	return nil
}

func configureProfiles(cmd *Cmd) error {
	cont := true
	for cont {
		name, err := cmd.Input.Prompt("Profile name: ")
		if err != nil {
			return err
		}

		roleARN, err := cmd.Input.Prompt("Role ARN (e.g. arn:aws:iam::123456789001:role/EngineeringRole): ")
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

		more, err := cmd.Input.Prompt("Do you want to configure more profiles? [y/N]: ")
		if err != nil {
			return err
		}
		fmt.Print("\n")

		cont = strings.ToLower(more)[0] == 'y'
	}
	return nil
}
