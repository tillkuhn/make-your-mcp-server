package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"mcp-curl/internal"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"mcp-time",
		"0.0.7",
	)

	// Add a tool
	tool := mcp.NewTool("use_time",
		mcp.WithDescription("fetch the current date and time"),
		//mcp.WithString("url",
		//	mcp.Required(),
		//	mcp.Description("url of the webpage to fetch"),
		//),
	)

	// Add a tool handler
	s.AddTool(tool, timeHandler)

	fmt.Println("🚀 Server started")
	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("😡 Server error: %v\n", err)
	}
	fmt.Println("👋 Server stopped")
}

func timeHandler(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	internal.LogRequest("mcp-time", req)

	//url, ok := request.Params.Arguments["url"].(string)
	//if !ok {
	//	return mcp.NewToolResultError("url must be a string"), nil
	//}
	//cmd := exec.Command("curl", "-s", url)
	//output, err := cmd.Output()
	//if err != nil {
	//	return mcp.NewToolResultError(err.Error()), nil
	//}
	currentTime := time.Now().Format("Monday, 02 Jan 2006, 15:04:05 MST")

	horst, err := os.Hostname()
	if err != nil {
		internal.LogError("mcp-time", err)
	}
	content := "It's now " + currentTime + " on " + horst
	res := mcp.NewToolResultText(content)
	internal.LogResponse("mcp-time", res)
	return res, nil
}
