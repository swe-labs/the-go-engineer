package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

type SpanContext struct {
	TraceID      string
	SpanID       string
	ParentSpanID string
}

type Span struct {
	Context    SpanContext
	Operation  string
	StartTime  time.Time
	EndTime    time.Time
	Attributes map[string]string
	Events     []SpanEvent
	Children   []*Span
	parent     *Span
}

type SpanEvent struct {
	Timestamp time.Time
	Name      string
}

func generateID(bytes int) string {
	b := make([]byte, bytes)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func NewTrace() SpanContext {
	return SpanContext{
		TraceID: generateID(16),
		SpanID:  generateID(8),
	}
}

func newChildContext(parent SpanContext) SpanContext {
	return SpanContext{
		TraceID:      parent.TraceID,
		SpanID:       generateID(8),
		ParentSpanID: parent.SpanID,
	}
}

func StartSpan(ctx SpanContext, operation string) *Span {
	return &Span{
		Context:    ctx,
		Operation:  operation,
		StartTime:  time.Now(),
		Attributes: make(map[string]string),
	}
}

func (s *Span) End() {
	s.EndTime = time.Now()
}

func (s *Span) SetAttribute(key, value string) {
	s.Attributes[key] = value
}

func (s *Span) AddEvent(name string) {
	s.Events = append(s.Events, SpanEvent{Timestamp: time.Now(), Name: name})
}

func (s *Span) StartChild(operation string) *Span {
	child := StartSpan(newChildContext(s.Context), operation)
	child.parent = s
	s.Children = append(s.Children, child)
	return child
}

func (s *Span) Duration() time.Duration {
	if s.EndTime.IsZero() {
		return time.Since(s.StartTime)
	}
	return s.EndTime.Sub(s.StartTime)
}

func (s *Span) String() string {
	return s.format("")
}

func (s *Span) format(indent string) string {
	result := fmt.Sprintf("%s%s [%s] (%v)", indent, s.Operation, s.Context.SpanID[:6], s.Duration())
	for k, v := range s.Attributes {
		result += fmt.Sprintf("\n%s  %s: %s", indent, k, v)
	}
	for _, e := range s.Events {
		result += fmt.Sprintf("\n%s  ! %s", indent, e.Name)
	}
	for _, c := range s.Children {
		result += "\n" + c.format(indent+"  ")
	}
	return result
}

func simulateRequest() *Span {
	rootCtx := NewTrace()
	root := StartSpan(rootCtx, "GET /api/orders")
	root.SetAttribute("http.method", "GET")
	root.SetAttribute("http.path", "/api/orders")

	auth := root.StartChild("auth.verify_token")
	auth.SetAttribute("token.source", "header")
	time.Sleep(5 * time.Millisecond)
	auth.AddEvent("token.cached")
	auth.End()

	db := root.StartChild("db.query_orders")
	db.SetAttribute("db.statement", "SELECT * FROM orders WHERE user_id = ?")
	time.Sleep(20 * time.Millisecond)
	db.End()

	payment := root.StartChild("payment.check_status")
	payment.SetAttribute("payment.provider", "stripe")
	time.Sleep(15 * time.Millisecond)
	payment.End()

	root.End()
	return root
}

func main() {
	root := simulateRequest()
	fmt.Println(root)
}
