# Config files

## Learning objective

Load application configuration from JSON and YAML files into Go structs, merge config from multiple sources, and apply file-based configuration patterns used in production Go services.

## Why this matters

Every non-trivial application needs configuration: database URLs, ports, log levels, feature flags. Hard-coding these values makes deployments brittle. Config files provide a standard, versionable, auditable way to parameterize applications across environments. Understanding config file patterns is essential for building deployable Go services.

## Mental model

Think of your Go struct as a control panel with labeled dials and switches. A config file is a set of instructions for setting each control. Loading config is like having an operator read the instructions and adjust the panel accordingly. Multiple config files are like having a stack of instruction sheets — newer sheets override older ones for the same control.

## Core idea

A config struct holds all configurable parameters for an application. A loading function reads a file (JSON, YAML, TOML), unmarshals its contents into the struct, and returns the populated config.

Key patterns:
- **Single source**: Load one config file at startup.
- **Merge with defaults**: Hard-code default values in the struct, then overlay file values.
- **Multi-source merge**: Load default config, then environment-specific config, then local overrides — each layer overwrites the previous.
- **File discovery**: Search well-known paths (`.` `./config`, `/etc/app/`, `$HOME/.config/app/`).

For YAML support, use the `gopkg.in/yaml.v3` package (not in stdlib).

## Under the hood

Config loading is a three-step pipeline:
1. **Discovery**: Find the config file by searching paths or using a `--config` flag.
2. **Read**: `os.ReadFile(path)` reads the file bytes.
3. **Parse**: `json.Unmarshal` or `yaml.Unmarshal` decodes bytes into the config struct.

For merging, a simple approach is:
- Start with a struct pre-populated with defaults.
- Load the first file and unmarshal into the struct (overwrites defaults).
- Load subsequent files and unmarshal into the same struct (overwrites earlier values).

Go's `json.Unmarshal` and `yaml.Unmarshal` do not zero out fields that are missing from the file — they only overwrite fields that appear in the data. This is the key behavior that makes merging work.

## How Go uses it

- Database connection settings: host, port, user, password, database name
- HTTP server config: bind address, TLS cert paths, rate limits, CORS origins
- Logging config: level, format, output path
- Feature flags: boolean toggles for gradual rollouts
- Monitoring: metrics export config, tracing endpoint

## Go example

```go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

type AppConfig struct {
	ServerPort int       `json:"server_port"`
	LogLevel   string    `json:"log_level"`
	Database   DBConfig  `json:"database"`
	Debug      bool      `json:"debug"`
}

func DefaultConfig() AppConfig {
	return AppConfig{
		ServerPort: 8080,
		LogLevel:   "info",
		Database: DBConfig{
			Host: "localhost",
			Port: 5432,
		},
		Debug: false,
	}
}

func LoadConfig(path string) (AppConfig, error) {
	cfg := DefaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func LoadConfigWithOverride(basePath, overridePath string) (AppConfig, error) {
	cfg, err := LoadConfig(basePath)
	if err != nil {
		return cfg, err
	}
	if _, err := os.Stat(overridePath); os.IsNotExist(err) {
		return cfg, nil
	}
	data, err := os.ReadFile(overridePath)
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func main() {
	dir, _ := os.Getwd()
	base := filepath.Join(dir, "config.base.json")
	over := filepath.Join(dir, "config.override.json")

	// Create sample config files
	os.WriteFile(base, []byte(`{"server_port":9090,"database":{"host":"prod.example.com"}}`), 0644)
	os.WriteFile(over, []byte(`{"debug":true,"database":{"port":5433}}`), 0644)
	defer os.Remove(base)
	defer os.Remove(over)

	cfg, err := LoadConfigWithOverride(base, over)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Port: %d\n", cfg.ServerPort)          // 9090 (from base)
	fmt.Printf("LogLevel: %s\n", cfg.LogLevel)         // "info" (default, not overridden)
	fmt.Printf("Debug: %v\n", cfg.Debug)               // true (from override)
	fmt.Printf("DB Host: %s\n", cfg.Database.Host)     // "prod.example.com" (from base)
	fmt.Printf("DB Port: %d\n", cfg.Database.Port)     // 5433 (overridden)
}
```

## Step-by-step execution

1. `DefaultConfig()` creates a struct with `ServerPort: 8080`, `LogLevel: "info"`, etc.
2. `base` file is read: `{"server_port":9090,"database":{"host":"prod.example.com"}}`.
3. `json.Unmarshal` overwrites `ServerPort` to 9090 and `Database.Host` to `"prod.example.com"`. Other fields (`LogLevel`, `Database.Port`) stay at defaults.
4. `overridePath` file is read: `{"debug":true,"database":{"port":5433}}`.
5. `json.Unmarshal` overwrites `Debug` to `true` and `Database.Port` to 5433. `ServerPort` stays at 9090 (not in override).
6. Result: merged config with defaults → base → override layering.

## Common mistakes

- Mistake: Config file not found silently — using an empty path or ignoring the error.
  - Fix: Always check file read errors with `os.IsNotExist` if the file is optional, otherwise fail.

- Mistake: Using `os.ReadFile` on a directory — produces a confusing error.
  - Fix: Validate the path with `os.Stat` and reject directories.

- Mistake: Overwriting the entire struct on merge instead of letting `Unmarshal` overlay.
  - Fix: Unmarshal into the same struct instance (not a new one).

- Mistake: Storing secrets (passwords, API keys) in config files committed to version control.
  - Fix: Use environment variables or secret stores for secrets; keep only non-sensitive defaults in files.

## Debugging walkthrough

Config loads but values are wrong:

```go
cfg := DefaultConfig()
data, _ := os.ReadFile("config.json")
json.Unmarshal(data, &cfg)
fmt.Println(cfg.ServerPort) // 0, not 8080
```

**Investigation**: Check `config.json`:
```json
{
  "server_port": 0
}
```

**Root cause**: The file explicitly sets `server_port` to 0. `json.Unmarshal` overwrites the default. Zero is a valid JSON number.

**Fix**: Remove `"server_port"` from the file, or use pointer types (`*int`) to distinguish "not set" from "set to zero".

## Production notes

- **Version control your base configs** — they document the available parameters and their defaults.
- **Do not commit secrets** — use `.gitignore` for override files containing secrets.
- **Validate config after loading** — even with defaults, check for required fields and valid ranges.
- **Use a config library** like `spf13/viper` for complex multi-source config merging (file + env + flags).
- **Hot reloading**: Some applications watch config files for changes and reload without restart using `fsnotify`.

## Performance implications

- Config loading happens once at startup — performance is rarely a concern.
- File I/O for a small config file takes microseconds. The main cost is JSON/YAML parsing, which is ~1-10 ms for typical configs.
- Unmarshal into a struct is faster than into a `map[string]any` because the types are known at compile time.
- Config merge operations are negligible — a few struct field assignments.

## Practice task

Write a function `LoadAppConfig(paths ...string) (AppConfig, error)` that loads multiple config files in order, merging each into the same struct (defaults first). If a file doesn't exist, skip it (don't error). Returns the merged config. Test with two files where the second overrides specific fields.

## Tests / verification

```bash
go run ./curriculum/modules/07-cli-files-json-config/lessons/16-config-files
go test ./curriculum/modules/07-cli-files-json-config/lessons/16-config-files
```

## Review questions

1. Why does `json.Unmarshal` into an existing struct work as a merge operation?
2. What happens if a config file sets a field to its zero value (e.g., `"port": 0`)?
3. How would you handle optional config files without failing?
4. What's the risk of committing config files with secrets to version control?
5. What is the typical order of precedence for config sources?

## NEXT UP

Environment variables for configuration — reading configuration from the environment using 12-factor app principles.
