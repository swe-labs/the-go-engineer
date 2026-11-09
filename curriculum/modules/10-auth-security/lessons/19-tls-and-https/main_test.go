package main

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerateCert(t *testing.T) {
	cert, err := generateCert()
	if err != nil {
		t.Fatalf("generateCert() error = %v", err)
	}

	if len(cert.Certificate) == 0 {
		t.Fatal("expected at least one certificate")
	}

	c, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		t.Fatalf("ParseCertificate error = %v", err)
	}

	if c.Subject.CommonName != "localhost" {
		t.Errorf("CN = %q, want %q", c.Subject.CommonName, "localhost")
	}

	if len(c.DNSNames) == 0 || c.DNSNames[0] != "localhost" {
		t.Errorf("DNSNames = %v, want [localhost]", c.DNSNames)
	}
}

func TestTLSServer(t *testing.T) {
	tlsConfig := serverConfig()
	mux := http.NewServeMux()
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	ts := httptest.NewUnstartedServer(mux)
	ts.TLS = tlsConfig
	ts.StartTLS()
	defer ts.Close()

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: clientConfig(),
		},
	}

	resp, err := client.Get(ts.URL + "/test")
	if err != nil {
		t.Fatalf("TLS request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	if resp.TLS == nil {
		t.Fatal("expected TLS connection state")
	}

	if resp.TLS.Version < tls.VersionTLS12 {
		t.Errorf("TLS version = %d, want >= %d", resp.TLS.Version, tls.VersionTLS12)
	}
}

func TestTLSConfigMinVersion(t *testing.T) {
	cfg := serverConfig()

	tests := []struct {
		name    string
		version uint16
		wantOK  bool
	}{
		{name: "TLS 1.2", version: tls.VersionTLS12, wantOK: true},
		{name: "TLS 1.3", version: tls.VersionTLS13, wantOK: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.version < cfg.MinVersion {
				t.Errorf("version %d is below min %d", tc.version, cfg.MinVersion)
			}
		})
	}
}
