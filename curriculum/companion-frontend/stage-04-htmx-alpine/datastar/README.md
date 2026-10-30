# Stage 4 Variant: Real-Time Streaming with Datastar & Go

> **Companion Track Stage:** 4 (Modern Alternative Variant)
> **Core Module Alignment:** Module 11 (_Concurrency & Lifecycle_) and Module 12 (_Observability & Diagnostics_)
> **Key Technologies:** `Go 1.24+`, `Datastar v1.0.4` (~13 KB single script), `Server-Sent Events (SSE)`, Standard Library `http.Flusher`.
> **Prerequisites:** [Stage 4 Primary (HTMX &amp; Alpine.js)](../README.md), Go goroutines and buffered channels.

---

## 1. Overview & Architectural Shift

While the primary Stage 4 track teaches **HTMX & Alpine.js** (the current industry standard for hypermedia), this directory demonstrates the emerging **Datastar** paradigm — a hypermedia framework built specifically around **Server-Sent Events (SSE)** and native Go concurrency.

```
┌────────────────────────────────────────────────────────┐
│                   Browser Client                       │
│  - Single vendor script: datastar.js (~13 KB)          │
│  - Unified Signals (data-signals, data-bind, data-show)│
│  - SSE Morphing Engine (datastar-patch-elements)       │
│  - Zero client state desynchronization                 │
└───────────────────────────▲────────────────────────────┘
                            │
               Continuous SSE Stream (text/event-stream)
                            │
┌────────────────────────────────────────────────────────┐
│                  Go Dashboard Server                   │
│  - Native http.Flusher streaming                       │
│  - Pure standard library (zero external dependencies)  │
│  - Direct targeting of any DOM selector via SSE        │
│  - Worker pool goroutines broadcasting runtime metrics │
└────────────────────────────────────────────────────────┘
```

---

## 2. The Core Problems Datastar Solves

In the primary HTMX + Alpine.js architecture:

1. **The "Two Minds" Problem:** HTMX controls the DOM from the server via AJAX swaps, while Alpine controls local client state via `x-data`. When HTMX swaps elements containing Alpine components, Alpine state can desync, requiring manual lifecycle hooks (`@htmx:after-request`).
2. **Out-of-Band (OOB) Hacks:** To update multiple disconnected elements (e.g., adding a table row _and_ updating a metrics card), HTMX requires appending dummy elements marked with `hx-swap-oob="true"`.
3. **Clunky SSE Extensions:** HTMX treats SSE as an auxiliary extension (`hx-ext="sse"`), which requires awkward selector mappings and custom listeners.

**Datastar eliminates all three:**

- **Single ~13 KB Script:** Replaces both HTMX (14 KB) and Alpine (15 KB) with zero build steps.
- **Native SSE First:** Every action (`@post`, `@get`) and stream can emit `datastar-patch-elements` targeting _any_ CSS selector without OOB workarounds.
- **Unified Reactive Signals:** Client UI state (`$modalOpen`, `$terminalOpen`) is declared via `data-signals` and can be inspected or mutated both on the client and directly by server SSE events (`datastar-patch-signals`).

---

## 3. Side-by-Side Rosetta Stone: HTMX + Alpine vs. Datastar

| Interaction Pattern        | HTMX + Alpine.js (Primary Stage 4)                                                                                                 | Datastar (This Variant)                                                                                 |
| :------------------------- | :--------------------------------------------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------ |
| **Vendor Scripts**         | `<script src="/static/htmx.min.js"></script><script src="/static/sse.js"></script>``<script src="/static/alpine.min.js"></script>` | `<script type="module" src="/static/datastar.js"></script>`                                             |
| **State Declaration**      | `<body x-data="{ modalOpen: false }">`                                                                                             | `<body data-signals="{ modalOpen: false }">`                                                            |
| **Button Action**          | `<button hx-post="/tasks?type=Ingest" hx-target="#tasks" hx-swap="afterbegin">`                                                    | `<button data-on-click="@post('/tasks?type=Ingest')">`                                                  |
| **Two-Way Input Binding**  | `<input x-model="duration">`                                                                                                       | `<input data-bind="customDuration">`                                                                    |
| **Element Visibility**     | `<div x-show="modalOpen">`                                                                                                         | `<div data-show="$modalOpen">`                                                                          |
| **Dynamic Button Text**    | `<button x-text="open ? 'Close' : 'Open'">`                                                                                        | `<button data-text="$open ? 'Close' : 'Open'">`                                                         |
| **Server Closing a Modal** | Client-side hook:`@htmx:after-request="modalOpen = false"`                                                                         | Server-driven SSE patch:`event: datastar-patch-signalsdata: signals {"modalOpen": false}`               |
| **Updating Metrics Card**  | `<div id="metrics" hx-swap-oob="true">...</div>`                                                                                   | Server-driven SSE patch:`event: datastar-patch-elementsdata: selector #metrics-panel``data: mode outer` |

---

## 4. The Datastar SSE Wire Protocol (Zero-Magic)

Under the hood, Datastar operates over standard `text/event-stream`. When Go writes to `http.ResponseWriter` via `http.Flusher`, it formats events like this:

### 1. Patching DOM Elements (`datastar-patch-elements`)

```http
event: datastar-patch-elements
data: selector #tasks-table-body
data: mode prepend
data: elements <tr id="task-5"><td>#5</td><td>Data Ingestion</td></tr>
```

_Modes supported:_ `outer` (default morph), `inner`, `prepend`, `append`, `before`, `after`, `replace`.

### 2. Patching Client Reactive Signals (`datastar-patch-signals`)

```http
event: datastar-patch-signals
data: signals {"modalOpen": false}
```

Datastar automatically deserializes this JSON and updates client signals, closing modal dialogs or resetting form fields without writing client JavaScript.

---

## 5. The GOTD Pipeline: Go + Templ + Datastar

This starter comes equipped with a first-class **GOTD (Go + Orchestration/Tailwind + Templ + Datastar)** pipeline in [`components/`](./components). 

Instead of manual string formatting, components are defined in type-safe `.templ` files and streamed directly to client targets:

```go
// PatchComponent renders a compile-time checked templ.Component directly to an SSE patch
func (d *DatastarSSE) PatchComponent(ctx context.Context, c templ.Component, selector string, mode string) error {
    var buf bytes.Buffer
    if err := c.Render(ctx, &buf); err != nil {
        return err
    }
    return d.PatchElements(buf.String(), selector, mode)
}
```

### Type-Safe Component Example (`task_row.templ`):
```templ
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
        <td class="font-mono text-xs">{ t.DurationStr }</td>
        <td class="text-right">
            if t.Status == "QUEUED" || t.Status == "RUNNING" {
                <button
                    type="button"
                    class="btn-action-cancel"
                    data-on-click={ fmt.Sprintf("@post('/tasks/%d/cancel')", t.ID) }
                    title="Cancel execution"
                >
                    Cancel
                </button>
            }
        </td>
    </tr>
}
```

When initialized with `tmpl = nil`, the server operates in **pure compile-time type-safe mode** with zero template parsing overhead.

---

## 6. Directory Structure

```
stage-04-htmx-alpine/datastar/
├── components/           # Datastar-native templ component library (GOTD stack)
│   ├── components_test.go # Comprehensive unit tests for templ components
│   ├── dashboard.templ   # Root HTML dashboard layout in templ with Datastar signals
│   ├── dashboard_templ.go # Generated Go source code
│   ├── metrics_card.templ # Reactive metrics card
│   ├── metrics_card_templ.go
│   ├── task_row.templ    # Granular partial table row with data-on-click="@post(...)"
│   ├── task_row_templ.go
│   └── types.go          # Domain types (Task, MetricsSnapshot)
├── static/
│   ├── datastar.js       # Vendored datastar.js v1.0.4 (~13 KB single script)
│   └── style.css         # Modern dark dashboard & status styles
├── templates/            # Standard html/template files (for comparison)
│   ├── dashboard.html    # Full dashboard view with Datastar attributes
│   ├── metrics_card.html # Metrics partial
│   └── task_row.html     # Task row fragment
├── main.go               # Server streaming datastar-patch-elements (dual-mode)
├── main_test.go          # Datastar test suite (html/template + templ modes)
└── README.md             # This guide
```

---

## 7. Official Go SDK Harmony

While this starter implements the SSE protocol with pure Go standard library code (`DatastarSSE`), the Datastar project provides an official first-party Go SDK:

```bash
go get github.com/starfederation/datastar/sdk/go
```

With the official SDK and [`a-h/templ`](https://templ.guide), handlers become exceptionally concise:

```go
import datastar "github.com/starfederation/datastar/sdk/go"

func handleTask(w http.ResponseWriter, r *http.Request) {
    sse := datastar.NewSSE(w, r)

    // Render and stream a compile-time checked templ component
    _ = sse.MergeFragmentTempl(components.TaskRow(task),
        datastar.WithMergeTarget("#tasks-table-body"),
        datastar.WithMergeMode(datastar.FragmentMergeModePrepend),
    )

    // Patch signals to dismiss client modal
    _ = sse.MergeSignals(map[string]any{"modalOpen": false})
}
```

---

## 8. Running the Datastar Starter

### 1. Start the Datastar Server

```bash
go run main.go
```

The server will start listening at `http://localhost:8081` (defaulted to port 8081 so you can run it concurrently with Stage 4's HTMX server on port 8080).

### 2. Compare Both Side-by-Side

Open two browser windows:

1. **HTMX + Alpine.js:** `http://localhost:8080/`
2. **Datastar SSE Edition:** `http://localhost:8081/`

Observe:

- In Datastar, clicking **Custom Task Dispatch** and submitting sends `@post('/tasks')` which returns an SSE stream that prepends the table row, updates the metrics card, and resets the client modal signal in **a single continuous HTTP stream**.
- Inspect Chrome DevTools Network tab: notice that only 1 script tag (`datastar.js`, 13 KB) is transferred.

### 3. Run the Unit Test Suite

```bash
go test -v ./...
```

All tests verify root rendering, SSE wire formatting, signal patching, task cancellation, and real-time event streaming.
