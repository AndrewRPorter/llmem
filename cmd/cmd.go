package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var Version = "dev"

func Run(args []string, stdout, stderr io.Writer) int {
	flag.Usage = func() {
		fmt.Fprintf(stderr, "llmem - LLM memory tracker\n\n")
		fmt.Fprintf(stderr, "Usage: llmem <command> [flags]\n\n")
		fmt.Fprintf(stderr, "Commands:\n")
		fmt.Fprintf(stderr, "  add        Add a new memory\n")
		fmt.Fprintf(stderr, "  help       Show this help message\n")
		fmt.Fprintf(stderr, "  init       Initialize a new .llmem directory\n")
		fmt.Fprintf(stderr, "  read       Read memories\n")
		fmt.Fprintf(stderr, "  remove     Remove a memory by ID\n")
		fmt.Fprintf(stderr, "  update     Update llmem to latest version\n")
		fmt.Fprintf(stderr, "  version    Print the version\n")
	}

	if len(args) < 1 {
		flag.Usage()
		return 1
	}

	switch args[0] {
	case "version":
		fmt.Fprintln(stdout, "llmem "+Version)
		return 0
	case "init":
		return cmdInit(args[1:], stdout, stderr)
	case "add":
		return cmdAdd(args[1:], stdout, stderr)
	case "read":
		return cmdRead(args[1:], stdout, stderr)
	case "remove":
		return cmdRemove(args[1:], stdout, stderr)
	case "update":
		return cmdUpdate(args[1:], stdout, stderr)
	case "help":
		flag.Usage()
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n", args[0])
		flag.Usage()
		return 1
	}
}

func realRun() {
	code := Run(os.Args[1:], os.Stdout, os.Stderr)
	if code != 0 {
		os.Exit(code)
	}
}
