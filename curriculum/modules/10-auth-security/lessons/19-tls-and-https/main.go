package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"time"
)

func generateCert() (tls.Certificate, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Go Engineer Lesson"},
			CommonName:   "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return tls.Certificate{}, err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})

	return tls.X509KeyPair(certPEM, keyPEM)
}

func serverConfig() *tls.Config {
	cert, err := generateCert()
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}
}

func clientConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
	}
}

func main() {
	fmt.Println("=== TLS Server Demo ===")
	fmt.Println("Generating self-signed certificate...")

	tlsConfig := serverConfig()
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"Hello over TLS!"}`))
	})

	ts := httptest.NewUnstartedServer(mux)
	ts.TLS = tlsConfig
	ts.StartTLS()
	defer ts.Close()

	fmt.Println("TLS server at:", ts.URL)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: clientConfig(),
		},
	}

	resp, err := client.Get(ts.URL + "/hello")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Response:", string(body))
	fmt.Println("TLS version:", resp.TLS.Version)
	fmt.Println("Server cert CN:", resp.TLS.PeerCertificates[0].Subject.CommonName)

	fmt.Println("\n=== TLS config check ===")
	cfg := tlsConfig
	fmt.Printf("MinVersion: %d (TLS 1.2=%d, TLS 1.3=%d)\n", cfg.MinVersion, tls.VersionTLS12, tls.VersionTLS13)
	fmt.Printf("Certificates: %d\n", len(cfg.Certificates))
}
