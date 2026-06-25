package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
			errs = append(errs, ValidationError{fmt.Sprintf("jobs[%d].name", i), "job name is required"})
		}
		if job.RunsOn == "" {
			errs = append(errs, ValidationError{fmt.Sprintf("jobs[%d].runs-on", i), "runs-on is required"})
		}
		for j, step := range job.Steps {
			if step.Name == "" {
				errs = append(errs, ValidationError{fmt.Sprintf("jobs[%d].steps[%d].name", i, j), "step name is required"})
			}
			if step.Run == "" && step.Uses == "" {
				errs = append(errs, ValidationError{fmt.Sprintf("jobs[%d].steps[%d]", i, j), "step must have run or uses"})
			}
			if step.Timeout > 360 {
				errs = append(errs, ValidationError{fmt.Sprintf("jobs[%d].steps[%d].timeout", i, j), "timeout exceeds 360 minutes"})
			}
		}
	}
	return errs
}

func GoCacheKey(goVersion, osName string, sumHash string) string {
	return fmt.Sprintf("go-cache-%s-%s-%s", goVersion, osName, sumHash[:8])
}

func GenerateCacheStep(key, restoreKey, path string) Step {
	return Step{
		Name: "Cache Go modules",
		Uses: "actions/cache@v4",
		With: map[string]string{
			"key":          key,
			"restore-keys": restoreKey,
			"path":         path,
		},
	}
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

func SuggestRunner(goVersion string, osSuffix string) string {
	if osSuffix == "" {
		return "ubuntu-latest"
	}
	if strings.Contains(osSuffix, "win") {
		return "windows-latest"
	}
	if strings.Contains(osSuffix, "mac") {
		return "macos-latest"
	}
	return "ubuntu-latest"
}

func init() {
	_ = filepath.Join
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
					{Name: "Setup Go", Uses: "actions/setup-go@v5", With: map[string]string{"go-version": "1.22"}},
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

	cacheStep := GenerateCacheStep(
		GoCacheKey("1.22", "ubuntu", "abcdef1234567890"),
		"go-cache-1.22-ubuntu-",
		"~/.cache/go",
	)
	fmt.Printf("Cache step: %s (key=%s)\n", cacheStep.Name, cacheStep.With["key"])

	matrix := map[string][]string{
		"go": {"1.21", "1.22"},
		"os": {"ubuntu", "windows"},
	}
	combos := ExpandMatrix(matrix)
	fmt.Printf("Matrix combinations: %d\n", len(combos))
	for _, c := range combos {
		fmt.Printf("  go=%s os=%s runner=%s\n", c["go"], c["os"], SuggestRunner(c["go"], c["os"]))
	}
}
