package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type RunbookPhase string

const (
	PhasePreDeploy  RunbookPhase = "pre-deploy"
	PhaseDeploy     RunbookPhase = "deploy"
	PhasePostDeploy RunbookPhase = "post-deploy"
	PhaseRollback   RunbookPhase = "rollback"
)

type RunbookStep struct {
	Phase      RunbookPhase
	Name       string
	Command    string
	Timeout    time.Duration
	Required   bool
	CheckAfter string
}

type Runbook struct {
	ServiceName string
	Version     string
	Steps       []RunbookStep
	Contacts    []string
}

func NewRunbook(serviceName, version string) *Runbook {
	return &Runbook{
		ServiceName: serviceName,
		Version:     version,
		Contacts:    []string{},
	}
}

func (rb *Runbook) AddStep(step RunbookStep) {
	rb.Steps = append(rb.Steps, step)
}

func (rb *Runbook) AddContact(contact string) {
	rb.Contacts = append(rb.Contacts, contact)
}

func (rb *Runbook) ExecutePhase(phase RunbookPhase) (bool, []string) {
	var executed []string
	for _, step := range rb.Steps {
		if step.Phase != phase {
			continue
		}
		executed = append(executed, fmt.Sprintf("[%s] %s...", phase, step.Name))
		if step.CheckAfter != "" {
			executed = append(executed, fmt.Sprintf("  -> verify: %s", step.CheckAfter))
		}
	}
	return true, executed
}

func (rb *Runbook) GenerateChecklist() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("=== Deployment Runbook: %s v%s ===\n", rb.ServiceName, rb.Version))
	b.WriteString(fmt.Sprintf("Contacts: %s\n\n", strings.Join(rb.Contacts, ", ")))
	phases := []RunbookPhase{PhasePreDeploy, PhaseDeploy, PhasePostDeploy, PhaseRollback}
	for _, phase := range phases {
		b.WriteString(fmt.Sprintf("--- %s ---\n", strings.ToUpper(string(phase))))
		for _, step := range rb.Steps {
			if step.Phase != phase {
				continue
			}
			req := ""
			if step.Required {
				req = " [REQUIRED]"
			}
			b.WriteString(fmt.Sprintf("  [] %s%s\n", strings.ReplaceAll(step.Name, "_", " "), req))
			if step.Command != "" {
				b.WriteString(fmt.Sprintf("     $ %s\n", step.Command))
			}
			if step.CheckAfter != "" {
				b.WriteString(fmt.Sprintf("     verify: %s\n", step.CheckAfter))
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (rb *Runbook) Validate() []error {
	var errs []error
	if rb.ServiceName == "" {
		errs = append(errs, fmt.Errorf("service name is required"))
	}
	if rb.Version == "" {
		errs = append(errs, fmt.Errorf("version is required"))
	}
	if len(rb.Contacts) == 0 {
		errs = append(errs, fmt.Errorf("at least one incident contact is required"))
	}
	phaseCounts := make(map[RunbookPhase]int)
	for _, step := range rb.Steps {
		phaseCounts[step.Phase]++
		if step.Name == "" {
			errs = append(errs, fmt.Errorf("step name is required"))
		}
	}
	if phaseCounts[PhasePreDeploy] == 0 {
		errs = append(errs, fmt.Errorf("at least one pre-deploy step is required"))
	}
	if phaseCounts[PhaseDeploy] == 0 {
		errs = append(errs, fmt.Errorf("at least one deploy step is required"))
	}
	if phaseCounts[PhasePostDeploy] == 0 {
		errs = append(errs, fmt.Errorf("at least one post-deploy step is required"))
	}
	return errs
}

func (rb *Runbook) Print() {
	fmt.Println(rb.GenerateChecklist())
}

func main() {
	rb := NewRunbook("payment-service", "v2.1.0")
	rb.AddContact("oncall@company.com")
	rb.AddContact("#ops-slack")

	rb.AddStep(RunbookStep{
		Phase: PhasePreDeploy, Name: "Verify CI pipeline passed",
		Command:    "gh run list --workflow ci.yml --branch main --limit 1 --json conclusion",
		CheckAfter: "conclusion == success", Required: true,
	})
	rb.AddStep(RunbookStep{
		Phase: PhasePreDeploy, Name: "Check database migration status",
		Command:    "go run ./cmd/migrate --dry-run",
		CheckAfter: "no pending migrations", Required: true,
	})
	rb.AddStep(RunbookStep{
		Phase: PhaseDeploy, Name: "Build and push Docker image",
		Command:  "docker build -t payment-service:v2.1.0 . && docker push",
		Required: true,
	})
	rb.AddStep(RunbookStep{
		Phase: PhaseDeploy, Name: "Deploy to staging",
		Command:    "kubectl set image deployment/payment-service app=payment-service:v2.1.0 -n staging",
		CheckAfter: "kubectl rollout status deployment/payment-service -n staging",
		Required:   true,
	})
	rb.AddStep(RunbookStep{
		Phase: PhasePostDeploy, Name: "Verify health endpoint in staging",
		Command:  "curl -f http://staging.payment/healthz",
		Required: true,
	})
	rb.AddStep(RunbookStep{
		Phase: PhasePostDeploy, Name: "Smoke test critical paths",
		Command:  "go run ./cmd/smoke-test --env staging",
		Required: true,
	})
	rb.AddStep(RunbookStep{
		Phase: PhaseRollback, Name: "Rollback to previous version",
		Command:    "kubectl rollout undo deployment/payment-service -n staging",
		CheckAfter: "kubectl rollout status deployment/payment-service -n staging",
		Required:   true,
	})

	fmt.Println("=== Deployment Runbook ===")
	rb.Print()

	errs := rb.Validate()
	if len(errs) > 0 {
		fmt.Fprintf(os.Stderr, "Validation errors:\n")
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "  - %v\n", e)
		}
		os.Exit(1)
	}
	fmt.Println("Runbook validation: OK")
}
