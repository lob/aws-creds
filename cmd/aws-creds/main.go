package main

import (
	"os"

	"github.com/lob/aws-creds/pkg/cmd"
	"github.com/lob/aws-creds/pkg/input"
)

func main() {
	i := input.New(os.Stdin, os.Stdout)
	cmd.Execute(i)
}
