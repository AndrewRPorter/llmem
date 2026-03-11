# llmem

A CLI tool that tracks LLM interactions as NDJSON event logs — edits, mistakes, new files, and more — and lets agents retrieve relevant memories via tool calls.

## Install

```sh
curl -fsSL https://raw.githubusercontent.com/AndrewRPorter/llmem/main/install.sh | bash
```

## CLI Usage

```sh
# Initialize a .llmem directory in the current folder
llmem init

# Initialize in a specific directory
llmem init /path/to/project

# Print the version
llmem version
```

## LLM Integration

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

## License

MIT
