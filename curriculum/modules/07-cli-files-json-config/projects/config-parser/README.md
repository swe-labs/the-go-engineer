# Config Parser

## Learning objective

Build a robust, layered configuration parser in Go. Master JSON and key-value file parsing with `encoding/json`, schema validation, environment variable override precedence (Twelve-Factor App methodology), and mockable environment lookups.

## Why this matters

Modern containerized microservices must dynamically adapt their runtime configuration based on local development files, Kubernetes config maps, and production secret injections. Establishing a clean configuration hierarchy (Defaults < Files < Environment Overrides) with strict validation prevents services from launching with broken ports, missing database targets, or insecure defaults.

## Prerequisites

Complete all Module 07 lessons (01–21) on JSON serialization/deserialization, file reading with `os.ReadFile`, environment variable management with `os.LookupEnv`, and struct validation.

## Available concepts

This project draws on structs, `encoding/json`, `os`, `path/filepath`, `strconv`, dependency injection for `envLookup`, custom validation error aggregation, and table-driven testing.

## Project overview

The `_starter/` directory provides type signatures, baseline default structs, and stub loader methods. Your job is to implement the layered configuration loading and validation engine so that:

- Defaults are applied automatically when fields are unspecified.
- Configuration files (`.json` or `.env`) merge into defaults without overwriting unspecified fields.
- Environment variables (`APP_*`) take highest precedence, overriding file and default values.
- Invariants are verified before returning: `server.port` must be a valid TCP port (1–65535), `database.name` must be non-empty, and `database.max_conns` must be positive.
- Unsupported file formats return `ErrUnsupportedFormat`.

## Precedence Hierarchy

```
[ Defaults ] (Lowest)
     |
     v
[ Config File (.json / .env) ]
     |
     v
[ Environment Variables (APP_*) ] (Highest)
     |
     v
[ Validation Engine ] ---> Final AppConfig
```

## Environment Variable Mappings

| Environment Variable | Target Struct Field |
|---|---|
| `APP_ENV` | `Environment` |
| `APP_SERVER_HOST` | `Server.Host` |
| `APP_SERVER_PORT` | `Server.Port` |
| `APP_DB_HOST` | `Database.Host` |
| `APP_DB_PORT` | `Database.Port` |
| `APP_DB_NAME` | `Database.Name` |
| `APP_DB_USER` | `Database.User` |
| `APP_DB_PASSWORD` | `Database.Password` |
| `APP_DB_MAX_CONNS` | `Database.MaxConns` |
| `APP_METRICS_ENABLED` | `Metrics.Enabled` |

## Tasks

### Task 1 — Implement file loading & format detection

In `_starter/main.go`:
- Inspect file extension:
  - `.json`: unmarshal into `AppConfig` using `encoding/json`.
  - `.env` / `.kv`: parse line-by-line `KEY=VALUE` pairs via `parseKeyValueConfig`.
  - Unrecognized extension: return `ErrUnsupportedFormat`.

### Task 2 — Implement environment variable overrides

- Implement `applyEnvOverrides(cfg *AppConfig)`:
  - Query environment variables through `p.envLookup` (mockable).
  - Convert strings to integers or booleans where appropriate (`strconv.Atoi`, `strings.EqualFold`).

### Task 3 — Implement schema validation

- Implement `Validate(cfg AppConfig) error`:
  - Verify `cfg.Server.Port` is between 1 and 65535.
  - Verify `cfg.Database.Name` is non-empty.
  - Verify `cfg.Database.MaxConns` > 0.
  - Return aggregated error wrapping `ErrValidation` if any checks fail.

## Verification

Run all automated unit tests:

```bash
go test -v ./curriculum/modules/07-cli-files-json-config/projects/config-parser/_starter
```

Expected output:
```
=== RUN   TestConfigParser_Defaults
--- PASS: TestConfigParser_Defaults (0.00s)
=== RUN   TestConfigParser_JSONFileMerge
--- PASS: TestConfigParser_JSONFileMerge (0.01s)
=== RUN   TestConfigParser_EnvOverridePrecedence
--- PASS: TestConfigParser_EnvOverridePrecedence (0.01s)
=== RUN   TestConfigParser_ValidationErrors
--- PASS: TestConfigParser_ValidationErrors (0.00s)
=== RUN   TestConfigParser_UnsupportedFormat
--- PASS: TestConfigParser_UnsupportedFormat (0.01s)
PASS
```

## Rubric

| Criterion | Excellent (5) | Good (3) | Needs Work (0) |
|---|---|---|---|
| **Precedence Ordering** | Env vars strictly override file configs, which strictly override defaults | Overrides work but order of precedence is reversed or incomplete | Overrides fail to update nested struct fields |
| **Schema Validation** | Port range, required strings, and positive bounds enforced with clean error messages | Basic validation present with partial checks | Invalid configurations accepted silently |
| **Format Flexibility** | JSON and KV formats supported; format detection based on extension is clean | Only JSON supported | Parsing crashes on valid JSON or malformed files |
| **Testability & Design** | `envLookup` dependency injected allowing deterministic testing without touching process env | Tests mutate global process environment (`os.Setenv`) | Untested override logic |

## NEXT UP

[Module 08: HTTP REST APIs](../../08-http-rest-apis/README.md)
