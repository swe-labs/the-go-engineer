package main

import (
	"fmt"
)

type BoundedContext struct {
	Name     string
	OwnsData []string
	TeamSize int
	APICalls []string
}

type SplittingAnalysis struct {
	Context           string
	Reasons           []string
	RecommendedAction string
	RiskLevel         string
}

func AnalyzeSplitContext(ctx BoundedContext) SplittingAnalysis {
	analysis := SplittingAnalysis{Context: ctx.Name}
	reasons := []string{}

	if ctx.TeamSize > 8 {
		reasons = append(reasons, "team size exceeds Dunbar limit")
	}
	if len(ctx.OwnsData) > 5 {
		reasons = append(reasons, "context owns too many data entities")
	}
	analysis.Reasons = reasons

	if len(reasons) >= 2 {
		analysis.RecommendedAction = fmt.Sprintf("split %s into separate service", ctx.Name)
		analysis.RiskLevel = "medium"
	} else if len(reasons) >= 1 {
		analysis.RecommendedAction = fmt.Sprintf("monitor %s for further growth", ctx.Name)
		analysis.RiskLevel = "low"
	} else {
		analysis.RecommendedAction = "keep as part of monolith"
		analysis.RiskLevel = "none"
	}
	analysis.Reasons = reasons
	return analysis
}

type ServiceCoupling struct {
	ServiceA   string
	ServiceB   string
	SharedData []string
	CallPerSec int
}

func EvaluateCoupling(c ServiceCoupling) string {
	if len(c.SharedData) > 0 && c.CallPerSec > 100 {
		return fmt.Sprintf("high coupling between %s and %s: consider merging", c.ServiceA, c.ServiceB)
	}
	if len(c.SharedData) > 0 {
		return fmt.Sprintf("medium coupling: shared data %v needs ownership resolution", c.SharedData)
	}
	return fmt.Sprintf("low coupling: %s and %s are independent", c.ServiceA, c.ServiceB)
}

type SplittingStrategy struct {
	Name       string
	Steps      []string
	Downtime   bool
	DurationMs int
}

var Strategies = map[string]SplittingStrategy{
	"strangler-fig": {
		Name: "Strangler Fig",
		Steps: []string{
			"identify bounded context boundary",
			"add proxy layer routing traffic to new service",
			"migrate one endpoint at a time",
			"remove old code when all traffic routed",
		},
		Downtime:   false,
		DurationMs: 5000,
	},
	"database-separation": {
		Name: "Database Separation",
		Steps: []string{
			"extract schema for bounded context",
			"set up new database instance",
			"run dual-writes for consistency check",
			"switch reads to new database",
			"remove old schema",
		},
		Downtime:   true,
		DurationMs: 10000,
	},
	"event-driven-split": {
		Name: "Event-driven Split",
		Steps: []string{
			"publish events for the bounded context's data",
			"build new service consuming events",
			"route new requests to new service",
			"backfill old data",
			"deprecate old service",
		},
		Downtime:   false,
		DurationMs: 15000,
	},
}

func main() {
	orderCtx := BoundedContext{
		Name:     "Order Management",
		OwnsData: []string{"orders", "order_items", "payments", "refunds", "shipments", "inventory_reservations"},
		TeamSize: 10,
	}
	analysis := AnalyzeSplitContext(orderCtx)
	fmt.Printf("Analysis for %s:\n", analysis.Context)
	for _, r := range analysis.Reasons {
		fmt.Printf("  - %s\n", r)
	}
	fmt.Printf("  Action: %s (risk: %s)\n\n", analysis.RecommendedAction, analysis.RiskLevel)

	coupling := ServiceCoupling{
		ServiceA:   "Orders",
		ServiceB:   "Inventory",
		SharedData: []string{"product_sku"},
		CallPerSec: 250,
	}
	fmt.Println(EvaluateCoupling(coupling))

	strat := Strategies["strangler-fig"]
	fmt.Printf("\n%s strategy:\n", strat.Name)
	for _, s := range strat.Steps {
		fmt.Printf("  - %s\n", s)
	}
	fmt.Printf("Downtime required: %v\n", strat.Downtime)
}
