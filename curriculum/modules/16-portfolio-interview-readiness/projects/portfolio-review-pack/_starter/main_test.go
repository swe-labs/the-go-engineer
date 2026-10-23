package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func findRepoRoot(t *testing.T) string {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("could not find repository root containing go.mod from %s", dir)
		}
		dir = parent
	}
}

func TestPortfolioAuditor_RealShowcaseAudit(t *testing.T) {
	root := findRepoRoot(t)
	auditor := NewPortfolioAuditor(root)

	report := auditor.AuditPortfolio()

	if report.TotalProjects != len(ShowcaseProjects) {
		t.Errorf("expected %d total projects, got %d", len(ShowcaseProjects), report.TotalProjects)
	}

	for _, p := range report.Projects {
		if p.Score < 80 {
			t.Errorf("project %s scored %d (<80%%). Gaps: %v", p.Title, p.Score, p.Gaps)
		}
	}

	if !report.ReadyForReview {
		t.Errorf("expected portfolio to be ready for review, summary: %s", report.Summary)
	}
}

func TestPortfolioAuditor_SyntheticGapsDetection(t *testing.T) {
	tmpDir := t.TempDir()

	projDir := filepath.Join(tmpDir, "sample-project", "_solution")
	if err := os.MkdirAll(projDir, 0755); err != nil {
		t.Fatalf("failed to create temp project dir: %v", err)
	}
	_ = os.WriteFile(filepath.Join(projDir, "main.go"), []byte("package main"), 0644)
	_ = os.WriteFile(filepath.Join(projDir, "main_test.go"), []byte("package main"), 0644)

	auditor := NewPortfolioAuditor(tmpDir)
	status := auditor.AuditProject(PortfolioShowcaseDef{
		Slug:        "sample-project",
		Title:       "Sample Project",
		Module:      "module-99",
		RelativeDir: "sample-project",
	})

	if status.HasReadme {
		t.Errorf("expected HasReadme to be false")
	}
	if !status.HasCode || !status.HasTests {
		t.Errorf("expected HasCode and HasTests to be true")
	}
	if status.Score >= 100 {
		t.Errorf("expected score to be docked for missing README, got %d", status.Score)
	}

	foundGap := false
	for _, g := range status.Gaps {
		if strings.Contains(g, "README") {
			foundGap = true
			break
		}
	}
	if !foundGap {
		t.Errorf("expected gap for missing README, got: %v", status.Gaps)
	}
}

func TestPortfolioAuditor_MarkdownReportGeneration(t *testing.T) {
	root := findRepoRoot(t)
	auditor := NewPortfolioAuditor(root)
	report := auditor.AuditPortfolio()

	var buf bytes.Buffer
	err := auditor.WriteMarkdownReport(&buf, report)
	if err != nil {
		t.Fatalf("failed to write markdown report: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "# Go Engineer Portfolio Audit Scorecard") {
		t.Errorf("missing report title")
	}
	if !strings.Contains(out, "READY FOR INTERVIEWS") {
		t.Errorf("missing ready status in report:\n%s", out)
	}
	if !strings.Contains(out, "Book Catalog REST API") {
		t.Errorf("missing showcase item row")
	}
}
