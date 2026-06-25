package main

import (
	"math"
	"testing"
)

func TestEstimateStorage_Basic(t *testing.T) {
	design := SystemDesign{Name: "test"}
	design.AddRequirement("data", 1, Gigabytes, "test data")
	design.EstimateStorage()
	if len(design.Estimates) < 1 {
		t.Fatal("expected at least one estimate")
	}
	if design.Estimates[0].EstimatedValue < 0.9 || design.Estimates[0].EstimatedValue > 1.1 {
		t.Errorf("expected ~1 GB storage, got %.2f", design.Estimates[0].EstimatedValue)
	}
}

func TestEstimateStorage_MultipleUnits(t *testing.T) {
	design := SystemDesign{Name: "test"}
	design.AddRequirement("small", 500, Megabytes, "test")
	design.AddRequirement("medium", 2, Gigabytes, "test")
	design.EstimateStorage()
	for _, e := range design.Estimates {
		if e.EstimatedUnit != Gigabytes && e.EstimatedUnit != Terabytes {
			continue
		}
		if e.EstimatedValue < 2 || e.EstimatedValue > 3 {
			t.Errorf("expected ~2.5 GB, got %.2f %s", e.EstimatedValue, e.EstimatedUnit)
		}
	}
}

func TestEstimateThroughput_HasAssumptions(t *testing.T) {
	design := SystemDesign{Name: "test"}
	design.AddRequirement("reads", 100, RequestsPerSecond, "test reads")
	design.EstimateThroughput()
	for _, e := range design.Estimates {
		if len(e.Assumptions) == 0 {
			t.Errorf("expected throughput estimate to have assumptions")
		}
		if len(e.Recommendations) == 0 {
			t.Errorf("expected throughput estimate to have recommendations")
		}
	}
}

func TestBytesToUnit(t *testing.T) {
	tests := []struct {
		input    float64
		wantUnit DataUnit
		tol      float64
	}{
		{500, Bytes, 1},
		{2048, Kilobytes, 1},
		{1048576, Megabytes, 1},
		{1073741824, Gigabytes, 1},
		{1099511627776, Terabytes, 1},
	}
	for _, tc := range tests {
		val, unit := bytesToUnit(tc.input)
		if unit != tc.wantUnit {
			t.Errorf("bytesToUnit(%f) unit = %v, want %v", tc.input, unit, tc.wantUnit)
		}
		if math.Abs(val-math.Round(val)) > tc.tol && val > tc.tol {
			t.Errorf("bytesToUnit(%f) = %f, want ~%f", tc.input, val, math.Round(val))
		}
	}
}

func TestSystemDesign_RequirementsStored(t *testing.T) {
	design := SystemDesign{Name: "test"}
	design.AddRequirement("users", 1000, RequestsPerSecond, "test users")
	if len(design.Requirements) != 1 {
		t.Fatalf("expected 1 requirement, got %d", len(design.Requirements))
	}
	if design.Requirements[0].Name != "users" {
		t.Errorf("expected name 'users', got '%s'", design.Requirements[0].Name)
	}
}
