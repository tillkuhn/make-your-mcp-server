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

.PHONY: run
run: mcp-curl mcp-time mcp-random ## runs mcphost with mcp.kson config and specified ollama model
	mcphost --config ./mcp.json --model ollama:$(MODEL)

.PHONY: models
models:  ## list ollama models
	ollama list
