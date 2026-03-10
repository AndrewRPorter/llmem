package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func cmdInit(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintf(stderr, "Usage: llmem init [directory]\n\n")
		fmt.Fprintf(stderr, "Initialize a new .llmem directory with an empty events file.\n")
		fmt.Fprintf(stderr, "Defaults to the current directory if no argument is given.\n")
	}
	if err := fs.Parse(args); err != nil {
		return 1
	}

	if fs.NArg() > 0 && fs.Arg(0) == "help" {
		fs.Usage()
		return 0
	}

	dir := "."
	if fs.NArg() > 0 {
		dir = fs.Arg(0)
	}

	llmemDir := filepath.Join(dir, ".llmem")
	eventsFile := filepath.Join(llmemDir, "events.ndjson")

	if err := os.MkdirAll(llmemDir, 0o755); err != nil {
		fmt.Fprintf(stderr, "error creating directory: %v\n", err)
		return 1
	}

	f, err := os.OpenFile(eventsFile, os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		if os.IsExist(err) {
			fmt.Fprintf(stderr, "already initialized: %s\n", llmemDir)
			return 1
		}
		fmt.Fprintf(stderr, "error creating events file: %v\n", err)
		return 1
	}
	f.Close()

	fmt.Fprintf(stdout, "Initialized %s\n", llmemDir)
	return 0
}
