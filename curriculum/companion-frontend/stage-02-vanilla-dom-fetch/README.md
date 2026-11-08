# Stage 2: Browser JavaScript, DOM & Fetch API

> **Companion Track Stage:** 2 of 4  
> **Core Module Alignment:** Module 08 (*HTTP & REST APIs*) and Module 09 (*SQL & Persistence*)  
> **Estimated Time:** 12–16 Hours  
> **Prerequisites:** Stage 1, HTTP status codes, JSON serialization.

---

## 1. Overview

This starter builds an interactive browser client that communicates asynchronously with a Go REST API backend. It demonstrates modern vanilla JavaScript without frontend frameworks, bundlers, or npm dependencies.

```
┌────────────────────────────────────────────────────────┐
│                   Browser Client                       │
│  - Debounced Input Events (300ms delay)                │
│  - Event Delegation (Parent listener for delete)       │
│  - Fetch API (async/await, HTTP status management)     │
│  - Accessible Field-Level Validation Error Rendering   │
└───────────────────────────▲────────────────────────────┘
                            │
               JSON Payloads (REST API: /api/books)
                            │
┌───────────────────────────▼────────────────────────────┐
│                  Go App Server                         │
│  - Thread-safe repository (sync.RWMutex)               │
│  - Strict JSON validation & structured error responses │
│  - Embedded static SPA delivery via //go:embed         │
└────────────────────────────────────────────────────────┘
```

---

## 2. Core Frontend Patterns Demonstrated

### 1. Debounced Search Input
Without debouncing, every keystroke fires an HTTP request to the backend. With a 300ms debounce timer, requests are only dispatched after the learner pauses typing:

```javascript
function debounce(func, delayMs) {
    let timeoutId;
    return function (...args) {
        clearTimeout(timeoutId);
        timeoutId = setTimeout(() => func.apply(this, args), delayMs);
    };
}
```

### 2. Event Delegation
Because book cards are created and deleted dynamically, attaching individual event listeners to each button leads to memory leaks and broken listeners upon re-rendering. Instead, an event listener is placed on the parent container:

```javascript
booksContainer.addEventListener("click", e => {
    const deleteBtn = e.target.closest('[data-action="delete"]');
    if (!deleteBtn) return;
    const id = parseInt(deleteBtn.dataset.id, 10);
    deleteBook(id);
});
```

### 3. Asynchronous Error Deserialization
When Go handlers return `400 Bad Request` with structured field errors:
```json
{
  "error": "Validation failed",
  "details": [
    {"field": "isbn", "message": "ISBN must be 10 or 13 digits"}
  ]
}
```
The client maps the error messages directly onto the associated input elements with red highlight rings and screen-reader accessible error spans.

---

## 3. Directory Structure

```
stage-02-vanilla-dom-fetch/
├── static/
│   ├── app.js            # Vanilla JS: state, debounce, fetch, event delegation
│   ├── index.html        # Semantic HTML client interface
│   └── style.css         # Modern dark-mode styling, loading spinners & cards
├── main.go               # Go REST API backend & static file server
├── main_test.go          # API validation & static integration tests
└── README.md             # This guide
```

---

## 4. Running the Starter

### Start the Server
```bash
go run main.go
```
The server will start listening at `http://localhost:8080`.

Open your browser to:
`http://localhost:8080/`

Try:
1. Typing in the search bar — notice the loading spinner and debounced queries in your terminal log.
2. Submitting an empty book form to trigger backend validation errors.
3. Adding a new book and deleting an existing one asynchronously.

### Run the Test Suite
```bash
go test -v ./...
```

---

## 5. Learner Extension Challenges

1. **Inline Edit Modal**: Add an "Edit" button to book cards that opens a modal dialog and dispatches `PUT /api/books/{id}`.
2. **Sort by Price/Title**: Add client-side sort dropdowns to reorder the active collection in memory.
3. **Optimistic UI Updates**: Update the DOM immediately on deletion, with a rollback mechanism if the network call fails.

---

## Stage mission

Build an interactive browser client that talks asynchronously to a Go REST API using modern vanilla JavaScript — no frameworks, no bundler, no npm. The stage proves you can manage DOM state, debounce input, and render server-returned validation errors by hand.

## Backend alignment

Aligns with [Module 08: HTTP and REST APIs](../../modules/08-http-rest-apis/README.md) and [Module 09: SQL and PostgreSQL Persistence](../../modules/09-sql-postgres-persistence/README.md). The Go server mirrors those modules' REST semantics and persistence patterns; the browser is a second client of the same API shape.

## Verification and test suite

```bash
go test -v ./...
```

Then `go run main.go` and open `http://localhost:8080/`. Confirm debounced search, field-level validation error rendering, and asynchronous add/delete without a reload.

## Next stage

Continue to [Stage 3: Type-Safe Server Components with templ & Tailwind](../stage-03-templ-tailwind/README.md).

