package main

import (
	"testing"
)

func TestNewBlueGreenDeployment(t *testing.T) {
	d := NewBlueGreenDeployment("v1", "v2")
	if d.Strategy != StrategyBlueGreen {
		t.Errorf("expected blue-green strategy, got %s", d.Strategy)
	}
	if d.Current.Version != "v1" {
		t.Errorf("expected current v1, got %s", d.Current.Version)
	}
	if d.New.Version != "v2" {
		t.Errorf("expected new v2, got %s", d.New.Version)
	}
	if d.Done {
		t.Error("should not be done initially")
	}
}

func TestNewCanaryDeployment(t *testing.T) {
	d := NewCanaryDeployment("v1", "v2", 10)
	if d.Strategy != StrategyCanary {
		t.Errorf("expected canary strategy, got %s", d.Strategy)
	}
	if d.CanaryPercent != 10 {
		t.Errorf("expected 10%% canary, got %d", d.CanaryPercent)
	}
}

func TestBlueGreenDeploy(t *testing.T) {
	d := NewBlueGreenDeployment("v1", "v2")
	events, err := d.Deploy()
	if err != nil {
		t.Fatalf("deploy failed: %v", err)
	}
	if !d.Done {
		t.Error("expected done after deploy")
	}
	if !d.New.Healthy {
		t.Error("expected new target to be healthy after deploy")
	}
	if len(events) == 0 {
		t.Error("expected events from deploy")
	}
}

func TestCanaryDeploy(t *testing.T) {
	d := NewCanaryDeployment("v1", "v2", 10)
	events, err := d.Deploy()
	if err != nil {
		t.Fatalf("deploy failed: %v", err)
	}
	if !d.Done {
		t.Error("expected done after deploy")
	}
	if len(events) == 0 {
		t.Error("expected events from deploy")
	}
}

func TestRecreateDeploy(t *testing.T) {
	d := &Deployment{Strategy: StrategyRecreate, Current: &Target{Version: "v1", Healthy: true}, New: &Target{Version: "v2"}}
	events, err := d.Deploy()
	if err != nil {
		t.Fatalf("deploy failed: %v", err)
	}
	if !d.Done {
		t.Error("expected done")
	}
	if d.Current.Healthy {
		t.Error("expected current to be unhealthy after recreate")
	}
	if !d.New.Healthy {
		t.Error("expected new to be healthy after recreate")
	}
	if len(events) == 0 {
		t.Error("expected events")
	}
}

func TestUnknownStrategy(t *testing.T) {
	d := &Deployment{Strategy: "unknown"}
	_, err := d.Deploy()
	if err == nil {
		t.Error("expected error for unknown strategy")
	}
}

func TestBlueGreenRollback(t *testing.T) {
	d := NewBlueGreenDeployment("v1", "v2")
	d.Deploy()
	events := d.Rollback()
	if len(events) == 0 {
		t.Error("expected rollback events")
	}
}

func TestCanaryRollback(t *testing.T) {
	d := NewCanaryDeployment("v1", "v2", 10)
	d.Deploy()
	events := d.Rollback()
	if len(events) == 0 {
		t.Error("expected rollback events")
	}
}

func TestRollbackBeforeDeploy(t *testing.T) {
	d := NewBlueGreenDeployment("v1", "v2")
	events := d.Rollback()
	if len(events) != 1 || events[0] != "Deployment not complete, nothing to rollback" {
		t.Errorf("unexpected rollback before deploy: %v", events)
	}
}

func TestZeroDowntimeRequired(t *testing.T) {
	tests := []struct {
		strategy DeploymentStrategy
		want     bool
	}{
		{StrategyBlueGreen, true},
		{StrategyCanary, true},
		{StrategyRolling, true},
		{StrategyRecreate, false},
	}
	for _, tc := range tests {
		got := ZeroDowntimeRequired(tc.strategy)
		if got != tc.want {
			t.Errorf("ZeroDowntimeRequired(%q) = %v, want %v", tc.strategy, got, tc.want)
		}
	}
}

func TestTargetInitialState(t *testing.T) {
	target := &Target{Name: "test", Version: "v1", Healthy: true}
	if !target.Healthy {
		t.Error("expected healthy target")
	}
}

func TestDeploymentMultipleStrategies(t *testing.T) {
	strategies := []DeploymentStrategy{StrategyBlueGreen, StrategyCanary, StrategyRolling, StrategyRecreate}
	for _, s := range strategies {
		d := &Deployment{Strategy: s, Current: &Target{Version: "v1", Healthy: true}, New: &Target{Version: "v2"}}
		if s == StrategyCanary {
			d.CanaryPercent = 10
		}
		_, err := d.Deploy()
		if err != nil {
			t.Errorf("strategy %s failed: %v", s, err)
		}
	}
}
