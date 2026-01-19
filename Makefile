# Caution: Only models with "tools" support work for this demo, e.g. mistral or granite4:3b (https://ollama.com/library/)
#MODEL ?= mistral:latest
MODEL ?= granite4:3b

.PHONY: help
help: ## this help
	@# https://gist.github.com/prwhite/8168133#gistcomment-3291344
	@grep -E "^$$PFX[0-9a-zA-Z_-]+:.*?## .*$$" $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'; echo "";\

clean: ## clean build artifacts
	@rm -rf ./bin/*

mcp-curl: cmd/curl/main.go ## build mcp-curl
	go build -o bin/mcp-curl cmd/curl/main.go

mcp-time: cmd/time/main.go ## build mcp-time
	go build -o bin/mcp-time cmd/time/main.go

mcp-random: cmd/random/main.go ## build mcp-random
	go build -o bin/mcp-random cmd/random/main.go

install: ## install MCPHost if not already installed
	@if ! command -v mcphost >/dev/null 2>&1; then \
		echo "Installing MCPHost..."; \
		go install github.com/mark3labs/mcphost@latest; \
	else \
		echo "MCPHost already installed."; \
	fi

.PHONY: test
# Run all Go unit tests for all mcp command tools
# Usage: make test
# Description: Runs all Go tests recursively in ./cmd/*
test: ## run all unit tests for mcp tools
	go test ./cmd/...

.PHONY: run
run: install mcp-curl mcp-time mcp-random ## runs MCPHost with mcp.kson config and specified ollama model
	mcphost --config ./mcp.json --model ollama:$(MODEL)


.PHONY: ollama
ollama:  ## install ollama if not already installed, and run it as a service
	@if ! command -v ollama >/dev/null 2>&1; then \
		echo "Installing Ollama..."; \
		brew install ollama; \
	else \
		echo "MCPHost already installed."; \
	fi
	brew services info ollama
	brew services run ollama
	@ollama list


.PHONY: ollama-stop
ollama-stop:  ## stop ollama service
	brew services stop ollama

.PHONY: ollama-logs
ollama-logs:  ## show logs from ollama service (blocking)
	tail -f /opt/homebrew/var/log/ollama.log

.PHONY: models
models:  ## list ollama models
	ollama list
