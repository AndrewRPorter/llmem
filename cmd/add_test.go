package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupInitializedDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer
	code := Run([]string{"init", dir}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("init failed: %s", stderr.String())
	}

	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(dir)
	return dir
}

func TestAddCreatesMemory(t *testing.T) {
	dir := setupInitializedDir(t)

	var stdout, stderr bytes.Buffer
	code := Run([]string{"add", "--name", "test memory", "--memory", "something happened", "--paths", "main.go"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d; stderr: %s", code, stderr.String())
	}

	if !strings.Contains(stdout.String(), "Added memory") {
		t.Errorf("expected 'Added memory' output, got %q", stdout.String())
	}

	data, err := os.ReadFile(filepath.Join(dir, ".llmem", "events.ndjson"))
	if err != nil {
		t.Fatalf("could not read events file: %v", err)
	}

	var m Memory
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("could not parse memory: %v", err)
	}

	if m.Name != "test memory" {
		t.Errorf("got name %q, want %q", m.Name, "test memory")
	}
	if m.Memory != "something happened" {
		t.Errorf("got memory %q, want %q", m.Memory, "something happened")
	}
	if len(m.Paths) != 1 || m.Paths[0] != "main.go" {
		t.Errorf("got paths %v, want [main.go]", m.Paths)
	}
	if m.ID == "" {
		t.Error("expected non-empty UUID")
	}
	if m.UpdatedAt == "" {
		t.Error("expected non-empty updated_at")
	}
}

func TestAddMultiplePaths(t *testing.T) {
	dir := setupInitializedDir(t)

	var stdout, stderr bytes.Buffer
	code := Run([]string{"add", "--name", "multi", "--memory", "edited files", "--paths", "a.go, b.go, c.go"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d; stderr: %s", code, stderr.String())
	}

	data, err := os.ReadFile(filepath.Join(dir, ".llmem", "events.ndjson"))
	if err != nil {
		t.Fatalf("could not read events file: %v", err)
	}

	var m Memory
	json.Unmarshal(data, &m)
	if len(m.Paths) != 3 {
		t.Errorf("got %d paths, want 3", len(m.Paths))
	}
}

func TestAddMultipleMemories(t *testing.T) {
	dir := setupInitializedDir(t)

	for i := 0; i < 3; i++ {
		var stdout, stderr bytes.Buffer
		Run([]string{"add", "--name", "mem", "--memory", "entry", "--paths", "f.go"}, &stdout, &stderr)
	}

	data, err := os.ReadFile(filepath.Join(dir, ".llmem", "events.ndjson"))
	if err != nil {
		t.Fatalf("could not read events file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}

func TestAddMissingRequiredFlags(t *testing.T) {
	setupInitializedDir(t)

	tests := []struct {
		name string
		args []string
	}{
		{"missing name", []string{"add", "--memory", "x", "--paths", "y"}},
		{"missing memory", []string{"add", "--name", "x", "--paths", "y"}},
		{"missing paths", []string{"add", "--name", "x", "--memory", "y"}},
		{"all missing", []string{"add"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := Run(tt.args, &stdout, &stderr)
			if code != 1 {
				t.Errorf("expected exit code 1, got %d", code)
			}
		})
	}
}

func TestAddWithoutInit(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	var stdout, stderr bytes.Buffer
	code := Run([]string{"add", "--name", "x", "--memory", "y", "--paths", "z"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if !strings.Contains(stderr.String(), "llmem init") {
		t.Errorf("expected hint to run init, got %q", stderr.String())
	}
}

func TestAddHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"add", "help"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !strings.Contains(stderr.String(), "Usage: llmem add") {
		t.Errorf("expected add usage output, got %q", stderr.String())
	}
}
