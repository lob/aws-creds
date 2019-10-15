package input

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// Prompter defines methods used for getting user input.
type Prompter interface {
	Prompt(string) (string, error)
	PromptPassword(string) (string, error)
}

// Input implements Prompter using the given io.Reader and io.Writer.
type Input struct {
	in  io.Reader
	out io.Writer
}

// New creates a new Input struct with the given io.Reader and io.Writer.
func New(in io.Reader, out io.Writer) *Input {
	return &Input{in, out}
}

// Prompt prompts the user for input with the given message, using the provided io Reader and Writer.
func (i *Input) Prompt(msg string) (string, error) {
	_, err := fmt.Fprint(i.out, msg)
	if err != nil {
		return "", err
	}
	reader := bufio.NewReader(i.in)
	resp, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(resp), nil
}

// PromptPassword prompts the user for input that won't be printed back to them.
func (i *Input) PromptPassword(msg string) (string, error) {
	_, err := fmt.Fprint(i.out, msg)
	if err != nil {
		return "", err
	}
	bytePassword, err := terminal.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", err
	}
	fmt.Print("\n")
	return string(bytePassword), nil
}
