package okta

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lob/aws-creds/config"
	"github.com/lob/aws-creds/input"
)

// Auth represents the authentication response from Okta.
type Auth struct {
	StateToken   string `json:"stateToken"`
	Status       string `json:"status"`
	SessionToken string `json:"sessionToken"`
	Embedded     struct {
		Factors []*Factor `json:"factors"`
	} `json:"_embedded"`
}

// Factor contains information about a specific MFA factor.
type Factor struct {
	FactorType string `json:"factorType"`
	Profile    struct {
		CredentialID string `json:"credentialId"`
	} `json:"profile"`
	Links struct {
		Verify struct {
			Href string `json:"href"`
		} `json:"verify"`
	} `json:"_links"`
}

func login(c *Client, username, password string) (*Auth, error) {
	payload := []byte(fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password))
	resp, err := c.Post("/api/v1/authn", payload)
	if err != nil {
		return nil, err
	}

	auth := &Auth{}
	err = json.NewDecoder(resp).Decode(auth)
	if err != nil {
		return nil, err
	}
	return auth, nil
}

func (auth *Auth) verifyMFA(c *Client, conf *config.Config, p input.Prompter) error {
	var factor *Factor
	if conf.PreferredFactorType != "" {
		for _, f := range auth.Embedded.Factors {
			if f.FactorType == conf.PreferredFactorType {
				factor = f
				break
			}
		}
		if factor == nil {
			return fmt.Errorf("%s isn't available for MFA; reconfigure aws-creds or check your Okta settings", conf.PreferredFactorType)
		}
	} else {
		factorIndex, err := promptForFactor(auth.Embedded.Factors, p)
		if err != nil {
			return err
		}
		factor = auth.Embedded.Factors[factorIndex]

		save, err := p.Prompt("Do you want to remember to use this factor? [y/N]: ")
		if err != nil {
			return err
		}
		if strings.ToLower(save)[0] == 'y' {
			conf.PreferredFactorType = factor.FactorType
			err = conf.Save()
			if err != nil {
				return err
			}
		}
	}

	fmt.Printf("Using MFA factor %s\n", factor.FactorType)

	switch factor.FactorType {
	case totpFactorType:
		return verifyTOTP(c, factor, auth, p)
	default:
		return fmt.Errorf("%s factor not implemented", factor.FactorType)
	}
}
