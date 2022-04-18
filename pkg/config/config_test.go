package config

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lob/aws-creds/internal/test"
)

const badPermissions = 0200

func TestNew(t *testing.T) {
	conf, err := New("")
	if err != nil {
		t.Fatalf("unexpected error when creating config: %s", err)
	}
	if !strings.Contains(conf.CredentialsFilepath, "") {
		t.Errorf("expected %s to contain .aws/credentials", conf.CredentialsFilepath)
	}

	cfp := path.Join(os.TempDir(), "aws-creds", "config")
	defer test.Cleanup(t, cfp)
	err = os.Setenv(sharedCrendentialsFileEnv, cfp)
	if err != nil {
		t.Fatalf("unexpected error when setting environment variable %s: %s", sharedCrendentialsFileEnv, err)
	}
	conf, err = New("")
	if err != nil {
		t.Fatalf("unexpected error when creating config: %s", err)
	}
	if conf.CredentialsFilepath != cfp {
		t.Errorf("expected %s to contain .aws/credentials", conf.CredentialsFilepath)
	}
	err = os.Unsetenv(sharedCrendentialsFileEnv)
	if err != nil {
		t.Fatalf("unexpected error when unsetting environment variable %s: %s", sharedCrendentialsFileEnv, err)
	}
}

func TestConfig(t *testing.T) {
	path := path.Join(os.TempDir(), "aws-creds", "config")
	defer test.Cleanup(t, path)

	conf1, err := New(path)
	if err != nil {
		t.Fatalf("unexpected error when creating config: %s", err)
	}
	conf1.Username = "test_user"
	conf1.OktaHost = "https://test.okta.com"
	conf1.OktaAppPath = "/home/amazon_aws/0oa54k1gk2ukOJ9nGDt7/252"
	conf1.EnableKeyring = true
	conf1.Profiles = []*Profile{
		{"staging", "arn:staging", 3600},
		{"production", "arn:production", 3600},
	}
	if err := conf1.Save(); err != nil {
		t.Fatalf("unexpected error when saving config: %s", err)
	}

	conf2, err := New(path)
	if err != nil {
		t.Fatalf("unexpected error when creating config: %s", err)
	}
	if err := conf2.Load(); err != nil {
		t.Fatalf("unexpected error when loading config: %s", err)
	}

	cases := []struct {
		got, want string
	}{
		{conf2.Username, conf1.Username},
		{conf2.OktaHost, conf1.OktaHost},
		{conf2.OktaAppPath, conf1.OktaAppPath},
		{conf2.Profiles[0].Name, conf1.Profiles[0].Name},
		{conf2.Profiles[0].RoleARN, conf1.Profiles[0].RoleARN},
		{conf2.Profiles[1].Name, conf1.Profiles[1].Name},
		{conf2.Profiles[1].RoleARN, conf1.Profiles[1].RoleARN},
	}

	for _, tc := range cases {
		if tc.got != tc.want {
			t.Errorf("got %s, wanted %s", tc.got, tc.want)
		}
	}
}

func TestLoadErrors(t *testing.T) {
	path := path.Join(os.TempDir(), "aws-creds", "config")
	defer test.Cleanup(t, path)

	conf, err := New(path)
	if err != nil {
		t.Fatalf("unexpected error when creating config: %s", err)
	}

	if err := conf.Load(); err == nil {
		t.Errorf("expected error when loading non-existent config")
	}

	test.PrepTempFile(t, path)
	if err := ioutil.WriteFile(path, []byte(""), filePermissions); err != nil {
		t.Fatalf("unexpected error when writing file: %s", err)
	}
	if err := conf.Load(); err == nil {
		t.Errorf("expected error when loading empty config")
	}

	if err := os.Chmod(path, badPermissions); err != nil {
		t.Fatalf("unexpected error when changing permissions of file: %s", err)
	}
	if err := conf.Load(); err == nil {
		t.Errorf("expected error when loading config with bad permissions")
	}
}

func TestSaveErrors(t *testing.T) {
	path := path.Join(os.TempDir(), "parent", "aws-creds", "config")
	dir := filepath.Dir(path)
	defer test.Cleanup(t, dir)

	conf, err := New(path)
	if err != nil {
		t.Fatalf("unexpected error when creating config: %s", err)
	}

	if err := os.MkdirAll(filepath.Dir(dir), badPermissions); err != nil {
		t.Fatalf("unexpected error when making directory: %s", err)
	}
	if err := conf.Save(); err == nil {
		t.Errorf("expected error when saving config with bad permissions")
	}
}
