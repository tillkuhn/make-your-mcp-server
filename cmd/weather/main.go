package main

import (
	"context"
	"os/exec"

	"fmt"
	"mcp-curl/internal"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const toolName = "weather_info"

var execCommandContext = exec.CommandContext

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"mcp-curl",
		"1.0.0",
	)

	// Add a curlTool
	tool := mcp.NewTool(toolName,
		mcp.WithDescription("get the current weather for the specified location"),
		mcp.WithString("location",
			mcp.DefaultString(""), // empty will use current location based on IP
			mcp.Description("location for weather to fetch"),
		),
	)
	s.AddTool(tool, weatherHandler)

	fmt.Println("🚀 Server started for tool " + toolName)
	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("😡 Server error: %v\n", err)
	}
	fmt.Println("👋 Server stopped for tool " + toolName)
}

func weatherHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	param, ok := request.Params.Arguments["location"].(string)
	if !ok {
		return mcp.NewToolResultError("attribute must be a string"), nil
	}
	url := fmt.Sprintf("%s://%s/%s?format=3", "http", "wttr.in", param)
	internal.LogRequest(toolName, request)
	internal.LogRequest(toolName, "URL: "+url)
	cmd := execCommandContext(ctx, "curl", "-s", url)
	output, err := cmd.Output()
	if err != nil {
		internal.LogError("mcp-curl", err)
		return mcp.NewToolResultError(err.Error()), nil
	}
	content := string(output)
	res := mcp.NewToolResultText(content)
	internal.LogResponse(toolName, res)
	return res, nil
}
