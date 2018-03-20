// Package test defines structs useful for testing.
package test

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"runtime"
	"testing"
)

// NoopInput implements input.Prompter but doesn't have any side-effects.
type NoopInput struct{}

// NewNoopInput creates a new NoopInput struct.
func NewNoopInput() *NoopInput {
	return &NoopInput{}
}

// Prompt just returns the inputted message.
func (i NoopInput) Prompt(msg string) (string, error) {
	return msg, nil
}

// PromptPassword just returns the inputted message.
func (i NoopInput) PromptPassword(msg string) (string, error) {
	return msg, nil
}

// ArrayInput implements input.Prompter where every subsequent call to Prompt returns a new response.
type ArrayInput struct {
	responses []string
	count     int
}

// NewArrayInput creates a new ArrayInput struct.
func NewArrayInput(responses []string) *ArrayInput {
	return &ArrayInput{responses, 0}
}

// Prompt returns a new response based on how many calls have been made previously.
func (i *ArrayInput) Prompt(msg string) (string, error) {
	resp := i.responses[i.count]
	i.count = i.count + 1
	return resp, nil
}

// PromptPassword just returns the inputted message.
func (i *ArrayInput) PromptPassword(msg string) (string, error) {
	return msg, nil
}

// LoadTestFile fetches the contents of the file from the testdata directory as a string.
func LoadTestFile(t *testing.T, name string) string {
	_, b, _, _ := runtime.Caller(0)
	projectDir := filepath.Dir(b)
	p := path.Join(projectDir, "..", "testdata", name)
	contents, err := ioutil.ReadFile(p)
	if err != nil {
		t.Fatalf("unexpected error when reading file %s: %s", p, err)
	}
	return string(contents)
}
