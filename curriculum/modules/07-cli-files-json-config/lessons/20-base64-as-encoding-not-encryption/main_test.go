package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestEncodeDecodeText(t *testing.T) {
	original := []byte("Hello, Go!")
	enc, dec, err := EncodeDecode(original)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(enc) == 0 {
		t.Error("encoded string should not be empty")
	}
	if !bytes.Equal(original, dec) {
		t.Errorf("round-trip mismatch: got %v, want %v", dec, original)
	}
}

func TestEncodeDecodeBinary(t *testing.T) {
	original := []byte{0x00, 0xFF, 0x80, 0x40, 0x7F}
	enc, dec, err := EncodeDecode(original)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(enc) == 0 {
		t.Error("encoded string should not be empty")
	}
	if !bytes.Equal(original, dec) {
		t.Errorf("round-trip mismatch: got %v, want %v", dec, original)
	}
}

func TestEncodeDecodeEmpty(t *testing.T) {
	original := []byte{}
	enc, dec, err := EncodeDecode(original)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if enc != "" {
		t.Errorf("expected empty string, got %q", enc)
	}
	if len(dec) != 0 {
		t.Errorf("expected empty slice, got %v", dec)
	}
}

func TestURLEncodingSafe(t *testing.T) {
	original := []byte{0xFF, 0xFE, 0xFD, 0xFC, 0xFB}
	enc, _, _ := EncodeDecode(original)
	if strings.ContainsAny(enc, "+/") {
		t.Errorf("URLEncoding should not contain + or /, got %q", enc)
	}
}

func TestCompiles(t *testing.T) {
}
