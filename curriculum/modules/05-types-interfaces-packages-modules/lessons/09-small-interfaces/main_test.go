package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestUpperWriter(t *testing.T) {
	var buf bytes.Buffer
	uw := UpperWriter{w: &buf}
	uw.Write([]byte("hello"))
	if buf.String() != "HELLO" {
		t.Errorf("expected HELLO, got %s", buf.String())
	}
}

func TestCountingWriter(t *testing.T) {
	var buf bytes.Buffer
	cw := &CountingWriter{w: &buf}
	n, err := cw.Write([]byte("hello"))
	if err != nil {
		t.Fatal(err)
	}
	if n != 5 {
		t.Errorf("expected 5 bytes written, got %d", n)
	}
	if cw.BytesWritten() != 5 {
		t.Errorf("expected count 5, got %d", cw.BytesWritten())
	}
}

func TestCountingWriterMultipleWrites(t *testing.T) {
	var buf bytes.Buffer
	cw := &CountingWriter{w: &buf}
	cw.Write([]byte("a"))
	cw.Write([]byte("bb"))
	cw.Write([]byte("ccc"))
	if cw.BytesWritten() != 6 {
		t.Errorf("expected count 6, got %d", cw.BytesWritten())
	}
}

func TestCountingWriterImplementsWriter(t *testing.T) {
	var buf bytes.Buffer
	cw := &CountingWriter{w: &buf}
	var w io.Writer = cw
	w.Write([]byte("test"))
}

func TestStringReaderWithIOCopy(t *testing.T) {
	var buf bytes.Buffer
	r := strings.NewReader("hello world")
	n, err := io.Copy(&buf, r)
	if err != nil {
		t.Fatal(err)
	}
	if n != 11 {
		t.Errorf("expected 11 bytes copied, got %d", n)
	}
	if buf.String() != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", buf.String())
	}
}
