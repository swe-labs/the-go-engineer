package main

import (
	"strings"
	"testing"
	"time"
)

func TestNewTrace(t *testing.T) {
	ctx := NewTrace()
	if len(ctx.TraceID) != 32 {
		t.Errorf("expected trace ID length 32, got %d", len(ctx.TraceID))
	}
	if len(ctx.SpanID) != 16 {
		t.Errorf("expected span ID length 16, got %d", len(ctx.SpanID))
	}
	if ctx.ParentSpanID != "" {
		t.Errorf("expected empty parent span ID for root, got %s", ctx.ParentSpanID)
	}
}

func TestParentChildRelationship(t *testing.T) {
	ctx := NewTrace()
	root := StartSpan(ctx, "root")
	child := root.StartChild("child")

	if child.Context.TraceID != root.Context.TraceID {
		t.Error("child must share the same trace ID as parent")
	}
	if child.Context.ParentSpanID != root.Context.SpanID {
		t.Error("child's ParentSpanID must equal parent's SpanID")
	}
	if child.Context.SpanID == root.Context.SpanID {
		t.Error("child must have a different span ID than parent")
	}
	if len(root.Children) != 1 {
		t.Errorf("expected 1 child, got %d", len(root.Children))
	}
}

func TestSpanAttributes(t *testing.T) {
	ctx := NewTrace()
	span := StartSpan(ctx, "test")
	span.SetAttribute("key1", "value1")
	span.SetAttribute("key2", "value2")

	if span.Attributes["key1"] != "value1" {
		t.Errorf("expected attribute value1, got %s", span.Attributes["key1"])
	}
}

func TestSpanEvents(t *testing.T) {
	ctx := NewTrace()
	span := StartSpan(ctx, "test")
	span.AddEvent("cache.miss")

	if len(span.Events) != 1 {
		t.Errorf("expected 1 event, got %d", len(span.Events))
	}
	if span.Events[0].Name != "cache.miss" {
		t.Errorf("expected event name 'cache.miss', got '%s'", span.Events[0].Name)
	}
}

func TestSpanDuration(t *testing.T) {
	ctx := NewTrace()
	span := StartSpan(ctx, "test")
	time.Sleep(time.Nanosecond)
	span.End()

	if span.Duration() <= 0 {
		t.Error("expected positive duration after End()")
	}
}

func TestTraceTreeStructure(t *testing.T) {
	root := simulateRequest()

	if root.Operation != "GET /api/orders" {
		t.Errorf("expected root operation 'GET /api/orders', got '%s'", root.Operation)
	}
	if len(root.Children) != 3 {
		t.Errorf("expected 3 child spans, got %d", len(root.Children))
	}

	childOps := make(map[string]bool)
	for _, c := range root.Children {
		childOps[c.Operation] = true
	}
	if !childOps["auth.verify_token"] {
		t.Error("missing child span: auth.verify_token")
	}
	if !childOps["db.query_orders"] {
		t.Error("missing child span: db.query_orders")
	}
	if !childOps["payment.check_status"] {
		t.Error("missing child span: payment.check_status")
	}
}

func TestStringOutput(t *testing.T) {
	root := simulateRequest()
	output := root.String()

	if !strings.Contains(output, "GET /api/orders") {
		t.Error("string output missing root operation name")
	}
	if !strings.Contains(output, "auth.verify_token") {
		t.Error("string output missing child span")
	}
	if !strings.Contains(output, "http.method") {
		t.Error("string output missing attribute key")
	}
	if !strings.Contains(output, "token.cached") {
		t.Error("string output missing event name")
	}
}
