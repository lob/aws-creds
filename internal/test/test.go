// Package test defines structs useful for testing.
package test

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

const directoryPermissions = 0700

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
	i.count++
	return resp, nil
}

// PromptPassword just returns the inputted message.
func (i *ArrayInput) PromptPassword(msg string) (string, error) {
	resp := i.responses[i.count]
	i.count++
	return resp, nil
}

// MockSTS is a mock for the AWS STS service.
type MockSTS struct {
	stsiface.STSAPI
	Creds    *sts.Credentials
	Duration int64
}

// AssumeRoleWithSAML takes in an AssumeRoleWithSAMLInput and returns AssumeRoleWithSAMLOutput.
func (m *MockSTS) AssumeRoleWithSAML(in *sts.AssumeRoleWithSAMLInput) (*sts.AssumeRoleWithSAMLOutput, error) {
	m.Duration = *in.DurationSeconds
	return &sts.AssumeRoleWithSAMLOutput{Credentials: m.Creds}, nil
}

// Cleanup deletes the directory the given file is in. Usually meant to be used in a deferred call.
func Cleanup(t *testing.T, path string) {
	dir := filepath.Dir(path)
	if err := os.RemoveAll(dir); err != nil {
		t.Fatalf("unexpected error when cleaning up: %s", err)
	}
}

// PrepTempFile takes the directory of the given temp file and ensures that it exists.
func PrepTempFile(t *testing.T, p string) {
	dir := filepath.Dir(p)
	err := os.MkdirAll(dir, directoryPermissions)
	if err != nil {
		t.Fatalf("unexpected error when making %s: %s", dir, err)
	}
}

// LoadTestFile fetches the contents of the file from the testdata directory as a string.
func LoadTestFile(t *testing.T, name string) string {
	_, b, _, _ := runtime.Caller(0)
	projectDir := filepath.Dir(b)
	p := path.Join(projectDir, "..", "testdata", name)
	contents, err := ioutil.ReadFile(p) // nolint: gosec
	if err != nil {
		t.Fatalf("unexpected error when reading file %s: %s", p, err)
	}
	return string(contents)
}

// NewCredentials creates a new Credentials struct with populated values.
func NewCredentials() *sts.Credentials {
	return &sts.Credentials{
		AccessKeyId:     aws.String("AccessKeyId"),
		SecretAccessKey: aws.String("SecretAccessKey"),
		SessionToken:    aws.String("SessionToken"),
	}
}
