# Creating an MCP Server in Go and Serving it with Docker

## TKs Notebook

* Forked from Project Article https://github.com/ollama-tlms-golang/05-make-your-mcp-server
* Previous Article: https://github.com/ollama-tlms-golang/04-what-is-mcp
* Next Article https://github.com/ollama-tlms-golang/06-mcp-host-genaiapp
* Tips Ollama on MacOS: https://medium.com/@whyamit101/local-ai-using-ollama-with-agents-114c72182c97
* MCPHost 🤖: https://github.com/mark3labs/mcphost "A CLI host application that enables Large Language Models (LLMs) to interact with external tools through the Model Context Protocol (MCP). Currently supports Claude, OpenAI, Google Gemini, and Ollama models."
* LLMs with tool calling support: https://www.docker.com/blog/local-llm-tool-calling-a-practical-evaluation/, try mistral or granite4, see https://ollama.com/search?q=tools

Run ollama in shell to see log interaction (to run as service and enable after re-login: brew services start ollama);
```
$  OLLAMA_FLASH_ATTENTION="1" OLLAMA_KV_CACHE_TYPE="q8_0" /opt/homebrew/opt/ollama/bin/ollama serve

time=2026-01-19T16:44:42.838+01:00 level=INFO source=routes.go:1667 msg="Listening on 127.0.0.1:11434 (version 0.14.2)"
```

Build MCP Server and run mcphost with ollama tinyllama model:
```
$ go build -o mcp-curl main.go
$ mcphost --config ./mcp.json --model ollama:tinyllama:latest

Model loaded: ollama (tinyllama:latest)
Loaded 1 tools from MCP servers

/tools
Available tools:
1. mcp-curl__use_curl
```

Interact programmatically with [OpenCode Zen endpoints](https://opencode.ai/docs/zen#endpoints) Big Pickle Model, API Key can be obtained from [OpenCode Zen](https://opencode.ai/docs/zen/)

```
$ curl -sS https://opencode.ai/zen/v1/models |jq .
{ "object": "list",
  "data": [{
    "id": "big-pickle",
    "object": "model",
    "created": 1769531122,
    "owned_by": "opencode"}(...)
```


```
$ curl https://opencode.ai/zen/v1/chat/completions \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $(cat .opencode-api-key)" \
-d '{ "model": "big-pickle","messages": [{"role": "user", "content": "Say this is a test!"}], "temperature": 0.7, "max_tokens": 5}'
```


## Introduction

Today we'll look at how to create an MCP server in Go and serve it with Docker.

Prerequisites: having read [Understanding the Model Context Protocol (MCP)](https://k33g.hashnode.dev/understanding-the-model-context-protocol-mcp)

To write an MCP server, there are several official SDKs:

- [TypeScript SDK](https://github.com/modelcontextprotocol/typescript-sdk)
- [Python SDK](https://github.com/modelcontextprotocol/python-sdk)
- [Kotlin SDK](https://github.com/modelcontextprotocol/kotlin-sdk)

In recent years, I've developed an appetite for Go, so I looked around to see if there was a Go implementation. On this page [https://github.com/punkpeye/awesome-mcp-servers?tab=readme-ov-file#frameworks](https://github.com/punkpeye/awesome-mcp-servers?tab=readme-ov-file#frameworks), I found several, notably [https://github.com/mark3labs/mcp-go](https://github.com/mark3labs/mcp-go), by the creator of `mcphost` (which I discuss and use in the [previous blog post](https://k33g.hashnode.dev/understanding-the-model-context-protocol-mcp)).

The advantage of the `mcp-go` project is that it's simple and also provides tools to develop an MCP client, which will be very useful for a future blog post.

## Creating an MCP Server

I won't detail or implement all the possibilities of an MCP server here. I'll just implement the essentials:

- Provide a list of tools for the LLM
- Execute these tools when the LLM invokes them

When I use an LLM, there's one thing I'd like to be able to do: give it a website link so it can find information there, make a summary for me, etc.

So I created an MCP server that calls the `curl` utility with a URL parameter and returns its content. But let's look at the source code:

**`go.mod`**:
```go
module mcp-curl

go 1.23.4

require github.com/mark3labs/mcp-go v0.8.2
require github.com/google/uuid v1.6.0
```

**`main.go`**:
```go
package main

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"mcp-curl",
		"1.0.0",
	)

	// Add a tool
	tool := mcp.NewTool("use_curl",
		mcp.WithDescription("fetch this webpage"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("url of the webpage to fetch"),
		),
	)

	// Add a tool handler
	s.AddTool(tool, curlHandler)

	fmt.Println("🚀 Server started")
	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("😡 Server error: %v\n", err)
	}
	fmt.Println("👋 Server stopped")
}

func curlHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	url, ok := request.Params.Arguments["url"].(string)
	if !ok {
		return mcp.NewToolResultError("url must be a string"), nil
	}
	cmd := exec.Command("curl", "-s", url)
	output, err := cmd.Output()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	content := string(output)

	return mcp.NewToolResultText(content), nil
}
```

### Explanations

It's really very simple. My server will provide only one tool. The code implements a server that exposes a web "fetching" tool based on `curl` using the Model Control Protocol.

- First, I create an MCP server named `mcp-curl` version `1.0.0` that works with standard input/output (stdio)
- Then I define a tool named `use_curl` that takes a required `url` parameter
- Finally, I add the `handler` for this tool that executes the `curl` command with the `-s` option to retrieve the webpage content
- Of course, I handle errors and return the webpage content as text

## Packaging the MCP Server with Docker

To make my life easier, I'll use a `Dockerfile` to build a Docker image containing my MCP server. This way, I can more easily deploy it on different platforms or make it available to anyone who wants to use it.

**`Dockerfile`**:
```Dockerfile
FROM golang:1.23.4-alpine AS builder
WORKDIR /app
COPY go.mod .
COPY main.go .

RUN <<EOF
go mod tidy 
go build
EOF

FROM curlimages/curl:8.6.0
WORKDIR /app
COPY --from=builder /app/mcp-curl .
ENTRYPOINT ["./mcp-curl"]
```

My `Dockerfile` has two parts:

- The first part builds my MCP server in Go
- The second part uses the `curlimages/curl:8.6.0` image that contains `curl` and copies my `mcp-curl` server into the image (so `mcp-curl` is the executable that calls `curl`)

To build the Docker image, I run the following command:

```bash
docker build -t mcp-curl .
```

And now let's see how to use our new MCP server with `mcphost`:

## Using the MCP Server with `mcphost`

First, we need to create a configuration file `mcp.json` for `mcphost`:

```json
{
  "mcpServers": {
    "mcp-curl-with-docker" :{
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-i",
        "mcp-curl"
      ]
    }
  }
}
```

Then we can use `mcphost` to use our MCP server and an LLM with Ollama, like this:

```bash
mcphost --config ./mcp.json --model ollama:qwen2.5-coder:14b
```

The MCP server is recognized by `mcphost`:
![mcp](imgs/01-mcp.png)

You can request the list of available tools with the `/tools` command:
![mcp](imgs/02-mcp.png)

Now you can request the content of a webpage and analyze its content (in my example I retrieve Go code from GitHub):
![mcp](imgs/03-mcp.png)

Wait a little bit:
![mcp](imgs/04-mcp.png)

And here's the webpage content:
![mcp](imgs/05-mcp.png)

## Conclusion

You can see that with just a few lines, it becomes really easy to give "superpowers" to your LLMs. In a future blog post, I'll show you how to create a generative AI application with Ollama and an MCP client in Go to interact with our MCP server.
