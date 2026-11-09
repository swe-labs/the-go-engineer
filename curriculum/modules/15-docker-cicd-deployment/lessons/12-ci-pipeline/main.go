package main

import (
	"fmt"
	"os"
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
		result := PipelineResult{Stage: stage.Name, Passed: true, Output: fmt.Sprintf("[simulated] %s", stage.Command)}
		results = append(results, result)
		if !result.Passed && stage.Required {
			break
		}
	}
	return results
}

func LintPackage(dir string) (string, error) {
	cmd := exec.Command("go", "vet", ".")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func VetPackage(path string) (string, error) {
	return LintPackage(path)
}

func TestPackage(path string, race bool, shuffle bool) (string, error) {
	args := []string{"test"}
	if race {
		args = append(args, "-race")
	}
	if shuffle {
		args = append(args, "-shuffle=on")
	}
	args = append(args, path)
	cmd := exec.Command("go", args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func BuildPackage(dir, output string) (string, error) {
	args := []string{"build", "-o", output, "."}
	cmd := exec.Command("go", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func IntegrationTestAvailable() bool {
	_, err := os.Stat("./integration")
	return err == nil
}

func main() {
	fmt.Println("=== CI Pipeline ===")
	fmt.Println()

	pipeline := NewPipeline([]PipelineStage{
		{Name: "lint", Command: "vet ./...", Required: false},
		{Name: "test", Command: "test -race -shuffle=on ./...", Required: true},
		{Name: "build", Command: "build -o /dev/null ./...", Required: true},
	})

	fmt.Println("Pipeline stages:")
	for _, s := range pipeline.Stages {
		required := ""
		if s.Required {
			required = " (required)"
		}
		fmt.Printf("  %s: go %s%s\n", s.Name, s.Command, required)
	}

	results := pipeline.Simulate()
	fmt.Println()
	fmt.Println("Simulated results:")
	for _, r := range results {
		status := "PASS"
		if !r.Passed {
			status = "FAIL"
		}
		fmt.Printf("  %s: %s\n", r.Stage, status)
	}

	fmt.Println()
	fmt.Println("CI vs CD:")
	fmt.Println("  CI: Continuous Integration — merge code frequently, run automated checks")
	fmt.Println("  CD: Continuous Delivery/Deployment — automatically deploy to production")
}
