package main

import (
	"testing"
)

func TestNewRollbackPlanManual(t *testing.T) {
	p := NewRollbackPlan("v1.0.0", ReasonManual)
	if p.Version != "v1.0.0" {
		t.Errorf("expected v1.0.0, got %s", p.Version)
	}
	if p.Reason != ReasonManual {
		t.Errorf("expected manual reason, got %s", p.Reason)
	}
	if len(p.Steps) == 0 {
		t.Error("expected rollback steps")
	}
}

func TestNewRollbackPlanHighErrorRate(t *testing.T) {
	p := NewRollbackPlan("v1.0.0", ReasonHighErrorRate)
	if len(p.Steps) != 4 {
		t.Errorf("expected 4 steps for high error rate, got %d", len(p.Steps))
	}
	if p.DataLossRisk {
		t.Error("high error rate should not have data loss risk")
	}
}

func TestNewRollbackPlanHealthCheckFail(t *testing.T) {
	p := NewRollbackPlan("v1.0.0", ReasonHealthCheckFail)
	if len(p.Steps) != 3 {
		t.Errorf("expected 3 steps for health check fail, got %d", len(p.Steps))
	}
}

func TestNewRollbackPlanDatabaseIssue(t *testing.T) {
	p := NewRollbackPlan("v1.0.0", ReasonDatabaseIssue)
	if !p.DataLossRisk {
		t.Error("database rollback should have data loss risk")
	}
	if p.EstimatedTime.Minutes() < 10 {
		t.Errorf("database rollback should take > 10 minutes, got %v", p.EstimatedTime)
	}
}

func TestNewRollbackPlanLatencySpike(t *testing.T) {
	p := NewRollbackPlan("v1.0.0", ReasonLatencySpike)
	if len(p.Steps) != 3 {
		t.Errorf("expected 3 steps for latency spike, got %d", len(p.Steps))
	}
}

func TestShouldRollbackUnhealthy(t *testing.T) {
	needed, _ := ShouldRollback(
		DeploymentState{Healthy: false}, 0.05, 500)
	if !needed {
		t.Error("should rollback when unhealthy")
	}
}

func TestShouldRollbackFailedDeploy(t *testing.T) {
	needed, _ := ShouldRollback(
		DeploymentState{Healthy: true, FailedDeploy: true}, 0.05, 500)
	if !needed {
		t.Error("should rollback when deploy failed")
	}
}

func TestShouldRollbackHealthy(t *testing.T) {
	needed, _ := ShouldRollback(
		DeploymentState{Healthy: true, FailedDeploy: false, ErrorRate: 0.01}, 0.05, 500)
	if needed {
		t.Error("should not rollback when healthy")
	}
}

func TestRollbackAutomationEvaluateTriggered(t *testing.T) {
	ra := &RollbackAutomation{
		State: DeploymentState{
			Healthy:         false,
			FailedDeploy:    true,
			PreviousVersion: "v1.0.0",
		},
	}
	plan, ok := ra.Evaluate()
	if !ok {
		t.Error("expected rollback plan")
	}
	if plan == nil {
		t.Fatal("plan should not be nil")
	}
	if plan.Version != "v1.0.0" {
		t.Errorf("expected rollback to v1.0.0, got %s", plan.Version)
	}
}

func TestRollbackAutomationEvaluateNotNeeded(t *testing.T) {
	ra := &RollbackAutomation{
		State: DeploymentState{
			Healthy:      true,
			FailedDeploy: false,
			ErrorRate:    0.01,
		},
	}
	_, ok := ra.Evaluate()
	if ok {
		t.Error("should not need rollback")
	}
}

func TestNewRollbackPlanDefault(t *testing.T) {
	p := NewRollbackPlan("v1", "")
	if len(p.Steps) != 3 {
		t.Errorf("expected 3 default steps, got %d", len(p.Steps))
	}
}

func TestRollbackPlanTimeEstimate(t *testing.T) {
	p := NewRollbackPlan("v1", ReasonHealthCheckFail)
	if p.EstimatedTime <= 0 {
		t.Error("estimated time should be positive")
	}
}
