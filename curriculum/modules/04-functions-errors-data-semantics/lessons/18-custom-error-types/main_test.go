package main

import (
	"errors"
	"testing"
)

func TestQueryDB_NotFound(t *testing.T) {
	err := queryDB("")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Error("expected errors.Is to match ErrNotFound")
	}
	var dbErr *DBError
	if !errors.As(err, &dbErr) {
		t.Fatal("expected *DBError")
	}
	if dbErr.Code != "NOT_FOUND" {
		t.Errorf("expected code NOT_FOUND, got %s", dbErr.Code)
	}
}

func TestQueryDB_DuplicateKey(t *testing.T) {
	err := queryDB("dup")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if errors.Is(err, ErrNotFound) {
		t.Error("expected NOT to match ErrNotFound")
	}
	var dbErr *DBError
	if !errors.As(err, &dbErr) {
		t.Fatal("expected *DBError")
	}
	if dbErr.Code != "DUPLICATE_KEY" {
		t.Errorf("expected code DUPLICATE_KEY, got %s", dbErr.Code)
	}
	if dbErr.Err == nil {
		t.Error("expected underlying error, got nil")
	}
}

func TestQueryDB_Valid(t *testing.T) {
	err := queryDB("valid")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestDBError_IsMethod(t *testing.T) {
	e1 := &DBError{Code: "NOT_FOUND", Message: "a"}
	e2 := &DBError{Code: "NOT_FOUND", Message: "b"}
	e3 := &DBError{Code: "OTHER", Message: "c"}

	if !e1.Is(e2) {
		t.Error("expected e1.Is(e2) to be true (same code)")
	}
	if e1.Is(e3) {
		t.Error("expected e1.Is(e3) to be false (different code)")
	}
}

func TestDBError_Unwrap(t *testing.T) {
	inner := errors.New("inner error")
	e := &DBError{Code: "ERR", Err: inner}
	if unwrapped := e.Unwrap(); unwrapped != inner {
		t.Errorf("expected inner error, got %v", unwrapped)
	}
}

func TestDBError_NilUnwrap(t *testing.T) {
	e := &DBError{Code: "ERR"}
	if e.Unwrap() != nil {
		t.Error("expected nil unwrap")
	}
}
