package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestInitCreatesFiles(t *testing.T) {
	dir := t.TempDir()

	var stdout, stderr bytes.Buffer
	code := Run([]string{"init", dir}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d; stderr: %s", code, stderr.String())
	}

	llmemDir := filepath.Join(dir, ".llmem")
	eventsFile := filepath.Join(llmemDir, "events.ndjson")

	if _, err := os.Stat(llmemDir); os.IsNotExist(err) {
		t.Error(".llmem directory was not created")
	}
	if _, err := os.Stat(eventsFile); os.IsNotExist(err) {
		t.Error("events.ndjson was not created")
	}

	info, err := os.Stat(eventsFile)
	if err != nil {
		t.Fatalf("could not stat events file: %v", err)
	}
	if info.Size() != 0 {
		t.Errorf("events file should be empty, got %d bytes", info.Size())
	}
}

func TestInitAlreadyInitialized(t *testing.T) {
	dir := t.TempDir()

	var stdout, stderr bytes.Buffer
	code := Run([]string{"init", dir}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("first init failed: %s", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"init", dir}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("expected exit code 1 on duplicate init, got %d", code)
	}
	if !bytes.Contains(stderr.Bytes(), []byte("already initialized")) {
		t.Errorf("expected 'already initialized' error, got %q", stderr.String())
	}
}

func TestInitDefaultDir(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	var stdout, stderr bytes.Buffer
	code := Run([]string{"init"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d; stderr: %s", code, stderr.String())
	}

	if _, err := os.Stat(filepath.Join(dir, ".llmem", "events.ndjson")); os.IsNotExist(err) {
		t.Error("events.ndjson was not created in current directory")
	}
}

func TestInitHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"init", "help"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !bytes.Contains(stderr.Bytes(), []byte("Usage: llmem init")) {
		t.Errorf("expected init usage output, got %q", stderr.String())
	}
}
