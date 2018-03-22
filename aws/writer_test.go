package aws

import (
	"os"
	"path"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/go-ini/ini"
	"github.com/lob/aws-creds/config"
	"github.com/lob/aws-creds/test"
)

func TestWriteCreds(t *testing.T) {
	name := "staging"
	creds := test.NewCredentials()
	profile := &config.Profile{Name: name, RoleARN: "staging:arn"}
	cfp := path.Join(os.TempDir(), ".aws", "credentials")
	test.PrepTempFile(t, cfp)
	defer test.Cleanup(t, cfp)

	err := WriteCreds(creds, profile, cfp)
	if err != nil {
		t.Fatalf("unexpected error when writing credentials: %s", err)
	}

	file, err := ini.Load(cfp)
	if err != nil {
		t.Fatalf("unexpected error when reading %s as an ini: %s", cfp, err)
	}
	section := file.Section(name)

	cases := []struct {
		key, want string
	}{
		{accessKeyID, *creds.AccessKeyId},
		{secretAccessKey, *creds.SecretAccessKey},
		{sessionToken, *creds.SessionToken},
	}

	for _, tc := range cases {
		key, err := section.GetKey(tc.key)
		if err != nil {
			t.Fatalf("unexpected error when getting key %s: %s", tc.key, err)
		}
		got := key.Value()
		if got != tc.want {
			t.Errorf("got %s, wanted %s", got, tc.want)
		}
	}

	creds.AccessKeyId = aws.String("existing credentials file")
	err = WriteCreds(creds, profile, cfp)
	if err != nil {
		t.Fatalf("unexpected error when writing credentials to an existing file: %s", err)
	}

	file, err = ini.Load(cfp)
	if err != nil {
		t.Fatalf("unexpected error when reading %s as an ini: %s", cfp, err)
	}
	section = file.Section(name)
	key, err := section.GetKey(accessKeyID)
	if err != nil {
		t.Fatalf("unexpected error when getting key %s: %s", accessKeyID, err)
	}
	got := key.Value()
	if got != *creds.AccessKeyId {
		t.Errorf("got %s, wanted %s", got, *creds.AccessKeyId)
	}
}
