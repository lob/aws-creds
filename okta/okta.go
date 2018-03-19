package okta

import (
	"fmt"

	"github.com/lob/aws-creds/config"
	"github.com/lob/aws-creds/input"
)

// Login to Okta with the given username and password.
func Login(conf *config.Config, p input.Prompter, password string) (*SAMLResponse, error) {
	c, err := NewClient(conf.OktaHost)
	if err != nil {
		return nil, err
	}

	auth, err := login(c, conf.Username, password)
	if err != nil {
		return nil, err
	}
	fmt.Println("Authentication succeeded!")

	if auth.Status == "MFA_REQUIRED" {
		if err = auth.verifyMFA(c, conf, p); err != nil {
			return nil, err
		}
		fmt.Println("MFA confirmed!")
	}

	return getSAMLResponse(c, conf.OktaAppPath, auth.SessionToken)
}
