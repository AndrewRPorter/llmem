package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func seedMemories(t *testing.T, dir string, memories []struct{ name, memory, paths string }) {
	t.Helper()
	for _, m := range memories {
		var stdout, stderr bytes.Buffer
		code := Run([]string{"add", "--name", m.name, "--memory", m.memory, "--paths", m.paths}, &stdout, &stderr)
		if code != 0 {
			t.Fatalf("seed add failed: %s", stderr.String())
		}
	}
}

func TestReadAll(t *testing.T) {
	dir := setupInitializedDir(t)
	_ = dir

	seedMemories(t, dir, []struct{ name, memory, paths string }{
		{"mem1", "first", "a.go"},
		{"mem2", "second", "b.go"},
		{"mem3", "third", "c.go"},
	})

	var stdout, stderr bytes.Buffer
	code := Run([]string{"read"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d; stderr: %s", code, stderr.String())
	}

	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 memories, got %d", len(lines))
	}
}

func TestReadFilterByFile(t *testing.T) {
	dir := setupInitializedDir(t)
	_ = dir

	seedMemories(t, dir, []struct{ name, memory, paths string }{
		{"mem1", "first", "a.go"},
		{"mem2", "second", "b.go"},
		{"mem3", "third", "a.go"},
	})

	var stdout, stderr bytes.Buffer
	code := Run([]string{"read", "--path", "a.go"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d; stderr: %s", code, stderr.String())
	}

	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 memories for a.go, got %d", len(lines))
	}

	for _, line := range lines {
		var m Memory
		json.Unmarshal([]byte(line), &m)
		found := false
		for _, p := range m.Paths {
			if p == "a.go" {
				found = true
			}
		}
		if !found {
			t.Errorf("expected path a.go in %v", m.Paths)
		}
	}
}

func TestReadFilterByDirectory(t *testing.T) {
	dir := setupInitializedDir(t)
	_ = dir

	seedMemories(t, dir, []struct{ name, memory, paths string }{
		{"mem1", "first", "cmd/add.go"},
		{"mem2", "second", "main.go"},
		{"mem3", "third", "cmd/init.go"},
	})

	var stdout, stderr bytes.Buffer
	code := Run([]string{"read", "--path", "cmd/"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d; stderr: %s", code, stderr.String())
	}

	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 memories for cmd/, got %d", len(lines))
	}
}

func TestReadFilterByFileMatchesDirectory(t *testing.T) {
	dir := setupInitializedDir(t)
	_ = dir

	seedMemories(t, dir, []struct{ name, memory, paths string }{
		{"mem1", "first", "cmd/"},
		{"mem2", "second", "main.go"},
	})

	var stdout, stderr bytes.Buffer
	code := Run([]string{"read", "--path", "cmd/add.go"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d; stderr: %s", code, stderr.String())
	}

	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) != 1 {
		t.Errorf("expected 1 memory matching cmd/add.go via cmd/, got %d", len(lines))
	}
}

func TestReadWithN(t *testing.T) {
	dir := setupInitializedDir(t)
	_ = dir

	seedMemories(t, dir, []struct{ name, memory, paths string }{
		{"mem1", "first", "a.go"},
		{"mem2", "second", "b.go"},
		{"mem3", "third", "c.go"},
		{"mem4", "fourth", "d.go"},
		{"mem5", "fifth", "e.go"},
	})

	var stdout, stderr bytes.Buffer
	code := Run([]string{"read", "-n", "2"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d; stderr: %s", code, stderr.String())
	}

	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 memories, got %d", len(lines))
	}

	// Should be the last 2
	var m Memory
	json.Unmarshal([]byte(lines[0]), &m)
	if m.Name != "mem4" {
		t.Errorf("expected mem4, got %s", m.Name)
	}
	json.Unmarshal([]byte(lines[1]), &m)
	if m.Name != "mem5" {
		t.Errorf("expected mem5, got %s", m.Name)
	}
}

func TestReadWithPathAndN(t *testing.T) {
	dir := setupInitializedDir(t)
	_ = dir

	seedMemories(t, dir, []struct{ name, memory, paths string }{
		{"mem1", "first", "cmd/a.go"},
		{"mem2", "second", "main.go"},
		{"mem3", "third", "cmd/b.go"},
		{"mem4", "fourth", "cmd/c.go"},
	})

	var stdout, stderr bytes.Buffer
	code := Run([]string{"read", "--path", "cmd/", "-n", "1"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d; stderr: %s", code, stderr.String())
	}

	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) != 1 {
		t.Errorf("expected 1 memory, got %d", len(lines))
	}

	var m Memory
	json.Unmarshal([]byte(lines[0]), &m)
	if m.Name != "mem4" {
		t.Errorf("expected mem4 (last cmd/ match), got %s", m.Name)
	}
}

func TestReadEmpty(t *testing.T) {
	setupInitializedDir(t)

	var stdout, stderr bytes.Buffer
	code := Run([]string{"read"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d; stderr: %s", code, stderr.String())
	}

	if strings.TrimSpace(stdout.String()) != "" {
		t.Errorf("expected empty output, got %q", stdout.String())
	}
}

func TestReadWithoutInit(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	var stdout, stderr bytes.Buffer
	code := Run([]string{"read"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}

func TestReadHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"read", "help"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !strings.Contains(stderr.String(), "Usage: llmem read") {
		t.Errorf("expected read usage output, got %q", stderr.String())
	}
}
