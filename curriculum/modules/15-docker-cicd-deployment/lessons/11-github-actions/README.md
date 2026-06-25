# GitHub Actions

## Learning objective

Parse, validate, and generate GitHub Actions workflow configurations using Go, and explain how workflow YAML, jobs, steps, actions, matrix builds, runners, and caching work together in a CI pipeline.

## Why this matters

GitHub Actions is the most widely used CI/CD platform for open-source and enterprise Go projects. Every Go engineer publishing a library or deploying a service on GitHub interacts with workflow files. Understanding how workflows are structured, how matrix builds parallelize tests, and how caching speeds up module downloads is essential for shipping reliable Go software. Writing Go code to validate or generate workflows lets you treat your CI configuration as code, catch errors before committing, and automate repetitive workflow setup.

## Mental model

A GitHub Actions workflow is a directed acyclic graph of jobs. Each job runs on a virtual machine (runner) and contains sequential steps. Steps either run shell commands or invoke prebuilt actions (reusable units). A matrix strategy fans a single job definition into multiple parallel runs across combinations of operating systems, Go versions, or other variables. Caching is a key-value store keyed by a hash of your dependency files: if the key matches a previous run, the cache is restored instead of downloading fresh.

Think of a workflow as a factory assembly line: the trigger event (push, PR) is the start bell, each job is a workstation, each step is a task at that workstation, and the matrix is multiple parallel assembly lines running the same procedure with different tools.

## Core idea

GitHub Actions automates software workflows directly from your repository. The key building blocks:

| Block | Purpose | Example |
|---|---|---|
| Workflow | Top-level automation defined in `.github/workflows/*.yml` | `ci.yml` |
| Event | Trigger that starts the workflow | `push`, `pull_request`, `schedule` |
| Job | A group of steps running on one runner | `test`, `build` |
| Step | An individual task within a job | `Run go test ./...` |
| Action | A reusable step packaged as a GitHub repo | `actions/checkout@v4` |
| Runner | The VM that executes jobs | `ubuntu-latest`, `windows-latest` |
| Matrix | Cross-product of variables expanding into multiple jobs | `{go: [1.21, 1.22], os: [linux, windows]}` |
| Cache | Store/restore dependencies by key to avoid redownloads | `actions/cache@v4` for Go module cache |

## Under the hood

When a workflow is triggered, GitHub spins up a fresh runner VM, clones the repository, and executes each job in dependency order (defined by `needs`). Each step runs in a new shell process. If a step fails, remaining steps in that job are skipped. The runner sends logs and status back to GitHub's API in real time.

Caching works by computing a cache key (typically a hash of `go.sum`), uploading a tarball of `~/.cache/go` on the first run, and restoring it on subsequent runs when the key matches. Cache hits reduce `go mod download` from minutes to seconds.

Matrix builds create one job per combination. GitHub limits matrix size to 256 jobs per workflow. Each matrix job gets its own runner and reports results individually.

Go-specific points: Go compilation is fast enough that many projects skip caching and just run `go mod download` each time. But for monorepos or projects with many dependencies, caching the module cache and the build cache (`go build -cache`) significantly reduces CI time.

## How Go uses it

Go projects on GitHub Actions typically follow this pattern:

- **Lint**: `golangci-lint run` to catch style and correctness issues.
- **Vet**: `go vet ./...` for suspicious constructs.
- **Test**: `go test -race -shuffle=on ./...` with matrix across Go versions.
- **Build**: `go build ./...` or build the specific binary.
- **Cache**: `actions/cache@v4` with key based on `go.sum` hash.

The Go team provides `actions/setup-go@v5` which installs a specific Go version, sets `GOROOT` and `GOPATH`, and adds `$GOROOT/bin` and `$GOPATH/bin` to `PATH`.

Many open-source Go projects use a single workflow with a matrix of Go versions and operating systems. The `go.mod` file pins the module path and minimum Go version, while the workflow matrix tests across the N-2 latest Go releases.

## Go example

```go
package main

import (
	"fmt"
	"os"
)

type Workflow struct {
	Name string
	On   string
	Jobs []Job
}

type Job struct {
	Name   string
	RunsOn string
	Steps  []Step
	Needs  []string
	Matrix map[string][]string
}

type Step struct {
	Name    string
	Run     string
	Uses    string
	With    map[string]string
	Timeout int
}

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func ValidateWorkflow(w Workflow) []error {
	var errs []error
	if w.Name == "" {
		errs = append(errs, ValidationError{"name", "workflow name is required"})
	}
	if w.On == "" {
		errs = append(errs, ValidationError{"on", "trigger event is required"})
	}
	for i, job := range w.Jobs {
		if job.Name == "" {
			errs = append(errs, ValidationError{
				fmt.Sprintf("jobs[%d].name", i), "job name is required",
			})
		}
		if job.RunsOn == "" {
			errs = append(errs, ValidationError{
				fmt.Sprintf("jobs[%d].runs-on", i), "runs-on is required",
			})
		}
		for j, step := range job.Steps {
			if step.Name == "" {
				errs = append(errs, ValidationError{
					fmt.Sprintf("jobs[%d].steps[%d].name", i, j), "step name is required",
				})
			}
			if step.Run == "" && step.Uses == "" {
				errs = append(errs, ValidationError{
					fmt.Sprintf("jobs[%d].steps[%d]", i, j),
					"step must have run or uses",
				})
			}
		}
	}
	return errs
}

func ExpandMatrix(matrix map[string][]string) []map[string]string {
	if len(matrix) == 0 {
		return []map[string]string{{}}
	}
	var result []map[string]string
	keys := make([]string, 0, len(matrix))
	values := make([][]string, 0, len(matrix))
	for k, v := range matrix {
		keys = append(keys, k)
		values = append(values, v)
	}
	var recurse func(int, map[string]string)
	recurse = func(idx int, combo map[string]string) {
		if idx == len(keys) {
			cpy := make(map[string]string, len(combo))
			for k, v := range combo {
				cpy[k] = v
			}
			result = append(result, cpy)
			return
		}
		for _, v := range values[idx] {
			combo[keys[idx]] = v
			recurse(idx+1, combo)
		}
	}
	recurse(0, make(map[string]string))
	return result
}

func main() {
	w := Workflow{
		Name: "CI",
		On:   "push",
		Jobs: []Job{
			{
				Name:   "test",
				RunsOn: "ubuntu-latest",
				Steps: []Step{
					{Name: "Checkout", Uses: "actions/checkout@v4"},
					{Name: "Setup Go", Uses: "actions/setup-go@v5",
						With: map[string]string{"go-version": "1.22"}},
					{Name: "Run tests", Run: "go test ./..."},
				},
			},
		},
	}
	errs := ValidateWorkflow(w)
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "validation error: %v\n", e)
		}
		os.Exit(1)
	}
	fmt.Printf("Workflow %q validated OK (%d jobs)\n", w.Name, len(w.Jobs))

	matrix := map[string][]string{
		"go": {"1.21", "1.22"},
		"os": {"ubuntu", "windows"},
	}
	combos := ExpandMatrix(matrix)
	fmt.Printf("Matrix produces %d combinations:\n", len(combos))
	for _, c := range combos {
		fmt.Printf("  go=%s os=%s\n", c["go"], c["os"])
	}
}
```

## Step-by-step execution

For the workflow validation of the example above:

1. `ValidateWorkflow` receives a `Workflow` with `Name="CI"`, `On="push"`, one job.
2. It checks `Name` is not empty (passes).
3. It checks `On` is not empty (passes).
4. For job index 0, it checks `Name="test"` (passes), `RunsOn="ubuntu-latest"` (passes).
5. For steps, it iterates: step 0 has `Uses` (passes), step 1 has `Uses` (passes), step 2 has `Run` (passes).
6. No errors are returned.

For `ExpandMatrix` with `go: [1.21, 1.22]` and `os: [ubuntu, windows]`:

1. `keys = ["go", "os"]`, `values = [["1.21", "1.22"], ["ubuntu", "windows"]]`.
2. Recurse with idx=0, iterate over `["1.21", "1.22"]`:
   - Set `combo["go"]="1.21"`, recurse idx=1.
     - Set `combo["os"]="ubuntu"`, recurse idx=2 (base): copy and append.
     - Set `combo["os"]="windows"`, recurse idx=2 (base): copy and append.
   - Set `combo["go"]="1.22"`, recurse idx=1.
     - Same as above for `os`.
3. Result: 4 combinations.

## Common mistakes

- **Forgetting `runs-on`**: Every job must specify a runner. Missing it causes a workflow parsing error.
- **Step without `run` or `uses`**: Each step must either run a shell command or use an action. An empty step fails validation.
- **Cache key too specific**: If the cache key includes the OS but the workflow runs on multiple OS matrix entries, use separate keys per OS. Otherwise cache misses are common.
- **Matrix explosion**: A matrix with 4 variables each having 4 values produces 256 jobs, hitting GitHub's limit. Keep matrices focused.
- **Not caching the build cache**: Only caching `go mod download` misses the Go build cache. Cache both `~/.cache/go` (module download) and the build cache directory for faster rebuilds.

## Debugging walkthrough

Consider a workflow file that fails validation:

```yaml
name: CI
on: push
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run tests
```

**Symptom**: GitHub shows "Invalid workflow file" with no detailed error in the UI.

**Investigation**: Using the `ValidateWorkflow` function from the Go example:

```go
w := Workflow{
	Name: "CI",
	On:   "push",
	Jobs: []Job{{
		Name:   "test",
		RunsOn: "ubuntu-latest",
		Steps: []Step{
			{Uses: "actions/checkout@v4"},
			{Name: "Run tests"},
		},
	}},
}
errs := ValidateWorkflow(w)
for _, e := range errs {
	fmt.Println(e)
}
```

**Root cause**: The first step (`Uses`) has no `Name` (reported as `steps[0].name required`). The second step (`Name: "Run tests"`) has neither `run` nor `uses` (reported as `steps[1]: step must have run or uses`). Two validation errors surface both issues in one pass.

**Fix**: Add `name` to the checkout step and add `run: go test ./...` to the test step.

## Production notes

In production Go repositories, workflows follow additional conventions:

- **Pin action versions by SHA**: Instead of `actions/checkout@v4`, pin to the full commit SHA (`actions/checkout@a12f...`). Tags can be moved by a malicious actor; SHAs are immutable.
- **Separate lint, test, and deploy workflows**: Keep CI fast by running lint and test in parallel, and only running deploy after both succeed.
- **Use workflow_dispatch for manual triggers**: Add `workflow_dispatch` alongside `push`/`pull_request` to allow manual re-runs from the GitHub UI.
- **Concurrency groups**: Use `concurrency` to cancel duplicate workflow runs on the same branch, saving runner minutes.

## Performance implications

- **Matrix builds increase total CI time linearly** with the number of combinations but reduce wall-clock time by running tests in parallel across runners.
- **Caching cuts `go mod download` from 30-60 seconds to under 2 seconds** on cache hit. The Go build cache (`go build -cache`) saves even more time on incremental builds.
- **Self-hosted runners** can reduce costs for high-volume CI, but require maintenance. GitHub-hosted runners start in seconds and are zero-maintenance.
- **Workflow startup overhead** is typically 10-30 seconds per job (spinning up the runner, cloning the repo). Fewer jobs with matrix steps is faster than many sequential jobs.

## Practice task

Write a function `GenerateWorkflowYAML(w Workflow) string` that produces a minimal YAML representation of a GitHub Actions workflow in Go. Include:

- The workflow name and trigger event.
- Each job with its runner and steps.
- Steps that have a `Uses` field should produce `- uses: ...`; steps with `Run` should produce `- run: ...`.
- Steps with `With` entries should produce `with:` indented block.

Then write a `main()` that creates a workflow with at least two jobs (test and build), generates the YAML, and prints it.

## Tests / verification

```bash
go run ./curriculum/modules/15-docker-cicd-deployment/lessons/11-github-actions
go test ./curriculum/modules/15-docker-cicd-deployment/lessons/11-github-actions
```

The existing tests validate workflow structure, cache key generation, matrix expansion, and runner suggestion. After completing the practice task, add tests for your `GenerateWorkflowYAML` function covering workflows with multiple jobs, actions with parameters, and empty steps.

## Review questions

1. What three fields are required for every GitHub Actions step? What error does a step missing all three produce?
2. A matrix with `go: [1.21, 1.22, 1.23]` and `os: [ubuntu, windows, macos]` produces how many jobs? What is the GitHub limit?
3. Why does caching the Go module cache reduce CI time? What file is typically used to compute the cache key?
4. How does `actions/setup-go@v5` differ from manually downloading and installing Go in a step?
5. What happens when a job's step fails? Does the workflow continue, stop the job, or stop all jobs?

## NEXT UP

CI pipeline: linting, vetting, testing, and building Go code as automated pipeline stages.
