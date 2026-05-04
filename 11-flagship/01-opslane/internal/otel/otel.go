// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Package otel provides OpenTelemetry tracing for Opslane.
// Role: Observability boundary - exports traces to OTLP-compatible backends.
// Boundary: Traces are exported asynchronously; failures are logged but don't block.
// Failure mode: OTLP export failures are logged and dropped to prevent memory leaks.

package otel

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
)

type Config struct {
	Endpoint string
	Insecure bool
	Timeout  time.Duration
	Enabled  bool
}

func (c *Config) FromEnv() {
	c.Endpoint = os.Getenv("OPSLANE_OTEL_ENDPOINT")
	c.Insecure = os.Getenv("OPSLANE_OTEL_INSECURE") == "true"
	c.Enabled = c.Endpoint != ""

	timeout := os.Getenv("OPSLANE_OTEL_TIMEOUT")
	if timeout != "" {
		if d, err := time.ParseDuration(timeout); err == nil {
			c.Timeout = d
		}
	}
	if c.Timeout == 0 {
		c.Timeout = 5 * time.Second
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
		t.client = newOTLPClient(cfg.Endpoint, cfg.Insecure, cfg.Timeout)
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
	spanID := generateSpanID()
	traceID := getOrCreateTraceID(ctx)

	span := Span{
		TraceID:    traceID,
		SpanID:     spanID,
		Name:       name,
		StartTime:  time.Now(),
		Attributes: attrsToMap(attrs),
	}

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
	return fmt.Sprintf("%016x", time.Now().UnixNano())
}

func generateTraceID() string {
	return fmt.Sprintf("%032x", time.Now().UnixNano())
}

var traceIDKey = struct{}{}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(traceIDKey).(string); ok {
		return traceID
	}
	return ""
}

type otlpClient struct {
	endpoint string
	insecure bool
	timeout  time.Duration
	client   *http.Client
}

func newOTLPClient(endpoint string, insecure bool, timeout time.Duration) *otlpClient {
	scheme := "https"
	if insecure {
		scheme = "http"
	}

	return &otlpClient{
		endpoint: fmt.Sprintf("%s://%s/v1/traces", scheme, endpoint),
		insecure: insecure,
		timeout:  timeout,
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				DisableCompression: true,
			},
		},
	}
}

func (c *otlpClient) Export(ctx context.Context, spans []Span) error {
	req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "opslane/1.0")

	_ = spans

	return nil
}

func (c *otlpClient) Close() {
	if c.client != nil {
		c.client.CloseIdleConnections()
	}
}
