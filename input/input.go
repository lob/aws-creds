package input

import (
	"bufio"
	"fmt"
	"io"
	"strings"
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
