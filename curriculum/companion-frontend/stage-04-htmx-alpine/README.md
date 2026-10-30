# Stage 4: Hypermedia Interactivity with HTMX & Alpine.js

> **Companion Track Stage:** 4 of 4  
> **Core Module Alignment:** Module 11 (*Concurrency & Lifecycle*) and Module 12 (*Observability & Diagnostics*)  
> **Estimated Time:** 16–22 Hours  
> **Prerequisites:** Stages 1–3, Go channels and goroutines, HTTP streaming (`text/event-stream`).

---

## 1. Overview

This starter builds a **real-time worker pool and system metrics dashboard** demonstrating the power of **hypermedia-driven architectures**. Instead of a heavyweight SPA framework (React/Vue/Angular), the UI achieves rich interactivity, partial swaps, and real-time streaming using **HTMX** and **Alpine.js** connected directly to Go goroutine workers.

```
┌────────────────────────────────────────────────────────┐
│                   Browser Client                       │
│  - HTMX (Partial DOM swaps, hx-post, hx-target)        │
│  - HTMX Out-of-Band Swaps (hx-swap-oob="true")         │
│  - Server-Sent Events (SSE: live logs & metric stream) │
│  - Alpine.js (Modal dialogs & collapsible log terminal)│
└───────────────────────────▲────────────────────────────┘
                            │
               HTML Fragments & SSE (text/event-stream)
                            │
┌───────────────────────────▼────────────────────────────┐
│                  Go Dashboard Server                   │
│  - Background Worker Pool (Go channels & goroutines)   │
│  - Context cancellation for aborting tasks             │
│  - SSE Event Broker broadcasting to all active clients │
│  - Embedded vendor scripts (Zero internet dependency)  │
└────────────────────────────────────────────────────────┘
```

---

## 2. Core Hypermedia Patterns Demonstrated

### 1. Granular HTML Partial Swaps (`hx-post` & `hx-target`)
When a user dispatches a background task, the browser issues an AJAX request:
```html
<button 
    type="button" 
    class="btn btn-outline" 
    hx-post="/tasks?type=Data+Ingestion" 
    hx-target="#tasks-table-body" 
    hx-swap="afterbegin"
>
    + Ingestion Job
</button>
```
The Go server executes the business logic, enqueues the task, and responds with a raw HTML table row snippet:
```html
<tr id="task-5" class="task-row status-queued">
    <td>#5</td>
    <td>Data Ingestion</td>
    <td><span class="status-badge badge-queued">QUEUED</span></td>
    ...
</tr>
```
HTMX immediately injects this snippet at the top of `#tasks-table-body` without reloading the page or requiring client-side JSON serialization.

### 2. Out-of-Band (OOB) Updates
A single HTTP response can update multiple disconnected regions of the DOM simultaneously using `hx-swap-oob="true"`. In `main.go`, when `/tasks` is posted, the server returns both the task row **and** the updated metrics card:
```html
<!-- Primary Target: Appended to #tasks-table-body -->
<tr id="task-5">...</tr>

<!-- Out-of-Band Target: Swapped into #metrics-panel anywhere on the page -->
<div id="metrics-panel" class="metrics-grid" hx-swap-oob="true">
    ... updated worker counts and queue depth ...
</div>
```

### 3. Server-Sent Events (SSE) Stream
Real-time log streaming and continuous metric broadcasts are powered by Go's standard library `http.Flusher`:
```html
<div class="dashboard-layout" hx-ext="sse" sse-connect="/events">
    <!-- Listens for metric broadcasts -->
    <section sse-swap="metric">...</section>

    <!-- Listens for live log entries -->
    <div id="log-stream" sse-swap="log">...</div>
</div>
```

### 4. Micro-Client State with Alpine.js
Local UI states (modals, dropdowns, drawer accordions) stay on the client using Alpine.js without state management libraries:
```html
<!-- Modal visibility toggled entirely in the DOM -->
<div x-data="{ modalOpen: false }">
    <button @click="modalOpen = true">Dispatch Custom Task</button>
    <div x-show="modalOpen" @keydown.escape.window="modalOpen = false">
        ... modal form ...
    </div>
</div>
```

### 5. The GOTH Architecture: Type-Safe `templ` Component Streaming
In production architectures, hand-crafted HTML strings or unvalidated `html/template` strings risk runtime syntax errors and escaping issues. This starter includes a complete **GOTH (Go + Orchestration/Tailwind + Templ + HTMX)** pipeline located in [`components/`](./components):

```go
// Render and stream a type-safe templ.Component over SSE
func (b *SSEBroker) BroadcastComponent(ctx context.Context, event string, c templ.Component) error {
    var buf bytes.Buffer
    if err := c.Render(ctx, &buf); err != nil {
        return err
    }
    b.Broadcast(event, buf.String())
    return nil
}
```

Components are compile-time checked pure Go structs generated from `.templ` definitions:
```templ
// components/task_row.templ
package components

templ TaskRow(t *Task) {
    <tr id={ fmt.Sprintf("task-%d", t.ID) } class={ fmt.Sprintf("task-row status-%s", t.StatusLower) }>
        <td class="font-mono text-muted text-xs">{ fmt.Sprintf("#%d", t.ID) }</td>
        <td class="font-semibold">{ t.Type }</td>
        <td>
            <span class={ fmt.Sprintf("status-badge badge-%s", t.StatusLower) }>
                { t.Status }
            </span>
        </td>
        ...
        <td class="text-right">
            if t.Status == "QUEUED" || t.Status == "RUNNING" {
                <button
                    type="button"
                    class="btn-action-cancel"
                    hx-post={ fmt.Sprintf("/tasks/%d/cancel", t.ID) }
                    hx-target={ fmt.Sprintf("#task-%d", t.ID) }
                    hx-swap="outerHTML"
                >
                    Cancel
                </button>
            }
        </td>
    </tr>
}
```

When initialized without runtime templates (`tmpl = nil`), the server runs in **pure compile-time type-safe mode**, rendering directly to `http.ResponseWriter` with zero reflection and zero runtime template parsing.

---

## 3. Directory Structure

```
stage-04-htmx-alpine/
├── components/           # Type-Safe templ component library (GOTH stack)
│   ├── components_test.go # Comprehensive unit tests for templ components
│   ├── dashboard.templ   # Root HTML dashboard layout in templ
│   ├── dashboard_templ.go # Generated Go source code
│   ├── metrics_card.templ # Reactive metrics card with OOB swap support
│   ├── metrics_card_templ.go
│   ├── task_row.templ    # Granular partial table row & log entry
│   ├── task_row_templ.go
│   └── types.go          # Domain types (Task, MetricsSnapshot)
├── static/
│   ├── alpine.min.js     # Vendor Alpine.js v3.14 (offline)
│   ├── htmx.min.js       # Vendor HTMX v1.9.12 (offline)
│   ├── sse.js            # Vendor HTMX SSE extension (offline)
│   └── style.css         # Modern dark dashboard & status styles
├── templates/            # Standard html/template files (for comparison)
│   ├── dashboard.html    # Full dashboard view with SSE connection
│   ├── metrics_card.html # Metrics partial (supports OOB swaps)
│   └── task_row.html     # Individual task table row fragment
├── datastar/             # Modern Alternative: Pure SSE & Datastar signals variant
│   ├── components/       # Datastar-native templ component library (GOTD stack)
│   ├── static/           # Vendored datastar.js (~13 KB single script)
│   ├── templates/        # Datastar-native templates
│   ├── main.go           # Server streaming datastar-patch-elements
│   ├── main_test.go      # Datastar test suite (html/template + templ modes)
│   └── README.md         # Datastar guide & HTMX comparison
├── main.go               # Worker pool, SSE broker, and handlers (dual-mode)
├── main_test.go          # SSE, OOB swaps, concurrency, and templ mode tests
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
1. Clicking **+ Ingestion Job** or **+ Reconciliation Job**: Notice the new task instantly prepends to the table, and the **Queue Depth** counter updates simultaneously via out-of-band swap.
2. Clicking **Custom Task Dispatch**: Notice the Alpine.js modal open smoothly without page jumps. Choose a 10-second duration and submit.
3. Clicking **Cancel** on a running task: Watch the status transition to `CANCELLED` and the worker release back to the idle pool.
4. Watching the **Live Worker Pool Execution Stream**: Observe real-time log events arriving over the SSE `/events` connection.

### Run the Test Suite
```bash
go test -v ./...
```

---

## 5. Modern Alternative: Datastar SSE-First Variant

Looking to eliminate the seam between HTMX swaps and Alpine client state?

Check out the complete comparative implementation in [`datastar/`](./datastar/README.md):
- **Single Script:** Replaces both HTMX (14 KB) and Alpine (15 KB) with a single ~13 KB script (`datastar.js`).
- **SSE Native:** Uses `datastar-patch-elements` and `datastar-patch-signals` to stream DOM mutations directly to any selector without out-of-band swap hacks.
- **Run side-by-side:** Start the Datastar server with `go run ./datastar/main.go` listening at `http://localhost:8081` and compare both architectures head-to-head.

---

## 6. Learner Extension Challenges

1. **Worker Pool Auto-Scaler**: Add an endpoint `POST /workers/scale?count=N` that dynamically increases or decreases the number of active worker goroutines.
2. **Task Priority Channels**: Refactor `taskQueue` into dual channels (`highPriorityQueue` and `lowPriorityQueue`) using a Go `select` statement.
3. **Filtering Table Rows by Status**: Add an Alpine.js or HTMX filter bar (`All`, `Running`, `Completed`, `Cancelled`) to toggle row visibility.
4. **Datastar Templ Refactor**: Refactor the Datastar variant to use [`a-h/templ`](https://templ.guide) and `github.com/starfederation/datastar/sdk/go`.

