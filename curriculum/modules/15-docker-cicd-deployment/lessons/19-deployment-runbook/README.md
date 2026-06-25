# Deployment runbook

## Learning objective

Design, validate, and generate deployment runbooks in Go that structure pre-deploy checks, deploy steps, post-deploy verification, rollback procedures, and incident contacts for a production Go service.

## Why this matters

A deployment runbook is the canonical source of truth for how to deploy a service. Without a runbook, each deployment relies on tribal knowledge: one engineer knows the steps, but if they are unavailable, the deployment stalls or mistakes are made. A runbook standardizes the process, reduces human error, provides a checklist to follow under pressure, and is the first place to look during an incident. For Go services deployed across multiple environments, a well-structured runbook ensures every team member deploys the same way every time.

## Mental model

A deployment runbook is a recipe for deploying a specific service. Like a cooking recipe, it lists:

- **Ingredients**: What you need (artifacts, credentials, access).
- **Prep work**: What to check before starting (pre-deploy).
- **Cooking steps**: How to deploy (deploy).
- **Taste test**: How to verify it worked (post-deploy).
- **Emergency procedure**: What to do if something goes wrong (rollback).
- **Chef contact**: Who to call for help (incident contacts).

The runbook is versioned alongside the service code. When the service changes (new configuration, new dependencies), the runbook is updated in the same PR. This keeps the runbook accurate and reviewable.

## Core idea

A deployment runbook has four phases:

| Phase | Purpose | Example steps |
|---|---|---|
| Pre-deploy | Verify readiness before making changes | CI passed, DB migration status, config validation |
| Deploy | Execute the deployment | Build image, push to registry, update orchestrator |
| Post-deploy | Verify the deployment succeeded | Health check, smoke tests, metric check |
| Rollback | Revert if deployment fails | Switch traffic, redeploy previous version, run down migration |

Each step has:
- **Name**: What to do.
- **Command**: The exact command to run.
- **CheckAfter**: How to verify the step succeeded.
- **Required**: Whether the pipeline should stop if this step fails.
- **Phase**: Which phase the step belongs to.

Incident contacts are the people and channels to notify if any step fails unexpectedly.

## Under the hood

A runbook is more than a checklist. It is a decision tree:

```
Pre-deploy checks pass?
  YES → Continue to deploy
  NO  → Stop. Investigate. Do not proceed.

Deploy completes?
  YES → Continue to post-deploy checks
  NO  → Execute rollback phase

Post-deploy checks pass?
  YES → Deployment is complete
  NO  → Execute rollback phase

Rollback completes?
  YES → System is on previous version. Notify team.
  NO  → Escalate to incident response.
```

The runbook can be automated as a CI/CD pipeline or followed manually by an operator. In either case, the structure is the same. Automation executes the commands; manual execution prints them for the operator to run.

## How Go uses it

Go teams typically store runbooks in the repository under `deploy/runbook.md` or as a Go-based tool in `cmd/deploy-runbook`. The Go tool reads the runbook definition (from code or YAML), validates it, and prints a formatted checklist or executes steps interactively.

Key practices:

- The runbook is updated in the same PR as the code changes that affect deployment.
- Runbook changes require review, just like code changes.
- The runbook must be executable by any team member, not just the service owner.
- Incident contacts are reviewed quarterly and updated in the runbook.

## Go example

```go
package main

import (
	"fmt"
	"strings"
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
	}
}

func (rb *Runbook) AddStep(step RunbookStep) {
	rb.Steps = append(rb.Steps, step)
}

func (rb *Runbook) AddContact(contact string) {
	rb.Contacts = append(rb.Contacts, contact)
}

func (rb *Runbook) GenerateChecklist() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("=== Deployment: %s v%s ===\n",
		rb.ServiceName, rb.Version))
	b.WriteString(fmt.Sprintf("Contacts: %s\n\n",
		strings.Join(rb.Contacts, ", ")))

	phases := []RunbookPhase{
		PhasePreDeploy, PhaseDeploy,
		PhasePostDeploy, PhaseRollback,
	}
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
			b.WriteString(fmt.Sprintf("  [] %s%s\n", step.Name, req))
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
	phases := map[RunbookPhase]bool{}
	for _, step := range rb.Steps {
		phases[step.Phase] = true
		if step.Name == "" {
			errs = append(errs, fmt.Errorf("step name is required"))
		}
	}
	if !phases[PhasePreDeploy] {
		errs = append(errs, fmt.Errorf("pre-deploy phase is required"))
	}
	if !phases[PhaseDeploy] {
		errs = append(errs, fmt.Errorf("deploy phase is required"))
	}
	if !phases[PhasePostDeploy] {
		errs = append(errs, fmt.Errorf("post-deploy phase is required"))
	}
	return errs
}

func main() {
	rb := NewRunbook("payment-service", "v2.1.0")
	rb.AddContact("oncall@company.com")
	rb.AddContact("#ops-slack")

	rb.AddStep(RunbookStep{
		Phase: PhasePreDeploy, Name: "Verify CI passed",
		Command: "gh run list --workflow ci.yml --branch main",
		Required: true,
	})
	rb.AddStep(RunbookStep{
		Phase: PhaseDeploy, Name: "Deploy to production",
		Command: "kubectl set image deploy/payment app=v2.1.0",
		CheckAfter: "kubectl rollout status deploy/payment",
		Required: true,
	})
	rb.AddStep(RunbookStep{
		Phase: PhasePostDeploy, Name: "Verify health endpoint",
		Command: "curl -f https://payment.example.com/healthz",
		Required: true,
	})
	rb.AddStep(RunbookStep{
		Phase: PhaseRollback, Name: "Rollback to v2.0.0",
		Command: "kubectl rollout undo deploy/payment",
		Required: true,
	})

	fmt.Println(rb.GenerateChecklist())
}
```

## Step-by-step execution

For the runbook above when validated:

1. `NewRunbook("payment-service", "v2.1.0")` creates the runbook with service name and version.
2. `AddContact("oncall@company.com")` adds an email contact.
3. `AddStep` adds four steps across all four phases.
4. `Validate()` checks:
   - Service name not empty (passes).
   - Version not empty (passes).
   - At least one contact (passes).
   - All four phases have at least one step (passes).
   - All step names are non-empty (passes).
5. Returns empty errors.

Executing the pre-deploy phase via `GenerateChecklist()` produces:

```
--- PRE-DEPLOY ---
  [] Verify CI passed [REQUIRED]
     $ gh run list --workflow ci.yml --branch main
```

This checklist is printed to the operator's terminal or posted in a deployment channel.

## Common mistakes

- **Runbook not updated with service changes**: If a deploy command changes but the runbook is not updated, the runbook becomes dangerous. Inconsistent runbooks lead to mistakes.
- **No rollback phase**: Every deployment must have a rollback plan. If the runbook does not say how to rollback, it is incomplete.
- **Missing verification steps**: Deploying and assuming it worked is how incidents start. Always include health check and smoke test steps.
- **Outdated incident contacts**: If the on-call engineer changed but the runbook still lists the old contact, escalation is delayed. Review contacts quarterly.
- **Runbook too long**: A runbook with 50 steps is overwhelming and will be skipped. Keep deploy phases to 3-5 steps each. Detailed troubleshooting goes in a separate troubleshooting guide.

## Debugging walkthrough

Consider a deployment that failed during post-deploy verification:

```
$ curl -f https://payment.example.com/healthz
curl: (7) Failed to connect to payment.example.com port 443: Connection refused
```

**Symptom**: The health check command in the post-deploy step fails.

**Investigation**:
1. Check the deploy step: was the image tag correct? `kubectl get deployment payment -o jsonpath='{.spec.template.spec.containers[0].image}'` should show `v2.1.0`.
2. Check pod status: `kubectl get pods -l app=payment`. Are pods running or CrashLoopBackOff?
3. Check logs: `kubectl logs -l app=payment --tail=50`.

**Root cause**: The new version has a configuration error. The Go server starts but fails to bind to the port because the config file is missing a required field.

**Fix**: Execute the rollback step: `kubectl rollout undo deployment/payment`. This reverts to v2.0.0, and the health check passes again.

**Post-incident**: Add a pre-deploy step that validates the config before deploying:

```
  [] Validate config file [REQUIRED]
     $ go run ./cmd/validate-config --config config/prod.yaml
     verify: exit code == 0
```

Another scenario: The runbook says to run `kubectl set image` but the actual deployment uses Helm. The runbook is wrong.

**Root cause**: The runbook was not updated when the team switched from kubectl to Helm.

**Fix**: Update the runbook to reflect the current deployment tooling. Add a review step to the PR checklist: "Has the runbook been updated?"

## Production notes

In production deployment runbooks:

- **Automate what you can, document what you cannot**: Automated steps (CI, image build, rollout) should be executed by the pipeline. Manual steps (verification, rollback decision) should be clear checkboxes.
- **Runbook drill**: Practice the runbook quarterly. Simulate a deployment and a rollback. Time how long it takes. Improve slow or confusing steps.
- **Runbook as code**: Store the runbook in the repository as a Go program or YAML file. Version it alongside the service. Generate the human-readable output from the structured definition.
- **Blame-free postmortems**: When a runbook step fails or is incorrect, file a bug to fix the runbook. The goal is a runbook so good that anyone can deploy the service correctly.

## Performance implications

- **Runbook readability impacts deployment speed**: A well-organized runbook with clear commands reduces deployment time by 50-70% compared to ad-hoc deployment.
- **Automated runbooks execute in seconds**: A fully automated runbook (GitHub Actions + kubectl) takes 2-5 minutes. A manual runbook takes 15-30 minutes.
- **Checklist overhead**: Printing a checklist adds zero runtime overhead. The overhead is in the human reading and following it, which is the desired behavior.
- **Validating runbook structure is O(n)**: Validating N steps requires iterating over them once. Even a runbook with 100 steps validates in microseconds.

## Practice task

Write a function `ParseRunbookYAML(yamlContent string) (*Runbook, error)` that parses a YAML-like format into a `Runbook` struct:

```yaml
service: payment-service
version: v2.1.0
contacts:
  - oncall@company.com
steps:
  - phase: pre-deploy
    name: Verify CI
    command: gh run list
    required: true
  - phase: deploy
    name: Deploy
    command: kubectl set image deploy/payment app=v2.1.0
```

Then write a `main()` that reads the YAML content as a string literal, parses it, validates it, and prints the checklist. Handle errors gracefully.

## Tests / verification

```bash
go run ./curriculum/modules/15-docker-cicd-deployment/lessons/19-deployment-runbook
go test ./curriculum/modules/15-docker-cicd-deployment/lessons/19-deployment-runbook
```

The existing tests verify runbook creation, step and contact addition, checklist generation (all four phases present, service name, contacts), runbook validation (missing name, version, contacts, phases, empty step names), phase execution, and empty phase execution. After completing the practice task, add tests for `ParseRunbookYAML` covering valid YAML, missing fields, and malformed YAML.

## Review questions

1. What are the four required phases of a deployment runbook? What happens if one is missing?
2. Why should a runbook be stored in the repository alongside the service code rather than in a wiki or shared drive?
3. How does a runbook differ from a troubleshooting guide? When would you use each?
4. A deployment fails at the post-deploy health check. What does the runbook say to do? What if the rollback also fails?
5. What information should be included in each runbook step beyond the command to run?

## NEXT UP

Congratulations on completing Module 15! You now understand Docker, CI/CD, deployment strategies, and DevOps practices for Go services. Next up: Module 16 — Microservices, gRPC, and Service Mesh.
