package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const downloadBase = "https://github.com/AndrewRPorter/llmem/releases/download"

var repoAPI = "https://api.github.com/repos/AndrewRPorter/llmem/releases/latest"

func setRepoAPI(url string) {
	repoAPI = url
}

type githubRelease struct {
	TagName string `json:"tag_name"`
}

func cmdUpdate(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("update", flag.ContinueOnError)
	fs.SetOutput(stderr)

	fs.Usage = func() {
		fmt.Fprintf(stderr, "Usage: llmem update\n\n")
		fmt.Fprintf(stderr, "Update llmem to the latest version.\n")
	}

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if fs.NArg() > 0 && fs.Arg(0) == "help" {
		fs.Usage()
		return 0
	}

	latest, err := fetchLatestVersion()
	if err != nil {
		fmt.Fprintf(stderr, "error checking for updates: %v\n", err)
		return 1
	}

	current := strings.TrimPrefix(Version, "v")
	latest = strings.TrimPrefix(latest, "v")

	if current == latest {
		fmt.Fprintf(stdout, "Already up to date (v%s)\n", current)
		return 0
	}

	fmt.Fprintf(stdout, "Updating llmem v%s -> v%s...\n", current, latest)

	goos := runtime.GOOS
	goarch := runtime.GOARCH
	url := fmt.Sprintf("%s/v%s/llmem_%s_%s.tar.gz", downloadBase, latest, goos, goarch)

	execPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(stderr, "error finding current binary: %v\n", err)
		return 1
	}

	// Download and replace via shell to handle in-place replacement
	script := fmt.Sprintf(
		`set -e; TMPDIR=$(mktemp -d); trap 'rm -rf "$TMPDIR"' EXIT; `+
			`curl -fsSL "%s" -o "$TMPDIR/llmem.tar.gz"; `+
			`tar -xzf "$TMPDIR/llmem.tar.gz" -C "$TMPDIR"; `+
			`chmod +x "$TMPDIR/llmem"; `+
			`mv "$TMPDIR/llmem" "%s"`,
		url, execPath,
	)

	cmd := exec.Command("sh", "-c", script)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(stderr, "error updating binary: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "Updated to v%s\n", latest)
	return 0
}

func fetchLatestVersion() (string, error) {
	resp, err := http.Get(repoAPI)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	if release.TagName == "" {
		return "", fmt.Errorf("no tag found in latest release")
	}

	return release.TagName, nil
}
