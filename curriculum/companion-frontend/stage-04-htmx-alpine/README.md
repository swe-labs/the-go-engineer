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

---

## 3. Directory Structure

```
stage-04-htmx-alpine/
├── static/
│   ├── alpine.min.js     # Vendor Alpine.js v3.14 (offline)
│   ├── htmx.min.js       # Vendor HTMX v1.9.12 (offline)
│   ├── sse.js            # Vendor HTMX SSE extension (offline)
│   └── style.css         # Modern dark dashboard & status styles
├── templates/
│   ├── dashboard.html    # Full dashboard view with SSE connection
│   ├── metrics_card.html # Metrics partial (supports OOB swaps)
│   └── task_row.html     # Individual task table row fragment
├── main.go               # Worker pool, SSE broker, and handlers
├── main_test.go          # SSE, OOB swaps, and concurrency tests
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

## 5. Learner Extension Challenges

1. **Worker Pool Auto-Scaler**: Add an endpoint `POST /workers/scale?count=N` that dynamically increases or decreases the number of active worker goroutines.
2. **Task Priority Channels**: Refactor `taskQueue` into dual channels (`highPriorityQueue` and `lowPriorityQueue`) using a Go `select` statement.
3. **Filtering Table Rows by Status**: Add an Alpine.js or HTMX filter bar (`All`, `Running`, `Completed`, `Cancelled`) to toggle row visibility.
