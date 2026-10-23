# Stage 1: Semantic HTML5 & Modern CSS Wiki

> **Companion Track Stage:** 1 of 4  
> **Core Module Alignment:** Module 07 (*CLI, Files, JSON, and Templates*)  
> **Estimated Time:** 10–14 Hours  
> **Prerequisites:** Go standard library fundamentals, basic terminal usage.

---

## 1. Overview

This starter implements a **server-rendered engineering documentation wiki** using Go's standard library `html/template`, semantic HTML5 elements, and modern CSS without external UI frameworks or client-side JavaScript.

### Architectural Principles

```
┌────────────────────────────────────────────────────────┐
│                   Browser Client                       │
│  - Semantic HTML5: <header>, <nav>, <main>, <article>  │
│  - CSS Custom Properties (Theme Tokens)                │
│  - CSS Grid (Workspace) & Flexbox (Nav / Rail)         │
│  - Zero JavaScript Theme Toggle via Checkbox Hack      │
└───────────────────────────▲────────────────────────────┘
                            │
               HTTP / HTML (Fully Rendered)
                            │
┌───────────────────────────▼────────────────────────────┐
│                  Go Wiki Server                        │
│  - Go standard library net/http & html/template        │
│  - Embedded assets via //go:embed                      │
│  - Contextual metadata calculation (Word count, read)  │
└────────────────────────────────────────────────────────┘
```

---

## 2. Key Features

1. **Pure Semantic Markup**:
   - `<header class="wiki-header">` for brand identity and accessibility controls.
   - `<nav class="sidebar-nav">` with unordered lists (`<ul>`, `<li>`, `<a>`) for document navigation.
   - `<main class="wiki-main">` containing `<article class="wiki-article">` for primary content.
   - `<aside class="wiki-rail">` with `<dl>`, `<dt>`, `<dd>` definition lists for document specs.
   - `<footer class="wiki-footer">` for system provenance.

2. **Modern CSS Layout Engine**:
   - **CSS Grid Workspace**: Three-column responsive grid (`260px minmax(0, 1fr) 260px`).
   - **Responsive Breakpoints**: Gracefully collapses to two columns on tablets (`<= 1080px`) and a single vertical column on mobile (`<= 768px`).
   - **Sticky Rails**: Sidebar navigation and document specs rail stick to the viewport during long article scrolling (`position: sticky`).

3. **Zero-JS Dark/Light Theme Switching**:
   - Implemented via a hidden input element `<input type="checkbox" id="theme-toggle">` and the CSS sibling selector `#theme-toggle:checked ~ .wiki-layout`.
   - Modifies CSS custom properties (`--bg-canvas`, `--bg-surface`, `--text-main`, `--border-color`) in real time without a single line of JavaScript.

4. **Self-Contained Go Binary**:
   - Templates, static stylesheets, and initial markdown documents are embedded into the compiled binary via `//go:embed`.
   - Runs everywhere with zero runtime disk dependencies.

---

## 3. Directory Structure

```
stage-01-semantic-html-css/
├── content/              # Raw markdown engineering documentation
│   ├── goroutines.md
│   ├── memory-model.md
│   ├── tcp-handshake.md
│   └── welcome.md
├── static/
│   └── style.css         # Modern CSS styles with theme variables & Grid layout
├── templates/
│   ├── layout.html       # Primary semantic HTML5 layout template
│   └── not_found.html    # 404 error page template
├── main.go               # HTTP server, routing, template rendering
├── main_test.go          # Unit & integration test suite
└── README.md             # This guide
```

---

## 4. Running the Starter

### Start the Server
```bash
go run main.go
```
By default, the server listens at `http://localhost:8080`.

Open your browser to:
- `http://localhost:8080/wiki/welcome`
- `http://localhost:8080/wiki/goroutines`
- `http://localhost:8080/wiki/memory-model`
- `http://localhost:8080/wiki/tcp-handshake`

### Run the Test Suite
```bash
go test -v ./...
```

---

## 5. Learner Extension Challenges

1. **Add Table of Contents (TOC)**: Enhance `main.go` to scan markdown headings (`##`, `###`) and generate an in-page anchor navigation list inside `<aside class="wiki-rail">`.
2. **Add a New Article**: Create `content/distributed-consensus.md` and verify that the sidebar dynamically indexes and highlights it.
3. **Print Stylesheet**: Add an `@media print` rule to `style.css` that hides the sidebar and metadata rail, ensuring clean, high-contrast hardcopy documentation prints.
