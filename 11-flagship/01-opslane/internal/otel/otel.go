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
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Config struct {
	Endpoint    string
	Insecure    bool
	Timeout     time.Duration
	Enabled     bool
	SampleRate  float64
	ServiceName string
	Environment string
}

func (c *Config) FromEnv() {
	c.Endpoint = os.Getenv("OPSLANE_OTEL_ENDPOINT")
	c.Insecure = os.Getenv("OPSLANE_OTEL_INSECURE") == "true"
	c.Enabled = c.Endpoint != ""
	c.ServiceName = os.Getenv("OPSLANE_SERVICE_NAME")
	c.Environment = os.Getenv("OPSLANE_ENV")

	timeout := os.Getenv("OPSLANE_OTEL_TIMEOUT")
	if timeout != "" {
		if d, err := time.ParseDuration(timeout); err == nil {
			c.Timeout = d
		}
	}
	if c.Timeout == 0 {
		c.Timeout = 5 * time.Second
	}
	c.SampleRate = 1.0
	if sampleRate := os.Getenv("OPSLANE_OTEL_SAMPLE_RATE"); sampleRate != "" {
		if f, err := strconv.ParseFloat(sampleRate, 64); err == nil && f >= 0 && f <= 1 {
			c.SampleRate = f
		}
	}
}

type Tracer struct {
	config  Config
	client  *otlpClient
	spans   chan Span
	stopped chan struct{}
	wg      sync.WaitGroup
	logger  *slog.Logger
}

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

func (t *Tracer) Stop() {
	if t.client != nil {
		close(t.stopped)
		t.wg.Wait()
		t.client.Close()
	}
}

func (t *Tracer) Enabled() bool {
	return t.config.Enabled
}

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

var spanKey = struct{}{}

func getOrCreateTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(traceIDKey).(string); ok && traceID != "" {
		return traceID
	}
	return generateTraceID()
}

func attrsToMap(attrs []string) map[string]string {
	m := make(map[string]string)
	for i := 0; i < len(attrs)-1; i += 2 {
		m[attrs[i]] = attrs[i+1]
	}
	return m
}

func generateSpanID() string {
	return randomHex(8)
}

func generateTraceID() string {
	return randomHex(16)
}

var traceIDKey = struct{}{}
var spanIDKey = struct{}{}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(traceIDKey).(string); ok {
		return traceID
	}
	return ""
}

func WithSpanID(ctx context.Context, spanID string) context.Context {
	return context.WithValue(ctx, spanIDKey, spanID)
}

func GetSpanID(ctx context.Context) string {
	if spanID, ok := ctx.Value(spanIDKey).(string); ok {
		return spanID
	}
	return ""
}

func WithTraceParent(ctx context.Context, traceParent string) context.Context {
	traceID, parentID, ok := ParseTraceParent(traceParent)
	if !ok {
		return ctx
	}
	ctx = WithTraceID(ctx, traceID)
	return WithSpanID(ctx, parentID)
}

func ParseTraceParent(traceParent string) (traceID, parentID string, ok bool) {
	parts := strings.Split(strings.TrimSpace(traceParent), "-")
	if len(parts) != 4 {
		return "", "", false
	}
	if parts[0] != "00" || len(parts[1]) != 32 || len(parts[2]) != 16 || len(parts[3]) != 2 {
		return "", "", false
	}
	return parts[1], parts[2], true
}

func FormatTraceParent(traceID, spanID string) string {
	if len(traceID) != 32 || len(spanID) != 16 {
		return ""
	}
	return fmt.Sprintf("00-%s-%s-01", strings.ToLower(traceID), strings.ToLower(spanID))
}

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

func randomHex(n int) string {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("%0*x", n*2, time.Now().UnixNano())
	}
	return hex.EncodeToString(buf)
}

type otlpClient struct {
	endpoint     string
	insecure     bool
	timeout      time.Duration
	client       *http.Client
	serviceName  string
	environment  string
}

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

func (c *otlpClient) Close() {
	if c.client != nil {
		c.client.CloseIdleConnections()
	}
}
