package main

import "testing"

func TestTypedNilTrap(t *testing.T) {
	err := doWork()
	if err == nil {
		t.Error("doWork() returned nil interface — expected non-nil (typed nil trap)")
	}
}

func TestTrueNilInterface(t *testing.T) {
	err := doWorkSafe()
	if err != nil {
		t.Error("doWorkSafe() returned non-nil — expected true nil interface")
	}
}

func TestGetErrorNil(t *testing.T) {
	err := getError(false)
	if err != nil {
		t.Error("getError(false) should be nil")
	}
}

func TestGetErrorNonNil(t *testing.T) {
	err := getError(true)
	if err == nil {
		t.Fatal("getError(true) should not be nil")
	}
	if err.Error() != "something went wrong" {
		t.Errorf("unexpected message: %s", err.Error())
	}
}
