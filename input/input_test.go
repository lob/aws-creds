package input

import (
	"bytes"
	"errors"
	"fmt"
	"syscall"
	"testing"
)

type errWriter struct{}

func (w errWriter) Write(p []byte) (int, error) {
	return 0, errors.New("err")
}

func TestPrompt(t *testing.T) {
	var bin bytes.Buffer
	var bout bytes.Buffer
	wantIn := "Testing"
	wantOut := "Prompt: "

	bin.WriteString(fmt.Sprintf("%s\n", wantIn))
	gotIn, err := Prompt(wantOut, &bin, &bout)
	if err != nil {
		t.Fatalf("unexpected error when prompting for input: %s", err)
	}

	if gotIn != wantIn {
		t.Errorf("got %s, wanted %s", gotIn, wantIn)
	}
	gotOut := bout.String()
	if gotOut != wantOut {
		t.Errorf("got %s, wanted %s", gotOut, wantOut)
	}
}

func TestPromptErrors(t *testing.T) {
	var bin bytes.Buffer
	var bout bytes.Buffer

	_, err := Prompt("Prompt: ", &bin, &bout)
	if err == nil {
		t.Errorf("expected error when prompting with empty input")
	}

	_, err = Prompt("Prompt: ", &bin, errWriter{})
	if err == nil {
		t.Errorf("expected error when prompting with error writer")
	}
}

func TestPromptPassword(t *testing.T) {
	var bout bytes.Buffer
	wantOut := "Prompt: "

	_, err := PromptPassword(wantOut, syscall.Stdin, &bout)
	if err == nil {
		t.Fatalf("expected error when prompting for password with stdin: %s", err)
	}

	gotOut := bout.String()
	if gotOut != wantOut {
		t.Errorf("got %s, wanted %s", gotOut, wantOut)
	}
}
