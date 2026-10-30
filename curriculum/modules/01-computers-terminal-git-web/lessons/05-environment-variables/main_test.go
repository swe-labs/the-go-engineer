package main

import (
	"testing"
)

func TestGetEnv(t *testing.T) {
	// 1. Test unset variable
	unsetVar := GetEnv("DEFINITELY_UNSET_ENV_VAR_12345")
	if unsetVar.IsSet {
		t.Errorf("expected variable to be unset")
	}

	// 2. Test set variable
	testKey := "TEST_ENV_VAR_GO_ENGINEER"
	testVal := "super_secret_value_42"
	t.Setenv(testKey, testVal)

	setVar := GetEnv(testKey)
	if !setVar.IsSet {
		t.Errorf("expected variable to be set")
	}
	if setVar.Value != testVal {
		t.Errorf("expected value %q, got %q", testVal, setVar.Value)
	}

	// 3. Test empty but set variable
	emptyKey := "TEST_EMPTY_ENV_VAR"
	t.Setenv(emptyKey, "")
	emptyVar := GetEnv(emptyKey)
	if !emptyVar.IsSet {
		t.Errorf("expected empty variable to be marked as set")
	}
	if emptyVar.Value != "" {
		t.Errorf("expected empty string value, got %q", emptyVar.Value)
	}
}

func TestMainExecution(t *testing.T) {
	main()
}
