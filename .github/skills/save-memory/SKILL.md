---
name: save-memory
description: "Save and recall persistent memories using llmem MCP tools. USE WHEN: you learn something new about the codebase, make a mistake and discover the fix, find a non-obvious pattern or convention, encounter a gotcha or edge case, complete a multi-step task worth documenting, or the user asks you to remember something. ALSO USE WHEN: starting work on a file or directory to read relevant prior memories first. DO NOT USE FOR: trivial facts, temporary debugging notes, or information already in code comments."
---

# Save Memory

## Purpose

Use llmem to build a persistent knowledge base of lessons learned, mistakes, conventions, and discoveries. This helps you avoid repeating mistakes, recall context across conversations, and accumulate project-specific knowledge over time.

## When to Save a Memory

**Always save a memory when you:**

- **Make a mistake and fix it** — Record what went wrong, why, and the correct approach so you don't repeat it
- **Discover a non-obvious pattern** — Conventions, implicit rules, or "the way things work here" that aren't documented in code
- **Learn something new about the codebase** — Architecture decisions, data flow, dependencies between components, quirks
- **Find a gotcha or edge case** — Something that surprised you or would trip up someone unfamiliar
- **Complete a significant task** — Summarize what was done and key decisions made, especially for multi-step work
- **The user tells you to remember something** — Any explicit request to store information
- **Discover a user preference** — Coding style, tool preferences, workflow habits, communication preferences
- **Resolve a confusing error** — The error message, root cause, and fix for future reference
- **Identify a useful debugging technique** — Steps that helped diagnose an issue in this specific project

**Do NOT save memories for:**

- Information already clearly documented in code comments or README
- Trivial one-off facts with no future value
- Temporary debugging state that won't matter later

## When to Read Memories

**Before starting work**, read relevant memories to benefit from past experience:

- Starting work on a file → read memories filtered to that file path
- Starting work on a directory/module → read memories filtered to that directory
- Beginning a new task → read recent memories for general context
- Encountering an error → check if a memory about it already exists

## Procedure

### Reading Memories

Before diving into a task, check for relevant prior knowledge:

1. **For a specific file**: Use `mcp_llmem_read` with `path` set to the file (e.g., `"cmd/add.go"`)
2. **For a directory/module**: Use `mcp_llmem_read` with `path` set to the directory with trailing slash (e.g., `"cmd/"`)
3. **For recent context**: Use `mcp_llmem_read` with `n` set to 5-10 for the most recent memories
4. **For everything**: Use `mcp_llmem_read` with no parameters

### Saving a Memory

When something worth remembering happens:

1. **Choose a clear, descriptive name** — Short but specific (e.g., "Binary was outdated causing MCP failure", not "fix")
2. **Write the memory content** — Include:
    - What happened or what you learned
    - Why it matters or what to do differently
    - The root cause if it was a bug/mistake
3. **Associate the right paths** — Comma-separated file or directory paths this memory relates to. Use paths relative to the project root. Be specific (prefer `"cmd/mcp.go"` over `"cmd/"`) but include broader paths when the lesson applies broadly.

Use `mcp_llmem_add` with:

- `name`: Concise title of the memory (e.g., "NDJSON append requires newline separator")
- `memory`: Detailed content explaining what was learned and why it matters
- `paths`: Comma-separated relevant file paths (e.g., `"cmd/add.go, cmd/read.go"`)

### Removing a Memory

If a memory is outdated, incorrect, or no longer relevant:

1. Use `mcp_llmem_read` to find the memory and its UUID
2. Use `mcp_llmem_remove` with the `id` parameter set to the UUID

## Writing Good Memories

**Good memory example:**

- Name: `"Installed binary must be rebuilt after adding new CLI commands"`
- Memory: `"When adding a new subcommand (like 'mcp'), the installed binary at /usr/local/bin/llmem must be rebuilt and reinstalled. The old binary won't recognize new commands, causing 'unknown command' errors. Always run 'go build -o llmem . && cp llmem /usr/local/bin/llmem' after adding commands."`
- Paths: `"cmd/cmd.go, Makefile"`

**What makes it good:**

- Title describes the lesson, not just the symptom
- Content explains the _why_, not just the _what_
- Includes the fix/solution
- Paths point to the relevant files

**Bad memory example:**

- Name: `"bug fix"`
- Memory: `"fixed a bug"`
- Paths: `"."`

## Categories of Memories Worth Saving

| Category          | Example                                                                       |
| ----------------- | ----------------------------------------------------------------------------- |
| Mistake & fix     | "Off-by-one in slice when limiting with -n flag"                              |
| Convention        | "All commands accept io.Writer for testability, never use os.Stdout directly" |
| Architecture      | "Memories stored as append-only NDJSON; remove rewrites entire file"          |
| User preference   | "User prefers concise PR descriptions with bullet points"                     |
| Debugging insight | "MCP server hangs if stderr output is too large"                              |
| Build/deploy      | "VERSION file is embedded at compile time via go:embed"                       |
| Edge case         | "Path filter with trailing slash matches all files in directory"              |
