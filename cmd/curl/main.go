package main

import (
	"context"
	"os/exec"

	"fmt"
	"mcp-curl/internal"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"mcp-curl",
		"1.0.0",
	)

	// Add a curlTool
	curlTool := mcp.NewTool("use_curl",
		mcp.WithDescription("fetch this url or webpage"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("url of the webpage to fetch"),
		),
	)
	s.AddTool(curlTool, curlHandler)

	fmt.Println("🚀 Server started")
	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("😡 Server error: %v\n", err)
	}
	fmt.Println("👋 Server stopped")
}

func curlHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	internal.LogRequest("mcp-curl", request)

	url, ok := request.Params.Arguments["url"].(string)
	if !ok {
		internal.LogError("mcp-curl", fmt.Errorf("url parameter missing or not a string"))
		return mcp.NewToolResultError("url must be a string"), nil
	}
	cmd := exec.CommandContext(ctx, "curl", "-s", url)
	output, err := cmd.Output()
	if err != nil {
		internal.LogError("mcp-curl", err)
		return mcp.NewToolResultError(err.Error()), nil
	}
	content := string(output)
	res := mcp.NewToolResultText(content)
	internal.LogResponse("mcp-curl", res)
	return res, nil
}
