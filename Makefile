.PHONY: build test test-race test-verbose bench vet fmt fmt-check lint tidy deps-check deps-update clean run-hello run-env cover validate ci help

# ============================================================================
# The Go Engineer — Makefile
# ============================================================================
# Usage:
#   make help       — Show all available commands
#   make ci         — Run CI-equivalent local checks
#   make build      — Compile all packages
#   make test       — Run all tests
#   make validate   — Run curriculum validator
# ============================================================================

.DEFAULT_GOAL := help

## Build & Compile
build: ## Compile all packages
	@echo "Building all packages..."
	go build ./...
	@echo "Build successful"

## Testing
test: ## Run all tests
	@echo "Running tests..."
	go test ./...

test-race: ## Run all tests with race detector
	@echo "Running tests with race detector..."
	go test -race ./...

test-verbose: ## Run all tests with verbose output
	@echo "Running tests with verbose output..."
	go test -v ./...

bench: ## Run benchmark suite
	@echo "Running benchmarks..."
	go test -bench=. -benchmem -count=1 ./08-quality-test/01-quality-and-performance/testing/benchmarks/

cover: ## Run tests with coverage report
	@echo "Generating coverage report..."
	go test -coverprofile coverage.out ./...
	go tool cover -func coverage.out
	@echo "HTML report: go tool cover -html coverage.out -o coverage.html"

## Code Quality
vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...
	@echo "Vet passed"

fmt: ## Format all Go files
	@echo "Formatting code..."
	gofmt -w .
	@echo "Formatted"

fmt-check: ## Check formatting without writing files
	@echo "Checking formatting..."
	@test -z "$$(gofmt -l .)" || (echo "Unformatted files:" && gofmt -l . && exit 1)
	@echo "All files formatted"

lint: vet fmt-check ## Run vet and format check
	@echo "Lint checks passed"

## Modules
tidy: ## Sync go.mod and go.sum with imports
	@echo "Tidying modules..."
	go mod tidy
	@git diff --exit-code -- go.mod go.sum
	@echo "Modules tidy"

deps-check: ## List available module updates
	@echo "Checking dependencies..."
	go list -u -m all

deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy
	@echo "Dependencies updated"

## Curriculum
validate: ## Run curriculum validator
	@echo "Validating curriculum..."
	go run ./scripts/validate_curriculum.go
	@echo "Curriculum valid"

## Run Examples
run-hello: ## Run the Hello World example
	go run ./01-getting-started/2-hello-world

run-env: ## Run the environment check
	go run ./01-getting-started/4-dev-environment

## CI
ci: build vet fmt-check tidy test test-race cover validate ## Run CI-equivalent local checks
	@echo "CI-equivalent checks passed"

clean: ## Remove build artifacts and caches
	@echo "Cleaning..."
	rm -f coverage.out coverage.html
	go clean -cache -testcache
	@echo "Clean"

help: ## Show this help message
	@echo "The Go Engineer — Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Usage: make <command>"
