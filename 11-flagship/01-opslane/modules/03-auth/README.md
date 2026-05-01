# OPSL.3 Authentication and Tenant Isolation

## Mission

Add identity and tenant-scoped request trust so later modules can rely on verified context.

## What This Module Builds

- password hashing
- signed token creation and verification
- auth middleware
- tenant and user identity propagation through request context

## You Are Here If

- you can explain why protected handlers do not accept tenant IDs from request bodies
- you understand where token verification happens
- you can trace a request from bearer token to trusted identity

## Proof Surface

```bash
go test ./11-flagship/01-opslane/internal/auth/...
```

## Required Files and Boundaries

- `internal/auth/token.go`
- `internal/auth/password.go`
- `internal/auth/service.go`
- `internal/auth/middleware.go`
- `internal/auth/context.go`

Implemented code surface: [SURFACE.md](./SURFACE.md)

Identity becomes trusted only after middleware verification. Repository calls should receive that trusted identity, not raw request claims.

## Engineering Questions

- Where should tenant identity be enforced first: middleware, service, or repository?
- What error shape should a tampered token produce?
- Which auth behavior belongs in middleware and which belongs in business services?

## Next Step

Next: `OPSL.4` -> `11-flagship/01-opslane/modules/04-http-api`

Open `11-flagship/01-opslane/modules/04-http-api/README.md` to continue.
