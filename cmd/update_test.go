package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpdateAlreadyUpToDate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"tag_name": "v%s"}`, Version)
	}))
	defer server.Close()

	origAPI := repoAPI
	defer func() { setRepoAPI(origAPI) }()
	setRepoAPI(server.URL)

	var stdout, stderr bytes.Buffer
	code := Run([]string{"update"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d; stderr: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Already up to date") {
		t.Errorf("expected 'Already up to date', got %q", stdout.String())
	}
}

func TestUpdateHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"update", "help"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !strings.Contains(stderr.String(), "Usage: llmem update") {
		t.Errorf("expected update usage output, got %q", stderr.String())
	}
}

func TestUpdateAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	origAPI := repoAPI
	defer func() { setRepoAPI(origAPI) }()
	setRepoAPI(server.URL)

	var stdout, stderr bytes.Buffer
	code := Run([]string{"update"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if !strings.Contains(stderr.String(), "error checking for updates") {
		t.Errorf("expected error message, got %q", stderr.String())
	}
}
