package main

import (
	"testing"
)

func TestResolveDomain(t *testing.T) {
	ips, err := resolveDomain("localhost")
	if err != nil {
		t.Fatalf("resolveDomain(localhost) failed: %v", err)
	}
	if len(ips) == 0 {
		t.Fatal("expected at least one IP for localhost")
	}
}

func TestResolveDomainInvalid(t *testing.T) {
	_, err := resolveDomain("this-domain-does-not-exist-12345-test.com")
	if err == nil {
		t.Log("Note: domain resolved unexpectedly (network may be returning results)")
	}
}

func TestResolveDomainEmpty(t *testing.T) {
	_, err := resolveDomain("")
	if err == nil {
		t.Error("expected error for empty domain")
	}
}

func TestDialTCPInvalid(t *testing.T) {
	conn, err := dialTCP("0.0.0.0:1")
	if err == nil {
		if conn != nil {
			conn.Close()
		}
		t.Log("Note: connection succeeded unexpectedly")
	}
}
