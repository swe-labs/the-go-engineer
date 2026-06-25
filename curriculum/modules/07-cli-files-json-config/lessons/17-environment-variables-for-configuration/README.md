# Environment variables for configuration

## Learning objective

Read configuration from environment variables using `os.Getenv`, `os.LookupEnv`, and `os.Environ`, combine environment variables with flags and config files, and apply 12-factor app methodology.

## Why this matters

Environment variables are the standard way to inject configuration into deployed applications — they work identically across development, staging, and production environments, on bare metal, containers, and orchestrators like Kubernetes. The 12-factor app manifesto mandates env vars as the config source of choice because they are language-agnostic, immutable at runtime, and separate cleanly from code.

## Mental model

Think of environment variables as the operating system's bulletin board. Any program can post a notice (set a variable) and any child process can read it. Your Go application is a child of the shell (or container runtime), so it inherits the parent's environment block — a set of `KEY=VALUE` string pairs. Reading `os.Getenv("PORT")` is like walking to the bulletin board and looking for a notice labeled "PORT".

## Core idea

The `os` package provides three functions for environment access:

- `os.Getenv(key string) string` — returns the value, or empty string if not set (can't distinguish "empty" from "unset").
- `os.LookupEnv(key string) (string, bool)` — returns the value and a boolean indicating if the key exists.
- `os.Environ() []string` — returns all environment variables as `"KEY=VALUE"` strings.

The `flag` package supports `flag.String`, `flag.Int`, etc. which can be combined with env vars using a pattern: flag value as default, then override with env var if set.

## Under the hood

On Windows, `os.Getenv` calls `GetEnvironmentVariableW` (kernel32). Environment variables are stored in the Process Environment Block (PEB), a kernel-managed structure created when the process starts. Child processes inherit a copy of the parent's environment block.

Changing environment variables with `os.Setenv` modifies the current process's block and is inherited by child processes, but does NOT affect the parent process (shell) or other running processes. Environment changes are process-local.

## How Go uses it

- **Port configuration**: `os.Getenv("PORT")` — defaults to `"8080"` if empty.
- **Database URLs**: `os.Getenv("DATABASE_URL")` — full connection string.
- **Log levels**: `os.Getenv("LOG_LEVEL")` — `"debug"`, `"info"`, `"warn"`, `"error"`.
- **Feature flags**: `os.Getenv("ENABLE_BETA")` — `"true"` or `"false"`.
- **Secret injection**: API keys, tokens, passwords from environment (or secret stores mapped to env vars).

## Go example

```go
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port    int
	DBURL   string
	Verbose bool
}

func LoadConfig() Config {
	cfg := Config{
		Port:    8080,
		DBURL:   "postgres://localhost:5432/app?sslmode=disable",
		Verbose: false,
	}

	// Override from environment variables
	if v := os.Getenv("PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			cfg.Port = p
		}
	}
	if v := os.Getenv("DATABASE_URL"); v != "" {
		cfg.DBURL = v
	}
	if v, ok := os.LookupEnv("VERBOSE"); ok && v == "true" {
		cfg.Verbose = true
	} else if ok && v == "false" {
		cfg.Verbose = false
	}

	return cfg
}

// Flag + env pattern: flag default can be overridden by env var
func LoadConfigWithFlags() Config {
	port := flag.Int("port", 8080, "server port")
	dbURL := flag.String("db", "postgres://localhost:5432/app?sslmode=disable", "database URL")
	verbose := flag.Bool("verbose", false, "enable verbose logging")
	flag.Parse()

	// Environment variables override flag defaults
	if v := os.Getenv("PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			*port = p
		}
	}
	if v := os.Getenv("DATABASE_URL"); v != "" {
		*dbURL = v
	}
	if v := os.Getenv("VERBOSE"); v != "" {
		*verbose = v == "true"
	}

	return Config{Port: *port, DBURL: *dbURL, Verbose: *verbose}
}

func main() {
	// Show current env config
	cfg := LoadConfig()
	fmt.Printf("Port: %d\n", cfg.Port)
	fmt.Printf("DBURL: %s\n", cfg.DBURL)
	fmt.Printf("Verbose: %v\n", cfg.Verbose)

	// Set an env var for this process and re-read
	os.Setenv("PORT", "3000")
	os.Setenv("DATABASE_URL", "postgres://prod:5432/db?sslmode=require")
	os.Setenv("VERBOSE", "true")

	cfg2 := LoadConfig()
	fmt.Println("\nAfter Setenv:")
	fmt.Printf("Port: %d\n", cfg2.Port)
	fmt.Printf("DBURL: %s\n", cfg2.DBURL)
	fmt.Printf("Verbose: %v\n", cfg2.Verbose)

	// Demonstrate flag-based config (run without flags to see defaults)
	cfg3 := LoadConfigWithFlags()
	fmt.Println("\nFlag-based config:")
	fmt.Printf("Port: %d\n", cfg3.Port)
	fmt.Printf("DBURL: %s\n", cfg3.DBURL)
	fmt.Printf("Verbose: %v\n", cfg3.Verbose)

	// List all environment variables (truncated)
	envs := os.Environ()
	fmt.Printf("\nTotal env vars: %d\n", len(envs))
}
```

## Step-by-step execution

1. `LoadConfig()` initializes `Config` with hard-coded defaults.
2. `os.Getenv("PORT")` reads the `PORT` variable. If empty (not set), the default 8080 remains.
3. `os.Getenv("DATABASE_URL")` reads the URL. If set, it replaces the default.
4. `os.LookupEnv("VERBOSE")` returns value + true/false. We check both `"true"` and `"false"` explicitly.
5. `os.Setenv("PORT", "3000")` changes the environment for the current process.
6. `LoadConfig()` called again now reads `PORT=3000`.
7. `os.Environ()` returns all variables — useful for debugging or passing to child processes.

## Common mistakes

- Mistake: Using `os.Getenv` and treating empty string as "not set".
  - Fix: Use `os.LookupEnv` when you need to distinguish "empty" from "unset".

- Mistake: Forgetting to parse string env vars into proper types (int, bool, duration).
  - Fix: Always parse: `strconv.Atoi`, `strconv.ParseBool`, `time.ParseDuration`.

- Mistake: Expecting `os.Setenv` to persist after the program exits.
  - Fix: It doesn't. Environment changes are per-process and lost when the process exits.

- Mistake: Using env vars for complex structured data (nested config).
  - Fix: Use a config file for complex data; env vars for simple scalars.

## Debugging walkthrough

Config port stays 8080 even though `PORT=9090` is set:

```go
port := os.Getenv("PORT")
// ... later
cfg.Port = 8080
if port != "" {
	cfg.Port = port // port is string, but cfg.Port is int
}
```

**Symptom**: Compile error or wrong value.

**Root cause**: `os.Getenv` returns a string. Assigning a string to an int field requires explicit conversion. Even if it compiled, `port` would be the literal value (string), not parsed.

**Fix**:
```go
if v := os.Getenv("PORT"); v != "" {
	if p, err := strconv.Atoi(v); err == nil {
		cfg.Port = p
	}
}
```

## Production notes

- **12-factor app**: Store config in the environment. One deployable artifact (the binary) promotes across environments — only env vars change.
- **Naming convention**: `APP_COMPONENT_KEY` — e.g., `MYAPP_DB_PORT`, `MYAPP_LOG_LEVEL`. Prefix with the app name to avoid collisions.
- **Secret management**: In production, use a secrets manager (Vault, AWS Secrets Manager) that injects values as env vars into the container.
- **Document required vars**: Provide a `.env.example` file (never `.env` itself) with documented defaults.
- **Validate at startup**: Fail fast if required env vars are missing.

## Performance implications

- `os.Getenv` and `os.LookupEnv` are O(1) — the environment block is a pre-built table in process memory.
- `os.Environ()` allocates a new slice of strings each call — cache the result if called frequently.
- `os.Setenv` is relatively expensive (copies the entire environment block on some platforms) — call only during startup initialization, never in hot paths.

## Practice task

Write a function `ConfigFromEnv() Config` that reads `PORT`, `DATABASE_URL`, and `VERBOSE` from environment variables with the same defaults as above. Write a `main()` that:
1. Sets environment variables using `os.Setenv`.
2. Calls `ConfigFromEnv()`.
3. Prints the config.
4. Also demonstrate `os.LookupEnv` for a missing variable.

## Tests / verification

```bash
go run ./curriculum/modules/07-cli-files-json-config/lessons/17-environment-variables-for-configuration
go test ./curriculum/modules/07-cli-files-json-config/lessons/17-environment-variables-for-configuration
```

## Review questions

1. What's the difference between `os.Getenv` and `os.LookupEnv`?
2. How do you distinguish between "env var set to empty string" and "env var not set"?
3. What does `os.Environ()` return?
4. Does `os.Setenv` affect the parent shell after the program exits?
5. What is the recommended naming convention for environment variables in a Go application?

## NEXT UP

Config validation — validating parsed configuration values before using them.
