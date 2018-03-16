package cmd

import (
	"bufio"
	"fmt"
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
	os.Args = []string{"cmd", "-h"}
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

	reader := bufio.NewReader(strings.NewReader(fmt.Sprintf("username\n%s\nprofile\narn\nn\n", exampleEmbedLink)))
	defer cleanup(t, *configFilepath)
	if err := execute([]string{configureCommand}, reader, os.Stdout); err != nil {
		t.Errorf("unexpected error when configuring with invalid custom config: %s", err)
	}
}

func TestProfile(t *testing.T) {
	*configFilepath = path.Join(os.TempDir(), "aws-creds", "config")
	defer resetFlags()

	if err := execute([]string{refreshCommand}, os.Stdin, os.Stdout); err == nil {
		t.Errorf("expected error when executing refresh command without configuring")
	}

	reader := bufio.NewReader(strings.NewReader(fmt.Sprintf("username\n%s\nprofile\narn\nn\n", exampleEmbedLink)))
	defer cleanup(t, *configFilepath)
	if err := execute([]string{configureCommand}, reader, os.Stdout); err != nil {
		t.Fatalf("unexpected error when configuring for profile tests: %s", err)
	}

	if err := execute([]string{refreshCommand}, os.Stdin, os.Stdout); err == nil {
		t.Errorf("expected error when executing refresh command without a profile set")
	}

	*profile = "profile"

	if err := execute([]string{refreshCommand}, os.Stdin, os.Stdout); err == nil {
		t.Errorf("unexpected error when executing refresh command with a profile set")
	}
}
func TestUnknownCommand(t *testing.T) {
	if err := execute([]string{"foo"}, os.Stdin, os.Stdout); err == nil {
		t.Errorf("expected error when executing unknown command")
	}
}

func resetFlags() {
	*help = false
	*configFilepath = defaultConfigFilepath
	*profile = ""
}
