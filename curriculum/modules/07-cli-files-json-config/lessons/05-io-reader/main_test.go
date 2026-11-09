package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestReadFull(t *testing.T) {
	r := strings.NewReader("hello reader")
	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello reader" {
		t.Errorf("expected 'hello reader', got %q", data)
	}
}

func TestReadPartial(t *testing.T) {
	r := strings.NewReader("abcdefghij")
	buf := make([]byte, 4)
	n, err := r.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if n != 4 {
		t.Errorf("expected 4 bytes read, got %d", n)
	}
	if string(buf) != "abcd" {
		t.Errorf("expected 'abcd', got %q", buf)
	}
}

func TestReadEof(t *testing.T) {
	r := strings.NewReader("x")
	buf := make([]byte, 10)
	n, err := r.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Errorf("expected 1 byte, got %d", n)
	}
	// Second read should return EOF
	n, err = r.Read(buf)
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v (n=%d)", err, n)
	}
}

func TestMultiReader(t *testing.T) {
	r := io.MultiReader(
		strings.NewReader("part1,"),
		strings.NewReader("part2,"),
		strings.NewReader("part3"),
	)
	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "part1,part2,part3" {
		t.Errorf("expected concatenated parts, got %q", data)
	}
}

func TestLimitReader(t *testing.T) {
	r := io.LimitReader(strings.NewReader("too much data"), 3)
	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "too" {
		t.Errorf("expected 'too', got %q", data)
	}
}

type testReader struct{ data []byte }

func (t *testReader) Read(p []byte) (int, error) {
	if len(t.data) == 0 {
		return 0, io.EOF
	}
	n := copy(p, t.data)
	t.data = t.data[n:]
	return n, nil
}

func TestCustomReader(t *testing.T) {
	r := &testReader{data: []byte("custom reader works")}
	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "custom reader works" {
		t.Errorf("expected 'custom reader works', got %q", data)
	}
}

func TestBufferRead(t *testing.T) {
	data := []byte("buffer test")
	r := bytes.NewReader(data)
	buf := make([]byte, 4)

	n, err := r.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if n != 4 || string(buf[:n]) != "buff" {
		t.Errorf("expected 'buff', got %q", buf[:n])
	}

	n, err = r.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if string(buf[:n]) != "er t" {
		t.Errorf("expected 'er t', got %q", buf[:n])
	}
}

func TestReadZeroBytes(t *testing.T) {
	r := strings.NewReader("data")
	_, err := r.Read(nil)
	if err != nil {
		t.Errorf("reading nil buffer should not error: %v", err)
	}
}
