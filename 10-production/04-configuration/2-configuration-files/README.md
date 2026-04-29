# CFG.2 Configuration Files

## Mission

Master hierarchical configuration. Learn how to use **YAML** or **JSON** files to manage complex application settings that are too large or nested for environment variables. Understand the pattern of **Precedence**: how to combine file-based defaults with environment variable overrides to create a flexible, production-ready configuration system.

## Prerequisites

- CFG.1 Environment Variables
- Section 05: Packages and I/O (Basic file reading and JSON encoding)

## Mental Model

Think of Configuration Files as **The Owner's Manual**.

1. **Environment Variables**: Simple "Switches" (On/Off, Port Number).
2. **Configuration Files**: Detailed "Schematics" (Database Connection Pool settings, Retry Policies, Feature Flags).
3. **The Hybrid**: You use a file for the 50 settings that stay the same in every production cluster, and you use environment variables to override the 2 settings that are unique to *this* cluster (like the Database Password).

## Visual Model

```mermaid
graph TD
    File[config.yaml] -->|Load Defaults| App[Go Binary]
    Env[Environment Vars] -->|Override| App
    App -->|Final Result| ValidatedConfig[In-Memory Config Struct]

    subgraph "Precedence Rule"
    Env > File
    end
```

## Machine View

- **`yaml.Unmarshal` / `json.Unmarshal`**: The primary tools for turning a file on disk into a Go struct.
- **Precedence**: A common pattern is:
  1. Load hardcoded defaults.
  2. Overwrite with values from a config file.
  3. Overwrite with values from environment variables.
  4. Overwrite with command-line flags.
- **Immutability**: Once the configuration is loaded and validated on startup, it should be treated as **Read-Only** throughout the life of the application.

## Run Instructions

```bash
# Run the example to see how YAML is parsed into a Go struct
go run ./10-production/04-configuration/2-configuration-files
```

## Code Walkthrough

### The Config Struct
Shows how to define a Go struct with tags (`yaml:"port"`, `json:"port"`) that match the file format.

### Parsing Logic
Demonstrates reading a file from the disk and unmarshaling it into the struct.

### The Override Pattern
Shows how to manually check for an environment variable and overwrite a field in the struct if it exists.

## Try It

1. Look at `main.go`. Create a new `config.yaml` file and change the port number. Run the code.
2. Add a nested section `database` to the YAML and update the Go struct to support it.
3. Discuss: Why should you never commit a `config.yaml` that contains a real production password to Git? (See SEC.9).

## In Production
**Beware of hidden defaults.** If your YAML parser ignores missing fields, your application might start up with a `0` or `""` value without you realizing it. Use a library that supports **Required Fields** or perform manual validation (CFG.5) immediately after parsing. Always provide an `example.config.yaml` in your repository so other developers know what keys are available.

## Thinking Questions
1. When is a config file better than an environment variable?
2. What is "Precedence," and why does it matter?
3. Why is YAML more common than JSON for configuration files?

## Next Step

Files and Environment variables handle most cases, but sometimes you need to override a setting quickly from the command line. Continue to [CFG.3 Flag Parsing](../3-flag-parsing).
