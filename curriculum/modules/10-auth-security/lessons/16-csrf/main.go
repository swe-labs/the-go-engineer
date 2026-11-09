package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
)

type CSRFProtector struct {
	mu     sync.RWMutex
	tokens map[string]string
	key    []byte
}

func NewCSRFProtector() *CSRFProtector {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}
	return &CSRFProtector{
		tokens: make(map[string]string),
		key:    key,
	}
}

func (c *CSRFProtector) GenerateToken(sessionID string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		panic(err)
	}
	encoded := hex.EncodeToString(token)
	salted := append([]byte(encoded), c.key...)
	hash := sha256.Sum256(salted)
	c.tokens[string(hash[:])] = sessionID
	return encoded
}

func (c *CSRFProtector) ValidateToken(token, sessionID string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	salted := append([]byte(token), c.key...)
	hash := sha256.Sum256(salted)
	key := string(hash[:])
	storedSession, ok := c.tokens[key]
	if !ok {
		return false
	}
	delete(c.tokens, key)
	return storedSession == sessionID
}

type csrfHandler struct {
	protector *CSRFProtector
}

func (h *csrfHandler) login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		sessionID := "session-abc-123"
		token := h.protector.GenerateToken(sessionID)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `{"session":"%s","csrf_token":"%s"}`, sessionID, token)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *csrfHandler) transfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	sessionID := r.Header.Get("X-Session-ID")
	csrfToken := r.Header.Get("X-CSRF-Token")
	if !h.protector.ValidateToken(csrfToken, sessionID) {
		http.Error(w, "CSRF validation failed", http.StatusForbidden)
		return
	}
	w.Write([]byte(`{"status":"transferred"}`))
}

func middlewareCSRF(next http.Handler, protector *CSRFProtector) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
			sameSite := r.Header.Get("X-Same-Site")
			if sameSite == "strict" {
				next.ServeHTTP(w, r)
				return
			}
			origin := r.Header.Get("Origin")
			referer := r.Header.Get("Referer")
			if origin == "" && referer == "" {
				http.Error(w, "CSRF: missing origin or referer", http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	protector := NewCSRFProtector()
	h := &csrfHandler{protector: protector}

	mux := http.NewServeMux()
	mux.HandleFunc("/login", h.login)
	mux.HandleFunc("/transfer", h.transfer)

	handler := middlewareCSRF(mux, protector)

	ts := httptest.NewServer(handler)
	defer ts.Close()

	fmt.Println("=== CSRF Protection Demo ===")

	resp, _ := http.Post(ts.URL+"/login", "application/json", http.NoBody)
	fmt.Println("POST /login response:", resp.Status)
	resp.Body.Close()

	req, _ := http.NewRequest("POST", ts.URL+"/transfer", http.NoBody)
	req.Header.Set("X-Session-ID", "session-abc-123")
	req.Header.Set("X-CSRF-Token", "invalid-token")
	resp, _ = http.DefaultClient.Do(req)
	body := new(strings.Builder)
	io.Copy(body, resp.Body)
	fmt.Println("POST /transfer (bad token):", resp.Status, body.String())
	resp.Body.Close()

	fmt.Println("\n=== SameSite cookie note ===")
	fmt.Println("Real CSRF protection uses SameSite=Strict cookies + CSRF tokens.")
	fmt.Println("This example shows the token validation pattern statelessly.")
}
