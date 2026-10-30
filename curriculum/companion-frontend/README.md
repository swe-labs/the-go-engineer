# The Go Engineer — Frontend Companion Track

> **Status:** Released Companion Track  
> **Target Audience:** Backend Go engineers seeking autonomous full-stack browser capability without bloated JavaScript build pipelines.  
> **Core Architecture:** Progressive Server-Driven Web Applications (`Go` + `templ` + `HTMX` + `Tailwind CSS` + `Alpine.js`, with `Datastar` SSE comparison).

---

## 1. Architectural Philosophy

A capable backend Go engineer must be able to deliver complete, working software products to users through the browser. However, forcing learners to juggle heavy Single-Page Application (SPA) toolchains (TypeScript, React, Next.js, Webpack/Vite, state management libraries, and hydration mismatches) concurrently with low-level systems programming dilutes cognitive focus and introduces artificial complexity.

The **Frontend Companion Track** runs parallel to the core Go engineering curriculum. It teaches a modern, lightweight, **hypermedia-driven server architecture** that leverages Go's high-concurrency server strengths:

```
┌────────────────────────────────────────────────────────┐
│                   Browser Client                       │
│  - Semantic HTML5 & Responsive Tailwind CSS            │
│  - HTMX (AJAX swaps, polling, out-of-band updates)     │
│  - Alpine.js (Lightweight client-only UI state)        │
│  - Datastar (~13 KB unified SSE-first reactive engine) │
└───────────────────────────▲────────────────────────────┘
                            │
               HTTP / HTML Fragments & SSE Streams
                            │
┌────────────────────────────────────────────────────────┐
│                  Go Backend Server                     │
│  - Type-safe components (a-h/templ)                    │
│  - High-performance routing & middleware               │
│  - Real-time SSE streaming (http.Flusher / goroutines) │
│  - PostgreSQL persistence & transactional business logic│
└────────────────────────────────────────────────────────┘
```

### Why Hypermedia Over SPAs?
1. **Single Language for Core Logic:** Business logic, validation, data models, and template rendering stay in type-checked Go.
2. **Zero JavaScript Build Pipeline:** No `npm`, `node_modules`, or bundler configuration required for initial deployments.
3. **Instant Latency & Low Memory:** Server-rendered HTML streams immediately; client memory footprint is negligible.
4. **Natural Architectural Progression:** Builds incrementally on Go standard library skills (`html/template`, `net/http`) without context switching.

---

## 2. 52-Week Convergence Timeline

The companion track is organized into four progressive stages coordinated with milestones in the backend core path:

| Stage | Weeks | Focus Topics | Core Path Alignment | Practical Convergence Project | Executable Starter Directory |
| :---: | :---: | :--- | :--- | :--- | :--- |
| **1** | **10–14** | **Semantic HTML5 & Modern CSS** | Module 07 (CLI & Files) | Static Documentation & Template Renderer | [`stage-01-semantic-html-css/`](./stage-01-semantic-html-css) |
| **2** | **15–24** | **Browser JS, DOM & Fetch API** | Module 08 & 09 (HTTP & SQL) | Dynamic Book Catalog with REST Client | [`stage-02-vanilla-dom-fetch/`](./stage-02-vanilla-dom-fetch) |
| **3** | **25–33** | **Type-Safe UI with `templ` & Tailwind** | Module 10 (Auth & Security) | Authenticated Multi-Tenant Management Portal | [`stage-03-templ-tailwind/`](./stage-03-templ-tailwind) |
| **4** | **34–42** | **Hypermedia Interactivity: HTMX, Alpine & Datastar** | Modules 11–15 (Concurrency & CI/CD) | Live Metrics Dashboard with SSE (HTMX + Datastar Starters) | [`stage-04-htmx-alpine/`](./stage-04-htmx-alpine) |
| **5** | **43–50** | **Flagship Opslane Integration** | Module 18 (Flagship Opslane) | Production Multi-Tenant Opslane Console | `curriculum/modules/18-flagship-opslane` |

---

## 3. Detailed Stage Breakdown

### Stage 1: Semantic HTML5 & Modern CSS (Weeks 10–14)
*Coordinated with Module 07 (CLI, Files, and HTML Templates)*

#### Learning Objectives
- Write accessible, semantic HTML5 markup (`<header>`, `<nav>`, `<main>`, `<article>`, `<section>`, `<form>`, `<dialog>`).
- Implement responsive viewport layouts using CSS Grid and Flexbox without UI frameworks.
- Understand CSS custom properties (variables), box-sizing, typography scales, and media queries.
- Connect Go's `html/template` package to deliver server-rendered dynamic web pages.

#### Core Verification Task
Build a server-rendered personal engineering wiki using Go's `html/template` that displays markdown documents with responsive sidebar navigation, dark/light theme toggle via CSS variables, and zero JavaScript.

---

### Stage 2: Browser JavaScript, DOM & Fetch API (Weeks 15–24)
*Coordinated with Module 08 (HTTP REST APIs) and Module 09 (SQL & Persistence)*

#### Learning Objectives
- Inspect and manipulate the Document Object Model (DOM) using modern vanilla JavaScript (`querySelector`, `addEventListener`).
- Master event bubbling, capture phases, and event delegation patterns.
- Execute asynchronous HTTP network requests via the `fetch()` API with `async/await`.
- Handle client-side error states, parse JSON payloads, and manage HTTP status codes (200, 400, 401, 403, 500).
- Debug network traffic, CORS headers, and cookie credentials in Chrome DevTools.

#### Core Verification Task
Implement an interactive client for the Module 08 **Book Catalog REST API**:
- Search books with debounced input requests.
- Add and delete entries asynchronously without page reload.
- Display accessible inline validation error messages returned by Go HTTP handlers.

---

### Stage 3: Type-Safe Server Components with `a-h/templ`, Tailwind & `shadcn-templ` (Weeks 25–33)
*Coordinated with Module 10 (Auth & Security)*

#### Learning Objectives
- Replace error-prone string templates with [`a-h/templ`](https://templ.guide) — compile-time checked Go component templates.
- Compose reusable components with typed parameters, children blocks, and standard Go conditionals/loops.
- Rapidly build sleek, modern user interfaces using Tailwind CSS utility classes and [`shadcn-templ`](https://shadcn-templ.com) copy-and-own UI primitives (`@ui.Button`, `@ui.Card`, `@ui.Badge`, `@ui.Input`, `@ui.Separator`), configured via `components.json`.
- Eliminate Node.js/npm dependencies entirely using the Standalone Tailwind CLI and live reloading via `air` or Task.
- Handle secure authentication flows: session cookies vs JWT Bearer tokens, CSRF protection middleware, and role-based UI rendering.

#### Component Example (`profile.templ` with `@ui.Button` & `@ui.Badge`):
```templ
package components

import "github.com/rasel9t6/the-go-engineer/curriculum/companion-frontend/stage-03-templ-tailwind/components/ui"

templ RoleIndicator(role string) {
    if role == "admin" {
        @ui.Badge(ui.BadgeProps{Variant: "warning"}) {
            Admin (Full Control)
        }
    } else {
        @ui.Badge(ui.BadgeProps{Variant: "secondary"}) {
            Viewer (Read Only)
        }
    }
}
```

#### Core Verification Task
Construct an authenticated user settings and profile management portal connected to the Module 10 **Secure Authenticated API**, with role-guarded controls, type-safe `shadcn-templ` components, and CSRF-protected form submissions.

---

### Stage 4: Hypermedia-Driven Interactivity with HTMX & Alpine.js (Weeks 34–42)
*Coordinated with Modules 11–15 (Concurrency, Observability, and Deployment)*

#### Learning Objectives
- Drive complex user interactions purely through HTML attributes with [HTMX](https://htmx.org):
  - `hx-get`, `hx-post`, `hx-put`, `hx-delete` for AJAX requests.
  - `hx-target` and `hx-swap` (`innerHTML`, `outerHTML`, `beforebegin`) for granular DOM updates.
  - `hx-trigger` with modifiers (`keyup changed delay:500ms`, `every 5s` for live polling).
  - Server-Sent Events (SSE) via `hx-ext="sse"` for real-time background worker updates.
- Manage micro-client states using [Alpine.js](https://alpinejs.dev) (modals, dropdown toggles, multi-step tab switching) without a full framework.

#### HTMX Active Search & Inline Edit Pattern:
```html
<!-- Live search with server-rendered partial table swap -->
<input type="text"
       name="search"
       placeholder="Search active orders..."
       class="border rounded px-3 py-2 w-full mb-4"
       hx-post="/admin/orders/search"
       hx-trigger="keyup changed delay:300ms"
       hx-target="#order-table-body"
       hx-indicator="#search-spinner" />

<div id="search-spinner" class="htmx-indicator text-sm text-blue-500">Searching...</div>

<table>
    <tbody id="order-table-body">
        <!-- Rendered by Go templ handler -->
    </tbody>
</table>
```

#### Core Verification Task
Build a live log stream and background job runner dashboard connected to Module 11 (Worker Pools) and Module 12 (Observability) using Server-Sent Events (SSE) and HTMX out-of-band swaps (`hx-swap-oob`).

#### The GOTH & GOTD Architectural Pipelines:
Stage 4 includes dual-mode implementations supporting both standard `html/template` and compile-time type-checked `templ` components:
- **GOTH (Go + Orchestration/Tailwind + Templ + HTMX):** Located in [`stage-04-htmx-alpine/components/`](./stage-04-htmx-alpine/components), rendering type-safe components via `BroadcastComponent` over SSE.
- **GOTD (Go + Orchestration/Tailwind + Templ + Datastar):** Located in [`stage-04-htmx-alpine/datastar/components/`](./stage-04-htmx-alpine/datastar/components), streaming typed components directly into `datastar-patch-elements` wire protocol via `PatchComponent`.

#### Modern Go-Native Alternative: Datastar (SSE-First Paradigm)
While HTMX and Alpine.js teach essential, resume-ready hypermedia foundations, Go backend engineers frequently encounter friction at the boundary between HTMX DOM swaps and Alpine client state.

The companion track includes a complete comparative starter in [`stage-04-htmx-alpine/datastar/`](./stage-04-htmx-alpine/datastar/README.md) demonstrating **Datastar**:
- **Unified ~13 KB Footprint:** Replaces both HTMX and Alpine with a single script, eliminating state desync.
- **Server-Sent Events Native:** Every HTTP interaction operates over persistent or on-demand SSE streams (`datastar-patch-elements` and `datastar-patch-signals`).
- **No Out-of-Band Workarounds:** Any SSE event can target arbitrary DOM selectors (`mode: prepend`, `mode: outer`, etc.) directly from Go's standard `http.Flusher`.
- **First-Class Go SDK & Templ Integration:** Integrates natively with `github.com/starfederation/datastar/sdk/go` and compile-time type-checked `a-h/templ` components.

---

## 4. Final Capstone Convergence: Flagship Opslane (Weeks 43–50)

In Module 18, the learner brings together backend engineering and frontend companion mastery to construct the **Flagship Opslane Console**:

1. **Multi-Tenant Dashboard:** Responsive, accessible management interface with tenant switching, role permissions, and live system metrics.
2. **Order & Payment Pipeline UI:** Real-time order tracking, payment processing states, and transactional ledger audits.
3. **Operational Controls:** Graceful node draining, live health monitoring, and system metrics visualization.
4. **Zero Single-Page Overhead:** Delivered as a single compiled Go binary containing embedded assets via Go's standard `//go:embed` directive.
