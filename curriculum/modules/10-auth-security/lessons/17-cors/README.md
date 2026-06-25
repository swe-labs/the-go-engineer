# CORS

## Learning objective

Configure Cross-Origin Resource Sharing (CORS) headers in Go HTTP servers, handle preflight requests correctly, and understand the security boundaries of the browser's same-origin policy.

## Why this matters

Every modern web application makes cross-origin requests. A single-page app at `app.example.com` needs to call an API at `api.example.com`. Without CORS, the browser blocks these requests. Misconfiguring CORS is one of the most common security issues in web APIs: too permissive (`Access-Control-Allow-Origin: *`) exposes data to any site, while too restrictive breaks legitimate clients. Go engineers building HTTP APIs must understand CORS to both enable legitimate cross-origin access and prevent unauthorized data access.

## Mental model

CORS is not a server-side security mechanism. It is a browser-enforced policy on which origins can read the response to a cross-origin request. A malicious client like `curl` or Postman ignores CORS entirely and can read any response. CORS protects the user's browser session, not the API.

Think of CORS as a bouncer at a club. The bouncer checks your ID (Origin header). If you are on the guest list (allowed origins), the bouncer lets you in (browser shows the response to JavaScript). If you are not, the bouncer turns you away (browser hides the response and raises an error). But if you sneak in through the back door (direct HTTP client), there is no bouncer.

## Core idea

The browser's same-origin policy prevents JavaScript from reading responses from a different origin (protocol + host + port). CORS is a mechanism to selectively relax this policy.

A CORS request involves these HTTP headers:

| Header | Direction | Purpose |
|---|---|---|
| `Origin` | Request | Indicates the requesting origin |
| `Access-Control-Allow-Origin` | Response | Specifies which origins can read the response |
| `Access-Control-Allow-Methods` | Response | Lists allowed HTTP methods for the resource |
| `Access-Control-Allow-Headers` | Response | Lists allowed custom headers |
| `Access-Control-Allow-Credentials` | Response | Whether credentials (cookies, auth) are allowed |
| `Access-Control-Max-Age` | Response | How long preflight results can be cached |
| `Access-Control-Expose-Headers` | Response | Which headers JavaScript can access |

For simple requests (GET, HEAD, POST with content types `application/x-www-form-urlencoded`, `multipart/form-data`, or `text/plain`), the browser sends the request with the `Origin` header and checks the response for `Access-Control-Allow-Origin`.

For all other requests (PUT, DELETE, PATCH, JSON content type, custom headers), the browser first sends a preflight `OPTIONS` request to check permissions before sending the actual request.

## Under the hood

When the browser makes a cross-origin request:

1. The browser adds the `Origin` header to the request.
2. The server responds with CORS headers.
3. The browser checks whether the response's `Access-Control-Allow-Origin` matches the requesting origin.
4. If the header is missing or does not match, the browser blocks JavaScript from reading the response. The network request still completes, but the JavaScript runtime throws an error.

For preflighted requests:

1. Browser sends `OPTIONS` with `Origin` and `Access-Control-Request-Method` headers.
2. Server responds with `Access-Control-Allow-Origin`, `Access-Control-Allow-Methods`, `Access-Control-Allow-Headers`.
3. Browser checks the preflight response. If allowed, it sends the actual request.
4. The actual request also needs CORS headers on the response.

The `Vary: Origin` header is critical for caching. Without it, a CDN might serve a cached response with `ACAO: https://example.com` to a user from `https://evil.com`.

## How Go uses it

In Go, CORS is typically implemented as middleware. The standard `net/http` package does not include CORS middleware, but it is easy to write:

```go
func corsMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			// Check if origin is allowed
			// Set CORS headers
			// Handle preflight OPTIONS
			next.ServeHTTP(w, r)
		})
	}
}
```

Popular Go web frameworks provide CORS middleware:

- **Chi**: `github.com/go-chi/cors`
- **Gin**: `gin-contrib/cors`
- **Echo**: `echo/middleware.CORS()`
- **gorilla/handlers**: `handlers.CORS()`

## Go example

```go
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Max-Age", "86400")
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"Hello from the API"}`))
	})

	handler := corsMiddleware([]string{"https://example.com"})
	ts := httptest.NewServer(handler(mux))
	defer ts.Close()

	fmt.Println("Server at:", ts.URL)
}
```

## Step-by-step execution

For a cross-origin POST request with `Content-Type: application/json`:

1. Browser (at `https://app.example.com`) wants to POST JSON to `https://api.example.com/data`.
2. Browser checks: `Content-Type: application/json` is not a simple content type. This requires a preflight.
3. Browser sends OPTIONS to `https://api.example.com/data` with headers: `Origin: https://app.example.com`, `Access-Control-Request-Method: POST`, `Access-Control-Request-Headers: content-type`.
4. Server receives OPTIONS. CORS middleware checks: is `https://app.example.com` in the allowed origins list? Yes.
5. Server responds with: `Access-Control-Allow-Origin: https://app.example.com`, `Access-Control-Allow-Methods: GET, POST, PUT, DELETE, PATCH`, `Access-Control-Allow-Headers: Content-Type, Authorization`. Status 204 No Content.
6. Browser receives the preflight response. All headers match requirements.
7. Browser sends the actual POST request with `Origin: https://app.example.com`.
8. Server processes the request and responds with `Access-Control-Allow-Origin: https://app.example.com`.
9. Browser checks the CORS header. It matches. JavaScript can read the response.

## Common mistakes

- Mistake: Setting `Access-Control-Allow-Origin: *` with `Access-Control-Allow-Credentials: true`.
  - Why it happens: Developers need to support credentials across origins and use a wildcard for simplicity.
  - Fix: This combination is invalid per the CORS specification. When credentials are allowed, the origin must be explicit: `Access-Control-Allow-Origin: https://example.com`.

- Mistake: Reflecting the `Origin` header in `Access-Control-Allow-Origin` without validation (echo back).
  - Why it happens: Developers write `w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))` for convenience.
  - Fix: This allows any site to read the response. Always validate the origin against an allowlist.

- Mistake: Not handling preflight OPTIONS requests properly.
  - Why it happens: Developers only add CORS headers to actual endpoints, not the OPTIONS handler.
  - Fix: OPTIONS requests must get CORS headers and return 204, not 404 or the full response.

- Mistake: Using CORS as a security mechanism for server-side protection.
  - Why it happens: Developers think CORS prevents cross-origin requests entirely.
  - Fix: CORS only prevents the browser from showing the response to JavaScript. The request still reaches the server. Server-side authentication and authorization are still required.

- Mistake: Not setting `Vary: Origin` on CORS responses.
  - Why it happens: Developers set ACAO but forget Vary, causing CDN cache poisoning.
  - Fix: Always set `Vary: Origin` on responses that include `Access-Control-Allow-Origin`.

## Debugging walkthrough

Consider this broken CORS setup:

```go
func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
```

Symptom: The browser console shows "CORS error: The value of the 'Access-Control-Allow-Origin' header in the response must not be the wildcard '*' when the request's credentials mode is 'include'."

Investigation: The frontend JavaScript uses `fetch(url, { credentials: 'include' })`. The browser requires an explicit origin, not `*`, when credentials are included.

Root cause: The API uses `*` for ACAO but the frontend needs to send cookies (credentials).

Fix: Use an explicit origin:

```go
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "https://app.example.com" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Vary", "Origin")
		}
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
```

## Production notes

- Maintain an explicit allowlist of origins. Do not use `*` for any API that handles user-specific data or credentials.
- For public APIs that serve non-sensitive data (e.g., public weather data), `*` is acceptable.
- Use environment-specific origins: `CORS_ALLOWED_ORIGINS=https://app.example.com,https://staging.example.com`.
- Always set `Vary: Origin` to prevent CDN cache poisoning.
- For APIs behind a reverse proxy (Nginx, Cloudflare), CORS headers can be set at the proxy level.
- Test CORS configuration using browser dev tools, not curl. Curl ignores CORS entirely.

## Performance implications

- CORS adds no significant server-side overhead. A few header comparisons and string operations per request.
- Preflight requests double the number of HTTP requests for non-simple cross-origin requests. Use `Access-Control-Max-Age` to cache preflight responses (up to 86400 seconds).
- Headers add a small amount of response size (typically 100-300 bytes).
- The main performance cost is the network latency of the extra preflight round trip, which can be mitigated by `Max-Age`.

## Practice task

Write a function `corsHandler(allowedOrigins []string, next http.Handler) http.Handler` that implements a complete CORS middleware. It should:

1. Check the `Origin` header against the allowed origins list.
2. For matching origins, set `Access-Control-Allow-Origin` and `Vary: Origin`.
3. For credentials (when `Access-Control-Allow-Credentials` is needed), ensure the origin is explicit (not `*`).
4. Handle OPTIONS preflight requests with `Access-Control-Allow-Methods`, `Access-Control-Allow-Headers`, and return 204.
5. Return 403 for disallowed origins on OPTIONS (optional, but recommended).

Write a `main()` that creates a test server with this middleware, makes requests with different origins, and prints the response headers.

## Tests / verification

```bash
go run ./curriculum/modules/10-auth-security/lessons/17-cors
go test ./curriculum/modules/10-auth-security/lessons/17-cors
```

The existing tests verify that allowed origins receive ACAO headers, disallowed origins do not, preflight requests return proper headers, and credentials mode works correctly.

## Review questions

1. Why does CORS only apply to browser-based requests and not to curl or server-to-server requests?
2. What triggers a preflight OPTIONS request versus a simple CORS request?
3. Why is `Access-Control-Allow-Origin: *` incompatible with `Access-Control-Allow-Credentials: true`?
4. What is the purpose of the `Vary: Origin` header in CORS responses?
5. How would an attacker exploit an API that reflects the `Origin` header without validation?

## NEXT UP

Rate limiting -- protecting APIs from abuse with token bucket, leaky bucket, and sliding window algorithms implemented as HTTP middleware in Go.
