package cmd

import (
	"bufio"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lob/aws-creds/config"
)

func TestExecuteConfigure(t *testing.T) {
	path := path.Join(os.TempDir(), "aws-creds", "config")
	defer cleanup(t, path)
	conf := config.New(path)

	cmd := fakeCmd("test_user\ntest\nstaging\narn:staging\nn\n", conf)
	if err := executeConfigure(cmd); err != nil {
		t.Errorf("unexpected error when configuring with 1 profile: %s", err)
	}
	if conf.Username != "test_user" {
		t.Errorf("got %s, wanted %s", conf.Username, "test_user")
	}
	if len(conf.Profiles) != 1 {
		t.Errorf("got len(conf.Profiles) = %d, wanted %d", len(conf.Profiles), 1)
	}

	cmd = fakeCmd("test_user\ntest\nstaging\narn:staging\ny\nproduction\narn:production\nn\n", conf)
	if err := executeConfigure(cmd); err != nil {
		t.Errorf("unexpected error when configuring with 2 profiles: %s", err)
	}
	if len(conf.Profiles) != 2 {
		t.Errorf("got len(conf.Profiles) = %d, wanted %d", len(conf.Profiles), 2)
	}

	cmd = fakeCmd("test_user\ntest\nsandbox\narn:sandbox\nn\n", conf)
	if err := executeConfigure(cmd); err != nil {
		t.Errorf("unexpected error when configuring with an additional profile: %s", err)
	}
	if len(conf.Profiles) != 3 {
		t.Errorf("got len(conf.Profiles) = %d, wanted %d", len(conf.Profiles), 3)
	}
}

func fakeCmd(inStr string, conf *config.Config) *Cmd {
	in := bufio.NewReader(strings.NewReader(inStr))
	return &Cmd{
		Command: configureCommand,
		Config:  conf,
		In:      in,
		Out:     os.Stdout,
	}
}

func cleanup(t *testing.T, path string) {
	dir := filepath.Dir(path)
	if err := os.RemoveAll(dir); err != nil {
		t.Fatalf("unexpected error when cleaning up: %s", err)
	}
}
