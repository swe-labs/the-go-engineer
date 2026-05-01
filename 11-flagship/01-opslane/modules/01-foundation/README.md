# OPSL.1 Foundation and Configuration

## Mission

Make Opslane start as a real service with validated configuration instead of magic defaults.

## What This Module Builds

- the server entrypoint
- environment-aware configuration loading
- startup validation
- local container wiring for PostgreSQL

## You Are Here If

- you can run the Opslane server from `cmd/server`
- you understand where runtime configuration comes from
- you can explain which values are safe defaults and which must be explicit

## Proof Surface

```bash
go test ./11-flagship/01-opslane/internal/config/...
go run ./11-flagship/01-opslane/cmd/server
```

## Required Files and Boundaries

- `cmd/server/main.go`
- `internal/config/config.go`
- `internal/config/environment.go`
- `.env.example`
- `Dockerfile`
- `docker-compose.yml`

Implemented code surface: [SURFACE.md](./SURFACE.md)

Do not move tenant, auth, or business logic into the config package.

## Engineering Questions

- Which settings should fail fast at startup instead of surfacing after traffic arrives?
- Which secrets should never have production defaults?
- How do you keep local overrides explicit without hard-coding environment behavior in handlers?

## Next Step

Next: `OPSL.2` -> `11-flagship/01-opslane/modules/02-database`

Open `11-flagship/01-opslane/modules/02-database/README.md` to continue.
