package okta

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/lob/aws-creds/input"
)

const (
	totpFactorType = "token:software:totp"
)

func promptForFactor(factors []*Factor, p input.Prompter) (int, error) {
	fmt.Println("Available MFA Factors:")
	for i, f := range factors {
		var id string
		switch f.FactorType {
		case totpFactorType:
			id = f.Profile.CredentialID
		default:
			id = "CURRENTLY UNSUPPORTED"
		}
		fmt.Printf("[ %d ] %s: %s\n", i, f.FactorType, id)
	}
	indexStr, err := p.Prompt("Preferred Factor: ")
	if err != nil {
		return 0, err
	}
	index, err := strconv.Atoi(indexStr)
	switch err.(type) {
	case *strconv.NumError:
		fmt.Println("Invalid selection, please select again")
		return promptForFactor(factors, p)
	case nil:
		if index < 0 || index >= len(factors) {
			fmt.Println("Invalid selection, please select again")
			return promptForFactor(factors, p)
		}
		return index, nil
	default:
		return 0, err
	}
}

func verifyTOTP(c *Client, f *Factor, a *Auth, p input.Prompter) error {
	totp, err := p.Prompt("Enter TOTP: ")
	if err != nil {
		return err
	}
	payload := []byte(fmt.Sprintf(`{"stateToken":"%s","answer":"%s"}`, a.StateToken, totp))
	u, err := url.Parse(f.Links.Verify.Href)
	if err != nil {
		return err
	}
	resp, err := c.Post(u.Path, payload)
	if err != nil {
		return err
	}
	return json.NewDecoder(resp).Decode(a)
}
