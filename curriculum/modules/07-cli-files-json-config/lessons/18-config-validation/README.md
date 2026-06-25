# Config validation

## Learning objective

Validate configuration values at startup by checking required fields, type constraints, and range boundaries, and return comprehensive validation errors that fail fast before the application serves traffic.

## Why this matters

A misconfigured application is a time bomb. It may start, appear healthy, then crash or behave incorrectly when a code path uses the bad config value. Validating configuration at startup — before accepting any requests — catches these problems immediately. Production services that skip config validation inevitably discover missing values during incidents, making recovery harder and slower.

## Mental model

Config validation is like a pre-flight checklist for an airplane. Before takeoff, the pilot verifies fuel level, flaps, instruments, and hydraulics. Each check passes or fails. If any check fails, the plane stays grounded. Your config validation function is the same: before the HTTP server starts, verify that required fields are present, ports are in range, URLs are parseable, and timeouts are positive.

## Core idea

Validation is a function that takes a config struct and returns an error (or multiple errors). It checks:
- **Required fields**: string or pointer fields that must be non-zero.
- **Range checks**: integers within min/max bounds.
- **Format validation**: URLs parseable, emails contain `@`, paths exist.
- **Enum checks**: value belongs to a known set.
- **Cross-field validation**: e.g., `EndPort > StartPort`.

Use `errors.Join` (Go 1.20+) or a custom `ValidationError` type to collect multiple errors so the user sees all problems at once (fail-fast).

## Under the hood

Go has no built-in validation framework. The standard approach is a method on the config struct:

```go
func (c *Config) Validate() error {
    var errs []error
    if c.Port < 1 || c.Port > 65535 {
        errs = append(errs, fmt.Errorf("port %d out of range [1,65535]", c.Port))
    }
    if c.DBURL == "" {
        errs = append(errs, errors.New("database_url is required"))
    }
    return errors.Join(errs...)
}
```

For struct field validation, third-party packages like `go-playground/validator` use struct tags: `validate:"required,min=1,max=65535"`. These use reflection to iterate fields and run validators, similar to how `encoding/json` uses tags for marshaling.

## How Go uses it

- **Server startup**: Validate config before `http.ListenAndServe`.
- **CLI tools**: Validate parsed args and flags before executing.
- **Database migrations**: Validate connection config before running migrations.
- **Feature flags**: Validate flag combinations are legal.
- **Testing**: Validate test config structs to catch invalid test setup.

## Go example

```go
package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
)

type ServerConfig struct {
	Port    int
	Host    string
	DBURL   string
	LogLevel string
	Workers int
}

func (c *ServerConfig) Validate() error {
	var errs []error

	if c.Port < 1 || c.Port > 65535 {
		errs = append(errs, fmt.Errorf("port %d out of range [1, 65535]", c.Port))
	}

	if c.Host == "" {
		errs = append(errs, errors.New("host is required"))
	}

	if c.DBURL == "" {
		errs = append(errs, errors.New("database_url is required"))
	} else {
		if _, err := url.Parse(c.DBURL); err != nil {
			errs = append(errs, fmt.Errorf("database_url is not a valid URL: %w", err))
		}
	}

	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[c.LogLevel] {
		errs = append(errs, fmt.Errorf("log_level %q must be one of: debug, info, warn, error", c.LogLevel))
	}

	if c.Workers < 1 {
		errs = append(errs, fmt.Errorf("workers must be at least 1, got %d", c.Workers))
	}

	return errors.Join(errs...)
}

func LoadAndValidate() error {
	// Simulate loading config from env/file
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	cfg := ServerConfig{
		Port:     port,
		Host:     os.Getenv("HOST"),
		DBURL:    os.Getenv("DATABASE_URL"),
		LogLevel: os.Getenv("LOG_LEVEL"),
		Workers:  0,
	}
	return cfg.Validate()
}

func main() {
	// Set some invalid values
	os.Setenv("PORT", "99999")
	os.Setenv("HOST", "")
	os.Setenv("DATABASE_URL", "not-a-url")
	os.Setenv("LOG_LEVEL", "critical")

	if err := LoadAndValidate(); err != nil {
		fmt.Println("Config validation failed:")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Config is valid")
}
```

## Step-by-step execution

1. `ServerConfig` is populated from environment variables (or config file).
2. `Validate()` is called on the config pointer.
3. `Port` is 99999 → check `c.Port < 1 || c.Port > 65535` → true → append range error.
4. `Host` is `""` → check `c.Host == ""` → true → append required error.
5. `DBURL` is `"not-a-url"` → not empty, so URL validation runs → `url.Parse` fails → append URL format error.
6. `LogLevel` is `"critical"` → check valid levels → not found → append enum error.
7. `Workers` is 0 → check `< 1` → true → append min error.
8. `errors.Join(errs...)` returns a single error containing all 5 errors.
9. Main prints all errors and exits with code 1.

## Common mistakes

- Mistake: Validating only the first error and returning immediately — the user has to fix, re-run, hit the next error, repeat.
  - Fix: Collect all validation errors and return them together (fail-fast for detection, not for reporting).

- Mistake: Not validating URLs, paths, or addresses — accepting any non-empty string.
  - Fix: Use `net/url.Parse`, `os.Stat`, `net.ResolveTCPAddr` to validate formats.

- Mistake: Using `panic` for validation failures instead of returning errors.
  - Fix: Return errors from `Validate()` and let the caller decide how to handle.

- Mistake: Ignoring the validation error and starting the server anyway.
  - Fix: Check `Validate()` before `ListenAndServe` and `os.Exit(1)` on failure.

## Debugging walkthrough

This config passes validation but the server crashes:

```go
type Config struct {
	Port int `validate:"min=1,max=65535"`
}
cfg := Config{Port: 0}
// Port is 0, which fails "min=1", but if not validated...
http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil) // :0 → random port
```

**Symptom**: Server starts on a random port, not the expected one.

**Root cause**: Port 0 is valid for `net/http` (means "pick a random port"), but semantically wrong for the application. The validation should reject 0.

**Fix**: Add explicit check: `if c.Port == 0 { errs = append(errs, errors.New("port is required")) }`, or set a valid default before validation.

## Production notes

- **Fail fast at startup**: Validate before listening on ports or connecting to databases. A startup failure is immediately visible in logs and monitoring.
- **Include field names in errors**: `fmt.Errorf("config.port: must be between 1 and 65535, got %d", c.Port)` makes debugging faster.
- **Use structured validation**: Return a list of `{Field, Message}` pairs so automated tooling can format them.
- **Consider third-party validators**: For large config structs, `go-playground/validator` reduces boilerplate.
- **Validate cross-field constraints**: e.g., `ReadTimeout` <= `WriteTimeout`, `MinPoolSize` <= `MaxPoolSize`.

## Performance implications

- Config validation runs once at startup — performance is irrelevant.
- `errors.Join` allocates a slice of errors. For typical configs (< 50 fields), this is negligible.
- Third-party validators use reflection and are slower than hand-written checks, but again — startup only.
- URL parsing (`net/url.Parse`) allocates, but only for config fields typed as URLs.

## Practice task

Write a function `ValidateConfig(cfg ServerConfig) error` that checks the same rules above. In main(), create a config with at least two validation failures and print the combined error. Then create a valid config and confirm no error is returned.

## Tests / verification

```bash
go run ./curriculum/modules/07-cli-files-json-config/lessons/18-config-validation
go test ./curriculum/modules/07-cli-files-json-config/lessons/18-config-validation
```

## Review questions

1. Why should config validation happen at startup rather than lazily when values are first used?
2. How does `errors.Join` help with validation?
3. What's the difference between a required field check and a range check?
4. Should you use `panic` or return `error` from a validation function?
5. Give an example of a cross-field validation check.

## NEXT UP

Text templates — generating dynamic text output using Go's text/template package.
