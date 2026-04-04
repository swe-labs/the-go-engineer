.PHONY: build test vet fmt lint run clean help

# ============================================================================
# The Go Engineer — Makefile
# ============================================================================
# Usage:
#   make help     — Show all available commands
#   make build    — Compile all packages
#   make test     — Run all tests
#   make lint     — Run linter and vet
#   make fmt      — Format all Go files
#   make clean    — Remove build artifacts
# ============================================================================

# Default target
.DEFAULT_GOAL := help

## Build & Compile
build: ## Compile all packages (verifies everything compiles)
	@echo "🔨 Building all packages..."
	go build ./...
	@echo "✅ Build successful"

## Testing
test: ## Run all tests
	@echo "🧪 Running tests..."
	go test ./...

test-race: ## Run all tests with race detector
	@echo "🧪 Running tests with race detection..."
	go test -race ./...

test-verbose: ## Run all tests with verbose output
	@echo "🧪 Running tests (verbose)..."
	go test -v ./...

bench: ## Run benchmarks
	@echo "📊 Running benchmarks..."
	go test -bench=. -benchmem -count=1 ./13-quality-and-performance/testing/benchmarks/

## Code Quality
vet: ## Run go vet (find suspicious code)
	@echo "🔍 Running go vet..."
	go vet ./...
	@echo "✅ Vet passed"

fmt: ## Format all Go files
	@echo "✨ Formatting code..."
	gofmt -w .
	@echo "✅ Formatted"

fmt-check: ## Check if code is formatted (CI use)
	@echo "🔍 Checking formatting..."
	@test -z "$$(gofmt -l .)" || (echo "❌ Unformatted files:" && gofmt -l . && exit 1)
	@echo "✅ All files formatted"

lint: vet fmt-check ## Run all linters (vet + format check)
	@echo "✅ All lint checks passed"

## Utility
tidy: ## Sync go.mod with imports
	@echo "📦 Tidying modules..."
	go mod tidy
	@echo "✅ Modules tidied"

deps-check: ## Check module dependencies
	@echo "📋 Checking dependencies..."
	@echo "All dependencies:"
	go list -u -m all
	@echo "✅ Dependency check complete"

deps-update: ## Update all dependencies to latest patch version
	@echo "🔄 Updating dependencies..."
	go get -u ./...
	go mod tidy
	@echo "✅ Dependencies updated"

clean: ## Remove build artifacts
	@echo "🧹 Cleaning..."
	rm -f coverage.out
	go clean -cache -testcache
	@echo "✅ Clean"

## Run Examples
run-hello: ## Run the Hello World example
	go run ./01-core-foundations/getting-started/2-hello-world

run-env: ## Run the environment check
	go run ./01-core-foundations/getting-started/4-dev-environment

## Coverage
cover: ## Run tests with coverage report
	@echo "📊 Generating coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	@echo "📄 HTML report: go tool cover -html=coverage.out"

## Help
help: ## Show this help message
	@echo "The Go Engineer — Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Usage: make <command>"
