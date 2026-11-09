package main

import (
	"fmt"
	"math"
)

type ComplexityLevel int

const (
	Trivial ComplexityLevel = iota + 1
	Easy
	Moderate
	Complex
	VeryComplex
)

func (c ComplexityLevel) String() string {
	switch c {
	case Trivial:
		return "Trivial"
	case Easy:
		return "Easy"
	case Moderate:
		return "Moderate"
	case Complex:
		return "Complex"
	case VeryComplex:
		return "Very Complex"
	default:
		return "Unknown"
	}
}

type ProjectScope struct {
	Name            string
	Description     string
	NumFeatures     int
	NumIntegrations int
	HasAuth         bool
	HasDatabase     bool
	HasExternalAPI  bool
	NeedsUI         bool
	TeamSize        int
	EstimatedWeeks  float64
	Complexity      ComplexityLevel
	RiskFactors     []string
}

type ComplexityEstimator struct {
	BaseWeeks     float64
	FeatureWeight float64
	IntegrationW  float64
	AuthWeight    float64
	DatabaseW     float64
	APIWeight     float64
	UIWeight      float64
	TeamOverhead  float64
}

func DefaultEstimator() ComplexityEstimator {
	return ComplexityEstimator{
		BaseWeeks:     2.0,
		FeatureWeight: 1.5,
		IntegrationW:  2.0,
		AuthWeight:    2.0,
		DatabaseW:     1.5,
		APIWeight:     1.5,
		UIWeight:      3.0,
		TeamOverhead:  0.15,
	}
}

func EstimateProject(scope ProjectScope, est ComplexityEstimator) ProjectScope {
	weeks := est.BaseWeeks
	weeks += float64(scope.NumFeatures) * est.FeatureWeight
	weeks += float64(scope.NumIntegrations) * est.IntegrationW

	if scope.HasAuth {
		weeks += est.AuthWeight
		scope.RiskFactors = append(scope.RiskFactors, "authentication adds security complexity")
	}
	if scope.HasDatabase {
		weeks += est.DatabaseW
		scope.RiskFactors = append(scope.RiskFactors, "database schema design and migration overhead")
	}
	if scope.HasExternalAPI {
		weeks += est.APIWeight
		scope.RiskFactors = append(scope.RiskFactors, "external API integration has reliability dependencies")
	}
	if scope.NeedsUI {
		weeks += est.UIWeight
		scope.RiskFactors = append(scope.RiskFactors, "UI development adds frontend complexity")
	}

	if scope.TeamSize > 0 {
		weeks *= 1.0 + est.TeamOverhead*float64(scope.TeamSize-1)
	}

	estimatedWeeks := math.Ceil(weeks)
	scope.EstimatedWeeks = estimatedWeeks

	switch {
	case estimatedWeeks <= 2:
		scope.Complexity = Trivial
	case estimatedWeeks <= 4:
		scope.Complexity = Easy
	case estimatedWeeks <= 8:
		scope.Complexity = Moderate
	case estimatedWeeks <= 16:
		scope.Complexity = Complex
	default:
		scope.Complexity = VeryComplex
	}

	return scope
}

func (s ProjectScope) Summary() string {
	return fmt.Sprintf(
		"Project: %s\n  Description: %s\n  Features: %d, Integrations: %d\n  Auth: %v, DB: %v, API: %v, UI: %v\n  Team Size: %d\n  Estimated Complexity: %s\n  Estimated Timeline: %.0f weeks\n  Risk Factors: %v\n",
		s.Name, s.Description, s.NumFeatures, s.NumIntegrations,
		s.HasAuth, s.HasDatabase, s.HasExternalAPI, s.NeedsUI,
		s.TeamSize, s.Complexity, s.EstimatedWeeks, s.RiskFactors,
	)
}

func main() {
	est := DefaultEstimator()

	projects := []ProjectScope{
		{
			Name:        "URL Shortener API",
			Description: "A REST API that shortens URLs with redirect tracking and analytics",
			NumFeatures: 4, NumIntegrations: 1,
			HasAuth: true, HasDatabase: true, HasExternalAPI: false,
			NeedsUI: false, TeamSize: 1,
		},
		{
			Name:        "CLI Task Manager",
			Description: "A command-line task management tool with local storage",
			NumFeatures: 3, NumIntegrations: 0,
			HasAuth: false, HasDatabase: true, HasExternalAPI: false,
			NeedsUI: false, TeamSize: 1,
		},
		{
			Name:        "Full SaaS Dashboard",
			Description: "A full-featured SaaS dashboard with multi-tenant auth and analytics",
			NumFeatures: 8, NumIntegrations: 3,
			HasAuth: true, HasDatabase: true, HasExternalAPI: true,
			NeedsUI: true, TeamSize: 3,
		},
	}

	for _, p := range projects {
		result := EstimateProject(p, est)
		fmt.Println(result.Summary())
	}
}
