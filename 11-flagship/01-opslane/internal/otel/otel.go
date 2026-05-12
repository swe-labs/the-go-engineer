// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Package otel provides OpenTelemetry tracing for Opslane.
// Role: Observability boundary - exports traces to OTLP-compatible backends.
// Boundary: Traces are exported asynchronously; failures are logged but don't block.
// Failure mode: OTLP export failures are logged and dropped to prevent memory leaks.

package otel

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Config (Struct): configures the OpenTelemetry tracer.
type Config struct {
	Endpoint    string
	Insecure    bool
	Timeout     time.Duration
	Enabled     bool
	SampleRate  float64
	ServiceName string
	Environment string
}

// Tracer (Struct): provides OpenTelemetry tracing for the Opslane backend.
// It implements span creation, sampling, and async OTLP export.
type Tracer struct {
	config  Config
	client  *otlpClient
	spans   chan Span
	stopped chan struct{}
	wg      sync.WaitGroup
	logger  *slog.Logger
}

// Span (Struct): represents a single trace span with timing and attribute data.
type Span struct {
	TraceID       string
	SpanID        string
	ParentID      string
	Name          string
	StartTime     time.Time
	EndTime       time.Time
	Attributes    map[string]string
	Status        string
	StatusMessage string
}

// New (Constructor): creates a Tracer and starts the background export loop if
// tracing is enabled. Returns a no-op tracer when cfg.Enabled is false.
func New(cfg Config, logger *slog.Logger) *Tracer {
	if !cfg.Enabled {
		return &Tracer{config: cfg}
	}

	t := &Tracer{
		config:  cfg,
		spans:   make(chan Span, 1000),
		stopped: make(chan struct{}),
		logger:  logger,
	}

	if cfg.Endpoint != "" {
		t.client = newOTLPClient(cfg.Endpoint, cfg.Insecure, cfg.Timeout, cfg.ServiceName, cfg.Environment)
		t.wg.Add(1)
		go t.exportLoop()
	}

	return t
}

// exportLoop (Goroutine): drains the span channel on a 5-second tick or 100-span
// batch boundary, then sends the batch to the OTLP endpoint. Runs in a background
// goroutine started by New().
func (t *Tracer) exportLoop() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var batch []Span
	flush := func() {
		if len(batch) > 0 && t.client != nil {
			ctx, cancel := context.WithTimeout(context.Background(), t.config.Timeout)
			if err := t.client.Export(ctx, batch); err != nil {
				t.logger.Error("OTLP export failed", slog.Any("error", err))
			}
			cancel()
			batch = batch[:0]
		}
	}

	for {
		select {
		case <-t.stopped:
			flush()
			return
		case span := <-t.spans:
			batch = append(batch, span)
			if len(batch) >= 100 {
				flush()
			}
		case <-ticker.C:
			flush()
		}
	}
}

// StartSpan (Method): creates a new span with the given name and attributes,
// sampling according to the configured rate. Returns a context containing the
// span and a finish function that records the end time and enqueues the span.
func (t *Tracer) StartSpan(ctx context.Context, name string, attrs ...string) (context.Context, func()) {
	if t.config.Enabled && t.config.SampleRate < 1.0 {
		if shouldDropSpan(t.config.SampleRate) {
			return ctx, func() {}
		}
	}

	spanID := generateSpanID()
	traceID := getOrCreateTraceID(ctx)
	parentID := GetSpanID(ctx)

	span := Span{
		TraceID:    traceID,
		SpanID:     spanID,
		ParentID:   parentID,
		Name:       name,
		StartTime:  time.Now(),
		Attributes: attrsToMap(attrs),
		Status:     "ok",
	}

	ctx = WithTraceID(ctx, traceID)
	ctx = WithSpanID(ctx, spanID)
	return context.WithValue(ctx, spanKey, span), func() {
		span.EndTime = time.Now()
		if t.config.Enabled {
			select {
			case t.spans <- span:
			default:
				t.logger.Warn("OTLP span channel full, dropping span")
			}
		}
	}
}

// Stop (Method): signals the export loop to flush remaining spans and waits for
// the background goroutine to finish, then closes the OTLP HTTP client.
func (t *Tracer) Stop() {
	if t.client != nil {
		close(t.stopped)
		t.wg.Wait()
		t.client.Close()
	}
}

// Enabled (Method): reports whether the tracer was configured with a non-empty
// endpoint, allowing callers to skip tracing work when disabled.
func (t *Tracer) Enabled() bool {
	return t.config.Enabled
}

// HTTPMiddleware (Function): returns an HTTP middleware that wraps each request
// in a named span with method and path attributes, enabling per-request tracing.
func HTTPMiddleware(tracer *Tracer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, finish := tracer.StartSpan(r.Context(), "http.request",
				"http.method", r.Method,
				"http.target", r.URL.Path,
			)
			defer finish()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// contextKey (Type): uniquely typed key for context.WithValue lookups.
// Using a named int type guarantees keys never collide across packages or
// with other values in the same context chain.
type contextKey int

const (
	traceIDKey contextKey = iota
	spanIDKey
	spanKey
)

// getOrCreateTraceID (Function): retrieves the trace ID from context, generating a
// new one if none exists. Each trace represents one end-to-end request flow.
func getOrCreateTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(traceIDKey).(string); ok && traceID != "" {
		return traceID
	}
	return generateTraceID()
}

// attrsToMap (Function): converts a flat key-value string slice into a string map
// for efficient attribute storage on spans.
func attrsToMap(attrs []string) map[string]string {
	m := make(map[string]string)
	for i := 0; i < len(attrs)-1; i += 2 {
		m[attrs[i]] = attrs[i+1]
	}
	return m
}

// generateSpanID (Function): creates a random 8-byte hex span ID for OTLP compatibility.
func generateSpanID() string {
	return randomHex(8)
}

// generateTraceID (Function): creates a random 16-byte hex trace ID for OTLP compatibility.
func generateTraceID() string {
	return randomHex(16)
}

// WithTraceID (Function): stores the trace ID in context for propagation through
// / the call chain without explicit parameter passing.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

// GetTraceID (Function): extracts the trace ID from context, returning empty string
// if tracing has not been initialized for this request.
func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(traceIDKey).(string); ok {
		return traceID
	}
	return ""
}

// WithSpanID (Function): stores the span ID in context so child spans can reference
// their parent without direct parameter plumbing.
func WithSpanID(ctx context.Context, spanID string) context.Context {
	return context.WithValue(ctx, spanIDKey, spanID)
}

// GetSpanID (Function): extracts the current span ID from context, returning empty
// string if the request has no active span.
func GetSpanID(ctx context.Context) string {
	if spanID, ok := ctx.Value(spanIDKey).(string); ok {
		return spanID
	}
	return ""
}

// WithTraceParent (Function): parses a W3C Trace-Context header value and stores
// the extracted trace/span IDs in context for distributed trace propagation.
func WithTraceParent(ctx context.Context, traceParent string) context.Context {
	traceID, parentID, ok := ParseTraceParent(traceParent)
	if !ok {
		return ctx
	}
	ctx = WithTraceID(ctx, traceID)
	return WithSpanID(ctx, parentID)
}

// ParseTraceParent (Function): validates and splits a W3C traceparent header into
// its trace ID, parent span ID, and trace flags components.
func ParseTraceParent(traceParent string) (traceID, parentID string, ok bool) {
	parts := strings.Split(strings.TrimSpace(traceParent), "-")
	if len(parts) != 4 {
		return "", "", false
	}
	if parts[0] != "00" || len(parts[1]) != 32 || len(parts[2]) != 16 || len(parts[3]) != 2 {
		return "", "", false
	}

	traceID, parentID = parts[1], parts[2]
	flags := parts[3]

	if !isValidHex(traceID) || !isValidHex(parentID) || !isValidHex(flags) {
		return "", "", false
	}

	if isAllZero(traceID) || isAllZero(parentID) {
		return "", "", false
	}

	return traceID, parentID, true
}

// isValidHex (Function): checks that every character in the string is a valid
// hexadecimal digit, used to validate trace and span IDs.
func isValidHex(s string) bool {
	for _, r := range s {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
			return false
		}
	}
	return len(s) > 0
}

// isAllZero (Function): returns true when the string consists entirely of zero
// characters, used to reject invalid zero-value trace and span IDs.
func isAllZero(s string) bool {
	for _, r := range s {
		if r != '0' {
			return false
		}
	}
	return true
}

// FormatTraceParent (Function): builds a W3C traceparent header string from the
// given trace and span IDs for downstream propagation.
func FormatTraceParent(traceID, spanID string) string {
	if len(traceID) != 32 || len(spanID) != 16 {
		return ""
	}
	return fmt.Sprintf("00-%s-%s-01", strings.ToLower(traceID), strings.ToLower(spanID))
}

// shouldDropSpan (Function): probabilistically decides whether to sample a span
// based on the configured sample rate, using crypto/rand for unbiased selection.
func shouldDropSpan(sampleRate float64) bool {
	if sampleRate <= 0 {
		return true
	}
	if sampleRate >= 1 {
		return false
	}
	b := make([]byte, 2)
	if _, err := rand.Read(b); err != nil {
		return false
	}
	v := float64(int(b[0])<<8|int(b[1])) / 65535.0
	return v > sampleRate
}

// randomHex (Function): generates n random bytes and returns them as a hex string,
// falling back to nanosecond timestamps if crypto/rand fails.
func randomHex(n int) string {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("%0*x", n*2, time.Now().UnixNano())
	}
	return hex.EncodeToString(buf)
}

// otlpClient (Struct): manages the HTTP connection and serialization for sending
// OTLP trace data to an OpenTelemetry-compatible backend.
type otlpClient struct {
	endpoint    string
	insecure    bool
	timeout     time.Duration
	client      *http.Client
	serviceName string
	environment string
}

// newOTLPClient (Constructor): creates an OTLP HTTP client configured with the
// given endpoint, security mode, and timeout. Uses HTTPS by default.
func newOTLPClient(endpoint string, insecure bool, timeout time.Duration, serviceName, environment string) *otlpClient {
	scheme := "https"
	if insecure {
		scheme = "http"
	}

	if serviceName == "" {
		serviceName = "opslane"
	}

	return &otlpClient{
		endpoint:    fmt.Sprintf("%s://%s/v1/traces", scheme, endpoint),
		insecure:    insecure,
		timeout:     timeout,
		serviceName: serviceName,
		environment: environment,
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				DisableCompression: true,
			},
		},
	}
}

// Export (Method): serializes a batch of spans into OTLP JSON format and sends
// them to the configured endpoint via HTTP POST.
func (c *otlpClient) Export(ctx context.Context, spans []Span) error {
	if len(spans) == 0 {
		return nil
	}

	attrs := []map[string]any{
		{"key": "service.name", "value": map[string]string{"stringValue": c.serviceName}},
	}
	if c.environment != "" {
		attrs = append(attrs, map[string]any{
			"key":   "deployment.environment",
			"value": map[string]string{"stringValue": c.environment},
		})
	}

	payload := map[string]any{
		"resourceSpans": []map[string]any{
			{
				"resource": map[string]any{
					"attributes": attrs,
				},
				"scopeSpans": []map[string]any{
					{
						"scope": map[string]string{"name": "opslane.internal.otel"},
						"spans": toOTLPSpans(spans),
					},
				},
			},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal otlp payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, strings.NewReader(string(body)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "opslane/1.0")
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("send otlp request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("otlp exporter status=%d body=%q", resp.StatusCode, string(b))
	}
	return nil
}

// toOTLPSpans (Function): converts internal Span structs into the OTLP JSON wire
// format for transmission to the telemetry backend.
func toOTLPSpans(spans []Span) []map[string]any {
	out := make([]map[string]any, 0, len(spans))
	for _, span := range spans {
		attrs := make([]map[string]any, 0, len(span.Attributes))
		for key, value := range span.Attributes {
			attrs = append(attrs, map[string]any{
				"key":   key,
				"value": map[string]string{"stringValue": value},
			})
		}

		name := span.Name
		if name == "" {
			name = "unnamed"
		}
		entry := map[string]any{
			"traceId":           span.TraceID,
			"spanId":            span.SpanID,
			"parentSpanId":      span.ParentID,
			"name":              name,
			"startTimeUnixNano": strconv.FormatInt(span.StartTime.UnixNano(), 10),
			"endTimeUnixNano":   strconv.FormatInt(span.EndTime.UnixNano(), 10),
			"attributes":        attrs,
		}
		out = append(out, entry)
	}
	return out
}

// Close (Method): closes idle HTTP connections held by the OTLP client's transport.
func (c *otlpClient) Close() {
	if c.client != nil {
		c.client.CloseIdleConnections()
	}
}
