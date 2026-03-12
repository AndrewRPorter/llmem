# llmem

A CLI tool that tracks LLM interactions as NDJSON event logs — edits, mistakes, new files, and more — and lets agents retrieve relevant memories via tool calls.

## Install

```sh
curl -fsSL https://raw.githubusercontent.com/AndrewRPorter/llmem/main/install.sh | bash
```

## Usage

```sh
# Initialize a .llmem directory in the current folder
llmem init

# Initialize in a specific directory
llmem init /path/to/project

# Print the version
llmem version
```

## Integration

### Agent Skills

For an example skill that teaches agents when and how to save/read memories, see [`.github/skills/save-memory/SKILL.md`](.github/skills/save-memory/SKILL.md).

### MCP

In an MCP file like `.vscode/mcp.json` add the llmem tool:

```json
{
    "servers": {
        "llmem": {
            "command": "llmem",
            "args": ["mcp"]
        }
    }
}
```

## Development

```sh
# Build the binary
make build

# Run it
./llmem
```

## Example Memories

```json
{"id":"0d707e9b-785a-40db-9dae-a8aae6e48bea","paths":["cmd/add.go","cmd/read.go","cmd/remove.go"],"name":"Data model and storage format","memory":"Memories are stored as append-only NDJSON in .llmem/events.ndjson. Each line is a JSON object with fields: id (UUID v4), paths ([]string), name (string), memory (string), updated_at (RFC3339 UTC). Adding appends a line; removing reads all lines, filters out the target, and rewrites the entire file. There is no update command — modify by remove + add. Path filtering supports exact match, directory prefix matching (trailing slash), and reverse directory matching.","updated_at":"2026-03-11T23:49:33Z"}
{"id":"8c1aad9a-80d0-430f-b322-059f8ef2d6c0","paths":["main.go","VERSION","cmd/mcp.go"],"name":"Version is embedded at compile time via go:embed","memory":"VERSION file is embedded at compile time via go:embed in main.go and set on cmd.Version. The MCP server in cmd/mcp.go uses this version when creating the server. When releasing, update the VERSION file — the binary reads it at build time, not runtime.","updated_at":"2026-03-11T23:49:58Z"}
```
