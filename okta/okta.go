package okta

import (
	"fmt"

	"github.com/lob/aws-creds/config"
	"github.com/lob/aws-creds/input"
)

// Login to Okta with the given username and password.
func Login(conf *config.Config, p input.Prompter, sessionCookie, password string) (*SAMLResponse, string, error) {
	c, err := NewClient(conf.OktaHost, sessionCookie)
	if err != nil {
		return nil, "", err
	}

	if sessionCookie != "" {
		saml, err := getSAMLResponse(c, conf.OktaAppPath, "")
		if err == nil {
			return saml, "", nil
		}
	}

	auth, err := login(c, conf.Username, password)
	if err != nil {
		return nil, "", err
	}
	fmt.Println("Authentication succeeded!")

	if auth.Status == "MFA_REQUIRED" {
		if err = auth.verifyMFA(c, conf, p); err != nil {
			return nil, "", err
		}
		fmt.Println("MFA confirmed!")
	}

	saml, err := getSAMLResponse(c, conf.OktaAppPath, auth.SessionToken)
	if err != nil {
		return nil, "", err
	}

	cookies := c.http.Jar.Cookies(c.url)
	for _, cookie := range cookies {
		if cookie.Name == "sid" {
			sessionCookie = cookie.Value
		}
	}

	return saml, sessionCookie, nil
}
