# REST design principles

## Learning objective

Design resource-oriented RESTful APIs with stateless semantics, idempotent operations, consistent URL conventions, and apply the REST Maturity Model to evaluate API quality.

## Why this matters

REST is the dominant API architecture on the web. Every public API you consume (GitHub, Stripe, Twilio) follows REST principles. A well-designed REST API is intuitive to use, easy to evolve, and self-documenting. A poorly designed one causes confusion, versioning nightmares, and client bugs. Understanding REST principles separates a professional API designer from someone who just mounts handlers on URLs.

## Mental model

REST treats your API as a _system of resources_. A resource is anything you can name -- a user, an order, a product. Each resource has a URL (its address) and supports standard operations via HTTP methods (its verbs).

Think of a library:
- Books are resources at `/books`.
- A specific book is at `/books/42`.
- You GET a book (read it), POST to create a new book (add it to the shelf), PUT to replace a book, PATCH to update a page, DELETE to remove it.
- The system is stateless -- the library doesn't remember what you read last time.

REST is not a protocol; it's a set of architectural constraints that, when followed, produce a system with desirable properties: scalability, cacheability, evolvability.

## Core idea

The five REST constraints (from Fielding's dissertation):

1. **Client-Server**: Separation of concerns. The client (frontend) and server (backend) evolve independently.
2. **Stateless**: Each request contains all the information needed to process it. No server-side session state.
3. **Cacheable**: Responses must explicitly indicate cacheability (via `Cache-Control`, `ETag`).
4. **Uniform Interface**: Resources are identified in requests; representations are transferred via standard methods (GET, POST, PUT, DELETE, PATCH); self-descriptive messages; HATEOAS.
5. **Layered System**: Intermediaries (proxies, load balancers, gateways) can exist between client and server.

The **REST Maturity Model** (Richardson) grades APIs:

| Level | Name | Description |
|---|---|---|
| 0 | The Swamp of POX | One URL, one method (POST), all actions in the body |
| 1 | Resources | Multiple URLs for resources |
| 2 | HTTP Verbs | Uses GET, POST, PUT, DELETE appropriately |
| 3 | HATEOAS | Responses include links to related actions |

## Under the hood

HTTP methods map to CRUD operations:

| Method | Action | Idempotent | Safe | Example |
|---|---|---|---|---|
| GET | Read | Yes | Yes | `GET /books/42` |
| POST | Create | No | No | `POST /books` |
| PUT | Replace | Yes | No | `PUT /books/42` |
| DELETE | Delete | Yes | No | `DELETE /books/42` |
| PATCH | Partial update | No | No | `PATCH /books/42` |

- **Idempotent**: Multiple identical requests produce the same result as one. DELETE is idempotent because deleting the same resource twice returns the same outcome (the resource no longer exists).
- **Safe**: The operation has no side effects. GET and HEAD are safe.

Resource naming conventions:
- Use nouns, not verbs: `/orders` not `/createOrder`.
- Plural for collections: `/users`, `/products`.
- Nest sub-resources: `/users/42/orders`.
- Use query params for filtering/pagination, not path segments: `/orders?status=pending&page=2`.

## How Go uses it

Go's `net/http` allows building RESTful APIs directly. The standard library's `http.ServeMux` supports method-based routing (from Go 1.22+): `mux.HandleFunc("GET /api/books/{id}", handler)`. Before that, explicit method checking was the norm (as shown in the example).

Frameworks like `chi`, `gorilla/mux`, and `gin` provide richer routing for REST APIs with path parameters, method patterns, and middleware. The principles remain the same regardless of framework.

## Go example

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type store struct {
	mu    sync.RWMutex
	books map[int]Book
	next  int
}

func (s *store) create(b Book) Book {
	s.mu.Lock()
	defer s.mu.Unlock()
	b.ID = s.next
	s.books[b.ID] = b
	s.next++
	return b
}

func (s *store) getAll() []Book {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Book, 0, len(s.books))
	for _, b := range s.books {
		result = append(result, b)
	}
	return result
}

// GET /api/books and POST /api/books
func booksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(bookStore.getAll())
	case http.MethodPost:
		var b Book
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
			return
		}
		b = bookStore.create(b)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(b)
	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

var bookStore = &store{books: make(map[int]Book), next: 1}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/books", booksHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## Step-by-step execution

A client creates a book and retrieves it:

1. Client sends `POST /api/books` with body `{"title":"Go in Action","author":"Kennedy"}`.
2. Server receives the request. `booksHandler` checks `r.Method` -- it's POST.
3. JSON body is decoded into a `Book` struct. `Title` and `Author` are populated.
4. `store.create` assigns ID 1 and stores the book.
5. Response: `201 Created` with body `{"id":1,"title":"Go in Action","author":"Kennedy"}`.
6. Client sends `GET /api/books`.
7. `booksHandler` sees GET, calls `store.getAll()`, returns `200 OK` with the list.
8. Client sends `DELETE /api/books/1` (if implemented) to remove the resource.
9. DELETE is idempotent -- the second DELETE returns the same `404 Not Found`.

## Common mistakes

- Using verbs in URLs: `/createBook` instead of `POST /books`.
- Not returning proper status codes. Use 201 for creation, 204 for deletion, 404 for not found.
- Ignoring idempotency. PUT and DELETE must be idempotent; POST is not.
- Mixing plural and singular: `/user` vs `/users`. Pick one convention (plural is standard).
- Deep nesting: `/users/42/orders/5/items/7`. Beyond 2-3 levels, use query params or sub-resources at the top level.
- Not handling content negotiation. Return JSON (set `Content-Type: application/json`) and check `Accept` header if supporting multiple formats.

## Debugging walkthrough

An API that returns 200 for missing resources:

```go
func getBook(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/books/"))
    b, ok := bookStore.getByID(id)
    if !ok {
        // BUG: missing return after http.Error
        http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
    }
    json.NewEncoder(w).Encode(b)
}
```

**Symptom**: Requesting a non-existent book (`/books/999`) returns 200 with an empty JSON body `{}`.

**Investigation**: The `http.Error` call writes the error and sets status code, but execution continues because there's no `return`. The code falls through to `json.NewEncoder(w).Encode(b)`, which writes `null` and appends to the already-written 404 body. But since `http.Error` already wrote headers, the response code is 404, and the body is `{"error":"not found"}null`. Wait -- `http.Error` writes the body, so the encoder writes after it, double-writing. The client sees a 404 with `{"error":"not found"}null`.

**Fix**: Add `return` after `http.Error`:

```go
if !ok {
    http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
    return
}
```

## Production notes

- Always return consistent error shapes: `{"error":"message","code":"RESOURCE_NOT_FOUND"}`.
- Document all endpoints with status codes and request/response shapes.
- Use `Accept` header versioning for major breaking changes.
- Rate limiting and authentication should be orthogonal to resource design.
- Consider HATEOAS for Level 3 REST -- include `self`, `next`, `prev` links in responses.
- Use UUIDs for resource IDs instead of sequential integers to avoid information leakage.

## Performance implications

- RESTful design has negligible performance overhead. URL routing and method dispatch are O(1) or O(log n) in any Go router.
- JSON encoding/decoding is the bottleneck, not the REST architecture.
- Statelessness enables horizontal scaling. Any server can handle any request.
- Cacheable responses (GET) can be served from CDN or reverse proxy, dramatically reducing origin load.

## Practice task

1. Design a RESTful API for a `Product` resource with fields: `id`, `name`, `price`, `category`.
2. Implement handlers for `GET /products`, `POST /products`, `GET /products/{id}`, `PUT /products/{id}`, `DELETE /products/{id}`.
3. Use proper HTTP status codes (200, 201, 204, 400, 404, 405).
4. Make POST reject missing `name` with 400.
5. Make PUT return 404 if the product doesn't exist.
6. Write table-driven tests for all endpoints.

## Tests / verification

```bash
go test ./curriculum/modules/08-http-rest-apis/lessons/17-rest-design-principles -v
```

Tests cover CRUD operations, input validation, idempotent DELETE, not-found handling, and method validation.

## Review questions

1. What is the difference between a safe HTTP method and an idempotent one? Give examples of each.
2. Why is POST not idempotent? Can you make POST idempotent with idempotency keys?
3. What are the five REST constraints? Which one is often the most challenging to implement?
4. What does Richardson's REST Maturity Model describe? What distinguishes Level 2 from Level 3?
5. Why should resource URLs use nouns instead of verbs?

## NEXT UP

API versioning -- strategies for evolving your API without breaking clients, including URL path versioning and header-based versioning.
