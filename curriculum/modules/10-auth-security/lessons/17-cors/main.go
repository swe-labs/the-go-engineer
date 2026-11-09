package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

func corsMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			allowed := ""
			for _, o := range allowedOrigins {
				if o == "*" || o == origin {
					allowed = o
					break
				}
			}
			if allowed != "" {
				w.Header().Set("Access-Control-Allow-Origin", allowed)
				w.Header().Set("Vary", "Origin")
			}

			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-CSRF-Token")
				w.Header().Set("Access-Control-Max-Age", "86400")
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func corsMiddlewareWithCredentials(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			allowed := ""
			for _, o := range allowedOrigins {
				if o == origin {
					allowed = o
					break
				}
			}

			if allowed != "" {
				w.Header().Set("Access-Control-Allow-Origin", allowed)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Vary", "Origin")
			}

			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-CSRF-Token")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Max-Age", "86400")
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func validateOriginHeader(origin string, allowedOrigins []string) string {
	if origin == "" {
		return ""
	}
	for _, o := range allowedOrigins {
		if o == "*" {
			return "*"
		}
		if strings.EqualFold(o, origin) {
			return origin
		}
	}
	return ""
}

func main() {
	fmt.Println("=== CORS Middleware Demo ===")

	mux := http.NewServeMux()
	mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"Hello from the API"}`))
	})

	handler := corsMiddleware([]string{"https://example.com", "https://app.example.com"})

	ts := httptest.NewServer(handler(mux))
	defer ts.Close()

	fmt.Println("\nTest server started at:", ts.URL)

	req, _ := http.NewRequest("GET", ts.URL+"/api/data", nil)
	req.Header.Set("Origin", "https://example.com")
	resp, _ := http.DefaultClient.Do(req)
	fmt.Println("Allowed origin response status:", resp.Status)
	fmt.Println("ACAO header:", resp.Header.Get("Access-Control-Allow-Origin"))
	resp.Body.Close()

	req, _ = http.NewRequest("GET", ts.URL+"/api/data", nil)
	req.Header.Set("Origin", "https://evil.com")
	resp, _ = http.DefaultClient.Do(req)
	fmt.Println("Disallowed origin response status:", resp.Status)
	fmt.Println("ACAO header:", resp.Header.Get("Access-Control-Allow-Origin"))
	resp.Body.Close()

	req, _ = http.NewRequest("OPTIONS", ts.URL+"/api/data", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	resp, _ = http.DefaultClient.Do(req)
	fmt.Println("Preflight response status:", resp.Status)
	fmt.Println("ACAM header:", resp.Header.Get("Access-Control-Allow-Methods"))
	fmt.Println("ACAC header:", resp.Header.Get("Access-Control-Allow-Credentials"))
	resp.Body.Close()

	fmt.Println("\nNote: CORS is enforced by the browser, not by curl or Go clients.")
}
