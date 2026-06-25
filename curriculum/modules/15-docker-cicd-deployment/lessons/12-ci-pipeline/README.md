# CI pipeline

## Learning objective

Design, model, and validate a continuous integration pipeline for a Go project, comprising lint, vet, test, and build stages, and distinguish CI from CD in the context of pipeline-as-code practices.

## Why this matters

A CI pipeline is the automated gate that every code change must pass before reaching production. Without CI, teams rely on manual testing and code review alone, which misses integration issues, race conditions, and build failures. For Go projects, CI catches vet warnings, test regressions, and cross-compilation problems before they affect users. A well-designed pipeline gives every engineer confidence that main is always deployable.

## Mental model

Think of a CI pipeline as a series of quality gates on an assembly line. Each gate checks a specific property of the code:

- **Lint gate**: Checks code style and common mistakes (golangci-lint).
- **Vet gate**: Checks suspicious constructs the compiler allows (go vet).
- **Test gate**: Runs unit tests and reports coverage.
- **Build gate**: Verifies the code compiles correctly.
- **Integration gate** (optional): Runs tests that require real dependencies.

If any required gate fails, the assembly line stops, and the change is rejected. If all gates pass, the code proceeds to the packaging and deployment stage (CD). The pipeline is defined as code (YAML or a Go script) so it is versioned, reviewable, and reproducible.

## Core idea

Continuous Integration means merging code changes into a shared mainline frequently (multiple times per day) and verifying each merge with an automated build and test process. The core stages for a Go CI pipeline:

| Stage | Tool/Command | Purpose |
|---|---|---|
| Lint | `golangci-lint run` | Enforce style, detect bugs early |
| Vet | `go vet ./...` | Find suspicious constructs |
| Unit test | `go test -race -shuffle=on ./...` | Verify correctness, data races |
| Build | `go build ./...` | Ensure compilation succeeds |
| Integration test | `go test -tags=integration ./...` | Test with real dependencies |

Pipeline as code means the pipeline definition lives in the repository (e.g., `.github/workflows/ci.yml`), is versioned alongside the application code, and goes through the same review process. This ensures pipeline changes are traceable and auditable.

## Under the hood

When a CI pipeline runs, each stage executes in an isolated environment (a fresh runner or container). The pipeline runner:

1. Clones the repository at the specific commit SHA.
2. Sets up the toolchain (Go version specified in `go.mod` or workflow).
3. Executes each stage sequentially or in parallel based on dependency declarations.
4. Collects stdout, stderr, and exit codes.
5. Reports status (pass/fail) back to the platform.

For Go specifically, `go vet` analyzes source code AST and reports suspicious patterns like unreachable code, nil comparisons on interface values, and incorrect `printf` format strings. `go test -race` compiles the test binary with the Go race detector, which instruments all memory accesses and detects unsynchronized concurrent reads and writes.

The build stage runs `go build`, which compiles all packages in the module. A successful build verifies that all dependencies resolve, types are correct, and the binary links correctly for the target OS and architecture.

## How Go uses it

Go's toolchain is designed for CI. `go vet`, `go test`, and `go build` produce consistent, machine-parseable output across platforms. Typical Go CI pipeline script:

```bash
go vet ./...
go test -race -shuffle=on -count=1 ./...
go build -o /dev/null ./...
```

For monorepos, run these commands for each changed module to keep CI fast. Use `go list -f '{{.Dir}}' ./...` to find all module roots and run pipeline stages per module.

Many Go projects combine linting (`golangci-lint`), static analysis (`staticcheck`), and vulnerability scanning (`govulncheck`) as additional pipeline stages. These tools can be cached between runs using the Go build cache.

## Go example

```go
package main

import (
	"fmt"
	"os/exec"
	"strings"
)

type PipelineStage struct {
	Name     string
	Command  string
	Required bool
}

type PipelineResult struct {
	Stage  string
	Passed bool
	Output string
	Error  string
}

type Pipeline struct {
	Stages []PipelineStage
}

func NewPipeline(stages []PipelineStage) *Pipeline {
	return &Pipeline{Stages: stages}
}

func (p *Pipeline) Run() []PipelineResult {
	var results []PipelineResult
	for _, stage := range p.Stages {
		result := PipelineResult{Stage: stage.Name}
		out, err := exec.Command("go", strings.Fields(stage.Command)...).CombinedOutput()
		result.Output = string(out)
		if err != nil {
			result.Error = err.Error()
			result.Passed = false
		} else {
			result.Passed = true
		}
		results = append(results, result)
		if !result.Passed && stage.Required {
			break
		}
	}
	return results
}

func (p *Pipeline) Simulate() []PipelineResult {
	var results []PipelineResult
	for _, stage := range p.Stages {
		result := PipelineResult{
			Stage:  stage.Name,
			Passed: true,
			Output: fmt.Sprintf("[simulated] go %s", stage.Command),
		}
		results = append(results, result)
	}
	return results
}

func main() {
	pipeline := NewPipeline([]PipelineStage{
		{Name: "lint", Command: "vet ./...", Required: false},
		{Name: "test", Command: "test -race -shuffle=on ./...", Required: true},
		{Name: "build", Command: "build -o /dev/null ./...", Required: true},
	})
	fmt.Println("CI Pipeline stages:")
	for _, s := range pipeline.Stages {
		req := ""
		if s.Required {
			req = " [required]"
		}
		fmt.Printf("  - %s: go %s%s\n", s.Name, s.Command, req)
	}
	results := pipeline.Simulate()
	fmt.Println("\nResults:")
	for _, r := range results {
		status := "PASS"
		if !r.Passed {
			status = "FAIL"
		}
		fmt.Printf("  %s: %s\n", r.Stage, status)
	}
}
```

## Step-by-step execution

For the pipeline above when `Run()` is called (if Go is installed):

1. Stage `lint`: Execute `go vet ./...`. If vet finds suspicious code, output contains warnings and exit code is non-zero. Since lint is not required, the pipeline continues even if it fails.
2. Stage `test`: Execute `go test -race -shuffle=on ./...`. This compiles test binaries with the race detector, shuffles test execution order to detect flaky tests, and runs all tests. Exit code is non-zero if any test fails.
3. If test fails, the pipeline stops (test is required). Result contains the failing test output.
4. Stage `build`: Execute `go build -o /dev/null ./...`. Verifies all packages compile. Exit code is non-zero if any package has a compilation error.

The `Simulate()` method returns mock results for demonstration without needing Go installed, useful for testing pipeline logic offline.

## Common mistakes

- **Running lint after test**: Lint is fast and fails early. Run lint and vet first to fail fast, then run the slower test and build stages.
- **Not using `-race` in CI**: Race conditions are only detectable at runtime. Always run `go test -race` in CI to catch data races before they reach production.
- **Skipping `-count=1`**: Go's test cache causes subsequent runs to return cached results. Use `-count=1` to disable caching and force real test execution in CI.
- **No test output on failure**: Use `-v` flag to see individual test results. Without it, CI output only shows PASS/FAIL per package.
- **Building to a path that may not exist**: Use `-o /dev/null` (Linux/macOS) or `-o NUL` (Windows) for build-only verification without creating a binary.

## Debugging walkthrough

Consider a CI pipeline that fails with:

```
go: errors parsing go.mod: unexpected end of JSON input
```

**Symptom**: The pipeline fails at the build stage with no compilation errors in the source code.

**Investigation**: The error is from `go.mod`, not from `.go` files. Check if `go.mod` was edited manually with invalid syntax.

**Root cause**: A merge conflict was resolved incorrectly, leaving merge conflict markers (`<<<<<<<`, `=======`) in `go.mod`. Run `go mod tidy` to regenerate a valid `go.mod`.

**Fix**:
```bash
git checkout main -- go.mod go.sum
go mod tidy
```

Then re-run the pipeline.

Another common failure:

```
=== RUN   TestParseConfig
    config_test.go:15: expected nil, got error: open /etc/app/config.json: no such file or directory
--- FAIL: TestParseConfig (0.00s)
```

**Root cause**: Test depends on a file that exists on the developer's machine but not on the CI runner. Filesystem-dependent tests must create temporary files in the test itself using `t.TempDir()`.

## Production notes

Production CI pipelines for Go services follow these practices:

- **Parallelize independent stages**: Run lint, vet, and test across multiple operating systems in parallel using matrix builds. Wait only for the build stage if subsequent steps depend on the binary.
- **Cache Go module and build cache**: Use `actions/cache@v4` with a key based on `go.sum`. This reduces CI time from minutes to seconds for dependency-heavy projects.
- **Fail fast, fail loud**: Required stages stop the pipeline immediately. Notifications (Slack, email) alert the team when a required stage fails.
- **Pipeline as code review**: Treat workflow YAML changes like application code changes. Require PR review for pipeline modifications.
- **GitHub status checks**: Configure required status checks in branch protection rules. A PR cannot merge unless all required CI stages pass.

## Performance implications

- **Sequential vs parallel stages**: Running lint, vet, and test sequentially adds latency. Run lint and vet in parallel with test and build when possible.
- **Race detector overhead**: `go test -race` is 5-20x slower than regular tests. Consider running race-detected tests on a subset (e.g., packages that deal with concurrency) and regular tests on the rest.
- **Cache hit ratio**: A well-tuned cache (correct keys, good restore-keys) saves 30-90 seconds per pipeline run. Monitor cache hit rates in GitHub Actions analytics.
- **Cross-compilation in CI**: Building for multiple targets (linux/amd64, linux/arm64, darwin/amd64) adds time proportional to the number of targets. Use parallel matrix jobs for different GOOS/GOARCH combinations.

## Practice task

Write a function `PipelineFromYAML(yamlContent string) (*Pipeline, error)` that parses a minimal YAML-like format:

```
stages:
  - name: lint
    command: vet ./...
    required: false
  - name: test
    command: test -race ./...
    required: true
```

Parse the lines starting with `name:`, `command:`, `required:` and build a `Pipeline` struct. Return an error if any line has invalid format. Then write a `main()` that calls `Simulate()` on the parsed pipeline and prints results.

## Tests / verification

```bash
go run ./curriculum/modules/15-docker-cicd-deployment/lessons/12-ci-pipeline
go test ./curriculum/modules/15-docker-cicd-deployment/lessons/12-ci-pipeline
```

The existing tests verify pipeline simulation, result formatting, empty pipeline behavior, and the ability to build and vet Go packages. After completing the practice task, add tests for your `PipelineFromYAML` function covering valid input, missing fields, and malformed lines.

## Review questions

1. What is the difference between a required and a non-required stage in a CI pipeline? When would you use each?
2. Why should `go test -race` be run in CI even if all tests pass locally?
3. What does `-count=1` do in `go test` and why is it important for CI?
4. How does pipeline as code differ from configuring CI through a web UI? What are the advantages?
5. A build stage passes but the binary is not deployable. What kinds of issues could cause this?

## NEXT UP

Release artifacts: how to produce, version, checksum, and store Go binaries as deployable artifacts.
