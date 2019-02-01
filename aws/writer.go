package aws

import (
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/go-ini/ini"
	"github.com/lob/aws-creds/config"
)

const (
	accessKeyID     = "aws_access_key_id"
	secretAccessKey = "aws_secret_access_key" // nolint: gosec
	sessionToken    = "aws_session_token"
)

var mutex = &sync.Mutex{}

// WriteCreds takes AWS credentials and writes it to ~/.aws/credentials.
func WriteCreds(creds *sts.Credentials, profile *config.Profile, filepath string) error {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := ini.Load(filepath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		file = ini.Empty()
	}

	section := file.Section(profile.Name)
	for _, k := range section.Keys() {
		section.DeleteKey(k.Name())
	}

	_, err = section.NewKey(accessKeyID, *creds.AccessKeyId)
	if err != nil {
		return err
	}
	_, err = section.NewKey(secretAccessKey, *creds.SecretAccessKey)
	if err != nil {
		return err
	}
	_, err = section.NewKey(sessionToken, *creds.SessionToken)
	if err != nil {
		return err
	}

	err = file.SaveTo(filepath)
	if err != nil {
		return err
	}
	fmt.Printf("Credentials for %s written to %s.\n", profile.Name, filepath)
	return nil
}
