# Stage 3: Type-Safe Server Components with templ & Tailwind

> **Companion Track Stage:** 3 of 4  
> **Core Module Alignment:** Module 10 (*Authentication, Authorization, and Security*)  
> **Estimated Time:** 14–18 Hours  
> **Prerequisites:** Stages 1 & 2, HTTP cookie semantics, session authentication, CSRF attack models.

---

## 1. Overview

This starter builds an **authenticated multi-tenant console** using [`a-h/templ`](https://templ.guide) — compile-time type-checked Go components — paired with Tailwind CSS utilities.

```
┌────────────────────────────────────────────────────────┐
│                   Browser Client                       │
│  - Tailwind CSS Utility Styling (Dark-themed UI)       │
│  - Session Cookie Authentication (auth_session)        │
│  - CSRF Form Protection (Synchronizer Token Pattern)   │
│  - Role-Guarded UI (Admin controls vs Viewer badges)   │
└───────────────────────────▲────────────────────────────┘
                            │
               HTML Pages (Type-Safe Templ Components)
                            │
┌───────────────────────────▼────────────────────────────┐
│                  Go Auth Server                        │
│  - Compiled a-h/templ components                       │
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
            <span class="bg-amber-500/20 text-amber-400 text-xs px-2 py-0.5 rounded">Admin</span>
        } else {
            <span class="bg-slate-700 text-slate-300 text-xs px-2 py-0.5 rounded">Viewer</span>
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
│   ├── header.templ      # Navigation bar and role indicator
│   ├── header_templ.go
│   ├── layout.templ      # Master HTML document wrapper
│   ├── layout_templ.go
│   ├── login.templ       # Sign-in form
│   ├── login_templ.go
│   ├── profile.templ     # User profile and preferences form
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

## 4. Running the Starter

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
1. Signing in as **Alice**: Notice the amber `Admin` badge. Navigate to **Tenant Settings** and change the organization name.
2. Signing out and signing in as **Bob**: Notice the `Viewer` badge. Navigate to **Tenant Settings** and observe the read-only notice with disabled inputs.
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

---

## 5. Learner Extension Challenges

1. **Password Change Component**: Add a CSRF-protected password modification card in `profile.templ` requiring old password verification.
2. **Flash Message Auto-Dismissal**: Add an accessible toast notification component for successful saves.
3. **Audit Log Component**: Render a table of recent tenant activity with timestamps and actor emails in `settings.templ`.
