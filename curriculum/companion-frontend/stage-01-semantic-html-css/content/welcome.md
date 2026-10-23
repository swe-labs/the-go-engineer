# Engineering Wiki

Welcome to the **Go Engineer Wiki**. This system is rendered entirely using Go's standard library `html/template` package, semantic HTML5, and modern responsive CSS with **zero JavaScript**.

## Key Design Principles

1. **Semantic HTML5 First**: Structure determines meaning. We use `<nav>` for navigation trees, `<main>` for primary article content, `<article>` for documentation sections, and `<aside>` for contextual meta-information.
2. **Modern CSS Architecture**: Grid for two-dimensional layout orchestration (sidebar + content + metadata rail), Flexbox for one-dimensional component alignment.
3. **Pure CSS Dark/Light Theming**: Controlled by CSS Custom Properties (`--bg-primary`, `--text-primary`, `--border-color`) toggled via an accessible hidden checkbox state.
4. **Zero-JavaScript Resilience**: All core capabilities (navigation, page switching, theme toggling, responsive collapsing) function seamlessly with JavaScript disabled.

## Explore the Knowledge Base

- [Goroutines & Concurrency Model](/wiki/goroutines)
- [Go Memory Model & Allocation](/wiki/memory-model)
- [TCP Handshake & Lifecycle](/wiki/tcp-handshake)
