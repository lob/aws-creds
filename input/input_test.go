package input

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
)

type errWriter struct{}

func (w errWriter) Write(p []byte) (int, error) {
	return 0, errors.New("err")
}

func TestPrompt(t *testing.T) {
	var bin bytes.Buffer
	var bout bytes.Buffer
	i := New(&bin, &bout)
	wantIn := "Testing"
	wantOut := "Prompt: "

	bin.WriteString(fmt.Sprintf("%s\n", wantIn))
	gotIn, err := i.Prompt(wantOut)
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
	i := New(&bin, &bout)

	_, err := i.Prompt("Prompt: ")
	if err == nil {
		t.Errorf("expected error when prompting with empty input")
	}

	i = New(&bin, errWriter{})
	_, err = i.Prompt("Prompt: ")
	if err == nil {
		t.Errorf("expected error when prompting with error writer")
	}
}

func TestPromptPassword(t *testing.T) {
	var bin bytes.Buffer
	var bout bytes.Buffer
	i := New(&bin, &bout)
	wantOut := "Prompt: "

	_, err := i.PromptPassword(wantOut)
	if err == nil {
		t.Fatalf("expected error when prompting for password with stdin: %s", err)
	}

	gotOut := bout.String()
	if gotOut != wantOut {
		t.Errorf("got %s, wanted %s", gotOut, wantOut)
	}
}

func TestPromptPasswordErrors(t *testing.T) {
	var bin bytes.Buffer
	i := New(&bin, errWriter{})

	_, err := i.PromptPassword("Prompt: ")
	if err == nil {
		t.Errorf("expected error when prompting with error writer")
	}
}
