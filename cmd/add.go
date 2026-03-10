package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Memory struct {
	ID        string   `json:"id"`
	Paths     []string `json:"paths"`
	Name      string   `json:"name"`
	Memory    string   `json:"memory"`
	UpdatedAt string   `json:"updated_at"`
}

func cmdAdd(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var paths string
	var name string
	var memory string

	fs.StringVar(&paths, "paths", "", "Comma-separated file or directory paths")
	fs.StringVar(&name, "name", "", "Name of the memory")
	fs.StringVar(&memory, "memory", "", "What to record")

	fs.Usage = func() {
		fmt.Fprintf(stderr, "Usage: llmem add [flags]\n\n")
		fmt.Fprintf(stderr, "Add a new memory to the events log.\n\n")
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

	if name == "" || memory == "" || paths == "" {
		fmt.Fprintf(stderr, "error: --name, --memory, and --paths are required\n")
		fs.Usage()
		return 1
	}

	pathList := strings.Split(paths, ",")
	for i := range pathList {
		pathList[i] = strings.TrimSpace(pathList[i])
	}

	m := Memory{
		ID:        uuid.New().String(),
		Paths:     pathList,
		Name:      name,
		Memory:    memory,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	eventsFile := filepath.Join(".llmem", "events.ndjson")
	f, err := os.OpenFile(eventsFile, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Fprintf(stderr, "error opening events file: %v\n", err)
		fmt.Fprintf(stderr, "have you run 'llmem init'?\n")
		return 1
	}
	defer f.Close()

	line, err := json.Marshal(m)
	if err != nil {
		fmt.Fprintf(stderr, "error encoding memory: %v\n", err)
		return 1
	}
	line = append(line, '\n')

	if _, err := f.Write(line); err != nil {
		fmt.Fprintf(stderr, "error writing memory: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "Added memory %s\n", m.ID)
	return 0
}
