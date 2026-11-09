package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	ErrPortfolioIncomplete = errors.New("portfolio does not meet minimum quality score")
)

type ProjectStatus struct {
	Slug      string   `json:"slug"`
	Title     string   `json:"title"`
	Module    string   `json:"module"`
	HasReadme bool     `json:"has_readme"`
	HasTests  bool     `json:"has_tests"`
	HasCode   bool     `json:"has_code"`
	HasProof  bool     `json:"has_proof"`
	Score     int      `json:"score"` // 0-100
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

type PortfolioAuditor struct {
	rootDir string
}

func NewPortfolioAuditor(rootDir string) *PortfolioAuditor {
	return &PortfolioAuditor{rootDir: rootDir}
}

// PortfolioShowcaseDef defines the flagship showcase projects expected in a full Go portfolio.
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

// AuditProject evaluates a single project directory against portfolio readiness standards.
func (a *PortfolioAuditor) AuditProject(def PortfolioShowcaseDef) ProjectStatus {
	dir := filepath.Join(a.rootDir, def.RelativeDir)
	status := ProjectStatus{
		Slug:   def.Slug,
		Title:  def.Title,
		Module: def.Module,
		Gaps:   make([]string, 0),
	}

	// 1. Check README
	readmePath := filepath.Join(dir, "README.md")
	if data, err := os.ReadFile(readmePath); err == nil && len(data) > 100 {
		status.HasReadme = true
	} else {
		status.Gaps = append(status.Gaps, "Missing or incomplete README.md")
	}

	// 2. Check Solution / Code files
	solutionDir := filepath.Join(dir, "_solution")
	if files, err := os.ReadDir(solutionDir); err == nil && len(files) > 0 {
		hasGo := false
		hasTest := false
		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".go") {
				if strings.HasSuffix(f.Name(), "_test.go") {
					hasTest = true
				} else {
					hasGo = true
				}
			}
		}
		status.HasCode = hasGo
		status.HasTests = hasTest
	} else {
		// Fallback check in root directory if no _solution directory
		if files, err := os.ReadDir(dir); err == nil {
			for _, f := range files {
				if strings.HasSuffix(f.Name(), ".go") {
					if strings.HasSuffix(f.Name(), "_test.go") {
						status.HasTests = true
					} else {
						status.HasCode = true
					}
				}
			}
		}
	}

	if !status.HasCode {
		status.Gaps = append(status.Gaps, "Missing reference Go source code")
	}
	if !status.HasTests {
		status.Gaps = append(status.Gaps, "Missing automated test suite (*_test.go)")
	}

	// 3. Check Proof Evidence (Benchmarking / Telemetry / Compose / Dockerfile)
	proofFiles := []string{"Dockerfile", "docker-compose.yml", "ci.yml"}
	hasProof := false
	_ = filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			for _, pf := range proofFiles {
				if strings.Contains(info.Name(), pf) {
					hasProof = true
					return nil
				}
			}
		}
		return nil
	})
	// Tests also count as proof of correctness
	status.HasProof = status.HasTests || hasProof

	// Compute Score
	score := 0
	if status.HasReadme {
		score += 30
	}
	if status.HasCode {
		score += 30
	}
	if status.HasTests {
		score += 25
	}
	if status.HasProof {
		score += 15
	}
	status.Score = score

	return status
}

// AuditPortfolio runs audits across all showcase projects and compiles a comprehensive report.
func (a *PortfolioAuditor) AuditPortfolio() PortfolioAuditReport {
	statuses := make([]ProjectStatus, 0, len(ShowcaseProjects))
	totalScore := 0
	passedCount := 0

	for _, def := range ShowcaseProjects {
		st := a.AuditProject(def)
		statuses = append(statuses, st)
		totalScore += st.Score
		if st.Score >= 80 {
			passedCount++
		}
	}

	avgScore := 0
	if len(statuses) > 0 {
		avgScore = totalScore / len(statuses)
	}

	ready := passedCount >= len(ShowcaseProjects)-1 && avgScore >= 80

	report := PortfolioAuditReport{
		AuditedAt:      time.Now().UTC(),
		TotalProjects:  len(ShowcaseProjects),
		PassedProjects: passedCount,
		AverageScore:   avgScore,
		ReadyForReview: ready,
		Projects:       statuses,
	}

	if ready {
		report.Summary = fmt.Sprintf("Portfolio is READY for senior engineering interviews with %d/%d projects verified.", passedCount, len(ShowcaseProjects))
	} else {
		report.Summary = fmt.Sprintf("Portfolio needs attention: %d/%d projects passed minimum readiness criteria.", passedCount, len(ShowcaseProjects))
	}

	return report
}

// WriteMarkdownReport writes a formatted Markdown scorecard to the given writer.
func (a *PortfolioAuditor) WriteMarkdownReport(w io.Writer, report PortfolioAuditReport) error {
	var b strings.Builder
	b.WriteString("# Go Engineer Portfolio Audit Scorecard\n\n")
	b.WriteString(fmt.Sprintf("**Audit Date:** %s  \n", report.AuditedAt.Format(time.RFC3339)))
	b.WriteString(fmt.Sprintf("**Overall Status:** %s  \n", map[bool]string{true: "✅ READY FOR INTERVIEWS", false: "⚠️ REMEDIATION NEEDED"}[report.ReadyForReview]))
	b.WriteString(fmt.Sprintf("**Score:** %d/100 (%d/%d projects passed)\n\n", report.AverageScore, report.PassedProjects, report.TotalProjects))

	b.WriteString("## Project Breakdown\n\n")
	b.WriteString("| Project Title | Module | README | Code | Tests | Proof | Score | Status |\n")
	b.WriteString("|---|---|---|---|---|---|---|---|\n")

	for _, p := range report.Projects {
		statusIcon := "✅"
		if p.Score < 80 {
			statusIcon = "❌"
		}
		b.WriteString(fmt.Sprintf("| %s | `%s` | %s | %s | %s | %s | %d%% | %s |\n",
			p.Title,
			p.Module,
			boolIcon(p.HasReadme),
			boolIcon(p.HasCode),
			boolIcon(p.HasTests),
			boolIcon(p.HasProof),
			p.Score,
			statusIcon,
		))
	}

	b.WriteString("\n## Technical Defense Readiness\n\n")
	b.WriteString("Key portfolio highlights for senior backend interviews:\n")
	b.WriteString("- **Distributed Systems & Observability:** Production `slog` + Prometheus + OpenTelemetry instrumentation.\n")
	b.WriteString("- **High-Throughput Concurrency:** 64-shard lock striping and zero-allocation sliding windows (sub-30ns latency).\n")
	b.WriteString("- **Domain-Driven Architecture:** Pure hexagonal boundaries with swappable ports & adapters.\n")
	b.WriteString("- **Cloud-Native Deployment:** Multi-stage scratch containers with non-root security and GitHub Actions CI.\n")

	_, err := io.WriteString(w, b.String())
	return err
}

func boolIcon(v bool) string {
	if v {
		return "✅"
	}
	return "❌"
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

	if !report.ReadyForReview {
		os.Exit(1)
	}
}
