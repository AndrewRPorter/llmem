# llmem

A CLI tool that tracks LLM interactions as NDJSON event logs — edits, mistakes, new files, and more — and lets agents retrieve relevant memories via tool calls.

## Install

```sh
go install github.com/andrew/llmem@latest
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

## Development

```sh
# Build the binary
make build

# Run it
./llmem
```

## License

MIT
