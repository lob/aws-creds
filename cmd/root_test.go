package cmd

import (
	"bufio"
	"os"
	"path"
	"strings"
	"testing"
)

func TestExecute(t *testing.T) {
	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()
	os.Args = []string{"cmd", "--help"}
	Execute()
}

func TestHelp(t *testing.T) {
	*help = true
	defer resetFlags()

	if err := execute([]string{}, os.Stdin, os.Stdout); err != nil {
		t.Errorf("unexpected error when executing help")
	}
}

func TestConfig(t *testing.T) {
	_ = execute([]string{}, os.Stdin, os.Stdout)
	if !strings.Contains(*configFilepath, defaultConfigFilepath) {
		t.Errorf("expected %s to contain %s", *configFilepath, defaultConfigFilepath)
	}

	*configFilepath = path.Join(os.TempDir(), "aws-creds", "config")
	defer resetFlags()

	if err := execute([]string{}, os.Stdin, os.Stdout); err == nil {
		t.Errorf("expected error when executing with invalid custom config")
	}

	reader := bufio.NewReader(strings.NewReader("username\norg\nprofile\narn\nn\n"))
	if err := execute([]string{"configure"}, reader, os.Stdout); err != nil {
		t.Errorf("unexpected error when configuring with invalid custom config: %s", err)
	}
}

func TestUnknownCommand(t *testing.T) {
	if err := execute([]string{"foo"}, os.Stdin, os.Stdout); err == nil {
		t.Errorf("expected error when executing unknown command")
	}
}

func resetFlags() {
	*help = false
	*configFilepath = ""
}
