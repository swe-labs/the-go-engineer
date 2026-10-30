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
	go test -bench=. -benchmem -count=1 ./curriculum/modules/13-performance-memory-engineering/projects/profile-and-fix-lab/_solution

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

fmt: ## Format all Go files (excluding starter code)
	@echo "Formatting code..."
	@find . -name "*.go" -not -path "*/_starter/*" -exec gofmt -w {} + 2>/dev/null || gofmt -w .
	@echo "Formatted"

fmt-check: ## Check formatting without writing files (excluding starter code)
	@echo "Checking formatting..."
	@unformatted="$$(gofmt -l . | grep -v '_starter' || true)"; \
	test -z "$$unformatted" || (echo "Unformatted files:" && echo "$$unformatted" && exit 1)
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
	go run ./internal/tools/curriculum validate-all
	@echo "Curriculum validation completed"

## Run Examples
run-hello: ## Run the Hello World example
	go run ./curriculum/modules/00-orientation/lessons/01-how-to-use-this-repository/main.go

run-env: ## Run the environment check
	go run ./curriculum/modules/02-go-setup-tooling/lessons/01-toolchain-installation-and-path/main.go

## CI
ci: build vet fmt-check tidy test test-race cover validate ## Run CI-equivalent local checks
	@echo "CI-equivalent checks passed"

clean: ## Remove build artifacts and caches
	@echo "Cleaning..."
	rm -f coverage.out coverage.html curriculum.exe *.exe
	go clean -cache -testcache
	@echo "Clean"

help: ## Show this help message
	@echo "The Go Engineer — Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Usage: make <command>"
