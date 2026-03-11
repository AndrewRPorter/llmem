package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func cmdMCP(args []string, stdout, stderr io.Writer) int {
	s := server.NewMCPServer("llmem", Version)

	addTool := mcp.NewTool("add",
		mcp.WithDescription("Add a new memory to the events log. Requires a name, the memory text, and associated file/directory paths."),
		mcp.WithString("name", mcp.Required(), mcp.Description("Name of the memory")),
		mcp.WithString("memory", mcp.Required(), mcp.Description("What to record")),
		mcp.WithString("paths", mcp.Required(), mcp.Description("Comma-separated file or directory paths associated with this memory")),
	)

	readTool := mcp.NewTool("read",
		mcp.WithDescription("Read memories from the events log, optionally filtered by path and limited to the N most recent."),
		mcp.WithString("path", mcp.Description("Filter by file or directory path (e.g. 'cmd/' or 'main.go')")),
		mcp.WithNumber("n", mcp.Description("Number of most recent memories to return (0 or omit for all)")),
	)

	removeTool := mcp.NewTool("remove",
		mcp.WithDescription("Remove a memory from the events log by its UUID."),
		mcp.WithString("id", mcp.Required(), mcp.Description("UUID of the memory to remove")),
	)

	s.AddTool(addTool, handleAdd)
	s.AddTool(readTool, handleRead)
	s.AddTool(removeTool, handleRemove)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(stderr, "mcp server error: %v\n", err)
		return 1
	}
	return 0
}

func handleAdd(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	name, _ := args["name"].(string)
	memory, _ := args["memory"].(string)
	paths, _ := args["paths"].(string)

	var out, errBuf bytes.Buffer
	code := cmdAdd([]string{"--name", name, "--memory", memory, "--paths", paths}, &out, &errBuf)
	if code != 0 {
		return mcp.NewToolResultError(errBuf.String()), nil
	}
	return mcp.NewToolResultText(out.String()), nil
}

func handleRead(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	toolArgs := request.GetArguments()
	var args []string
	if path, ok := toolArgs["path"].(string); ok && path != "" {
		args = append(args, "--path", path)
	}
	if n, ok := toolArgs["n"].(float64); ok && n > 0 {
		args = append(args, "-n", strconv.Itoa(int(n)))
	}

	var out, errBuf bytes.Buffer
	code := cmdRead(args, &out, &errBuf)
	if code != 0 {
		return mcp.NewToolResultError(errBuf.String()), nil
	}
	return mcp.NewToolResultText(out.String()), nil
}

func handleRemove(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	id, _ := args["id"].(string)

	var out, errBuf bytes.Buffer
	code := cmdRemove([]string{"--id", id}, &out, &errBuf)
	if code != 0 {
		return mcp.NewToolResultError(errBuf.String()), nil
	}
	return mcp.NewToolResultText(out.String()), nil
}
