# Stage 3: Type-Safe Server Components with templ, Tailwind & shadcn-templ

> **Companion Track Stage:** 3 of 4  
> **Core Module Alignment:** Module 10 (*Authentication, Authorization, and Security*)  
> **Estimated Time:** 14–18 Hours  
> **Prerequisites:** Stages 1 & 2, HTTP cookie semantics, session authentication, CSRF attack models.  
> **Reference Documentation:** [shadcn-templ Official Docs](https://shadcn-templ.com/docs/introduction)

---

## 1. Overview & Architecture

This starter builds an **authenticated multi-tenant console** using [`a-h/templ`](https://templ.guide) — compile-time type-checked Go components — paired with Tailwind CSS utilities and canonical [`shadcn-templ`](https://shadcn-templ.com) (version 2.0) copy-and-own primitives.

```
┌────────────────────────────────────────────────────────────────────────┐
│                          Browser Client                                │
│  - Tailwind CSS Utility Styling (Dark-themed UI tokens)                │
│  - Session Cookie Authentication (HttpOnly auth_session)               │
│  - CSRF Form Protection (Synchronizer Token Pattern)                   │
│  - Role-Guarded UI (Admin actions vs Viewer indicators)                │
│  - Type-Safe UI Primitives (@ui.Card, @ui.Button, @ui.Badge, @ui.Input)│
└───────────────────────────────────▲────────────────────────────────────┘
                                    │
                       HTML Pages (Type-Safe Templ Components)
                                    │
┌───────────────────────────────────▼────────────────────────────────────┐
│                           Go Auth Server                               │
│  - Compiled a-h/templ components (components/ & ui/)                   │
│  - In-memory session store & crypto token generation                   │
│  - Strict CSRF verification middleware                                 │
│  - Role-based authorization enforcement (Admin vs Viewer)              │
│  - Zero-Node runtime execution (pure Go binary)                        │
└────────────────────────────────────────────────────────────────────────┘
```

---

## 2. The `shadcn-templ` Philosophy & Core Principles

[`shadcn-templ`](https://shadcn-templ.com) (formerly known as *templui*) is the official, community-led port of [shadcn/ui](https://ui.shadcn.com) engineered specifically for Go and `templ`.

> **"This is not a component library. It is how you build your component library."**

Traditional web frameworks force developers to install opaque NPM packages with black-box dependencies. When you need custom behaviors or brand tweaks, you end up writing brittle CSS overrides, wrapping components, or fighting incompatible upstream APIs.

`shadcn-templ` solves this by introducing five core design pillars:

1. **Open Code:**  
   The actual source code of every component is dropped directly into your project (`components/ui/`). You have 100% control to inspect, extend, refactor, or adjust classes without waiting for an upstream release.
2. **Composition:**  
   Every component shares a common, predictable, composable interface using `templ` child slots (`{ children... }`). Primitives assemble like Lego blocks (e.g. `Card` -> `Header` -> `Title` -> `Description` -> `Action` -> `Content` -> `Footer`).
3. **Distribution:**  
   A standardized flat-file schema (`components.json`) paired with the `shadcn-templ` CLI distributes, manages, and resolves component dependencies recursively while rewriting Go import paths to match your module.
4. **Beautiful Defaults:**  
   Components ship with carefully calibrated default typography, spacing, and dark-mode color tokens matching modern SaaS aesthetics.
5. **AI-Ready:**  
   Because component code lives unminified inside your repository as pure Go and `templ`, LLMs and coding assistants can read, understand, generate, and safely modify UI components with zero hallucinated abstractions.

---

## 3. Directory Structure

```
stage-03-templ-tailwind/
├── components.json       # Official shadcn-templ 2.0 project configuration
├── components/           # Type-safe templ components and generated Go code
│   ├── ui/               # shadcn-templ copy-and-own UI component primitives
│   │   ├── badge.templ   # Accessible status badge component (8 variants + link mode)
│   │   ├── badge_templ.go
│   │   ├── button.templ  # Button with 6 visual variants, 8 sizes, and link/button switching
│   │   ├── button_templ.go
│   │   ├── card.templ    # Container Card, Header, Title, Description, Action, Content, Footer
│   │   ├── card_templ.go
│   │   ├── input.templ   # Accessible form input with dark styling
│   │   ├── input_templ.go
│   │   ├── separator.templ # Semantic/decorative divider (horizontal & vertical)
│   │   ├── separator_templ.go
│   │   └── ui_test.go    # Comprehensive unit tests for UI component rendering
│   ├── header.templ      # Navigation bar and role indicator (uses @ui.Badge & @ui.Button)
│   ├── header_templ.go
│   ├── layout.templ      # Master HTML document wrapper
│   ├── layout_templ.go
│   ├── login.templ       # Sign-in form (uses @ui.Card, @ui.Badge, @ui.Input, @ui.Button)
│   ├── login_templ.go
│   ├── profile.templ     # User profile form (uses @ui.Card, @ui.Input, @ui.Button, @ui.Badge)
│   ├── profile_templ.go
│   ├── settings.templ    # Tenant configuration (uses @ui.Card, @ui.Input, @ui.Separator, @ui.Button)
│   └── settings_templ.go
├── static/
│   └── tailwind.css      # Pre-compiled Tailwind CSS classes
├── main.go               # Session store, CSRF validator & HTTP handlers
├── main_test.go          # Auth, CSRF, and RBAC integration tests
└── README.md             # This guide
```

---

## 4. The `components.json` Specification

The [`components.json`](./components.json) file holds configuration for your project and tells the CLI where to vendor components and how to rewrite imports:

```json
{
  "$schema": "https://shadcn-templ.com/schema/components.json",
  "style": "base-nova",
  "tailwind": {
    "css": "static/tailwind.css",
    "baseColor": "zinc",
    "cssVariables": true
  },
  "iconLibrary": "lucide",
  "rtl": false,
  "menuColor": "default",
  "menuAccent": "subtle",
  "aliases": {
    "utils": "github.com/rasel9t6/the-go-engineer/curriculum/companion-frontend/stage-03-templ-tailwind/utils",
    "components": "github.com/rasel9t6/the-go-engineer/curriculum/companion-frontend/stage-03-templ-tailwind/components",
    "ui": "github.com/rasel9t6/the-go-engineer/curriculum/companion-frontend/stage-03-templ-tailwind/components/ui"
  },
  "scripts": {
    "dir": "static",
    "path": "/static"
  }
}
```

### Key Configuration Keys:
- **`$schema`**: Points to the official JSON schema validator.
- **`style`**: The preset design system (e.g. `base-nova`, `base-vega`, `base-maia`).
- **`tailwind`**: Specifies the CSS entry file, base color palette (`zinc`, `neutral`, `slate`), and enables CSS variables.
- **`aliases`**: Maps Go import paths. When you run `shadcn-templ add <component>`, the CLI automatically rewrites registry imports to your configured package alias.
- **`scripts`**: Defines where client JavaScript bundles are written (`dir`) and the public URL path (`path`) served to the browser.

---

## 5. Universal Component Props

Every component in the `shadcn-templ` ecosystem accepts three universal props that provide platform-native flexibility:

| Prop | Type | Description |
| :--- | :--- | :--- |
| **`ID`** | `string` | HTML `id` attribute rendered on the root container element. |
| **`Class`** | `string` | Additional utility classes merged with the component's default classes. |
| **`Attributes`** | `templ.Attributes` | Spread HTML attributes (`{ props.Attributes... }`) for ARIA, HTMX, Datastar, or test IDs. |

### Example: Spreading Attributes
```templ
@ui.Button(ui.ButtonProps{
    Variant: ui.VariantOutline,
    Attributes: templ.Attributes{
        "data-testid": "submit-btn",
        "aria-label":  "Confirm Deployment",
        "hx-post":     "/api/deploy",
        "hx-target":   "#status",
    },
}) {
    Deploy
}
```

---

## 6. Component Catalog & API Reference

All Stage 3 primitives are implemented in [`components/ui/`](./components/ui/):

### 1. Button (`@ui.Button`)

A polymorphic component that renders either a native `<button>` or a semantic `<a>` anchor styled as a button when `Href` is provided.

- **Variants:** `VariantDefault`, `VariantSecondary`, `VariantDestructive`, `VariantOutline`, `VariantGhost`, `VariantLink`
- **Sizes:** `SizeDefault` (h-9), `SizeXs` (h-7), `SizeSm` (h-8), `SizeLg` (h-10), `SizeIcon` (h-9 w-9)
- **Types:** `TypeButton`, `TypeSubmit`, `TypeReset`
- **Helper:** `ui.Variants(props)` returns raw class strings for styling arbitrary elements.

```templ
@ui.Button(ui.ButtonProps{
    Variant: ui.VariantDestructive,
    Size:    ui.SizeSm,
    Type:    ui.TypeSubmit,
}) {
    Revoke API Key
}

@ui.Button(ui.ButtonProps{
    Variant: ui.VariantOutline,
    Href:    "/audit-logs",
    Target:  "_blank",
}) {
    View Audit Logs
}
```

### 2. Card (`@ui.Card`)

The container component for cards, following the official composition hierarchy:

```
ui.Card
├── ui.CardHeader
│   ├── ui.CardTitle
│   ├── ui.CardDescription
│   └── ui.CardAction (top-right slot for buttons/badges)
├── ui.CardContent
└── ui.CardFooter
```

- **Sizes:** `CardSizeDefault`, `CardSizeSm` (compact padding)
- **Props:** Each sub-component accepts `Class`, `ID`, and `Attributes`.

```templ
@ui.Card(ui.CardProps{ID: "metrics-card"}) {
    @ui.CardHeader(ui.CardHeaderProps{}) {
        @ui.CardTitle(ui.CardTitleProps{}) { Database Metrics }
        @ui.CardDescription(ui.CardDescriptionProps{}) { Active connection pool status }
        @ui.CardAction(ui.CardActionProps{}) {
            @ui.Badge(ui.BadgeProps{Variant: ui.BadgeVariantSuccess}) { Healthy }
        }
    }
    @ui.CardContent(ui.CardContentProps{}) {
        <p class="text-sm text-slate-300">Connections: 42 / 100</p>
    }
    @ui.CardFooter(ui.CardFooterProps{}) {
        <span class="text-xs text-slate-500">Updated 2m ago</span>
    }
}
```

### 3. Badge (`@ui.Badge`)

Status badge supporting pill styling and anchor link mode.

- **Variants:** `BadgeVariantDefault`, `BadgeVariantSecondary`, `BadgeVariantDestructive`, `BadgeVariantOutline`, `BadgeVariantGhost`, `BadgeVariantLink`, `BadgeVariantSuccess`, `BadgeVariantWarning`
- **Link Mode:** Pass `Href` to render an accessible `<a>` tag with badge aesthetics.

```templ
@ui.Badge(ui.BadgeProps{Variant: ui.BadgeVariantWarning}) {
    Admin (Full Access)
}

@ui.Badge(ui.BadgeProps{Variant: ui.BadgeVariantOutline, Href: "/docs/security"}) {
    Read Security Spec &rarr;
}
```

### 4. Input (`@ui.Input`)

Accessible dark-mode form input supporting native HTML input types and states.

- **Supported Types:** `text`, `email`, `password`, `file`, `number`, `search`
- **State Props:** `Disabled`, `ReadOnly`, `Required`, `Placeholder`, `Value`, `Form`, `Accept`

```templ
@ui.Input(ui.InputProps{
    Type:        "email",
    ID:          "email",
    Name:        "email",
    Value:       vm.User.Email,
    Required:    true,
    Placeholder: "name@company.com",
})
```

### 5. Separator (`@ui.Separator`)

A visual or semantic divider for grouping related controls.

- **Orientations:** `OrientationHorizontal` (default, `h-[1px] w-full`), `OrientationVertical` (`h-full w-[1px]`)
- **Accessibility:** Supports `Decorative: true` (renders `role="none"`) or semantic `role="separator"` with `aria-orientation`.

```templ
@ui.Separator(ui.SeparatorProps{Class: "my-4"})
```

---

## 7. Client-Side JavaScript Bundling Architecture

While basic components (Button, Card, Badge, Input, Separator) are pure server-rendered HTML requiring **zero JavaScript**, advanced components (Dialogs, Dropdowns, Sheets, Popovers) use vanilla JavaScript ports of headless Base UI behaviors.

`shadcn-templ` handles client scripts through a **single build-time bundle**:

1. **No Per-Component Script Tags:**  
   The layout renders a single script tag in `<head>` via `@components.Scripts()`.
2. **Deterministic Build Hashing:**  
   The CLI scans `components/*/*.js` and bundles them into `static/shadcn-templ-<hash>.js`. The hash is stored in `components/scripts_bundle.go`.
3. **Immutable Caching:**  
   Production deployments can serve the hashed bundle with `Cache-Control: public, max-age=31536000, immutable`.
4. **DOM Watching (Hypermedia Compatibility):**  
   The bundled scripts use `MutationObserver` on the `document.body`. When new HTML fragments are swapped into the DOM by **HTMX**, **Datastar**, or **Alpine.js**, components automatically bind their event listeners without requiring custom re-initialization scripts!

---

## 8. CLI Tooling & Workflows

To use the official `shadcn-templ` CLI:

```bash
# Install or update the CLI
go install github.com/axadrn/shadcn-templ/v2/cmd/shadcn-templ@latest

# Initialize a project with an official preset (e.g. nova, vega, maia)
shadcn-templ init --preset nova

# Add components (resolves dependencies recursively)
shadcn-templ add button
shadcn-templ add card
shadcn-templ add dialog

# Bundle client-side JavaScript for components
shadcn-templ bundle

# Watch scripts during development
shadcn-templ bundle --watch

# Apply a new design system preset to an existing codebase
shadcn-templ apply vega
```

---

## 9. Zero-Node Production & Local Toolchain

You can develop, build, and deploy this entire stack with **zero Node.js or npm installed on your machine**:

### 1. Download Standalone Tailwind CSS v4 CLI
Download the official single-binary executable from [Tailwind CSS Releases](https://github.com/tailwindlabs/tailwindcss/releases/latest):
```bash
# Example (Linux):
curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64
chmod +x tailwindcss-linux-x64 && sudo mv tailwindcss-linux-x64 /usr/local/bin/tailwindcss
```

### 2. Development Task Runner (`Taskfile.yml`)
```yaml
version: "3"

tasks:
  tailwind:
    desc: Watch Tailwind CSS changes
    cmds:
      - tailwindcss -i ./static/input.css -o ./static/tailwind.css --watch

  templ:
    desc: Generate Go code from .templ files with live reload
    cmds:
      - go run github.com/a-h/templ/cmd/templ@v0.3.1070 generate --watch --proxy="http://localhost:8080" --cmd="go run main.go"

  dev:
    desc: Run Tailwind and templ concurrently
    cmds:
      - task --parallel tailwind templ
```

---

## 10. Running the Starter

### Start the Server
```bash
go run main.go
```
The server will start listening at `http://localhost:8080`.

Open your browser to:
`http://localhost:8080/login`

### Default Test Credentials:
| Account | Email | Password | Role | Access Level |
| :--- | :--- | :--- | :--- | :--- |
| **Admin** | `alice@acme.corp` | `opslane123` | `admin` | Full access (Can edit tenant settings and save profile) |
| **Viewer** | `bob@acme.corp` | `opslane123` | `viewer` | Read-only access (Tenant settings inputs disabled) |

### Verification Checklist:
1. **Alice (Admin):** Sign in and observe the amber `@ui.Badge` in the header. Navigate to **Tenant Settings**, edit the Organization Name, and submit. Observe flash confirmation.
2. **Bob (Viewer):** Sign in and observe the slate `@ui.Badge` in the header. Navigate to **Tenant Settings** and observe the read-only notice with disabled `@ui.Input` elements.
3. **CSRF Validation:** Verify that submitting form data without `csrf_token` yields an immediate `403 Forbidden` response.

### Modifying and Regenerating Templates
```bash
go run github.com/a-h/templ/cmd/templ@v0.3.1070 generate
```

### Run the Test Suite
```bash
go test -v ./...
```
This runs both the auth/security tests in `main_test.go` and the unit test suite in `components/ui/ui_test.go` covering all component variants, sizes, links, and attribute spreading.
