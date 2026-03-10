package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRemoveMemory(t *testing.T) {
	dir := setupInitializedDir(t)

	// Add a memory and capture its ID
	var stdout, stderr bytes.Buffer
	Run([]string{"add", "--name", "to-remove", "--memory", "delete me", "--paths", "a.go"}, &stdout, &stderr)

	data, _ := os.ReadFile(filepath.Join(dir, ".llmem", "events.ndjson"))
	var m Memory
	json.Unmarshal(data, &m)
	id := m.ID

	stdout.Reset()
	stderr.Reset()
	code := Run([]string{"remove", "--id", id}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d; stderr: %s", code, stderr.String())
	}

	if !strings.Contains(stdout.String(), "Removed memory") {
		t.Errorf("expected 'Removed memory' output, got %q", stdout.String())
	}

	// Verify file is now empty
	data, _ = os.ReadFile(filepath.Join(dir, ".llmem", "events.ndjson"))
	if strings.TrimSpace(string(data)) != "" {
		t.Errorf("expected empty events file, got %q", string(data))
	}
}

func TestRemoveKeepsOtherMemories(t *testing.T) {
	dir := setupInitializedDir(t)

	// Add 3 memories
	for _, name := range []string{"keep1", "remove-me", "keep2"} {
		var stdout, stderr bytes.Buffer
		Run([]string{"add", "--name", name, "--memory", "data", "--paths", "a.go"}, &stdout, &stderr)
	}

	// Find the ID of the middle one
	data, _ := os.ReadFile(filepath.Join(dir, ".llmem", "events.ndjson"))
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	var target Memory
	json.Unmarshal([]byte(lines[1]), &target)

	var stdout, stderr bytes.Buffer
	code := Run([]string{"remove", "--id", target.ID}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d; stderr: %s", code, stderr.String())
	}

	// Read remaining
	stdout.Reset()
	stderr.Reset()
	Run([]string{"read"}, &stdout, &stderr)
	remaining := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(remaining) != 2 {
		t.Errorf("expected 2 remaining memories, got %d", len(remaining))
	}

	for _, line := range remaining {
		var m Memory
		json.Unmarshal([]byte(line), &m)
		if m.Name == "remove-me" {
			t.Error("removed memory should not be present")
		}
	}
}

func TestRemoveNotFound(t *testing.T) {
	setupInitializedDir(t)

	var stdout, stderr bytes.Buffer
	code := Run([]string{"remove", "--id", "nonexistent-uuid"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if !strings.Contains(stderr.String(), "memory not found") {
		t.Errorf("expected 'memory not found' error, got %q", stderr.String())
	}
}

func TestRemoveMissingID(t *testing.T) {
	setupInitializedDir(t)

	var stdout, stderr bytes.Buffer
	code := Run([]string{"remove"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}

func TestRemoveWithoutInit(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	var stdout, stderr bytes.Buffer
	code := Run([]string{"remove", "--id", "some-id"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}

func TestRemoveHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"remove", "help"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !strings.Contains(stderr.String(), "Usage: llmem remove") {
		t.Errorf("expected remove usage output, got %q", stderr.String())
	}
}
