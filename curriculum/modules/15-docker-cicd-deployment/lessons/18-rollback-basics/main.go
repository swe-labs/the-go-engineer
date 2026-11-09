package main

import (
	"fmt"
	"time"
)

type RollbackReason string

const (
	ReasonHighErrorRate   RollbackReason = "high_error_rate"
	ReasonLatencySpike    RollbackReason = "latency_spike"
	ReasonHealthCheckFail RollbackReason = "health_check_failure"
	ReasonManual          RollbackReason = "manual_trigger"
	ReasonDatabaseIssue   RollbackReason = "database_issue"
)

type RollbackPlan struct {
	Version       string
	Reason        RollbackReason
	Steps         []string
	EstimatedTime time.Duration
	DataLossRisk  bool
}

type DeploymentState struct {
	CurrentVersion  string
	PreviousVersion string
	Healthy         bool
	ErrorRate       float64
	FailedDeploy    bool
}

func NewRollbackPlan(prevVersion string, reason RollbackReason) *RollbackPlan {
	plan := &RollbackPlan{
		Version:       prevVersion,
		Reason:        reason,
		EstimatedTime: 3 * time.Minute,
		DataLossRisk:  false,
	}
	switch reason {
	case ReasonHighErrorRate:
		plan.EstimatedTime = 2 * time.Minute
		plan.Steps = []string{
			"Revert load balancer traffic to previous version",
			"Scale down new version instances to zero",
			"Verify error rate returns to baseline",
			"Post-mortem: investigate root cause of errors",
		}
	case ReasonLatencySpike:
		plan.EstimatedTime = 5 * time.Minute
		plan.Steps = []string{
			"Gradually shift traffic back to previous version",
			"Run performance comparison between old and new",
			"Profile new version for slow code paths",
		}
	case ReasonHealthCheckFail:
		plan.EstimatedTime = 1 * time.Minute
		plan.Steps = []string{
			"Stop deploy: new instances never passed health checks",
			"Revert load balancer to old target group",
			"Check new version startup logs for errors",
		}
	case ReasonDatabaseIssue:
		plan.EstimatedTime = 15 * time.Minute
		plan.Steps = []string{
			"Revert application code to previous version",
			"Run database rollback migration (if applicable)",
			"Verify data integrity",
			"Scale down new version instances",
		}
		if reason == ReasonDatabaseIssue {
			plan.DataLossRisk = true
		}
	default:
		plan.Steps = []string{
			fmt.Sprintf("Redeploy version %s", prevVersion),
			"Verify health checks pass",
			"Monitor for 5 minutes post-rollback",
		}
	}
	return plan
}

func ShouldRollback(state DeploymentState, thresholdErrorRate float64, thresholdLatencyMs float64) (bool, RollbackReason) {
	if !state.Healthy {
		return true, ReasonHealthCheckFail
	}
	if state.FailedDeploy {
		return true, ReasonManual
	}
	return false, ""
}

type RollbackAutomation struct {
	State      DeploymentState
	Thresholds struct {
		ErrorRate    float64
		LatencyRatio float64
	}
}

func (ra *RollbackAutomation) Evaluate() (*RollbackPlan, bool) {
	shouldRollback, reason := ShouldRollback(ra.State, ra.Thresholds.ErrorRate, ra.Thresholds.LatencyRatio)
	if !shouldRollback {
		return nil, false
	}
	plan := NewRollbackPlan(ra.State.PreviousVersion, reason)
	return plan, true
}

func main() {
	fmt.Println("=== Rollback Basics ===")
	fmt.Println()

	state := DeploymentState{
		CurrentVersion:  "v2.0.0",
		PreviousVersion: "v1.0.0",
		Healthy:         false,
		ErrorRate:       0.15,
	}

	shouldRollback, reason := ShouldRollback(state, 0.05, 500)
	if shouldRollback {
		rp := NewRollbackPlan(state.PreviousVersion, reason)
		fmt.Printf("Rollback triggered: %s\n", rp.Reason)
		fmt.Printf("Estimated time: %s\n", rp.EstimatedTime)
		fmt.Printf("Data loss risk: %v\n", rp.DataLossRisk)
		fmt.Println("Steps:")
		for _, s := range rp.Steps {
			fmt.Printf("  - %s\n", s)
		}
	} else {
		fmt.Println("No rollback needed")
	}

	fmt.Println()
	ra := &RollbackAutomation{
		State: DeploymentState{
			CurrentVersion:  "v2.0.0",
			PreviousVersion: "v1.0.0",
			Healthy:         true,
			FailedDeploy:    true,
		},
	}
	rap, ok := ra.Evaluate()
	if ok {
		fmt.Printf("Automation rolled back: %s (reason: %s)\n", rap.Version, rap.Reason)
	} else {
		fmt.Println("Automation: no rollback needed")
	}
}
