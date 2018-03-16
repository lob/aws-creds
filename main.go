package main

import (
	"os"

	"github.com/lob/aws-creds/cmd"
	"github.com/lob/aws-creds/input"
)

func main() {
	i := input.New(os.Stdin, os.Stdout)
	cmd.Execute(i)
}
