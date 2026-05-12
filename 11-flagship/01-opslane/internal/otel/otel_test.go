package otel

import (
	"context"
	"log/slog"
	"sync"
	"testing"
	"time"
)

func TestNewDisabled(t *testing.T) {
	t.Parallel()

	cfg := Config{Enabled: false}
	tracer := New(cfg, slog.Default())

	if tracer.Enabled() {
		t.Fatal("expected disabled tracer")
	}
}

func TestNewEnabledNoEndpointNoBlock(t *testing.T) {
	t.Parallel()

	cfg := Config{
		Enabled:     true,
		Timeout:     time.Second,
		SampleRate:  1.0,
		ServiceName: "test",
		Environment: "test",
	}
	tracer := New(cfg, slog.Default())
	defer tracer.Stop()

	if !tracer.Enabled() {
		t.Fatal("expected enabled tracer")
	}
}

func TestSpanLifecycleWithoutExport(t *testing.T) {
	t.Parallel()

	cfg := Config{
		Enabled:     true,
		Timeout:     time.Second,
		SampleRate:  1.0,
		ServiceName: "test",
		Environment: "test",
	}
	tracer := New(cfg, slog.Default())

	ctx, finish := tracer.StartSpan(context.Background(), "test.operation",
		"key1", "value1",
		"key2", "value2",
	)
	if ctx == nil {
		t.Fatal("expected non-nil context")
	}

	traceID := GetTraceID(ctx)
	if traceID == "" {
		t.Fatal("expected non-empty trace ID")
	}

	spanID := GetSpanID(ctx)
	if spanID == "" {
		t.Fatal("expected non-empty span ID")
	}

	finish()
	tracer.Stop()
}

func TestTraceIDPropagation(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctx = WithTraceID(ctx, "abcdef0123456789abcdef0123456789")

	got := GetTraceID(ctx)
	if got != "abcdef0123456789abcdef0123456789" {
		t.Fatalf("trace ID = %q, want %q", got, "abcdef0123456789abcdef0123456789")
	}
}

func TestSpanIDPropagation(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctx = WithSpanID(ctx, "deadbeefdeadbeef")

	got := GetSpanID(ctx)
	if got != "deadbeefdeadbeef" {
		t.Fatalf("span ID = %q, want %q", got, "deadbeefdeadbeef")
	}
}

func TestTraceAndSpanIDKeysDoNotCollide(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctx = WithTraceID(ctx, "abcdef0123456789abcdef0123456789")
	ctx = WithSpanID(ctx, "deadbeefdeadbeef")

	traceID := GetTraceID(ctx)
	if traceID != "abcdef0123456789abcdef0123456789" {
		t.Fatalf("trace ID = %q, want abcdef... (keys collided)", traceID)
	}

	spanID := GetSpanID(ctx)
	if spanID != "deadbeefdeadbeef" {
		t.Fatalf("span ID = %q, want deadbeef... (keys collided)", spanID)
	}
}

func TestGetTraceIDFromEmptyContext(t *testing.T) {
	t.Parallel()

	if got := GetTraceID(context.Background()); got != "" {
		t.Fatalf("expected empty trace ID, got %q", got)
	}
}

func TestGetSpanIDFromEmptyContext(t *testing.T) {
	t.Parallel()

	if got := GetSpanID(context.Background()); got != "" {
		t.Fatalf("expected empty span ID, got %q", got)
	}
}

func TestWithTraceParentValid(t *testing.T) {
	t.Parallel()

	parent := "00-abcdef0123456789abcdef0123456789-deadbeefdeadbeef-01"
	ctx := WithTraceParent(context.Background(), parent)

	traceID := GetTraceID(ctx)
	if traceID != "abcdef0123456789abcdef0123456789" {
		t.Fatalf("trace ID = %q, want abcdef...", traceID)
	}

	spanID := GetSpanID(ctx)
	if spanID != "deadbeefdeadbeef" {
		t.Fatalf("span ID = %q, want deadbeef...", spanID)
	}
}

func TestWithTraceParentInvalid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
	}{
		{name: "empty", input: ""},
		{name: "wrong version", input: "01-abc-def-01"},
		{name: "short trace ID", input: "00-abc-def-01"},
		{name: "short span ID", input: "00-abcdef0123456789abcdef0123456789-abc-01"},
		{name: "all zero trace", input: "00-00000000000000000000000000000000-deadbeefdeadbeef-01"},
		{name: "all zero span", input: "00-abcdef0123456789abcdef0123456789-0000000000000000-01"},
		{name: "non-hex chars", input: "00-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx-deadbeefdeadbeef-01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := WithTraceParent(context.Background(), tt.input)
			if traceID := GetTraceID(ctx); traceID != "" {
				t.Fatalf("expected empty trace ID for invalid input, got %q", traceID)
			}
		})
	}
}

func TestParseTraceParent(t *testing.T) {
	t.Parallel()

	traceID, spanID, ok := ParseTraceParent("00-abcdef0123456789abcdef0123456789-deadbeefdeadbeef-01")
	if !ok {
		t.Fatal("expected valid parse")
	}
	if traceID != "abcdef0123456789abcdef0123456789" {
		t.Fatalf("trace ID = %q", traceID)
	}
	if spanID != "deadbeefdeadbeef" {
		t.Fatalf("span ID = %q", spanID)
	}
}

func TestParseTraceParentInvalid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
	}{
		{name: "empty", input: ""},
		{name: "too few parts", input: "00-abc-def"},
		{name: "too many parts", input: "00-abc-def-01-extra"},
		{name: "non-hex", input: "00-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx-yyyyyyyyyyyyyyyy-01"},
		{name: "all zero trace", input: "00-00000000000000000000000000000000-deadbeefdeadbeef-01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, ok := ParseTraceParent(tt.input)
			if ok {
				t.Fatalf("expected parse failure for %q", tt.name)
			}
		})
	}
}

func TestFormatTraceParent(t *testing.T) {
	t.Parallel()

	result := FormatTraceParent("abcdef0123456789abcdef0123456789", "deadbeefdeadbeef")
	want := "00-abcdef0123456789abcdef0123456789-deadbeefdeadbeef-01"
	if result != want {
		t.Fatalf("FormatTraceParent = %q, want %q", result, want)
	}
}

func TestFormatTraceParentShortID(t *testing.T) {
	t.Parallel()

	if result := FormatTraceParent("short", "also-short"); result != "" {
		t.Fatalf("expected empty result for short IDs, got %q", result)
	}
}

func TestIsValidHex(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  bool
	}{
		{input: "abcdef0123456789", want: true},
		{input: "ABCDEF0123456789", want: true},
		{input: "xyz", want: false},
		{input: "", want: false},
	}

	for _, tt := range tests {
		if got := isValidHex(tt.input); got != tt.want {
			t.Fatalf("isValidHex(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestIsAllZero(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  bool
	}{
		{input: "0000000000000000", want: true},
		{input: "abcdef0123456789", want: false},
		{input: "", want: true},
	}

	for _, tt := range tests {
		if got := isAllZero(tt.input); got != tt.want {
			t.Fatalf("isAllZero(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestAttrsToMap(t *testing.T) {
	t.Parallel()

	result := attrsToMap([]string{"key1", "value1", "key2", "value2"})
	if len(result) != 2 {
		t.Fatalf("expected 2 attributes, got %d", len(result))
	}
	if result["key1"] != "value1" {
		t.Fatalf("key1 = %q", result["key1"])
	}
	if result["key2"] != "value2" {
		t.Fatalf("key2 = %q", result["key2"])
	}
}

func TestAttrsToMapOddCount(t *testing.T) {
	t.Parallel()

	result := attrsToMap([]string{"key1", "value1", "key2"})
	if len(result) != 1 {
		t.Fatalf("expected 1 attribute for odd input, got %d", len(result))
	}
}

func TestRandomHexLength(t *testing.T) {
	t.Parallel()

	result := randomHex(8)
	if len(result) != 16 {
		t.Fatalf("randomHex(8) length = %d, want 16", len(result))
	}
}

func TestShouldDropSpan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		sampleRate float64
		wantDrop   bool
	}{
		{name: "zero rate drops all", sampleRate: 0, wantDrop: true},
		{name: "negative rate drops all", sampleRate: -1, wantDrop: true},
		{name: "full rate keeps all", sampleRate: 1.0, wantDrop: false},
		{name: "above one keeps all", sampleRate: 2.0, wantDrop: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldDropSpan(tt.sampleRate); got != tt.wantDrop {
				t.Fatalf("shouldDropSpan(%f) = %v, want %v", tt.sampleRate, got, tt.wantDrop)
			}
		})
	}
}

func TestGenerateSpanID(t *testing.T) {
	t.Parallel()

	id := generateSpanID()
	if len(id) != 16 {
		t.Fatalf("span ID length = %d, want 16", len(id))
	}
	if !isValidHex(id) {
		t.Fatalf("span ID contains non-hex characters: %q", id)
	}
}

func TestGenerateTraceID(t *testing.T) {
	t.Parallel()

	id := generateTraceID()
	if len(id) != 32 {
		t.Fatalf("trace ID length = %d, want 32", len(id))
	}
	if !isValidHex(id) {
		t.Fatalf("trace ID contains non-hex characters: %q", id)
	}
}

func TestOTLPSpanSerialization(t *testing.T) {
	t.Parallel()

	now := time.Now()
	spans := []Span{
		{
			TraceID:       "abcdef0123456789abcdef0123456789",
			SpanID:        "deadbeefdeadbeef",
			ParentID:      "",
			Name:          "test.op",
			StartTime:     now,
			EndTime:       now.Add(100 * time.Millisecond),
			Attributes:    map[string]string{"key": "value"},
			Status:        "ok",
			StatusMessage: "",
		},
		{
			TraceID:   "abcdef0123456789abcdef0123456789",
			SpanID:    "cafebabecafebabe",
			ParentID:  "deadbeefdeadbeef",
			Name:      "",
			StartTime: now,
			EndTime:   now.Add(50 * time.Millisecond),
		},
	}

	result := toOTLPSpans(spans)
	if len(result) != 2 {
		t.Fatalf("expected 2 OTLP spans, got %d", len(result))
	}

	if result[0]["name"] != "test.op" {
		t.Fatalf("span[0] name = %q", result[0]["name"])
	}

	if result[1]["name"] != "unnamed" {
		t.Fatalf("span[1] name should default to 'unnamed', got %q", result[1]["name"])
	}
}

func TestDisabledTracerStartSpanReturnsNoopContext(t *testing.T) {
	t.Parallel()

	cfg := Config{Enabled: false}
	tracer := New(cfg, slog.Default())

	ctx, finish := tracer.StartSpan(context.Background(), "noop")
	if ctx == nil {
		t.Fatal("expected non-nil context even for disabled tracer")
	}

	finish()
}

func TestOTLPClientExportEmptySpans(t *testing.T) {
	t.Parallel()

	client := newOTLPClient("otlp.example.com:4318", true, time.Second, "test", "test")
	defer client.Close()

	if err := client.Export(context.Background(), nil); err != nil {
		t.Fatalf("expected nil for empty export, got %v", err)
	}

	if err := client.Export(context.Background(), []Span{}); err != nil {
		t.Fatalf("expected nil for empty export, got %v", err)
	}
}

func TestGenerateIDsAreUnique(t *testing.T) {
	t.Parallel()

	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := generateSpanID()
		if ids[id] {
			t.Fatal("duplicate span ID generated")
		}
		ids[id] = true
	}
}

func TestClientCloseIdempotent(t *testing.T) {
	t.Parallel()

	client := newOTLPClient("otlp.example.com:4318", true, time.Second, "test", "test")
	client.Close()
	client.Close()
}

func TestMultipleGoroutinesUseIndependentContexts(t *testing.T) {
	t.Parallel()

	cfg := Config{
		Enabled:     true,
		Timeout:     time.Second,
		SampleRate:  1.0,
		ServiceName: "test",
		Environment: "test",
	}
	tracer := New(cfg, slog.Default())
	defer tracer.Stop()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, finish := tracer.StartSpan(context.Background(), "goroutine")
			if GetTraceID(ctx) == "" {
				t.Error("expected non-empty trace ID in goroutine")
			}
			if GetSpanID(ctx) == "" {
				t.Error("expected non-empty span ID in goroutine")
			}
			finish()
		}()
	}
	wg.Wait()
}

func TestEnabledTracerStopWithoutExportLoopIsSafe(t *testing.T) {
	t.Parallel()

	cfg := Config{
		Enabled:     true,
		Timeout:     time.Second,
		SampleRate:  1.0,
		ServiceName: "test",
	}
	tracer := New(cfg, slog.Default())
	tracer.Stop()
}

func TestStartSpanSamplingDropsBelowThreshold(t *testing.T) {
	t.Parallel()

	cfg := Config{
		Enabled:     true,
		Timeout:     time.Second,
		SampleRate:  0.0,
		ServiceName: "test",
	}
	tracer := New(cfg, slog.Default())
	defer tracer.Stop()

	ctx, finish := tracer.StartSpan(context.Background(), "dropped")
	traceID := GetTraceID(ctx)
	if traceID != "" {
		t.Fatalf("expected empty trace ID when sample rate is 0, got %q", traceID)
	}
	finish()
}
