package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Attributes struct {
	UserDepartment     string `json:"user_department"`
	UserClearance      int    `json:"user_clearance"`
	UserRole           string `json:"user_role"`
	ResourceDepartment string `json:"resource_department"`
	ResourceOwner      string `json:"resource_owner"`
	ResourceClass      int    `json:"resource_class"`
	EnvTime            string `json:"env_time"`
	EnvMethod          string `json:"env_method"`
	EnvIP              string `json:"env_ip"`
	EnvTenantID        string `json:"env_tenant_id"`
}

type Condition struct {
	Attribute string      `json:"attribute"`
	Operator  string      `json:"operator"`
	Value     interface{} `json:"value"`
}

type PolicyRule struct {
	Name       string      `json:"name"`
	Effect     string      `json:"effect"`
	Conditions []Condition `json:"conditions"`
}

type PolicyEngine struct {
	mu    sync.RWMutex
	rules []PolicyRule
}

func NewPolicyEngine() *PolicyEngine {
	return &PolicyEngine{}
}

func (pe *PolicyEngine) AddRule(rule PolicyRule) {
	pe.mu.Lock()
	defer pe.mu.Unlock()
	pe.rules = append(pe.rules, rule)
}

func (pe *PolicyEngine) Evaluate(attrs *Attributes) (bool, string, error) {
	pe.mu.RLock()
	defer pe.mu.RUnlock()
	for _, rule := range pe.rules {
		matches, err := evaluateConditions(rule.Conditions, attrs)
		if err != nil {
			return false, rule.Name, err
		}
		if matches {
			switch rule.Effect {
			case "allow":
				return true, rule.Name, nil
			case "deny":
				return false, rule.Name, nil
			}
		}
	}
	return false, "default-deny", nil
}

func evaluateConditions(conditions []Condition, attrs *Attributes) (bool, error) {
	for _, c := range conditions {
		attrValue, err := getAttributeValue(c.Attribute, attrs)
		if err != nil {
			return false, err
		}
		resolved := resolveValue(c.Value, attrs)
		ok, err := evaluateCondition(c.Operator, attrValue, resolved)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

func getAttributeValue(attr string, attrs *Attributes) (interface{}, error) {
	switch attr {
	case "user.department":
		return attrs.UserDepartment, nil
	case "user.clearance":
		return attrs.UserClearance, nil
	case "user.role":
		return attrs.UserRole, nil
	case "resource.department":
		return attrs.ResourceDepartment, nil
	case "resource.owner":
		return attrs.ResourceOwner, nil
	case "resource.class":
		return attrs.ResourceClass, nil
	case "env.time":
		return attrs.EnvTime, nil
	case "env.method":
		return attrs.EnvMethod, nil
	case "env.ip":
		return attrs.EnvIP, nil
	case "env.tenant_id":
		return attrs.EnvTenantID, nil
	default:
		return nil, fmt.Errorf("unknown attribute: %s", attr)
	}
}

// resolveValue resolves attribute references like "${user.department}".
func resolveValue(v interface{}, attrs *Attributes) interface{} {
	s, ok := v.(string)
	if !ok || !strings.HasPrefix(s, "${") || !strings.HasSuffix(s, "}") {
		return v
	}
	attr := strings.TrimSuffix(strings.TrimPrefix(s, "${"), "}")
	val, err := getAttributeValue(attr, attrs)
	if err != nil {
		return v
	}
	return val
}

func evaluateCondition(op string, actual, expected interface{}) (bool, error) {
	switch op {
	case "eq":
		return fmt.Sprintf("%v", actual) == fmt.Sprintf("%v", expected), nil
	case "neq":
		return fmt.Sprintf("%v", actual) != fmt.Sprintf("%v", expected), nil
	case "gt":
		return toFloat(actual) > toFloat(expected), nil
	case "gte":
		return toFloat(actual) >= toFloat(expected), nil
	case "lt":
		return toFloat(actual) < toFloat(expected), nil
	case "lte":
		return toFloat(actual) <= toFloat(expected), nil
	case "in":
		return inList(actual, expected)
	case "contains":
		return strings.Contains(fmt.Sprintf("%v", actual), fmt.Sprintf("%v", expected)), nil
	default:
		return false, fmt.Errorf("unknown operator: %s", op)
	}
}

func toFloat(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case int:
		return float64(val)
	}
	return 0
}

func inList(actual interface{}, expected interface{}) (bool, error) {
	list, ok := expected.([]interface{})
	if !ok {
		return false, errors.New("in operator requires a list value")
	}
	for _, item := range list {
		if fmt.Sprintf("%v", actual) == fmt.Sprintf("%v", item) {
			return true, nil
		}
	}
	return false, nil
}

func EvaluatePolicy(pe *PolicyEngine, userDept, userClearance string, resourceDept string, resourceClass int, method string) (bool, string) {
	clearance := 0
	if userClearance == "admin" {
		clearance = 99
	} else {
		fmt.Sscanf(userClearance, "%d", &clearance)
	}
	attrs := &Attributes{
		UserDepartment:     userDept,
		UserClearance:      clearance,
		UserRole:           userClearance,
		ResourceDepartment: resourceDept,
		ResourceClass:      resourceClass,
		EnvMethod:          method,
		EnvTime:            time.Now().Format("15:04"),
		EnvIP:              "10.0.0.1",
	}
	allowed, rule, _ := pe.Evaluate(attrs)
	return allowed, rule
}

func main() {
	pe := NewPolicyEngine()

	// Admin override must come first (break-glass policy).
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

	scenarios := []struct {
		name            string
		userDept        string
		userClearance   string
		resourceDept    string
		resourceClass   int
		method          string
		expectedAllowed bool
	}{
		{"Same dept GET", "engineering", "3", "engineering", 2, "GET", true},
		{"Different dept low clearance", "engineering", "1", "finance", 3, "GET", false},
		{"Admin override", "engineering", "admin", "finance", 5, "DELETE", true},
		{"High clearance different dept", "engineering", "5", "finance", 3, "GET", true},
		{"DELETE denied", "engineering", "3", "engineering", 2, "DELETE", false},
	}

	for _, s := range scenarios {
		allowed, rule := EvaluatePolicy(pe, s.userDept, s.userClearance, s.resourceDept, s.resourceClass, s.method)
		status := "PASS"
		if allowed != s.expectedAllowed {
			status = "FAIL"
		}
		fmt.Printf("[%s] %s: allowed=%v (rule: %s)\n", status, s.name, allowed, rule)
	}

	// HTTP
	http.HandleFunc("/api/documents/", func(w http.ResponseWriter, r *http.Request) {
		attrs := &Attributes{
			UserDepartment:     r.Header.Get("X-Department"),
			UserClearance:      3,
			UserRole:           r.Header.Get("X-Role"),
			ResourceDepartment: r.URL.Query().Get("department"),
			ResourceClass:      2,
			EnvMethod:          r.Method,
			EnvTime:            time.Now().Format("15:04"),
			EnvIP:              r.RemoteAddr,
		}
		allowed, rule, _ := pe.Evaluate(attrs)
		if !allowed {
			w.Header().Set("X-Deny-Reason", rule)
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "access granted"})
	})

	fmt.Println("\nABAC server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
