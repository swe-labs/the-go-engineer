package main

import (
	"os"
	"testing"
)

func TestAnalyzeFile_DetectsUndocumented(t *testing.T) {
	content := `package test
func Exported() {}
func unexported() {}`
	tmpFile := writeTempGoFile(t, content)
	defer os.Remove(tmpFile)

	report, err := AnalyzeFile(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	if report.Total != 1 {
		t.Fatalf("expected 1 exported symbol, got %d", report.Total)
	}
	if report.Items[0].Status != Missing {
		t.Errorf("expected doc comment to be missing")
	}
}

func TestAnalyzeFile_DetectsDocumented(t *testing.T) {
	content := `package test
// Exported does something.
func Exported() {}`
	tmpFile := writeTempGoFile(t, content)
	defer os.Remove(tmpFile)

	report, err := AnalyzeFile(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	if report.Total != 1 {
		t.Fatalf("expected 1 exported symbol, got %d", report.Total)
	}
	if report.Items[0].Status != Present {
		t.Errorf("expected doc comment to be present")
	}
}

func TestAnalyzeFile_CoverageCalculation(t *testing.T) {
	content := `package test
// Documented is documented.
func Documented() {}
func Undocumented() {}
// SomeType is a type.
type SomeType struct{}`
	tmpFile := writeTempGoFile(t, content)
	defer os.Remove(tmpFile)

	report, err := AnalyzeFile(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	if report.Total != 3 {
		t.Fatalf("expected 3 exported symbols, got %d", report.Total)
	}
	if report.Covered != 2 {
		t.Errorf("expected 2 documented, got %d", report.Covered)
	}
	if report.Coverage() != 66.66666666666666 {
		t.Errorf("expected ~66.67%% coverage, got %.1f%%", report.Coverage())
	}
}

func TestAnalyzeFile_NoExportedSymbols(t *testing.T) {
	content := `package test
func unexported() {}`
	tmpFile := writeTempGoFile(t, content)
	defer os.Remove(tmpFile)

	report, err := AnalyzeFile(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	if report.Total != 0 {
		t.Errorf("expected 0 exported symbols, got %d", report.Total)
	}
}

func writeTempGoFile(t *testing.T, content string) string {
	t.Helper()
	tmpFile := t.TempDir() + "/testfile.go"
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return tmpFile
}
