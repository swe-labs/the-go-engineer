# CFG.3 Flag Parsing

## Mission

Master command-line arguments. Learn how to use the standard library `flag` package to allow users to override configuration settings at launch time. Understand the role of flags in local development and CLI tools, and how they sit at the top of the **Configuration Precedence** hierarchy.

## Prerequisites

- CFG.2 Configuration Files

## Mental Model

Think of Flag Parsing as **The Steering Wheel**.

1. **Environment Variables & Files**: These are the "Engine Settings" (Static, defined before you get in the car).
2. **Flags**: These are the "Steering Inputs" (Active, defined by the driver right now).
3. **The Priority**: Even if the engine setting says "Go 60 MPH," if the driver turns the wheel (sets a flag like `--port=9000`), the car turns. Flags always represent the **Final Word** because they are the most explicit user input.

## Visual Model

```mermaid
graph TD
    User[User Command] -->|--port 9999| Flag[flag.Parse]
    Flag -->|Override| Config[Final Config Struct]

    subgraph "The Hierarchy"
    Flag > Env[Env Vars]
    Env > File[Config File]
    end
```

## Machine View

- **`flag.String` / `flag.Int`**: Functions that define a flag and return a *pointer* to the value.
- **`flag.Parse()`**: Must be called after all flags are defined and before the values are accessed. It parses the `os.Args[1:]` slice.
- **Usage**: Flags automatically provide a `-h` or `--help` output for your application, making it self-documenting.

## Run Instructions

```bash
# Run with default flags
go run ./10-production/04-configuration/3-flag-parsing

# Run with a custom flag
go run ./10-production/04-configuration/3-flag-parsing -port=7000
```

## Code Walkthrough

### Defining Flags
Shows how to define a flag with a name, a default value, and a help description.

### Parsing and Accessing
Demonstrates calling `flag.Parse()` and dereferencing the pointers to get the actual values.

### The precedence Merge
Shows a production pattern: loading a config file first, then checking if the flag value was changed by the user (not just the default).

## Try It

1. Run the code. Observe the default port.
2. Run with `-port=8888`. Verify the output.
3. Add a new boolean flag `-verbose` and use it to toggle extra log output (Track SL).
4. Discuss: Why are pointers used for flag values instead of direct values?

## In Production
**Don't over-use flags for service config.** Flags are hard to manage at scale in Kubernetes or Cloud platforms. If you have 200 settings, putting 200 flags in a Docker `ENTRYPOINT` is fragile and unreadable. Use flags for **Action Overrides** (e.g., `-migrate-db`) or **Critical Launch Parameters** (e.g., `-config-path`), and use Files or Environment variables for everything else.

## Thinking Questions
1. Why must `flag.Parse()` be called after all flags are defined?
2. What is the difference between `-port 80` and `-port=80`? (Hint: Both work in Go's flag package).
3. How do you handle "Subcommands" (like `git push` vs `git pull`) using the `flag` package?

## Next Step

Next: `CFG.4` -> `10-production/04-configuration/4-twelve-factor-principles`

Open `10-production/04-configuration/4-twelve-factor-principles/README.md` to continue.
