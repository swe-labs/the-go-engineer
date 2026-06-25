package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPipelineSimulate(t *testing.T) {
	p := NewPipeline([]PipelineStage{
		{Name: "lint", Command: "vet ./...", Required: false},
		{Name: "test", Command: "test ./...", Required: true},
	})
	results := p.Simulate()
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.Passed {
			t.Errorf("stage %s should pass in simulation", r.Stage)
		}
	}
}

func TestPipelineSimulateAlwaysPasses(t *testing.T) {
	p := &Pipeline{
		Stages: []PipelineStage{
			{Name: "stage1", Command: "test ./...", Required: true},
			{Name: "stage2", Command: "nonexistent", Required: false},
		},
	}
	results := p.Simulate()
	if len(results) != 2 {
		t.Fatalf("expected 2 results from simulation (always passes), got %d", len(results))
	}
	for _, r := range results {
		if !r.Passed {
			t.Errorf("expected all simulated stages to pass, got fail on %s", r.Stage)
		}
	}
}

func TestPipelineSimulateResultFormat(t *testing.T) {
	p := NewPipeline([]PipelineStage{
		{Name: "build", Command: "build -o /dev/null .", Required: true},
	})
	results := p.Simulate()
	r := results[0]
	if r.Stage != "build" {
		t.Errorf("expected stage build, got %s", r.Stage)
	}
	if !r.Passed {
		t.Errorf("expected pass")
	}
	if r.Output != "[simulated] build -o /dev/null ." {
		t.Errorf("unexpected output: %s", r.Output)
	}
}

func TestPipelineEmptyStages(t *testing.T) {
	p := NewPipeline(nil)
	results := p.Simulate()
	if len(results) != 0 {
		t.Errorf("expected 0 results for empty pipeline, got %d", len(results))
	}
}

func TestIntegrationTestAvailable(t *testing.T) {
	if IntegrationTestAvailable() {
		t.Log("integration directory exists")
	}
}

func TestStageRequiredField(t *testing.T) {
	stage := PipelineStage{Name: "test", Command: "go test", Required: true}
	if !stage.Required {
		t.Error("stage should be required")
	}
	stage2 := PipelineStage{Name: "lint", Command: "go vet", Required: false}
	if stage2.Required {
		t.Error("stage should not be required")
	}
}

func initGoModule(t *testing.T, dir string) {
	t.Helper()
	err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module testmodule\ngo 1.22\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLintPackage(t *testing.T) {
	dir := t.TempDir()
	initGoModule(t, dir)
	err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\nfunc main() {}\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	_, err = LintPackage(dir)
	if err != nil {
		t.Fatalf("expected no lint error on valid code: %v", err)
	}
}

func TestVetPackage(t *testing.T) {
	dir := t.TempDir()
	initGoModule(t, dir)
	err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\nfunc main() {}\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	_, err = VetPackage(dir)
	if err != nil {
		t.Fatalf("expected no vet error on valid code: %v", err)
	}
}

func TestBuildPackage(t *testing.T) {
	dir := t.TempDir()
	initGoModule(t, dir)
	err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\nimport \"fmt\"\nfunc main() { fmt.Println(\"hi\") }\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	outPath := filepath.Join(dir, "output")
	_, err = BuildPackage(dir, outPath)
	if err != nil {
		t.Fatalf("expected build to succeed: %v", err)
	}
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		t.Error("expected binary to exist after build")
	}
}
