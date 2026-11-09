package main

import (
	"testing"
)

func setupTestEngine() *PolicyEngine {
	pe := NewPolicyEngine()
	pe.AddRule(PolicyRule{
		Name: "admin-override", Effect: "allow",
		Conditions: []Condition{
			{Attribute: "user.role", Operator: "eq", Value: "admin"},
		},
	})
	pe.AddRule(PolicyRule{
		Name: "same-dept", Effect: "allow",
		Conditions: []Condition{
			{Attribute: "user.department", Operator: "eq", Value: "${resource.department}"},
			{Attribute: "env.method", Operator: "eq", Value: "GET"},
		},
	})
	pe.AddRule(PolicyRule{
		Name: "clearance", Effect: "allow",
		Conditions: []Condition{
			{Attribute: "user.clearance", Operator: "gte", Value: 3},
			{Attribute: "env.method", Operator: "eq", Value: "GET"},
		},
	})
	pe.AddRule(PolicyRule{
		Name: "deny-delete", Effect: "deny",
		Conditions: []Condition{
			{Attribute: "env.method", Operator: "eq", Value: "DELETE"},
		},
	})
	return pe
}

func TestEvaluatePolicy_SameDept(t *testing.T) {
	pe := setupTestEngine()
	allowed, rule := EvaluatePolicy(pe, "engineering", "2", "engineering", 1, "GET")
	if !allowed {
		t.Errorf("expected allowed, rule=%s", rule)
	}
	if rule != "same-dept" {
		t.Errorf("expected same-dept rule, got %s", rule)
	}
}

func TestEvaluatePolicy_HighClearance(t *testing.T) {
	pe := setupTestEngine()
	allowed, rule := EvaluatePolicy(pe, "engineering", "5", "finance", 3, "GET")
	if !allowed {
		t.Errorf("expected allowed, rule=%s", rule)
	}
	if rule != "clearance" {
		t.Errorf("expected clearance rule, got %s", rule)
	}
}

func TestEvaluatePolicy_DenyDelete(t *testing.T) {
	pe := setupTestEngine()
	allowed, rule := EvaluatePolicy(pe, "engineering", "3", "engineering", 2, "DELETE")
	if allowed {
		t.Errorf("expected denied, rule=%s", rule)
	}
	if rule != "deny-delete" {
		t.Errorf("expected deny-delete rule, got %s", rule)
	}
}

func TestEvaluatePolicy_AdminOverride(t *testing.T) {
	pe := setupTestEngine()
	allowed, rule := EvaluatePolicy(pe, "engineering", "admin", "finance", 5, "DELETE")
	if !allowed {
		t.Errorf("expected allowed for admin, rule=%s", rule)
	}
	if rule != "admin-override" {
		t.Errorf("expected admin-override rule, got %s", rule)
	}
}

func TestEvaluatePolicy_LowClearanceDifferentDept(t *testing.T) {
	pe := setupTestEngine()
	allowed, rule := EvaluatePolicy(pe, "engineering", "1", "finance", 3, "GET")
	if allowed {
		t.Errorf("expected denied, rule=%s", rule)
	}
	_ = rule
}

func TestEvaluateCondition(t *testing.T) {
	tests := []struct {
		op       string
		actual   interface{}
		expected interface{}
		want     bool
	}{
		{"eq", "hello", "hello", true},
		{"eq", "hello", "world", false},
		{"neq", "hello", "world", true},
		{"gt", 5, 3, true},
		{"gte", 5, 5, true},
		{"lt", 2, 5, true},
		{"contains", "hello world", "world", true},
	}
	for _, tt := range tests {
		got, err := evaluateCondition(tt.op, tt.actual, tt.expected)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}
		if got != tt.want {
			t.Errorf("evaluateCondition(%q, %v, %v) = %v, want %v", tt.op, tt.actual, tt.expected, got, tt.want)
		}
	}
}
