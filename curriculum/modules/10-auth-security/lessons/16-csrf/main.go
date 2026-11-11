package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"net/url"
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

// OriginConfig lets callers behind TLS termination specify explicit public
// origins or reverse-proxy CIDRs whose X-Forwarded-Proto header is trusted.
type OriginConfig struct {
	TrustedOrigins []string
	TrustedProxies []netip.Prefix
}

type originPair struct {
	scheme string
	host   string
}

func parseTrustedOrigins(raw []string) []originPair {
	var out []originPair
	for _, entry := range raw {
		u, err := url.Parse(strings.TrimSpace(entry))
		if err != nil || u.User != nil || u.Host == "" {
			continue
		}
		if !strings.EqualFold(u.Scheme, "http") && !strings.EqualFold(u.Scheme, "https") {
			continue
		}
		if (u.Path != "" && u.Path != "/") || u.RawQuery != "" || u.Fragment != "" {
			continue
		}
		out = append(out, originPair{scheme: u.Scheme, host: u.Host})
	}
	return out
}

func isTrustedProxy(remoteAddr string, proxies []netip.Prefix) bool {
	if len(proxies) == 0 || remoteAddr == "" {
		return false
	}
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		host = remoteAddr
	}
	addr, err := netip.ParseAddr(host)
	if err != nil {
		return false
	}
	addr = addr.Unmap()
	for _, p := range proxies {
		if p.Contains(addr) {
			return true
		}
	}
	return false
}

func forwardedProto(r *http.Request) string {
	raw, _, _ := strings.Cut(r.Header.Get("X-Forwarded-Proto"), ",")
	proto := strings.TrimSpace(raw)
	if strings.EqualFold(proto, "https") {
		return "https"
	}
	if strings.EqualFold(proto, "http") {
		return "http"
	}
	return ""
}

func expectedOrigin(r *http.Request, proxies []netip.Prefix) originPair {
	if isTrustedProxy(r.RemoteAddr, proxies) {
		if proto := forwardedProto(r); proto != "" {
			return originPair{scheme: proto, host: r.Host}
		}
	}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return originPair{scheme: scheme, host: r.Host}
}

func matchesOrigin(rawURL string, allowed []originPair, allowPath bool) bool {
	if rawURL == "" || rawURL == "null" {
		return false
	}
	u, err := url.Parse(rawURL)
	if err != nil || u.User != nil {
		return false
	}
	if !allowPath && (u.Path != "" || u.RawQuery != "" || u.Fragment != "") {
		return false
	}
	for _, a := range allowed {
		if a.host != "" && strings.EqualFold(u.Scheme, a.scheme) && strings.EqualFold(u.Host, a.host) {
			return true
		}
	}
	return false
}

func middlewareCSRF(next http.Handler, cfg ...OriginConfig) http.Handler {
	var policy OriginConfig
	for _, c := range cfg {
		policy.TrustedOrigins = append(policy.TrustedOrigins, c.TrustedOrigins...)
		policy.TrustedProxies = append(policy.TrustedProxies, c.TrustedProxies...)
	}
	explicitMode := len(policy.TrustedOrigins) > 0
	configuredOrigins := parseTrustedOrigins(policy.TrustedOrigins)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			next.ServeHTTP(w, r)
			return
		}
		// Origin/Referer comparison against configured public origins or the
		// request scheme/host (honoring X-Forwarded-Proto only from trusted
		// proxies) is a secondary defense alongside CSRF tokens and SameSite cookies.
		allowed := configuredOrigins
		if !explicitMode {
			allowed = []originPair{expectedOrigin(r, policy.TrustedProxies)}
		}
		origin := r.Header.Get("Origin")
		referer := r.Header.Get("Referer")
		switch {
		case origin != "":
			if !matchesOrigin(origin, allowed, false) {
				http.Error(w, "CSRF: untrusted origin", http.StatusForbidden)
				return
			}
		case referer != "":
			if !matchesOrigin(referer, allowed, true) {
				http.Error(w, "CSRF: untrusted referer", http.StatusForbidden)
				return
			}
		default:
			http.Error(w, "CSRF: missing origin or referer", http.StatusForbidden)
			return
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

	handler := middlewareCSRF(mux)

	ts := httptest.NewServer(handler)
	defer ts.Close()

	fmt.Println("=== CSRF Protection Demo ===")

	// The demo client sends Origin like a same-origin browser would.
	postWithOrigin := func(path string, headers map[string]string) (int, string) {
		req, _ := http.NewRequest("POST", ts.URL+path, http.NoBody)
		req.Header.Set("Origin", ts.URL)
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return 0, err.Error()
		}
		defer resp.Body.Close()
		body := new(strings.Builder)
		io.Copy(body, resp.Body)
		return resp.StatusCode, body.String()
	}

	_, loginBody := postWithOrigin("/login", nil)
	fmt.Println("POST /login:", loginBody)
	var loginResp struct {
		Session   string `json:"session"`
		CSRFToken string `json:"csrf_token"`
	}
	json.Unmarshal([]byte(loginBody), &loginResp)

	status, badBody := postWithOrigin("/transfer", map[string]string{
		"X-Session-ID": loginResp.Session,
		"X-CSRF-Token": "invalid-token",
	})
	fmt.Printf("POST /transfer (bad token): %d %s\n", status, badBody)

	status, goodBody := postWithOrigin("/transfer", map[string]string{
		"X-Session-ID": loginResp.Session,
		"X-CSRF-Token": loginResp.CSRFToken,
	})
	fmt.Printf("POST /transfer (valid token): %d %s\n", status, goodBody)

	fmt.Println("\n=== SameSite cookie note ===")
	fmt.Println("Real CSRF protection uses SameSite=Strict cookies + CSRF tokens.")
	fmt.Println("SameSite is enforced via the cookie attribute set by the server,")
	fmt.Println("never via a request header (clients control their own headers).")
}
