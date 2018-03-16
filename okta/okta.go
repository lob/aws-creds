package okta

import (
	"fmt"

	"github.com/lob/aws-creds/config"
)

// Login to Okta with the given username and password.
func Login(conf *config.Config, password string) error {
	c, err := NewClient(conf.OktaHost)
	if err != nil {
		return err
	}

	_, err = login(c, conf.Username, password)
	if err != nil {
		return err
	}
	fmt.Println("Authentication succeeded!")
	return nil
}
