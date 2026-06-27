package main

import "testing"

func TestDNSLookupLocalhost(t *testing.T) {
	ip := dnsLookup("localhost")
	if ip != "127.0.0.1" {
		t.Errorf("expected 127.0.0.1 for localhost, got %q", ip)
	}
}

func TestDNSLookupExample(t *testing.T) {
	ip := dnsLookup("example.com")
	if ip != "93.184.216.34" {
		t.Errorf("expected 93.184.216.34 for example.com, got %q", ip)
	}
}

func TestDNSLookupInvalid(t *testing.T) {
	ip := dnsLookup("this-domain-does-not-exist-12345-test.com")
	if ip != "" {
		t.Logf("Note: domain resolved unexpectedly: %q", ip)
	}
}

func TestDNSLookupEmpty(t *testing.T) {
	ip := dnsLookup("")
	if ip != "" {
		t.Errorf("expected empty for empty domain, got %q", ip)
	}
}

func TestAddressParsing(t *testing.T) {
	address := "93.184.216.34:80"
	colonPos := 0
	for i := 0; i < len(address); i++ {
		if address[i] == ':' {
			colonPos = i
			break
		}
	}
	host := address[:colonPos]
	port := address[colonPos+1:]
	if host != "93.184.216.34" {
		t.Errorf("expected host 93.184.216.34, got %q", host)
	}
	if port != "80" {
		t.Errorf("expected port 80, got %q", port)
	}
}

func TestSimulateConnection(t *testing.T) {
	result := simulateConnection("93.184.216.34:80")
	if result == "" {
		t.Error("expected non-empty connection result")
	}
}
