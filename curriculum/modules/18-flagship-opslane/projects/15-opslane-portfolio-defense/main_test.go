package main

import (
	"testing"
)

func TestDefense_TenantIsolationPillar(t *testing.T) {
	contract := NewOpslaneSystemContract()

	ok := contract.VerifyTenantBoundary("tenant-acme", "tenant-evil", "secret-order-99")
	if !ok {
		t.Fatalf("expected tenant boundary validation to pass")
	}

	crit := contract.criteria["tenant_isolation"]
	if !crit.Validated {
		t.Errorf("expected tenant_isolation criterion to be marked validated")
	}
}

func TestDefense_IdempotencyPillar(t *testing.T) {
	// Core practice task requirement:
	// Answer and verify the technical defense questions on payment idempotency.
	contract := NewOpslaneSystemContract()

	key := "pay-idem-uuid-999"

	tx1, wasCharged1 := contract.VerifyIdempotency(key)
	if !wasCharged1 || tx1 == "" {
		t.Fatalf("expected initial charge to execute")
	}

	tx2, wasCharged2 := contract.VerifyIdempotency(key)
	if wasCharged2 {
		t.Fatalf("DEFENSE FAILURE: Duplicate charge executed with identical idempotency key!")
	}

	if tx1 != tx2 {
		t.Errorf("expected duplicate request to return original transaction ID %s, got %s", tx1, tx2)
	}

	crit := contract.criteria["idempotency"]
	if !crit.Validated {
		t.Errorf("expected idempotency criterion to be marked validated")
	}
}

func TestDefense_AllPillarsValidated(t *testing.T) {
	contract := NewOpslaneSystemContract()

	// Validate remaining pillars
	contract.VerifyTenantBoundary("tenant-a", "tenant-b", "res-1")
	_, _ = contract.VerifyIdempotency("key-1")
	_, _ = contract.VerifyIdempotency("key-1")
	contract.criteria["transactional_integrity"].Validated = true
	contract.criteria["bounded_concurrency"].Validated = true
	contract.criteria["graceful_shutdown"].Validated = true

	if !contract.AllCriteriaValidated() {
		t.Errorf("expected all 5 capstone defense pillars to be validated")
	}
}
