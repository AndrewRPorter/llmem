package main

import (
	_ "embed"
	"os"
	"strings"

	"github.com/andrew/llmem/cmd"
)

//go:embed VERSION
var version string

func main() {
	cmd.Version = strings.TrimSpace(version)
	code := cmd.Run(os.Args[1:], os.Stdout, os.Stderr)
	if code != 0 {
		os.Exit(code)
	}
}
