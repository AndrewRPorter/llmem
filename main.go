package main

import (
	"os"

	"github.com/andrew/llmem/cmd"
)

func main() {
	code := cmd.Run(os.Args[1:], os.Stdout, os.Stderr)
	if code != 0 {
		os.Exit(code)
	}
}
