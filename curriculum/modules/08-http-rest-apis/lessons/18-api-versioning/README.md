# API versioning

## Learning objective

Implement API versioning using URL path prefixes (`/v1/`, `/v2/`) and `Accept` header content negotiation, manage backward compatibility, and design deprecation strategies.

## Why this matters

APIs evolve. Fields are added, renamed, removed; endpoints change behavior. Without versioning, every change breaks existing clients. Versioning allows you to release new capabilities while old clients continue using the previous version. Every public API (Stripe, GitHub, Twilio) uses versioning. It is not optional for production services.

## Mental model

Think of API versions like app store releases. Users on the old version continue using it; new users get the latest. You maintain both until the old version is deprecated and eventually sunset. The API version is the contract between client and server. Changing the contract without versioning is like force-updating everyone's app without warning.

Versioning strategies:
- **URL path versioning**: `/v1/users`, `/v2/users`. Simple, visible, cacheable. The version is part of the URL.
- **Header versioning**: `Accept: application/vnd.api.v2+json`. Cleaner URLs but harder to test manually.
- **Query parameter versioning**: `/users?version=2`. Easy to implement but pollutes query strings.
- **Subdomain versioning**: `v2.api.example.com`. Infrastructure-level separation.

## Core idea

The core versioning concern is backward compatibility. A change is backward-compatible if:
- Existing fields are not removed or renamed (additions are fine).
- Existing behavior for the same inputs is unchanged.
- Existing error codes and formats are unchanged.

When you must make a breaking change, create a new version. The old version continues to work for existing clients.

**Deprecation policy**: Announce deprecation, set a sunset date, and return a `Sunset` or `Deprecation` HTTP header pointing to migration docs.

```http
Deprecation: true
Sunset: Sat, 31 Dec 2025 23:59:59 GMT
```

## Under the hood

URL path versioning is straightforward: route to different handlers based on the prefix.

Header versioning (content negotiation) uses the `Accept` header to indicate the desired version. The server parses the header and selects the appropriate handler. The response includes a `Content-Type` with the version.

```go
func versionHandler(w http.ResponseWriter, r *http.Request) {
    accept := r.Header.Get("Accept")
    if strings.Contains(accept, "application/vnd.api.v2") {
        handleV2(w, r)
        return
    }
    handleV1(w, r) // default
}
```

The `vnd.` prefix in the media type stands for "vendor-specific". The format is `application/vnd.{org}.{resource}+{format}`.

## How Go uses it

Go's standard library does not prescribe a versioning strategy. However, popular Go routers handle path-based versioning natively:

```go
r := chi.NewRouter()
r.Route("/v1", func(r chi.Router) {
    r.Get("/users", listUsersV1)
})
r.Route("/v2", func(r chi.Router) {
    r.Get("/users", listUsersV2)
})
```

The `content-type` package in the Go ecosystem can parse media type parameters for header-based versioning.

## Go example

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// V1 user: single Name field
type UserV1 struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// V2 user: structured name with email
type UserV2 struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

var usersV1 = []UserV1{{ID: 1, Name: "Alice"}}
var usersV2 = []UserV2{{ID: 1, FirstName: "Alice", LastName: "Smith", Email: "alice@example.com"}}

// URL path versioning: /v1/users and /v2/users
func usersV1Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-API-Version", "v1")
	json.NewEncoder(w).Encode(usersV1)
}

func usersV2Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-API-Version", "v2")
	json.NewEncoder(w).Encode(usersV2)
}

// Header-based versioning: single endpoint, Accept header drives version
type VersionHandler struct {
	handlers map[string]http.HandlerFunc
}

func (vh *VersionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	accept := r.Header.Get("Accept")
	version := "v1"

	if strings.Contains(accept, "application/vnd.api.v2") {
		version = "v2"
	}

	if h, ok := vh.handlers[version]; ok {
		h(w, r)
		return
	}
	http.Error(w, `{"error":"unsupported version"}`, http.StatusNotImplemented)
}

func main() {
	mux := http.NewServeMux()

	// URL path versioning
	mux.HandleFunc("/v1/users", usersV1Handler)
	mux.HandleFunc("/v2/users", usersV2Handler)

	// Header versioning (unversioned URL)
	vh := &VersionHandler{
		handlers: map[string]http.HandlerFunc{
			"v1": usersV1Handler,
			"v2": usersV2Handler,
		},
	}
	mux.Handle("/api/users", vh)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## Step-by-step execution

Client requests `/v2/users`:

1. Request arrives at server. `http.ServeMux` matches `/v2/users` to `usersV2Handler`.
2. Handler sets `Content-Type: application/json` and `X-API-Version: v2` headers.
3. Handler marshals `usersV2` (with `first_name`, `last_name`, `email`) and writes the response.

Client requests `/api/users` with `Accept: application/vnd.api.v2+json`:

1. `VersionHandler.ServeHTTP` is called.
2. It reads `Accept` header, detects `v2` substring, and dispatches to `usersV2Handler`.
3. Same response as above, but the URL remains `/api/users`.

Deprecation lifecycle:

1. v1 endpoint returns header: `Deprecation: true`, `Sunset: Sat, 31 Dec 2025`.
2. Clients see the deprecation header and migrate to v2.
3. After the sunset date, the v1 handler returns `410 Gone`.

## Common mistakes

- Not having a versioning plan at all. Every breaking change becomes a crisis.
- Versioning too granularly (e.g., `/v1.0.1/users`). Major versions only.
- Forgetting to include version info in responses via headers so clients can confirm which version they're using.
- Removing old versions without a deprecation period and migration guide.
- Using versioning as an excuse for poor API design. Versioning is for breaking changes, not for every minor iteration.
- Versioning internal behavior (bug fixes) that shouldn't require a version bump.

## Debugging walkthrough

A client reports that v2 returns stale data:

```go
func usersV2Handler(w http.ResponseWriter, r *http.Request) {
    // BUG: accidentally sends v1 data
    json.NewEncoder(w).Encode(usersV1)
}
```

**Symptom**: `GET /v2/users` returns `{"name":"Alice"}` instead of `{"first_name":"Alice","last_name":"Smith","email":"..."}`.

**Investigation**: Check the handler code -- it references `usersV1` instead of `usersV2`. The URL routing is correct, but the handler body has a copy-paste error.

**Fix**: Change `usersV1` to `usersV2` in `usersV2Handler`.

## Production notes

- Prefer URL path versioning for simplicity. It's the most common approach (GitHub, Stripe).
- Use header versioning when you want clean, unchanging URLs.
- Return a `Sunset` header on deprecated versions with a date at least 6 months in the future.
- Log the version used by each request for analytics and deprecation tracking.
- Keep the old version running for at least one major release cycle.
- Automate migration testing: ensure that existing clients against the old version produce the same results after each deployment.

## Performance implications

- Versioning itself has zero performance overhead -- it's just routing.
- Maintaining multiple versions adds code complexity but no runtime cost beyond the handler dispatch.
- The real cost is cognitive: developers must understand multiple versions. Keep version surface areas minimal and clearly separated.
- Header versioning requires parsing the `Accept` header, which is a string operation (nanoseconds). Negligible.

## Practice task

1. Create a v1 handler for `Product` that returns `{id, name, price}`.
2. Create a v2 handler that returns `{id, name, price_eur, currency}`.
3. Implement URL path versioning: `/v1/products`, `/v2/products`.
4. Implement header versioning: `Accept: application/vnd.api.v2+json` routes to v2.
5. Add a deprecation warning header to the v1 handler.
6. Write tests confirming both versions work and return the correct data shapes.

## Tests / verification

```bash
go test ./curriculum/modules/08-http-rest-apis/lessons/18-api-versioning -v
```

Tests verify v1 and v2 handlers, URL routing, header-based version dispatch, the version middleware, and the v1-to-v2 converter.

## Review questions

1. What are the four common API versioning strategies? What are the tradeoffs of each?
2. Why is URL path versioning simpler than header versioning? When would you choose header versioning?
3. What HTTP headers indicate a deprecated API version and its sunset date?
4. What constitutes a breaking change in an API? When do you need a new version?
5. How would you phase out v1 of an API without breaking existing clients?

## NEXT UP

Pagination and filtering -- implementing cursor and offset pagination, query parameters for filtering and sorting RESTful list endpoints.
