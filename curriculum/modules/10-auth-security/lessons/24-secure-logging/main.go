package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

func (l LogLevel) String() string {
	return [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[l]
}

type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields"`
}

type SecureLogger struct {
	MinLevel      LogLevel
	RedactedKeys  map[string]bool
	AuditEnabled  bool
	auditCallback func(LogEntry)
}

func NewSecureLogger(minLevel LogLevel) *SecureLogger {
	return &SecureLogger{
		MinLevel: minLevel,
		RedactedKeys: map[string]bool{
			"password": true, "secret": true, "token": true,
			"api_key": true, "apiKey": true, "access_token": true,
			"refresh_token": true, "ssn": true, "credit_card": true,
		},
	}
}

func (sl *SecureLogger) redact(fields map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(fields))
	for k, v := range fields {
		if sl.RedactedKeys[strings.ToLower(k)] {
			result[k] = "[REDACTED]"
		} else {
			result[k] = v
		}
	}
	return result
}

func (sl *SecureLogger) log(level LogLevel, msg string, fields map[string]interface{}) {
	if level < sl.MinLevel {
		return
	}
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level.String(),
		Message:   msg,
		Fields:    sl.redact(fields),
	}
	data, _ := json.Marshal(entry)
	fmt.Println(string(data))

	if sl.AuditEnabled && level >= WARN {
		if sl.auditCallback != nil {
			sl.auditCallback(entry)
		}
	}
}

func (sl *SecureLogger) Debug(msg string, fields map[string]interface{}) {
	sl.log(DEBUG, msg, fields)
}

func (sl *SecureLogger) Info(msg string, fields map[string]interface{}) {
	sl.log(INFO, msg, fields)
}

func (sl *SecureLogger) Warn(msg string, fields map[string]interface{}) {
	sl.log(WARN, msg, fields)
}

func (sl *SecureLogger) Error(msg string, fields map[string]interface{}) {
	sl.log(ERROR, msg, fields)
}

type AuditTrail struct {
	Entries []LogEntry
}

func (at *AuditTrail) Append(entry LogEntry) {
	at.Entries = append(at.Entries, entry)
}

func (at *AuditTrail) Summary() string {
	return fmt.Sprintf("Audit trail: %d entries", len(at.Entries))
}

func main() {
	fmt.Println("=== Secure Logging Demo ===")

	trail := &AuditTrail{}
	logger := NewSecureLogger(INFO)
	logger.AuditEnabled = true
	logger.auditCallback = trail.Append

	logger.Info("User logged in", map[string]interface{}{
		"user_id": "u123",
		"email":   "alice@example.com",
	})

	logger.Info("Payment processed", map[string]interface{}{
		"order_id":    "ord-456",
		"amount":      49.99,
		"credit_card": "4111-1111-1111-1111",
		"cvv":         "123",
	})

	logger.Warn("Failed login attempt", map[string]interface{}{
		"user_id":  "u999",
		"password": "hunter2",
		"ip":       "192.168.1.100",
	})

	logger.Error("Database connection failed", map[string]interface{}{
		"database": "users_db",
		"error":    "connection refused",
	})

	fmt.Println("\n=== Audit Trail ===")
	fmt.Println(trail.Summary())
	for _, e := range trail.Entries {
		fmt.Printf("  [%s] %s: %s\n", e.Level, e.Timestamp, e.Message)
	}

	fmt.Println("\nNote: Passwords, credit cards, and tokens are redacted.")
	fmt.Println("Sensitive fields should never appear in logs.")
}
