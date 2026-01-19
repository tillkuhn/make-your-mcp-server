package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"mcp-random	",
		"0.0.7",
	)

	// Add a tool
	tool := mcp.NewTool("use_random",
		mcp.WithDescription("create random things"),
		mcp.WithString("thing",
			mcp.Required(),
			mcp.Description("what kind of random thing you would like to create, e.g. beer, job, food	 "),
		),
	)

	// Add a tool handler
	s.AddTool(tool, randomHandler)

	fmt.Println("🚀 Server started")
	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("😡 Server error: %v\n", err)
	}
	fmt.Println("👋 Server stopped")
}

func randomHandler(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	thing, ok := req.Params.Arguments["thing"].(string)
	if !ok {
		return mcp.NewToolResultError("thing must be a string"), nil
	}
	var content string
	switch strings.ToLower(thing) {
	case "beer":
		content = gofakeit.BeerName()
	case "job":
		content = gofakeit.JobTitle()
	case "food":
		content = gofakeit.MinecraftFood()
	case "hobby":
		content = gofakeit.Hobby()
	default:
		return mcp.NewToolResultError("invalid thing " + thing + " currently only support for beer, hobby, job and (minecraft) food"), nil
	}

	//cmd := exec.Command("curl", "-s", url)
	//output, err := cmd.Output()
	//if err != nil {
	//	return mcp.NewToolResultError(err.Error()), nil
	//}

	return mcp.NewToolResultText(content), nil
}
