# OPSL.4 HTTP API Layer

## Mission

Turn the auth and persistence foundation into a stable public JSON contract.

## What This Module Builds

- tenant and user setup endpoints
- login endpoint
- protected order and payment routes
- JSON error responses
- rate limiting and CORS middleware

## You Are Here If

- you can call the API and receive stable JSON responses
- you understand why authenticated routes read tenant identity from auth context
- you can explain where middleware stops and handlers begin

## Proof Surface

```bash
go test ./11-flagship/01-opslane/internal/handlers/...
go test ./11-flagship/01-opslane/internal/middleware/...
go run ./11-flagship/01-opslane/cmd/server
```

Manual checks:

```bash
curl http://localhost:8080/health
curl -X POST http://localhost:8080/api/v1/auth/login -H "Content-Type: application/json" -d "{\"tenant_id\":1,\"email\":\"admin@example.com\",\"password\":\"CorrectHorse7Battery\"}"
curl http://localhost:8080/api/v1/orders
```

The final request above must return `401` without a bearer token.

## Required Files and Boundaries

- `internal/handlers/handlers.go`
- `internal/handlers/api.go`
- `internal/middleware/middleware.go`

Implemented code surface: [SURFACE.md](./SURFACE.md)

Handlers parse and shape HTTP traffic. They should not become the order-processing engine.

## Engineering Questions

- Which requests should be rate-limited and which should stay exempt?
- How do you keep duplicate database conflicts from leaking raw driver errors to clients?
- Which timeouts and cleanup rules belong at the HTTP boundary?

## Next Step

Next: `OPSL.5` -> `11-flagship/01-opslane/modules/05-order-processing`

Open `11-flagship/01-opslane/modules/05-order-processing/README.md` to continue.
