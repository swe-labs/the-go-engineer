// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 09: Architecture & Security
// Title: TLS and HTTPS in Go
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn the transport-level rules that turn plain HTTP into encrypted, identity-checked HTTPS.
//
// WHY THIS MATTERS:
//   - TLS protects the channel by encrypting traffic and authenticating the server identity with certificates.
//
// RUN:
//   go run ./09-architecture/04-security/08-tls-and-https-in-go
//
// KEY TAKEAWAY:
//   - Learn the transport-level rules that turn plain HTTP into encrypted, identity-checked HTTPS.
// ============================================================================

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"
)

// handler (Function): responds to all HTTPS requests with a confirmation message.
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from the secure server! You reached: %s\n", r.URL.Path)
}

func main() {
	fmt.Println("=== SEC.8 TLS and HTTPS in Go ===")
	fmt.Println()

	// --- Generate self-signed certificate ---
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Failed to generate RSA key: %v", err)
	}

	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		log.Fatalf("Failed to generate serial number: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			Organization: []string{"The Go Engineer Demo"},
			CommonName:   "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(1 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %v", err)
	}

	// PEM-encode the certificate for display
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	fmt.Println("Generated self-signed certificate:")
	fmt.Printf("  Subject: %s\n", template.Subject.String())
	fmt.Printf("  Serial: %s\n", template.SerialNumber)
	fmt.Printf("  Valid until: %s\n", template.NotAfter.Format(time.RFC3339))
	fmt.Println("  DNS names: localhost")
	fmt.Println()

	// --- Start HTTPS server ---
	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		log.Fatalf("Failed to load key pair: %v", err)
	}

	server := &http.Server{
		Addr:    ":8443",
		Handler: http.HandlerFunc(handler),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{tlsCert},
			MinVersion:   tls.VersionTLS12,
		},
	}

	go func() {
		fmt.Printf("HTTPS server listening on https://localhost%s\n", server.Addr)
		if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Allow server to start
	time.Sleep(200 * time.Millisecond)

	// --- Test request with insecure skip (demo only) ---
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Get("https://localhost:8443/hello")
	if err != nil {
		log.Fatalf("Test request failed: %v", err)
	}
	defer resp.Body.Close()

	buf := make([]byte, 256)
	n, _ := resp.Body.Read(buf)
	fmt.Printf("Test response: %s", string(buf[:n]))
	fmt.Println()

	log.Println("TLS and HTTPS demo completed successfully.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: SEC.9 -> 09-architecture/04-security/09-secrets-management")
	fmt.Println("Current: SEC.8 (tls and https in go)")
	fmt.Println("---------------------------------------------------")
}
