.PHONY: help build run clean test deps

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Download dependencies
	go mod download
	go mod verify

build: ## Build the server binary
	go build -o bin/mcp-server main.go

run: ## Run the server
	go run main.go

clean: ## Clean build artifacts
	rm -rf bin/
	go clean

test: ## Run tests
	go test -v ./...

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

lint: fmt vet ## Run linters

dev: ## Run in development mode with auto-reload (requires air)
	@which air > /dev/null || (echo "air not found. Install with: go install github.com/cosmtrek/air@latest" && exit 1)
	air

install-dev-tools: ## Install development tools
	go install github.com/cosmtrek/air@latest



