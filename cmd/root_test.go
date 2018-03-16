package cmd

import (
	"os"
	"path"
	"strings"
	"testing"
)

type noopInput struct{}

func (i noopInput) Prompt(msg string) (string, error) {
	return msg, nil
}
func (i noopInput) PromptPassword(msg string) (string, error) {
	return msg, nil
}

func TestExecute(t *testing.T) {
	fakeInput := &noopInput{}
	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()
	os.Args = []string{"cmd", "-h"}
	Execute(fakeInput)
}

func TestHelp(t *testing.T) {
	fakeInput := &noopInput{}
	*help = true
	defer resetFlags()

	if err := execute([]string{}, fakeInput); err != nil {
		t.Errorf("unexpected error when executing help")
	}
}

func TestConfig(t *testing.T) {
	fakeInput := &noopInput{}

	_ = execute([]string{}, fakeInput)
	if !strings.Contains(*configFilepath, defaultConfigFilepath) {
		t.Errorf("expected %s to contain %s", *configFilepath, defaultConfigFilepath)
	}

	*configFilepath = path.Join(os.TempDir(), "aws-creds", "config")
	defer resetFlags()

	if err := execute([]string{}, fakeInput); err == nil {
		t.Errorf("expected error when executing with invalid custom config")
	}

	if err := execute([]string{configureCommand}, fakeInput); err != nil {
		t.Errorf("unexpected error when configuring with invalid custom config: %s", err)
	}
}

func TestProfile(t *testing.T) {
	fakeInput := &noopInput{}
	*configFilepath = path.Join(os.TempDir(), "aws-creds", "config")
	defer resetFlags()

	if err := execute([]string{refreshCommand}, fakeInput); err == nil {
		t.Errorf("expected error when executing refresh command without configuring")
	}

	if err := execute([]string{configureCommand}, fakeInput); err != nil {
		t.Fatalf("unexpected error when configuring for profile tests: %s", err)
	}

	*profile = "profile"

	if err := execute([]string{refreshCommand}, fakeInput); err == nil {
		t.Errorf("unexpected error when executing refresh command with a profile set")
	}
}
func TestUnknownCommand(t *testing.T) {
	fakeInput := &noopInput{}

	if err := execute([]string{"foo"}, fakeInput); err == nil {
		t.Errorf("expected error when executing unknown command")
	}
}

func resetFlags() {
	*help = false
	*configFilepath = defaultConfigFilepath
	*profile = ""
}
