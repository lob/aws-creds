package okta

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/lob/aws-creds/pkg/input"
)

const (
	totpFactorType = "token:software:totp"
	smsFactorType  = "sms"
)

func promptForFactor(factors []*Factor, p input.Prompter) (int, error) {
	fmt.Println("Available MFA Factors:")
	for i, f := range factors {
		var id string
		switch f.FactorType {
		case totpFactorType:
			id = f.Profile.CredentialID
		case smsFactorType:
			id = f.Profile.PhoneNumber
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
	var resp io.Reader
	for {
		totp, err := p.Prompt("Enter Okta One-Time Password (TOTP): ")
		if err != nil {
			return err
		}
		resp, err = verifyAnswer(totp, a.StateToken, c, f)
		if err == nil {
			break
		}
		fmt.Println(err)
	}
	return json.NewDecoder(resp).Decode(a)
}

func verifySMS(c *Client, f *Factor, a *Auth, p input.Prompter) error {
	_, err := verifyAnswer("", a.StateToken, c, f)
	if err != nil {
		return err
	}
	var resp io.Reader
	for {
		code, err := p.Prompt("Enter SMS Code: ")
		if err != nil {
			return err
		}
		resp, err = verifyAnswer(code, a.StateToken, c, f)
		if err == nil {
			break
		}
		fmt.Println(err)
	}
	return json.NewDecoder(resp).Decode(a)
}

func verifyAnswer(answer string, stateToken string, c *Client, f *Factor) (io.Reader, error) {
	payload := []byte(fmt.Sprintf(`{"stateToken":"%s","answer":"%s"}`, stateToken, answer))
	u, err := url.Parse(f.Links.Verify.Href)
	if err != nil {
		return nil, err
	}
	return c.Post(u.Path, payload)
}
