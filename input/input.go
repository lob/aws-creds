package input

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

// Prompt prompts the user for input with the given message, using the provided io Reader and Writer.
func Prompt(msg string, in io.Reader, out io.Writer) (string, error) {
	_, err := fmt.Fprint(out, msg)
	if err != nil {
		return "", err
	}
	reader := bufio.NewReader(in)
	resp, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(resp), nil
}

// PromptPassword prompts the user for input that won't be printed back to them.
func PromptPassword(message string, in int, out io.Writer) (string, error) {
	_, err := fmt.Fprint(out, message)
	if err != nil {
		return "", err
	}
	bytePassword, err := terminal.ReadPassword(in)
	if err != nil {
		return "", err
	}
	fmt.Print("\n")
	return string(bytePassword), nil
}
