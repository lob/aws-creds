package cmd

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/lob/aws-creds/config"
	"github.com/lob/aws-creds/test"
)

func TestExecuteConfigure(t *testing.T) {
	path := path.Join(os.TempDir(), "aws-creds", "config")
	defer cleanup(t, path)
	conf := config.New(path)

	cmd := fakeCmd([]string{"test_user", exampleEmbedLink, "staging", "arn:staging", "n"}, conf)
	if err := executeConfigure(cmd); err != nil {
		t.Errorf("unexpected error when configuring with 1 profile: %s", err)
	}
	if conf.Username != "test_user" {
		t.Errorf("got %s, wanted %s", conf.Username, "test_user")
	}
	if len(conf.Profiles) != 1 {
		t.Errorf("got len(conf.Profiles) = %d, wanted %d", len(conf.Profiles), 1)
	}

	cmd = fakeCmd([]string{"test_user", exampleEmbedLink, "staging", "arn:staging", "y", "production", "arn:production", "n"}, conf)
	if err := executeConfigure(cmd); err != nil {
		t.Errorf("unexpected error when configuring with 2 profiles: %s", err)
	}
	if len(conf.Profiles) != 2 {
		t.Errorf("got len(conf.Profiles) = %d, wanted %d", len(conf.Profiles), 2)
	}

	cmd = fakeCmd([]string{"test_user", exampleEmbedLink, "sandbox", "arn:sandbox", "n"}, conf)
	if err := executeConfigure(cmd); err != nil {
		t.Errorf("unexpected error when configuring with an additional profile: %s", err)
	}
	if len(conf.Profiles) != 3 {
		t.Errorf("got len(conf.Profiles) = %d, wanted %d", len(conf.Profiles), 3)
	}

	cmd = fakeCmd([]string{"test_user", "invalid", "sandbox", "arn:sandbox", "n"}, conf)
	if err := executeConfigure(cmd); err == nil {
		t.Errorf("expected error when configuring with bad app URL")
	}
}

func fakeCmd(resp []string, conf *config.Config) *Cmd {
	fakeInput := test.NewArrayInput(resp)
	return &Cmd{
		Command: configureCommand,
		Config:  conf,
		Input:   fakeInput,
	}
}

func cleanup(t *testing.T, path string) {
	dir := filepath.Dir(path)
	if err := os.RemoveAll(dir); err != nil {
		t.Fatalf("unexpected error when cleaning up: %s", err)
	}
}
