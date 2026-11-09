package main

import (
	"io"
	"strings"
	"testing"
)

func TestTypeAssertionString(t *testing.T) {
	v := interface{}("hello")
	s, ok := v.(string)
	if !ok {
		t.Fatal("expected ok for string assertion")
	}
	if s != "hello" {
		t.Errorf("expected 'hello', got '%s'", s)
	}
}

func TestTypeAssertionInt(t *testing.T) {
	v := interface{}(42)
	i, ok := v.(int)
	if !ok {
		t.Fatal("expected ok for int assertion")
	}
	if i != 42 {
		t.Errorf("expected 42, got %d", i)
	}
}

func TestTypeAssertionFail(t *testing.T) {
	v := interface{}("hello")
	_, ok := v.(int)
	if ok {
		t.Error("expected false for int assertion on string")
	}
}

func TestTypeAssertionBool(t *testing.T) {
	v := interface{}(true)
	b, ok := v.(bool)
	if !ok {
		t.Fatal("expected ok for bool assertion")
	}
	if b != true {
		t.Errorf("expected true, got %v", b)
	}
}

func TestTypeAssertionFloat64(t *testing.T) {
	v := interface{}(3.14)
	f, ok := v.(float64)
	if !ok {
		t.Fatal("expected ok for float64 assertion")
	}
	if f != 3.14 {
		t.Errorf("expected 3.14, got %f", f)
	}
}

func TestTypeAssertionToInterface(t *testing.T) {
	var r io.Reader = strings.NewReader("test")
	w, ok := r.(io.WriterTo)
	if ok {
		_ = w
	} else {
		// strings.Reader does not implement WriterTo
	}
}

func TestTypeAssertionToConcretePointer(t *testing.T) {
	var r io.Reader = strings.NewReader("test")
	sr, ok := r.(*strings.Reader)
	if !ok {
		t.Fatal("expected ok for *strings.Reader assertion")
	}
	if sr.Len() != 4 {
		t.Errorf("expected Len 4, got %d", sr.Len())
	}
}

func TestTypeAssertionPanicForm(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic from bad type assertion")
		}
	}()
	v := interface{}(42)
	_ = v.(string)
}

func TestClassify(t *testing.T) {
	items := []interface{}{"hello", 42, true, 3.14, nil}
	classify(items)
}

func TestPrintType(t *testing.T) {
	printType("hello")
	printType(42)
	printType(true)
}
