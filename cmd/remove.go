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

func cmdRemove(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("remove", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var id string
	fs.StringVar(&id, "id", "", "UUID of the memory to remove")

	fs.Usage = func() {
		fmt.Fprintf(stderr, "Usage: llmem remove --id <uuid>\n\n")
		fmt.Fprintf(stderr, "Remove a memory by its UUID.\n\n")
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

	if id == "" {
		fmt.Fprintf(stderr, "error: --id is required\n")
		fs.Usage()
		return 1
	}

	eventsFile := filepath.Join(".llmem", "events.ndjson")
	f, err := os.Open(eventsFile)
	if err != nil {
		fmt.Fprintf(stderr, "error opening events file: %v\n", err)
		fmt.Fprintf(stderr, "have you run 'llmem init'?\n")
		return 1
	}

	var kept []string
	found := false
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		var m Memory
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			kept = append(kept, line)
			continue
		}
		if m.ID == id {
			found = true
			continue
		}
		kept = append(kept, line)
	}
	if err := scanner.Err(); err != nil {
		f.Close()
		fmt.Fprintf(stderr, "error reading events file: %v\n", err)
		return 1
	}
	f.Close()

	if !found {
		fmt.Fprintf(stderr, "memory not found: %s\n", id)
		return 1
	}

	var content string
	if len(kept) > 0 {
		content = strings.Join(kept, "\n") + "\n"
	}
	if err := os.WriteFile(eventsFile, []byte(content), 0o644); err != nil {
		fmt.Fprintf(stderr, "error writing events file: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "Removed memory %s\n", id)
	return 0
}
