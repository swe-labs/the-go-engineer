package main

import (
	"bytes"
	"testing"
)

func TestEmbeddedLoggerPromotion(t *testing.T) {
	s := Server{
		Logger: Logger{prefix: "test"},
		Host:   "localhost",
	}
	// Verify that Log is promoted
	if s.Logger.prefix != "test" {
		t.Errorf("expected prefix test, got %s", s.Logger.prefix)
	}
}

func TestInterfaceEmbeddingReadWriter(t *testing.T) {
	var buf bytes.Buffer
	var rw ReadWriter = &buf
	rw.Write([]byte("hello world"))
	readBuf := make([]byte, 11)
	n, err := rw.Read(readBuf)
	if err != nil {
		t.Fatal(err)
	}
	if string(readBuf[:n]) != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", string(readBuf[:n]))
	}
}

func TestHealthCheckerPing(t *testing.T) {
	var hc HealthChecker = Monitor{}
	if err := hc.Ping(); err != nil {
		t.Errorf("unexpected ping error: %v", err)
	}
}

func TestHealthCheckerStatus(t *testing.T) {
	var hc HealthChecker = Monitor{}
	status := hc.Status()
	if status["status"] != "ok" {
		t.Errorf("expected status ok, got %v", status)
	}
}

func TestMonitorHealthCheckerSatisfaction(t *testing.T) {
	m := Monitor{}
	var _ HealthChecker = m
	_ = m
}
