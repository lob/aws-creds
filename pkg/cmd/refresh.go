package cmd

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"sync"

	"github.com/lob/aws-creds/pkg/aws"
	"github.com/lob/aws-creds/pkg/config"
	"github.com/lob/aws-creds/pkg/okta"
	"github.com/zalando/go-keyring"
)

const (
	keyringPasswordService = "aws-creds Password" // nolint: gosec
	keyringSessionService  = "aws-creds Session Cookie"
)

func executeRefresh(cmd *Cmd) error {
	if len(cmd.Profiles) == 0 {
		return errors.New("you must provide at least one profile with '-p', e.g. 'aws-creds -p production'")
	}

	var profiles []*config.Profile

	for _, selectedProfile := range cmd.Profiles {
		found := false
		for _, availableProfile := range cmd.Config.Profiles {
			if selectedProfile == availableProfile.Name {
				profiles = append(profiles, availableProfile)
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("profile %s not configured", selectedProfile)
		}
	}

	var password string
	var prompted bool
	var sessionCookie string

	sessionCookie, err := getSessionCookie(cmd)
	if err != nil {
		return err
	}

	if sessionCookie == "" {
		password, prompted, err = getPassword(cmd)
		if err != nil {
			return err
		}
	}

	saml, newSessionCookie, err := okta.Login(cmd.Config, cmd.Input, sessionCookie, password)
	if err != nil {
		// Check if error is due to a network problem. If it is a network error,
		// return error but do not delete the password from the keyring.
		if urlError, ok := err.(*url.Error); ok {
			if _, ok := urlError.Err.(*net.OpError); ok {
				fmt.Println("Unable to connect with Okta. Make sure you are connected to the internet and try again.")
				return err
			}
		}
	}

	if sessionCookie != "" && err != nil && cmd.Config.EnableKeyring {
		if deleteErr := keyring.Delete(keyringSessionService, cmd.Config.Username); deleteErr != nil {
			return deleteErr
		}

		fmt.Println("Invalid session token deleted from system keyring.")

		password, prompted, err = getPassword(cmd)
		if err != nil {
			return err
		}

		saml, newSessionCookie, err = okta.Login(cmd.Config, cmd.Input, "", password)
	}

	if err != nil {
		if cmd.Config.EnableKeyring {
			deleteErr := keyring.Delete(keyringPasswordService, cmd.Config.Username)
			switch deleteErr {
			case nil:
				fmt.Println("Invalid password deleted from system keyring.")
				return err
			case keyring.ErrNotFound:
				return err
			}
			if deleteErr != nil && deleteErr != keyring.ErrNotFound {
				return deleteErr
			}
		}

		return err
	}

	if prompted {
		err = promptSavePassword(cmd, password)
		if err != nil {
			return err
		}
	}

	if cmd.Config.EnableKeyring {
		err = keyring.Set(keyringSessionService, cmd.Config.Username, newSessionCookie)
		if err != nil {
			return err
		}
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(profiles))

	for _, profile := range profiles {
		wg.Add(1)
		go renewCredentials(cmd, profile, saml, &wg, errCh)
	}

	wg.Wait()

	close(errCh)

	for e := range errCh {
		if e != nil {
			return e
		}
	}

	return nil
}

func renewCredentials(cmd *Cmd, profile *config.Profile, saml *okta.SAMLResponse, wg *sync.WaitGroup, errCh chan error) {
	defer wg.Done()

	creds, err := aws.GetCreds(cmd.STS, saml, profile)
	if err != nil {
		errCh <- fmt.Errorf("error renewing creds for \"%s\" profile: %w", profile.Name, err)
		return
	}

	errCh <- aws.WriteCreds(creds, profile, cmd.Config.CredentialsFilepath)
}

func getSessionCookie(cmd *Cmd) (string, error) {
	if !cmd.Config.EnableKeyring {
		return "", nil
	}

	cookie, err := keyring.Get(keyringSessionService, cmd.Config.Username)
	if err != nil && err != keyring.ErrNotFound {
		return "", err
	}
	if err == nil {
		return cookie, nil
	}
	return "", nil
}

func getPassword(cmd *Cmd) (password string, prompted bool, err error) {
	if cmd.Config.EnableKeyring {
		password, err = keyring.Get(keyringPasswordService, cmd.Config.Username)
		if err != nil && err != keyring.ErrNotFound {
			return "", false, err
		}
		if err == nil {
			fmt.Println("Password fetched from keyring.")
			return password, false, nil
		}
	}

	fmt.Println("Password not found in keyring.")
	password, err = promptPassword(cmd)
	return password, true, err
}

func promptPassword(cmd *Cmd) (string, error) {
	msg := fmt.Sprintf("Enter password for %s: ", cmd.Config.Username)
	password, err := cmd.Input.PromptPassword(msg)
	if err != nil {
		return "", err
	}
	return password, nil
}

func promptSavePassword(cmd *Cmd, password string) error {
	if !cmd.Config.EnableKeyring {
		return nil
	}

	save, err := cmd.Input.Prompt("Do you want to securely save your password in your system keyring? [y/N]: ")
	if err != nil {
		return err
	}
	if strings.HasPrefix(strings.ToLower(save), "y") {
		err := keyring.Set(keyringPasswordService, cmd.Config.Username, password)
		if err != nil {
			return err
		}
		fmt.Println("Password saved!")
	}
	return nil
}
