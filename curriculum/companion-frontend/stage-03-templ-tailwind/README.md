# Stage 3: Type-Safe Server Components with templ, Tailwind & shadcn-templ

> **Companion Track Stage:** 3 of 4  
> **Core Module Alignment:** Module 10 (*Authentication, Authorization, and Security*)  
> **Estimated Time:** 14–18 Hours  
> **Prerequisites:** Stages 1 & 2, HTTP cookie semantics, session authentication, CSRF attack models.

---

## 1. Overview

This starter builds an **authenticated multi-tenant console** using [`a-h/templ`](https://templ.guide) — compile-time type-checked Go components — paired with Tailwind CSS utilities and modern [`shadcn-templ`](https://github.com/axadrn/shadcn-templ) (copy-and-own) primitives.

```
┌────────────────────────────────────────────────────────┐
│                   Browser Client                       │
│  - Tailwind CSS Utility Styling (Dark-themed UI)       │
│  - Session Cookie Authentication (auth_session)        │
│  - CSRF Form Protection (Synchronizer Token Pattern)   │
│  - Role-Guarded UI (Admin controls vs Viewer badges)   │
│  - Type-Safe UI Primitives (Button, Card, Badge)       │
└───────────────────────────▲────────────────────────────┘
                            │
               HTML Pages (Type-Safe Templ Components)
                            │
┌───────────────────────────▼────────────────────────────┐
│                  Go Auth Server                        │
│  - Compiled a-h/templ components (components/ & ui/)   │
│  - In-memory session store & crypto token generation   │
│  - Strict CSRF verification middleware                 │
│  - Role-based authorization enforcement                │
└────────────────────────────────────────────────────────┘
```

---

## 2. Key Architecture & Security Concepts

### 1. Type Safety with `a-h/templ`
Standard Go `html/template` parses strings at runtime. Typos in struct field names or passing incorrect types only fail when a user triggers the specific code path. 

In contrast, `templ` compiles `.templ` files into standard Go functions (`components.Profile(vm).Render(ctx, w)`). If you misname a struct property or pass an invalid type, `go build` fails immediately at compile time.

```templ
templ Header(user *UserContext) {
    <header class="flex justify-between items-center bg-slate-900 text-white p-4">
        <div>Opslane Console</div>
        if user.Role == "admin" {
            @ui.Badge(ui.BadgeProps{Variant: "warning"}) {
                Admin (Full Control)
            }
        } else {
            @ui.Badge(ui.BadgeProps{Variant: "secondary"}) {
                Viewer (Read Only)
            }
        }
    </header>
}
```

### 2. Session Cookies & CSRF Protection
- **HttpOnly Cookies:** `auth_session` cookie cannot be read by JavaScript `document.cookie`, neutralizing XSS token theft.
- **CSRF Token:** Every authenticated session generates a 256-bit cryptographically random token. All state-modifying requests (`POST /profile`, `POST /settings`) must submit matching tokens. Requests lacking the token are immediately rejected with `403 Forbidden`.

### 3. Role-Based Access Control (RBAC)
- **Viewer Role (`bob@acme.corp`):** Can view workspace and tenant parameters, but settings inputs are disabled (`disabled?={ vm.User.Role != "admin" }`) and POST submissions are blocked by server middleware.
- **Admin Role (`alice@acme.corp`):** Can toggle strict multi-tenant isolation and update organization identity.

---

## 3. Directory Structure

```
stage-03-templ-tailwind/
├── components/           # Type-safe templ components and generated Go code
│   ├── ui/               # shadcn-templ copy-and-own UI component primitives
│   │   ├── badge.templ   # Accessible status badge component
│   │   ├── badge_templ.go
│   │   ├── button.templ  # Button with 6 visual variants & size options
│   │   ├── button_templ.go
│   │   ├── card.templ    # Container Card, Header, Title, Content, Footer
│   │   ├── card_templ.go
│   │   └── ui_test.go    # Unit tests for UI component rendering & variants
│   ├── header.templ      # Navigation bar and role indicator
│   ├── header_templ.go
│   ├── layout.templ      # Master HTML document wrapper
│   ├── layout_templ.go
│   ├── login.templ       # Sign-in form
│   ├── login_templ.go
│   ├── profile.templ     # User profile and preferences form (using @ui.Button and @ui.Badge)
│   ├── profile_templ.go
│   ├── settings.templ    # Tenant configuration & role guard
│   └── settings_templ.go
├── static/
│   └── tailwind.css      # Pre-compiled Tailwind CSS classes
├── main.go               # Session store, CSRF validator & handlers
├── main_test.go          # Auth, CSRF, and RBAC integration tests
└── README.md             # This guide
```

---

## 4. Modern Production Architecture: The GOTH Stack & shadcn-templ

Modern production Go web development has converged around the **"GOTH" Stack**:
- **G**o (Standard library `net/http` or `go-chi/chi`)
- **O**rchestration / Zero-Node toolchains (`air` + Standalone Tailwind CLI)
- **T**empl (`a-h/templ` compile-time type-safe templates)
- **H**TMX (or Datastar) for server-driven hypermedia updates

### What is `shadcn-templ`?

[`shadcn-templ`](https://github.com/axadrn/shadcn-templ) (formerly known as *templui*) is a 1:1 open-source port of the React `shadcn/ui` design system engineered specifically for Go and `templ`.

#### Why the "Copy-and-Own" Model Fits Go
Traditional web frameworks force developers to import large, black-box npm packages that break on major version bumps. `shadcn-templ` completely rejects this:
1. **Source Code in Your Project:** Components live directly in your local `components/ui/` folder as `.templ` files.
2. **Zero Runtime Dependencies:** Components render to pure HTML via Go functions. No React, no Node.js, and zero JavaScript runtime dependencies required for basic components.
3. **Absolute Customization:** If you need to tweak padding, change variant colors, or add a custom data attribute, you simply edit `components/ui/button.templ` directly.

### Calling `shadcn-templ` Primitives in `.templ`

Once a component is placed in `components/ui/`, call it natively inside any template:

```templ
package components

import "github.com/rasel9t6/the-go-engineer/curriculum/companion-frontend/stage-03-templ-tailwind/components/ui"

templ QuickActions() {
    @ui.Card(ui.CardProps{}) {
        @ui.CardHeader(ui.CardProps{}) {
            @ui.CardTitle(ui.CardProps{}) { Quick Deployment }
            @ui.CardDescription(ui.CardProps{}) { Trigger zero-downtime updates }
        }
        @ui.CardContent(ui.CardProps{}) {
            <p class="text-sm text-slate-300">Cluster: production-us-east-1</p>
        }
        @ui.CardFooter(ui.CardProps{}) {
            @ui.Button(ui.ButtonProps{Variant: "default", Size: "sm"}) {
                Deploy Container
            }
            @ui.Button(ui.ButtonProps{Variant: "outline", Size: "sm", Href: "/logs"}) {
                View Runbook
            }
        }
    }
}
```

### Complete Zero-Node Toolchain Guide

You can run, compile, and hot-reload this entire stack with **zero Node.js or npm**:

#### 1. Download the Standalone Tailwind CLI
Download the official single-binary executable for your OS (Windows, Linux, or macOS) from [Tailwind CSS Releases](https://github.com/tailwindlabs/tailwindcss/releases):
```bash
# Example (Linux/macOS):
curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-arm64
chmod +x tailwindcss-macos-arm64
mv tailwindcss-macos-arm64 /usr/local/bin/tailwindcss
```

#### 2. Live Development with Air
Configure a `Makefile` or task runner to orchestrate CSS watching, `templ` compilation, and live Go binary reloading simultaneously:

```makefile
# Run everything concurrently in development:
dev:
	# 1. Watch .templ files and compile Tailwind CSS output
	tailwindcss -i ./static/input.css -o ./static/tailwind.css --watch &
	# 2. Watch .templ files and re-generate Go code
	go run github.com/a-h/templ/cmd/templ@latest generate --watch --proxy="http://localhost:8080" &
	# 3. Live-reload the Go server binary on code changes
	air
```

---

## 5. Running the Starter

### Start the Server
```bash
go run main.go
```
The server will start listening at `http://localhost:8080`.

Open your browser to:
`http://localhost:8080/login`

### Default Test Credentials:
| Account | Email | Password | Role |
| :--- | :--- | :--- | :--- |
| **Admin** | `alice@acme.corp` | `opslane123` | Full access (Can edit tenant settings) |
| **Viewer** | `bob@acme.corp` | `opslane123` | Read-only access (Tenant settings disabled) |

Try:
1. Signing in as **Alice**: Notice the amber `@ui.Badge` for `Admin`. Navigate to **Tenant Settings** and change the organization name.
2. Signing out and signing in as **Bob**: Notice the `@ui.Badge` for `Viewer`. Navigate to **Tenant Settings** and observe the read-only notice with disabled inputs.
3. Submitting a form without a valid CSRF token to verify the `403 Forbidden` defense.

### Modifying and Regenerating Templates
If you edit any `.templ` file, re-run the generator:
```bash
go run github.com/a-h/templ/cmd/templ@latest generate -path ./components
```

### Run the Test Suite
```bash
go test -v ./...
```
This runs both the authentication/security integration tests in `main_test.go` and the `shadcn-templ` UI component tests in `components/ui/ui_test.go`.

---

## 6. Learner Extension Challenges

1. **Add an Alert / Dialog Component**: Port a `shadcn-templ` `alert.templ` or modal dialog into `components/ui/` and use it to confirm dangerous tenant deletion actions.
2. **Password Change Component**: Add a CSRF-protected password modification card in `profile.templ` using `@ui.Card` and `@ui.Button(ui.ButtonProps{Variant: "destructive"})`.
3. **Hypermedia Integration**: Combine `@ui.Button` with HTMX (`hx-post="/api/tenant/toggle" hx-target="#tenant-status"`) to perform instant inline status updates without page reloads.
