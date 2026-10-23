package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

type ProjectStatus struct {
	Slug      string   `json:"slug"`
	Title     string   `json:"title"`
	Module    string   `json:"module"`
	HasReadme bool     `json:"has_readme"`
	HasTests  bool     `json:"has_tests"`
	HasCode   bool     `json:"has_code"`
	HasProof  bool     `json:"has_proof"`
	Score     int      `json:"score"`
	Gaps      []string `json:"gaps,omitempty"`
}

type PortfolioAuditReport struct {
	AuditedAt      time.Time       `json:"audited_at"`
	TotalProjects  int             `json:"total_projects"`
	PassedProjects int             `json:"passed_projects"`
	AverageScore   int             `json:"average_score"`
	ReadyForReview bool            `json:"ready_for_review"`
	Projects       []ProjectStatus `json:"projects"`
	Summary        string          `json:"summary"`
}

type PortfolioShowcaseDef struct {
	Slug        string
	Title       string
	Module      string
	RelativeDir string
}

var ShowcaseProjects = []PortfolioShowcaseDef{
	{
		Slug:        "book-catalog-rest-api",
		Title:       "Book Catalog REST API",
		Module:      "module-08",
		RelativeDir: "curriculum/modules/08-http-rest-apis/projects/book-catalog-rest-api",
	},
	{
		Slug:        "task-api-with-postgres",
		Title:       "Task API with Postgres & Migrations",
		Module:      "module-09",
		RelativeDir: "curriculum/modules/09-sql-postgres-persistence/projects/task-api-with-postgres",
	},
	{
		Slug:        "secure-authenticated-api",
		Title:       "Secure Authenticated RBAC API",
		Module:      "module-10",
		RelativeDir: "curriculum/modules/10-auth-security/projects/secure-authenticated-api",
	},
	{
		Slug:        "bounded-worker-pool",
		Title:       "Bounded Worker Pool with Backpressure",
		Module:      "module-11",
		RelativeDir: "curriculum/modules/11-lifecycle-context-concurrency/projects/bounded-worker-pool",
	},
	{
		Slug:        "observability-retrofit",
		Title:       "Observability Retrofit (slog + Prometheus + OTel)",
		Module:      "module-12",
		RelativeDir: "curriculum/modules/12-observability-diagnostics/projects/observability-retrofit",
	},
	{
		Slug:        "performance-analysis-lab",
		Title:       "High-Throughput Lock-Striped Rate Limiter",
		Module:      "module-13",
		RelativeDir: "curriculum/modules/13-performance-memory-engineering/projects/performance-analysis-lab",
	},
	{
		Slug:        "modular-refactor-review",
		Title:       "Modular Hexagonal Architecture Refactor",
		Module:      "module-14",
		RelativeDir: "curriculum/modules/14-architecture-distributed-systems/projects/modular-refactor-review",
	},
	{
		Slug:        "dockerized-service-with-ci-cd",
		Title:       "Dockerized Service with Multi-Stage CI/CD",
		Module:      "module-15",
		RelativeDir: "curriculum/modules/15-docker-cicd-deployment/projects/dockerized-service-with-ci-cd",
	},
}

type PortfolioAuditor struct {
	rootDir string
}

func NewPortfolioAuditor(rootDir string) *PortfolioAuditor {
	return &PortfolioAuditor{rootDir: rootDir}
}

// TODO: Task 1 - Audit single project directory
// Check for presence of README.md, Go source files, tests (*_test.go), and proof artifacts (Dockerfile, compose, ci.yml).
func (a *PortfolioAuditor) AuditProject(def PortfolioShowcaseDef) ProjectStatus {
	// TODO: Implement file system checks and populate ProjectStatus
	return ProjectStatus{
		Slug:   def.Slug,
		Title:  def.Title,
		Module: def.Module,
		Score:  0,
		Gaps:   []string{"Audit not implemented"},
	}
}

// TODO: Task 2 - Audit all showcase projects and compute overall readiness score
func (a *PortfolioAuditor) AuditPortfolio() PortfolioAuditReport {
	// TODO: Iterate over ShowcaseProjects, evaluate status, compute average score,
	// and determine if ReadyForReview (score >= 80% across projects).
	return PortfolioAuditReport{
		AuditedAt:     time.Now().UTC(),
		TotalProjects: len(ShowcaseProjects),
	}
}

// TODO: Task 3 - Render Markdown scorecard
func (a *PortfolioAuditor) WriteMarkdownReport(w io.Writer, report PortfolioAuditReport) error {
	// TODO: Write markdown table with icons and score breakdown
	_, err := fmt.Fprintf(w, "# Portfolio Audit Report\nTotal Projects: %d\n", report.TotalProjects)
	return err
}

func main() {
	jsonOutput := flag.Bool("json", false, "output report as JSON")
	flag.Parse()

	auditor := NewPortfolioAuditor(".")
	report := auditor.AuditPortfolio()

	if *jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(report)
	} else {
		_ = auditor.WriteMarkdownReport(os.Stdout, report)
	}
}
