# 🛠️ AGENTS.md — Guidelines for Autonomous Code Agents

This document provides robust conventions, commands, and style rules for any agentic (or human) contributor operating in this repository. Adhering to these guidelines ensures safe, maintainable, and efficient delivery of changes.

---

## 📦 Build, Lint, and Test Commands

### Build
- **Go Build (per tool):**
  - Build a single tool:
    ```sh
    go build -o bin/<tool-binary> cmd/<tool>/main.go
    ```
- **All Binaries (via Makefile):**
  - Build all core binaries:
    ```sh
    make mcp-curl mcp-time mcp-random
    ```
- **Docker Image:**
  - Package as Docker (example: mcp-curl):
    ```sh
    docker build -t mcp-curl .
    ```

### Lint
- **Formatting all code:**
  ```sh
  go fmt ./...
  ```
- No linter enforced by default (add `golint`, `staticcheck` or others if needed).

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
- **Run a specific test:**
  ```sh
  go test -run <TestFunctionName> ./cmd/<tool>/main_test.go
  ```
- **Test coverage per package:**
  ```sh
  go test -cover ./cmd/<tool>
  ```
- **Clean builds:**
  ```sh
  make clean
  # or
  rm -rf ./bin/*
  ```

---

## 🗂️ Code Style Guidelines

### Imports
- Group stdlib imports, third-party packages, then local modules (separate with blank lines).
- Use explicit imports; avoid dot imports and aliases unless truly needed for clarity.

Example:
```go
import (
    "fmt"           // stdlib
    "os/exec"

    "github.com/mark3labs/mcp-go/mcp"   // 3rd-party
    "github.com/mark3labs/mcp-go/server"
    "github.com/brianvoe/gofakeit/v7"
    "mcp-curl/internal"                 // local
)
```

### Formatting
- Enforce idiomatic Go formatting via `gofmt`.
- Use tabs for indentation (4 spaces width visually).
- No strict line length (suggested 100-120 chars for agents).

### Types and Interfaces
- Use concrete types unless a handler requires flexibility.
- Pointer receivers for mutating methods or large structs.
- Use Go error idioms (return `error` as last return value; never panic unless truly fatal).

### Naming Conventions
- Files: snake_case (e.g. `main_test.go`).
- Functions: CamelCase for functions and methods; tests are `TestXxx`.
- Variables: concise but descriptive; avoid single-letter unless trivial.
- Constants: PascalCase or ALL_CAPS (for public or global constants).
- Packages: all lower case.

### Error Handling
- Always check error returns, especially for I/O, OS, and subprocess calls.
  ```go
  output, err := cmd.Output()
  if err != nil {
      internal.LogError("mcp-curl", err)
      return mcp.NewToolResultError(err.Error()), nil
  }
  ```
- Validate all interface and input types (esp. for tool params and requests).
- Use `internal/log.go` helpers where available for logging errors.
- Use `log.Fatalf` or os.Exit only on unrecoverable, fatal errors.
- Boundary check all user or LLM-facing input; return clear, validated errors at boundary.

### Logging
- Log all requests, significant responses, and errors to `mcp.log`, including UTC timestamps.
- Valid log types: `[REQUEST]`, `[RESPONSE]`, `[ERROR]`.
- Do NOT print sensitive info in logs or output.

### Tests
- Store all tests alongside code, use `main_test.go` naming pattern.
- Use table-driven or property-based tests where possible.
- Always skip tests gracefully if required binaries (e.g. `curl`) or infrastructure are missing.

### Miscellaneous Conventions
- Prefer atomic file operations for test and code changes (no partial/dirty writes).
- Use safe concurrency: protect shared resources, avoid global mutable state unless explicitly synchronized.
- Avoid premature optimizations; prioritize readability, then performance if justified.
- Comment all exported symbols and interfaces for generated docs and agentic use.
- Limit direct shell/script execution to well-checked contexts (agents: double-check shell escapes, unexpected command construction is dangerous).

---

## 🪛 Tooling and Environment

- Go version: 1.23.4 (see `go.mod`).
- Dependencies: see `go.mod` for canonical list.
- MCPHost CLI for running/serving models and MCP servers.
- Docker multi-stage for packaging and deployment.
- No Prettier/lint config, nor Cursor or Copilot rules as of 2026-01-19.

---

## 🚦 Workflow & Collaboration

- Build *and* test before and after significant changes. Commit only after success.
- Use branches for significant features or refactorings; PRs should have clear rationale.
- Prefer `make` targets unless explicit variant is needed for CLI.
- Document new dependencies and explain rationale in PRs (keep `go.mod`/`go.sum` updated).
- Log all non-trivial agent or user-facing events (requests, errors, boundary cases).
- Automated agents: NEVER output secrets or credentials, and never log them.
- Keep public API minimal and strongly typed. Use interfaces for extensibility only.
- Add comments to exported functions, methods, types, and packages.
- Always check project-specific references: `README.md`, `Makefile`, code in `cmd/`.

---

## 📝 References
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [mcphost documentation](https://github.com/mark3labs/mcphost)
- [Model Context Protocol (MCP)](https://github.com/modelcontextprotocol/)

*Document improved for agentic use — 2026-01-19*
