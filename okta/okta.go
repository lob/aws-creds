package okta

import (
	"fmt"

	"github.com/lob/aws-creds/config"
	"github.com/lob/aws-creds/input"
)

// Login to Okta with the given username and password.
func Login(conf *config.Config, p input.Prompter, password string) error {
	c, err := NewClient(conf.OktaHost)
	if err != nil {
		return err
	}

	auth, err := login(c, conf.Username, password)
	if err != nil {
		return err
	}
	fmt.Println("Authentication succeeded!")

	if auth.Status == "MFA_REQUIRED" {
		if err = auth.verifyMFA(c, conf, p); err != nil {
			return err
		}
		fmt.Println("MFA confirmed!")
	}
	return nil
}
