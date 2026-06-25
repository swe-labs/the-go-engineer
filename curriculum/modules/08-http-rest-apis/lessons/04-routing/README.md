# Routing

## Learning objective

Use `http.ServeMux` patterns for path matching, implement method-based routing, extract path parameters, and compare with third-party routers like Chi and gorilla/mux.

## Why this matters

Routing is the entry point of every HTTP API. A well-designed router maps clean URLs to handler code, enforces HTTP method semantics, and extracts path parameters without string manipulation. Go 1.22 added method-based routing and path wildcards to the standard library `ServeMux`, reducing the need for third-party routers in many cases. Understanding both the built-in and third-party options lets you choose the right tool for each project's complexity.

## Mental model

Think of the router as a dispatcher at a hotel front desk. A guest arrives with a request (method + path). The dispatcher checks the board: "GET /rooms/42 → room service, POST /checkin → front desk, DELETE /booking/7 → cancellations." Each board entry has a method and a path pattern. The dispatcher finds the best match and directs the guest to the right department. If no match is found, the guest gets a 404.

## Core idea

Go 1.22+ `http.ServeMux` supports:

- **Exact path**: `"/users"` matches only `GET /users`.
- **Prefix path**: `"/users/"` matches any path starting with `/users/`.
- **Method+path**: `"GET /users"` matches only GET requests to `/users`. Without a method prefix, any method matches.
- **Path parameters**: `"GET /users/{id}"` captures the path segment as `r.PathValue("id")`.
- **Wildcard**: `"GET /users/{id}/posts/{postId}"` for multiple parameters.

Third-party routers (Chi, gorilla/mux) add:
- Route groups and subrouters.
- Middleware per route group.
- Regex-based path constraints.
- Named routes for reverse URL generation.

## Under the hood

Go 1.22's `ServeMux` compiles patterns into a radix tree (a trie) for efficient matching. Each node in the tree corresponds to a path segment. When a request arrives, the mux walks the tree segment by segment, matching literals exactly and wildcards to any single segment. Method-restricted patterns check the request method before delegating. If no pattern matches, `ServeMux` redirects `"/users"` to `"/users/"` (with a 301) when a subtree pattern exists, or returns 404.

## How Go uses it

With Go 1.22+:

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /users/{id}", getUser)
mux.HandleFunc("POST /users", createUser)
mux.HandleFunc("DELETE /users/{id}", deleteUser)
http.ListenAndServe(":8080", mux)
```

The `{id}` parameter is extracted with `r.PathValue("id")`.

Third-party example with Chi:

```go
r := chi.NewRouter()
r.Get("/users/{id}", getUser)
r.Post("/users", createUser)
r.Delete("/users/{id}", deleteUser)
```

## Go example

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// In-memory store for demonstration.
var users = map[string]User{
	"1": {ID: "1", Name: "Alice"},
	"2": {ID: "2", Name: "Bob"},
}

func listUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	u, ok := users[id]
	if !ok {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}
	users[u.ID] = u
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", listUsers)
	mux.HandleFunc("GET /users/{id}", getUser)
	mux.HandleFunc("POST /users", createUser)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## Step-by-step execution

1. Client sends `GET /users/1`.
2. `ServeMux` walks its patterns. `GET /users/{id}` matches method `GET` and path `/users/1`.
3. `{id}` is bound to `"1"`.
4. `getUser` is called. `r.PathValue("id")` returns `"1"`.
5. `users["1"]` exists → encoded as JSON and returned.
6. Client sends `POST /users` with body `{"id":"3","name":"Charlie"}`.
7. `POST /users` matches. `createUser` decodes the body, stores user, returns 201.
8. Client sends `DELETE /users/1`. No `DELETE /users/{id}` pattern registered → 405 Method Not Allowed (or 404 if no overlap exists).

## Common mistakes

- Using `r.PathValue` on Go versions before 1.22. It will compile but always return empty string. Check your `go.mod` version.
- Not handling the trailing slash: `"/users"` does NOT match `"/users/"`. They are different patterns.
- Forgetting that method-prefixed patterns only match that method. A `"GET /users"` pattern won't match a POST request.
- Using third-party routers when the built-in is sufficient. For simple APIs, standard library routing avoids a dependency and its API churn.
- Missing fallback: when no pattern matches, `ServeMux` returns 404. For APIs, always return a consistent JSON 404 body.

## Debugging walkthrough

A developer registers `GET /users/{id}` but requests to `/users/abc` return 404:

```go
mux.HandleFunc("GET /users/{id}", getUser)
```

**Symptom**: `curl http://localhost:8080/users/abc` → 404.

**Investigation**: Check if any other pattern is catching it. The developer also registered:

```go
mux.HandleFunc("/users", listUsers)
```

Without a method prefix. Since `GET /users` is a prefix of `GET /users/{id}` — but wait, `/users` is exact, not `/users/`. Actually, the issue is that `/users` matches exactly `/users`, so `/users/abc` should go to `{id}`. Let me check the mux tree... Actually, the issue could be if there's a trailing-slash redirect.

**Root cause**: There's a pattern `"/"` catch-all registered elsewhere that matches before the specific pattern.

**Fix**: Ensure catch-all patterns are not registered, or check pattern priority. In `ServeMux`, the most specific pattern wins (longest path). Remove the `"/"` pattern or make it more specific.

## Production notes

- Prefer method-prefixed patterns in Go 1.22+ for self-documenting routing. `"GET /users/{id}"` clearly states intent.
- Keep paths flat (max 2-3 levels). Deeply nested paths like `/api/v1/organizations/{org}/projects/{proj}/tasks` are harder to route, cache, and reason about.
- Use `http.NewServeMux()` instead of `http.DefaultServeMux` to create isolated routers, especially in tests.
- When using third-party routers, treat middleware as a first-class concern. Chi's middleware stack is composable per-route-group.

## Performance implications

- Go 1.22 `ServeMux` uses a radix tree with O(path-length) lookup, making routing negligible (<1µs per request).
- Third-party routers use similar trie structures. Chi and gorilla/mux add ~200-500ns overhead per lookup due to middleware dispatch.
- Route registration with wildcards is O(1). There is no performance benefit to prioritizing specific routes.
- Middleware chaining in third-party routers creates handler wrappers at registration time, not at request time, so there is no per-request allocation for middleware structure.

## Practice task

Build a task management API with three endpoints using Go 1.22+ `ServeMux`:
- `GET /tasks` — returns list of tasks as JSON.
- `GET /tasks/{id}` — returns a single task by ID, or 404.
- `POST /tasks` — creates a task from JSON body, returns 201.
- Use an in-memory `map[string]Task` as storage.
- Register all patterns with method prefixes.

## Tests / verification

```bash
go run ./curriculum/modules/08-http-rest-apis/lessons/04-routing
go test ./curriculum/modules/08-http-rest-apis/lessons/04-routing
```

## Review questions

1. What is the difference between the patterns `"GET /users"` and `"GET /users/"` in Go 1.22 ServeMux?
2. How do you extract a path parameter like `{id}` from a request?
3. When would you choose Chi or gorilla/mux over the standard library ServeMux?
4. What happens if you register both `"GET /users/{id}"` and `"GET /users/{name}"`?
5. Why does `DELETE /users/1` return 405 Method Not Allowed if only `GET /users/{id}` is registered?

## NEXT UP

Request parsing — extracting data from query strings, path parameters, headers, and JSON/form bodies.
