package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func setupTestLogger() (*bytes.Buffer, *RedactionStats, *slog.Logger) {
	buf := &bytes.Buffer{}
	stats := &RedactionStats{}
	handler := NewPIIRedactorHandler(slog.NewJSONHandler(buf, nil), stats)
	logger := slog.New(handler)
	return buf, stats, logger
}

func parseLogJSON(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("failed to parse JSON log: %v (raw: %s)", err, buf.String())
	}
	buf.Reset()
	return m
}

func TestPIIRedaction_SensitiveKeys(t *testing.T) {
	buf, stats, logger := setupTestLogger()

	logger.Info("Auth attempt",
		"username", "bob",
		"password", "hunter2",
		"token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
	)

	m := parseLogJSON(t, buf)

	if m["password"] != "[REDACTED_SECRET]" {
		t.Fatalf("expected password to be redacted, got %v", m["password"])
	}
	if m["token"] != "[REDACTED_SECRET]" {
		t.Fatalf("expected token to be redacted, got %v", m["token"])
	}
	if m["username"] != "bob" {
		t.Fatalf("expected username 'bob' preserved, got %v", m["username"])
	}
	if stats.PIIRedacted.Load() < 2 {
		t.Fatalf("expected at least 2 PII redacting events, got %d", stats.PIIRedacted.Load())
	}
}

func TestPIIRedaction_CreditCardAndSSN(t *testing.T) {
	buf, _, logger := setupTestLogger()

	logger.Info("Payment processed",
		"card_number", "4111-2222-3333-1234",
		"ssn", "123-45-6789",
	)

	m := parseLogJSON(t, buf)

	ccVal, ok := m["card_number"].(string)
	if !ok || !strings.Contains(ccVal, "****-****-****-1234") {
		t.Fatalf("expected masked credit card with suffix 1234, got %v", m["card_number"])
	}

	ssnVal, ok := m["ssn"].(string)
	if !ok || !strings.Contains(ssnVal, "***-**-6789") {
		t.Fatalf("expected masked SSN with suffix 6789, got %v", m["ssn"])
	}
}

func TestPIIRedaction_Email(t *testing.T) {
	buf, _, logger := setupTestLogger()

	logger.Info("User notification sent",
		"recipient", "john.doe@example.com",
	)

	m := parseLogJSON(t, buf)

	emailVal, ok := m["recipient"].(string)
	if !ok || strings.Contains(emailVal, "john.doe") {
		t.Fatalf("expected email address to be masked, got %v", m["recipient"])
	}
	if !strings.HasSuffix(emailVal, "@example.com") {
		t.Fatalf("expected domain preserved, got %v", emailVal)
	}
}

func TestPIIRedaction_NestedGroups(t *testing.T) {
	buf, _, logger := setupTestLogger()

	logger.Info("Order checkout",
		slog.Group("customer",
			"name", "Alice",
			"email", "alice@wonderland.org",
			"secret_pin", "[REDACTED_SECRET]",
		),
	)

	m := parseLogJSON(t, buf)
	cust, ok := m["customer"].(map[string]any)
	if !ok {
		t.Fatalf("expected customer group in JSON log: %v", m)
	}

	if cust["name"] != "Alice" {
		t.Fatalf("expected name Alice preserved, got %v", cust["name"])
	}
	if strings.Contains(cust["email"].(string), "alice@") {
		t.Fatalf("expected nested email masked, got %v", cust["email"])
	}
}

func TestPIIRedaction_BenignLogUnchanged(t *testing.T) {
	buf, stats, logger := setupTestLogger()

	logger.Info("System status operational",
		"uptime_seconds", 3600,
		"active_nodes", 5,
	)

	m := parseLogJSON(t, buf)
	if m["msg"] != "System status operational" {
		t.Fatalf("message modified: %v", m["msg"])
	}
	if m["uptime_seconds"].(float64) != 3600 {
		t.Fatalf("attribute modified: %v", m["uptime_seconds"])
	}
	if stats.PIIRedacted.Load() != 0 {
		t.Fatalf("expected 0 PII redacted on benign log, got %d", stats.PIIRedacted.Load())
	}
}
