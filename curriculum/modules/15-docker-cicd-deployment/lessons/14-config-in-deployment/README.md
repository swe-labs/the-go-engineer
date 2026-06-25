# Config in deployment

## Learning objective

Load, validate, and manage application configuration from environment variables in a Go service, following 12-factor app principles, and mask sensitive values in log output.

## Why this matters

A Go binary compiled once must run in development, staging, and production without recompilation. Configuration is the mechanism that makes this possible. Every environment has different database URLs, port numbers, log levels, and feature flags. Hardcoding these values creates a fragile deployment that breaks when moving between environments. Professional Go services externalize all configuration to the environment, keeping the binary environment-agnostic and deployment-friendly.

## Mental model

Configuration is the separation between code and the environment it runs in. Think of your Go binary as a universal appliance and configuration as the plug adapter that makes it work in a specific outlet.

- **Code** = what the service does (business logic, routes, handlers).
- **Config** = where and how it does it (port, database, log level).
- **Secrets** = who it authenticates as (passwords, API keys).

The 12-factor app manifesto (factor 3: Config) states: store config in environment variables. Environment variables are language- and OS-agnostic, never accidentally committed to version control, and trivially changed between deployments without code changes.

## Core idea

The 12-factor app approach to configuration:

| Principle | Implementation in Go |
|---|---|
| Store config in environment | `os.Getenv("PORT")` |
| Never hardcode defaults for each env | Defaults for dev only, override for prod |
| Group related config into a struct | `type Config struct { ... }` |
| Validate on startup | Fail fast with meaningful errors |
| Log sanitized config | Mask passwords, don't log secrets |

A config struct aggregates all configuration values. The `LoadConfigFromEnv` function reads environment variables, applies defaults, and returns a parsed `Config` or errors. Validation runs after loading to catch misconfigurations before the server starts.

## Under the hood

`os.Getenv` reads from the process's environment block, which is inherited from the parent process. On Unix, this is a list of `KEY=VALUE` strings in the process's memory. On container orchestration platforms (Kubernetes, Docker Compose), the orchestrator sets environment variables before starting the container.

When you run `os.Setenv("PORT", "8080")`, the Go runtime calls the platform's `setenv` syscall, modifying the process's environment block for the current process and any child processes.

The `strconv` package converts string environment variables to typed Go values. `strconv.Atoi` converts to `int`, `strconv.ParseBool` to `bool`, `strconv.ParseFloat` to `float64`.

## How Go uses it

Every production Go service uses environment-based configuration. Popular Go config libraries include:

- `github.com/kelseyhightower/envconfig` — maps environment variables to struct fields using tags.
- `github.com/spf13/viper` — reads from env, files, and remote config stores with a unified interface.
- `github.com/caarlos0/env` — lightweight struct tag-based parser.

The standard library alone is sufficient for most services. The pattern:

```go
type Config struct {
    Port    int    `env:"PORT" default:"8080"`
    DBURL   string `env:"DATABASE_URL" required:"true"`
}
```

At startup, call `LoadConfigFromEnv()`, validate, and log the result (with secrets masked). If validation fails, exit with a non-zero code immediately.

## Go example

```go
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port           int
	DatabaseURL    string
	LogLevel       string
	MaxConnections int
	Environment    string
	FeatureFlags   []string
}

type ConfigError struct {
	Field string
	Err   error
}

func (e ConfigError) Error() string {
	return fmt.Sprintf("config %s: %v", e.Field, e.Err)
}

func LoadConfigFromEnv() (*Config, []error) {
	var errs []error
	cfg := &Config{
		Port:           8080,
		LogLevel:       "info",
		MaxConnections: 100,
		Environment:    "development",
	}
	if v := os.Getenv("PORT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			errs = append(errs, ConfigError{"PORT", err})
		} else {
			cfg.Port = p
		}
	}
	if v := os.Getenv("DATABASE_URL"); v != "" {
		cfg.DatabaseURL = v
	}
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		cfg.LogLevel = v
	}
	if v := os.Getenv("MAX_CONNECTIONS"); v != "" {
		m, err := strconv.Atoi(v)
		if err != nil {
			errs = append(errs, ConfigError{"MAX_CONNECTIONS", err})
		} else {
			cfg.MaxConnections = m
		}
	}
	if v := os.Getenv("ENVIRONMENT"); v != "" {
		cfg.Environment = v
	}
	if v := os.Getenv("FEATURE_FLAGS"); v != "" {
		cfg.FeatureFlags = strings.Split(v, ",")
	}
	return cfg, errs
}

func ValidateConfig(cfg *Config) []error {
	var errs []error
	if cfg.Port < 1 || cfg.Port > 65535 {
		errs = append(errs, ConfigError{"Port",
			fmt.Errorf("must be between 1 and 65535, got %d", cfg.Port)})
	}
	if cfg.DatabaseURL == "" {
		errs = append(errs, ConfigError{"DatabaseURL",
			fmt.Errorf("database URL is required")})
	}
	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[cfg.LogLevel] {
		errs = append(errs, ConfigError{"LogLevel",
			fmt.Errorf("invalid log level: %s", cfg.LogLevel)})
	}
	if cfg.MaxConnections < 1 {
		errs = append(errs, ConfigError{"MaxConnections",
			fmt.Errorf("must be at least 1")})
	}
	validEnvs := map[string]bool{"development": true, "staging": true, "production": true}
	if !validEnvs[cfg.Environment] {
		errs = append(errs, ConfigError{"Environment",
			fmt.Errorf("invalid environment: %s", cfg.Environment)})
	}
	return errs
}

func main() {
	os.Setenv("PORT", "3000")
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/app")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("MAX_CONNECTIONS", "50")
	os.Setenv("ENVIRONMENT", "development")
	os.Setenv("FEATURE_FLAGS", "new-checkout,dark-mode")

	cfg, errs := LoadConfigFromEnv()
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "config error: %v\n", e)
		}
		os.Exit(1)
	}

	valErrs := ValidateConfig(cfg)
	if len(valErrs) > 0 {
		for _, e := range valErrs {
			fmt.Fprintf(os.Stderr, "validation error: %v\n", e)
		}
		os.Exit(1)
	}
	fmt.Printf("Config loaded: %+v\n", cfg)
}
```

## Step-by-step execution

For `LoadConfigFromEnv` with environment `PORT=3000`, `DATABASE_URL=postgres://...`:

1. Create `Config` with defaults: `Port=8080`, `LogLevel="info"`, `MaxConnections=100`, `Environment="development"`.
2. Read `PORT`: `os.Getenv("PORT")` returns `"3000"`. `strconv.Atoi("3000")` returns `3000, nil`. Set `cfg.Port = 3000`.
3. Read `DATABASE_URL`: returns `"postgres://..."`. Set `cfg.DatabaseURL = "postgres://..."`.
4. Read `LOG_LEVEL`: returns `"debug"`. Set `cfg.LogLevel = "debug"`.
5. Read `MAX_CONNECTIONS`: returns `"50"`. Parse to int, set `cfg.MaxConnections = 50`.
6. Read `ENVIRONMENT`: returns `"development"`. Set `cfg.Environment = "development"`.
7. Read `FEATURE_FLAGS`: returns `"new-checkout,dark-mode"`. Split by comma: `["new-checkout", "dark-mode"]`.
8. Return `cfg` with no errors.

Then `ValidateConfig` checks:
- Port 3000 is in range 1-65535 (passes).
- DatabaseURL is not empty (passes).
- "debug" is a valid log level (passes).
- MaxConnections 50 >= 1 (passes).
- "development" is valid (passes).
- Returns empty errors slice.

## Common mistakes

- **Hardcoding config values in the binary**: Makes the binary environment-specific. Always use environment variables or config files.
- **Not validating on startup**: A missing config value causes a runtime panic when the server starts serving requests, not at process start. Validate everything before the first HTTP request.
- **Logging sensitive config values**: Printing `DATABASE_URL` in startup logs exposes passwords. Mask the password portion before logging.
- **Confusing config with secrets**: Config is non-sensitive (port, log level). Secrets are sensitive (passwords, tokens). Use a secrets manager for secrets, not environment variables directly.
- **Overriding dev defaults in production code**: Keep defaults safe for development. Production should set all values explicitly.

## Debugging walkthrough

Consider a service that fails to connect to the database in production but works locally.

**Symptom**: `connection refused` error on start, but the database is running.

**Investigation**: The `DATABASE_URL` environment variable is set to `localhost:5432` in production, but the database is in a different container or host.

```go
cfg, errs := LoadConfigFromEnv()
fmt.Printf("Connecting to: %s\n", cfg.DatabaseURL)
```

**Root cause**: The production deployment is missing the `DATABASE_URL` env var, so the default (or an empty string) is used. The service tries to connect to `localhost`, which is the container itself, not the database container.

**Fix**: Set `DATABASE_URL=postgres://user:pass@db-container:5432/app?sslmode=disable` in the production environment configuration (Kubernetes ConfigMap, Docker Compose environment, or CI/CD pipeline variable).

Another scenario:

```
config error: config PORT: strconv.Atoi: parsing "8080\n": invalid syntax
```

**Root cause**: The environment variable contains a trailing newline. This often happens when the environment is set from a file with a trailing newline.

**Fix**: Trim the value: `v := strings.TrimSpace(os.Getenv("PORT"))`.

## Production notes

In production Go services:

- **Config struct validation**: Add a `Validate() []error` method to your config struct. Call it before starting the HTTP server or worker.
- **Config from multiple sources**: Use a library like `viper` that reads from env, YAML files, and remote config stores with priority ordering: CLI flags > env vars > config file > defaults.
- **Hot reload**: For config that changes at runtime (feature flags, log levels), use a separate reload endpoint or signal handler rather than restarting the process.
- **Config struct grouping**: Nest related settings: `type DatabaseConfig struct { URL string; MaxConns int }`. Embed in the top-level `Config`.

## Performance implications

- **Reading environment variables is O(1)**: `os.Getenv` is a hash map lookup in the process's environment block. No measurable overhead.
- **Validating config at startup adds < 1ms**: Validation is a one-time cost. It prevents hours of debugging from misconfigured deployments.
- **Config parsing libraries add import overhead**: The `envconfig` library adds ~100KB to the binary. Use the standard library for simple cases.
- **Hot-reload config**: Polling for config changes adds negligible overhead if done every 30-60 seconds. Avoid polling faster than necessary.

## Practice task

Write a function `LoadConfigFromFile(path string) (*Config, error)` that reads a simple `KEY=VALUE` line-delimited file and populates a `Config` struct. Skip lines starting with `#` (comments) and blank lines. Return an error for duplicate keys or invalid values. Then write a `main()` that:

1. Creates a temp `.env` file with valid config.
2. Loads the config from the file.
3. Validates the config.
4. Prints the loaded config (with masked database URL).

## Tests / verification

```bash
go run ./curriculum/modules/15-docker-cicd-deployment/lessons/14-config-in-deployment
go test ./curriculum/modules/15-docker-cicd-deployment/lessons/14-config-in-deployment
```

The existing tests verify default values, custom environment variable loading, invalid port handling, config validation with all error combinations, password masking in config strings, and DSN masking. After completing the practice task, add tests for your `LoadConfigFromFile` function covering valid files, files with comments, missing files, and duplicate keys.

## Review questions

1. According to the 12-factor app, where should configuration be stored? Why not in the code?
2. What is the difference between config and secrets? Give two examples of each.
3. Why should you validate configuration at startup rather than lazily when the value is first used?
4. What happens when `os.Getenv("UNDEFINED_VAR")` is called? How do you distinguish an empty value from a missing variable?
5. How would you support multiple config sources (env vars, YAML file, CLI flags) with a priority order in Go?

## NEXT UP

Secrets in deployment: managing sensitive configuration like API keys, database passwords, and TLS certificates.
