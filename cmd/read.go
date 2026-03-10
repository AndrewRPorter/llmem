package cmd

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func cmdRead(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("read", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var path string
	var n int

	fs.StringVar(&path, "path", "", "Filter by file or directory path")
	fs.IntVar(&n, "n", 0, "Number of most recent memories to return (0 = all)")

	fs.Usage = func() {
		fmt.Fprintf(stderr, "Usage: llmem read [flags]\n\n")
		fmt.Fprintf(stderr, "Read memories from the events log.\n\n")
		fmt.Fprintf(stderr, "Examples:\n")
		fmt.Fprintf(stderr, "  llmem read                  Read all memories\n")
		fmt.Fprintf(stderr, "  llmem read --path cmd/       Filter by directory\n")
		fmt.Fprintf(stderr, "  llmem read --path main.go   Filter by file\n")
		fmt.Fprintf(stderr, "  llmem read -n 5             Last 5 memories\n")
		fmt.Fprintf(stderr, "  llmem read --path cmd/ -n 3 Last 3 for cmd/\n\n")
		fmt.Fprintf(stderr, "Flags:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if fs.NArg() > 0 && fs.Arg(0) == "help" {
		fs.Usage()
		return 0
	}

	eventsFile := filepath.Join(".llmem", "events.ndjson")
	f, err := os.Open(eventsFile)
	if err != nil {
		fmt.Fprintf(stderr, "error opening events file: %v\n", err)
		fmt.Fprintf(stderr, "have you run 'llmem init'?\n")
		return 1
	}
	defer f.Close()

	var memories []Memory
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		var m Memory
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			fmt.Fprintf(stderr, "warning: skipping malformed line: %v\n", err)
			continue
		}

		if path != "" {
			if !matchesPath(m.Paths, path) {
				continue
			}
		}

		memories = append(memories, m)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(stderr, "error reading events file: %v\n", err)
		return 1
	}

	if n > 0 && len(memories) > n {
		memories = memories[len(memories)-n:]
	}

	for _, m := range memories {
		line, _ := json.Marshal(m)
		fmt.Fprintln(stdout, string(line))
	}

	return 0
}

func matchesPath(memPaths []string, filter string) bool {
	for _, p := range memPaths {
		if p == filter {
			return true
		}
		// directory filter: "cmd/" matches "cmd/add.go"
		if strings.HasSuffix(filter, "/") && strings.HasPrefix(p, filter) {
			return true
		}
		// file filter: "cmd/add.go" matches directory "cmd/"
		if strings.HasSuffix(p, "/") && strings.HasPrefix(filter, p) {
			return true
		}
	}
	return false
}
