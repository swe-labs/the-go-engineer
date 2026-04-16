.PHONY: build test vet fmt lint run clean help

# ============================================================================
# The Go Engineer â€” Makefile
# ============================================================================
# Usage:
#   make help     â€” Show all available commands
#   make build    â€” Compile all packages
#   make test     â€” Run all tests
#   make lint     â€” Run linter and vet
#   make fmt      â€” Format all Go files
#   make clean    â€” Remove build artifacts
# ============================================================================

# Default target
.DEFAULT_GOAL := help

## Build & Compile
build: ## Compile all packages (verifies everything compiles)
	@echo "ðŸ”¨ Building all packages..."
	go build ./...
	@echo "âœ… Build successful"

## Testing
test: ## Run all tests
	@echo "ðŸ§ª Running tests..."
	go test ./...

test-race: ## Run all tests with race detector
	@echo "ðŸ§ª Running tests with race detection..."
	go test -race ./...

test-verbose: ## Run all tests with verbose output
	@echo "ðŸ§ª Running tests (verbose)..."
	go test -v ./...

bench: ## Run benchmarks
	@echo "ðŸ“Š Running benchmarks..."
	go test -bench=. -benchmem -count=1 ./12-quality-and-performance/testing/benchmarks/

## Code Quality
vet: ## Run go vet (find suspicious code)
	@echo "ðŸ” Running go vet..."
	go vet ./...
	@echo "âœ… Vet passed"

fmt: ## Format all Go files
	@echo "âœ¨ Formatting code..."
	gofmt -w .
	@echo "âœ… Formatted"

fmt-check: ## Check if code is formatted (CI use)
	@echo "ðŸ” Checking formatting..."
	@test -z "$$(gofmt -l .)" || (echo "âŒ Unformatted files:" && gofmt -l . && exit 1)
	@echo "âœ… All files formatted"

lint: vet fmt-check ## Run all linters (vet + format check)
	@echo "âœ… All lint checks passed"

## Utility
tidy: ## Sync go.mod with imports
	@echo "ðŸ“¦ Tidying modules..."
	go mod tidy
	@echo "âœ… Modules tidied"

deps-check: ## Check module dependencies
	@echo "ðŸ“‹ Checking dependencies..."
	@echo "All dependencies:"
	go list -u -m all
	@echo "âœ… Dependency check complete"

deps-update: ## Update all dependencies to latest patch version
	@echo "ðŸ”„ Updating dependencies..."
	go get -u ./...
	go mod tidy
	@echo "âœ… Dependencies updated"

clean: ## Remove build artifacts
	@echo "ðŸ§¹ Cleaning..."
	rm -f coverage.out
	go clean -cache -testcache
	@echo "âœ… Clean"

## Run Examples
run-hello: ## Run the Hello World example
	go run ./01-foundations/01-getting-started/2-hello-world

run-env: ## Run the environment check
	go run ./01-foundations/01-getting-started/4-dev-environment

## Coverage
cover: ## Run tests with coverage report
	@echo "ðŸ“Š Generating coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	@echo "ðŸ“„ HTML report: go tool cover -html=coverage.out"

## Help
help: ## Show this help message
	@echo "The Go Engineer â€” Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Usage: make <command>"
