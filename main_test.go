package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	origArgs := os.Args
	os.Args = []string{"cmd", "--help"}
	defer func() {
		os.Args = origArgs
	}()
	main()
}
