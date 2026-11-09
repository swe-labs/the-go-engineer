package main

import (
	"fmt"
)

type DeploymentStrategy string

const (
	StrategyBlueGreen DeploymentStrategy = "blue-green"
	StrategyCanary    DeploymentStrategy = "canary"
	StrategyRolling   DeploymentStrategy = "rolling"
	StrategyRecreate  DeploymentStrategy = "recreate"
)

type Target struct {
	Name    string
	Version string
	Healthy bool
}

type Deployment struct {
	Strategy      DeploymentStrategy
	Current       *Target
	New           *Target
	CanaryPercent int
	Done          bool
}

func NewBlueGreenDeployment(currentVersion, newVersion string) *Deployment {
	return &Deployment{
		Strategy: StrategyBlueGreen,
		Current:  &Target{Name: "blue", Version: currentVersion, Healthy: true},
		New:      &Target{Name: "green", Version: newVersion, Healthy: false},
	}
}

func NewCanaryDeployment(currentVersion, newVersion string, canaryPercent int) *Deployment {
	return &Deployment{
		Strategy:      StrategyCanary,
		Current:       &Target{Name: "production", Version: currentVersion, Healthy: true},
		New:           &Target{Name: "canary", Version: newVersion, Healthy: false},
		CanaryPercent: canaryPercent,
	}
}

func (d *Deployment) Deploy() ([]string, error) {
	var events []string
	switch d.Strategy {
	case StrategyBlueGreen:
		events = d.blueGreenDeploy()
	case StrategyCanary:
		events = d.canaryDeploy()
	case StrategyRolling:
		events = d.rollingDeploy()
	case StrategyRecreate:
		events = d.recreateDeploy()
	default:
		return nil, fmt.Errorf("unknown strategy: %s", d.Strategy)
	}
	d.Done = true
	return events, nil
}

func (d *Deployment) blueGreenDeploy() []string {
	var events []string
	events = append(events, fmt.Sprintf("Spinning up %s environment (%s)", d.New.Name, d.New.Version))
	events = append(events, fmt.Sprintf("Running health checks on %s", d.New.Name))
	d.New.Healthy = true
	events = append(events, fmt.Sprintf("%s environment healthy, switching traffic", d.New.Name))
	events = append(events, fmt.Sprintf("Traffic switched to %s", d.New.Name))
	events = append(events, fmt.Sprintf("Keeping %s as rollback target", d.Current.Name))
	return events
}

func (d *Deployment) canaryDeploy() []string {
	var events []string
	events = append(events, fmt.Sprintf("Deploying %s as canary (%d%% traffic)", d.New.Version, d.CanaryPercent))
	events = append(events, "Monitoring error rates and latency for 5 minutes")
	d.New.Healthy = true
	events = append(events, "Canary healthy, gradually increasing traffic")
	events = append(events, fmt.Sprintf("Full rollout of %s complete", d.New.Version))
	return events
}

func (d *Deployment) rollingDeploy() []string {
	var events []string
	events = append(events, "Replacing instances one at a time")
	events = append(events, fmt.Sprintf("Instance 1: drained, updated to %s", d.New.Version))
	events = append(events, "Instance 2: drained, updated")
	events = append(events, "All instances updated to new version")
	d.Current.Healthy = false
	d.New.Healthy = true
	return events
}

func (d *Deployment) recreateDeploy() []string {
	var events []string
	events = append(events, "Stopping all current instances")
	d.Current.Healthy = false
	events = append(events, "Starting new instances")
	d.New.Healthy = true
	events = append(events, fmt.Sprintf("All traffic on %s", d.New.Version))
	return events
}

func (d *Deployment) Rollback() []string {
	var events []string
	if !d.Done {
		events = append(events, "Deployment not complete, nothing to rollback")
		return events
	}
	switch d.Strategy {
	case StrategyBlueGreen:
		events = append(events, fmt.Sprintf("Switching traffic back to %s (%s)", d.Current.Name, d.Current.Version))
	case StrategyCanary:
		events = append(events, "Reverting all traffic to production version")
		events = append(events, "Removing canary instances")
	default:
		events = append(events, fmt.Sprintf("Redeploying version %s", d.Current.Version))
	}
	d.New, d.Current = d.Current, d.New
	d.Done = false
	return events
}

func ZeroDowntimeRequired(strategy DeploymentStrategy) bool {
	return strategy == StrategyBlueGreen || strategy == StrategyCanary || strategy == StrategyRolling
}

func main() {
	fmt.Println("=== One Deployment Target ===")
	fmt.Println()

	fmt.Println("1. Blue-Green Deployment:")
	blueGreen := NewBlueGreenDeployment("v1.0.0", "v2.0.0")
	events, _ := blueGreen.Deploy()
	for _, e := range events {
		fmt.Printf("  - %s\n", e)
	}
	fmt.Println("  Zero downtime:", ZeroDowntimeRequired(blueGreen.Strategy))
	fmt.Println("  Rollback:")
	for _, e := range blueGreen.Rollback() {
		fmt.Printf("    - %s\n", e)
	}

	fmt.Println()
	fmt.Println("2. Canary Deployment (10%):")
	canary := NewCanaryDeployment("v1.0.0", "v2.0.0", 10)
	events, _ = canary.Deploy()
	for _, e := range events {
		fmt.Printf("  - %s\n", e)
	}
}
