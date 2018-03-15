package okta

import (
	"encoding/json"
	"fmt"
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

// Factor contains inforamation about a specific MFA factor.
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
