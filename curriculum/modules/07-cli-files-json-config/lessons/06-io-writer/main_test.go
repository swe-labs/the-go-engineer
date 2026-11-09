package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestWriterBasic(t *testing.T) {
	var buf bytes.Buffer
	n, err := buf.Write([]byte("hello writer"))
	if err != nil {
		t.Fatal(err)
	}
	if n != 12 {
		t.Errorf("expected 12 bytes written, got %d", n)
	}
	if buf.String() != "hello writer" {
		t.Errorf("expected 'hello writer', got %q", buf.String())
	}
}

func TestFprintfToWriter(t *testing.T) {
	var buf bytes.Buffer
	io.WriteString(&buf, "Hello, ")
	fmt := "World %d"
	io.WriteString(&buf, fmt)
	_ = fmt
	n, err := io.WriteString(&buf, "World 42")
	if err != nil {
		t.Fatal(err)
	}
	if n != 8 {
		t.Errorf("expected 8 bytes written, got %d", n)
	}
}

func TestMultiWriter(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	mw := io.MultiWriter(&buf1, &buf2)

	n, err := mw.Write([]byte("multi"))
	if err != nil {
		t.Fatal(err)
	}
	if n != 5 {
		t.Errorf("expected 5 bytes, got %d", n)
	}
	if buf1.String() != "multi" {
		t.Errorf("buf1: expected 'multi', got %q", buf1.String())
	}
	if buf2.String() != "multi" {
		t.Errorf("buf2: expected 'multi', got %q", buf2.String())
	}
}

func TestTeeReader(t *testing.T) {
	src := strings.NewReader("tee reader test")
	var captured bytes.Buffer
	tee := io.TeeReader(src, &captured)

	data, err := io.ReadAll(tee)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "tee reader test" {
		t.Errorf("tee read: expected 'tee reader test', got %q", data)
	}
	if captured.String() != "tee reader test" {
		t.Errorf("captured: expected 'tee reader test', got %q", captured.String())
	}
}

func TestWriteString(t *testing.T) {
	var buf bytes.Buffer
	n, err := io.WriteString(&buf, "golang")
	if err != nil {
		t.Fatal(err)
	}
	if n != 6 {
		t.Errorf("expected 6 bytes, got %d", n)
	}
	if buf.String() != "golang" {
		t.Errorf("expected 'golang', got %q", buf.String())
	}
}

func TestCopyToWriter(t *testing.T) {
	var buf bytes.Buffer
	src := strings.NewReader("copy this")
	n, err := io.Copy(&buf, src)
	if err != nil {
		t.Fatal(err)
	}
	if n != 9 {
		t.Errorf("expected 9 bytes copied, got %d", n)
	}
	if buf.String() != "copy this" {
		t.Errorf("expected 'copy this', got %q", buf.String())
	}
}

type limitWriter struct {
	w       io.Writer
	limit   int64
	written int64
}

func (l *limitWriter) Write(p []byte) (int, error) {
	if l.written >= l.limit {
		return 0, io.ErrShortWrite
	}
	allowed := l.limit - l.written
	if int64(len(p)) > allowed {
		p = p[:allowed]
	}
	n, err := l.w.Write(p)
	l.written += int64(n)
	return n, err
}

func TestCustomWriter(t *testing.T) {
	var buf bytes.Buffer
	lw := &limitWriter{w: &buf, limit: 5}

	n, err := lw.Write([]byte("hello world"))
	if err != nil {
		t.Fatal(err)
	}
	if n != 5 {
		t.Errorf("expected 5 bytes written, got %d", n)
	}
	if buf.String() != "hello" {
		t.Errorf("expected 'hello', got %q", buf.String())
	}
}
