package config

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

const badPermissions = 0200

func TestConfig(t *testing.T) {
	path := path.Join(os.TempDir(), "aws-creds", "config")
	defer cleanup(t, path)

	conf1 := New(path)
	conf1.Username = "test_user"
	conf1.OktaOrgURL = "https://test.okta.com"
	conf1.Profiles = []*Profile{
		{"staging", "arn:staging"},
		{"production", "arn:production"},
	}
	if err := conf1.Save(); err != nil {
		t.Fatalf("unexpected error when saving config: %s", err)
	}

	conf2 := New(path)
	if err := conf2.Load(); err != nil {
		t.Fatalf("unexpected error when loading config: %s", err)
	}

	cases := []struct {
		got, want string
	}{
		{conf2.Username, conf1.Username},
		{conf2.OktaOrgURL, conf1.OktaOrgURL},
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
	dir := filepath.Dir(path)
	defer cleanup(t, path)

	conf := New(path)

	if err := conf.Load(); err == nil {
		t.Errorf("expected error when loading non-existent config")
	}

	if err := os.MkdirAll(dir, directoryPermissions); err != nil {
		t.Fatalf("unexpected error when creating a directory: %s", err)
	}
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
	defer cleanup(t, dir)

	conf := New(path)

	if err := os.MkdirAll(filepath.Dir(dir), badPermissions); err != nil {
		t.Fatalf("unexpected error when making directory: %s", err)
	}
	if err := conf.Save(); err == nil {
		t.Errorf("expected error when saving config with bad permissions")
	}
}

func cleanup(t *testing.T, path string) {
	dir := filepath.Dir(path)
	if err := os.RemoveAll(dir); err != nil {
		t.Fatalf("unexpected error when cleaning up: %s", err)
	}
}
