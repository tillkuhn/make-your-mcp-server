# 🛠️ AGENTS.md — Guidelines for Autonomous Code Agents

This file provides best practices, commands, and conventions for agentic (and human) contributors working in this repository. Follow these rules and suggestions to maximize effective, robust, and maintainable modifications.

---

## 📦 Build, Lint, and Test Commands

### Build
- **Go Build (per tool):**
  - To build a single tool from `cmd/`, run:
    ```sh
    go build -o bin/<tool-binary> cmd/<tool>/main.go
    ```
  - Example:
    ```sh
    go build -o bin/mcp-curl cmd/curl/main.go
    ```
  
- **All Binaries (via Makefile):**
  - Build all standard binaries:
    ```sh
    make mcp-curl mcp-time mcp-random
    ```

- **Docker Image:**
  - To package e.g. `mcp-curl` as a Docker image:
    ```sh
    docker build -t mcp-curl .
    ```

### Lint
- **Go formatting:**
  - Format all code using:
    ```sh
    go fmt ./...
    ```
  - Lint tools (`golint`, `staticcheck`, etc.) are not configured by default. Add your preferred linter to workflow if needed.

### Test
- **All tests:**
    ```sh
    make test
    # or
    go test ./cmd/...
    ```
- **Single test file:**
    ```sh
    go test ./cmd/<tool>/main_test.go
    ```
- **Test coverage for a package:**
    ```sh
    go test -cover ./cmd/<tool>
    ```
- **Run a specific test:**
    ```sh
    go test -run <TestFunctionName> ./cmd/<tool>/main_test.go
    ```
  - Example:
    ```sh
    go test -run TestCurlHandler_Http ./cmd/curl/main_test.go
    ```

### Cleaning artifacts
- **Clean build outputs:**
    ```sh
    make clean
    # or
    rm -rf ./bin/*
    ```

---

## 🗂️ Code Style Guidelines

### Imports
- Group stdlib, then third-party, then local modules.
- Use blank lines to separate groups.
- Prefer explicit imports; avoid dot imports and aliases unless improving clarity.

```go
import (
    "fmt"           // stdlib
    "os/exec"

    "github.com/mark3labs/mcp-go/mcp"   // 3rd-party
    "github.com/mark3labs/mcp-go/server"
    "github.com/brianvoe/gofakeit/v7"
    "mcp-curl/internal"                 // local module
)
```

### Formatting
- Enforce `gofmt` (run `go fmt ./...`).
- Use tabs for indentation, 4 spaces visually.
- Typical max line length is 100-120, but not strictly enforced.

### Types and Interfaces
- Use concrete types unless handler flexibility is needed.
- Use pointer receivers when mutating state or for efficiency.
- Use Go error idioms: functions return an error value as the last argument.
- Prefer explicitly typed variables unless short declarations are self-explanatory.

### Naming Conventions
- Files: snake_case, all lower (e.g. `main_test.go`).
- Functions:
  - `CamelCase` (`curlHandler`, `LogRequest`).
  - Tests: `TestXxx`.
- Variables: short and meaningful (e.g. `res`, `err`, `thing`).
- Constants: `PascalCase` or `ALL_CAPS` if package-global.
- Packages: all lower case (`internal`, `mcp`, etc).

### Error Handling
- Check errors *immediately* after every operation that may fail (esp. I/O, execs, system calls):
  ```go
  output, err := cmd.Output()
  if err != nil {
      internal.LogError("mcp-curl", err)
      return mcp.NewToolResultError(err.Error()), nil
  }
  ```
- Use `internal/log.go` helpers for logging errors, requests, and responses. Log all incoming requests, all significant responses, and all errors.
- When accepting interface input (esp. via MCP requests), gracefully validate parameter types:
  ```go
  url, ok := request.Params.Arguments["url"].(string)
  if !ok {
      internal.LogError("mcp-curl", fmt.Errorf("url parameter missing or not a string"))
      return mcp.NewToolResultError("url must be a string"), nil
  }
  ```
- Fatal program errors (cannot open log file etc.) may use `log.Fatalf`.

### Logging
- Log files are written to `mcp.log` with UTC timestamps. Log type is one of `[REQUEST]`, `[RESPONSE]`, or `[ERROR]`.
- Use the provided singletons and don't create multiple loggers.
- All agent output should be idempotent and avoid printing secrets or credentials.

### Test Code
- Store tests alongside code, named `main_test.go`.
- Use table-driven tests when appropriate (not currently used, but encouraged for agent contributions).
- Use reflection with caution; it's mainly used for MCP test request injection.
- Tests should skip gracefully if external dependencies (e.g. curl binary) are missing.

---

## 🪛 Tooling and Environment

- Go Version: 1.23.4 (see go.mod)
- Dependencies: See `go.mod` — e.g., `github.com/mark3labs/mcp-go`, `github.com/brianvoe/gofakeit/v7`
- MCPHost CLI for running/serving models
- Docker can be used for packaging/running MCP servers
- No Lint/Prettier config files by default
- No Cursor, Copilot, or additional agentic rule files in this repo as of this version

## 🚦 Workflow Suggestions
- Build and test before and after major changes.
- Prefer `make` if unsure of build/test variant required.
- Log meaningful errors and response events; avoid returning silent failures.
- Automated agents should prefer atomic file operations and maintain code formatting.
- If introducing new dependencies, update `go.mod` and explain their purpose in PR descriptions.
- Keep public interfaces minimal and strongly typed.
- Comment all exported functions and structures.

---

## 📝 References
- [Official Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [`mcphost` documentation](https://github.com/mark3labs/mcphost)
- [Model Context Protocol (MCP)](https://github.com/modelcontextprotocol/)
- Project-specific usage: see `README.md`, `Makefile`, and code samples in `cmd/`


*Document updated automatically for agentic use — 2026-01-19*
