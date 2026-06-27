# Environment variables

## Learning objective

Understand and apply Environment variables in the context of professional Go software engineering.

## Why this matters

Applications need configuration that differs between development, testing, and production without modifying code. Environment variables provide a universal, language-agnostic mechanism for injecting configuration at runtime.

## Mental model

The Go toolchain includes build, test, fmt, vet, mod, and run. Each tool serves a specific purpose. Learning them replaces the need for complex build systems or task runners.

## Core idea

Without environment variables, configuration would require config files that must be managed across machines, or compile-time constants that require rebuilding. Environment variables provide simple, process-scoped configuration that works everywhere.

## Under the hood

On Linux, the environment block is stored at the top of the process's stack as an array of null-terminated strings: `environ`. The `execve` syscall receives the environment as a parameter. On Windows, each process has an environment block in its PEB (Process Environment Block), inherited from the parent and stored in kernel memory.

## How Go uses it

Go reads environment variables for toolchain configuration: `GOPATH` (module cache location), `GOROOT` (Go installation root), `GOOS`/`GOARCH` (cross-compilation target), `GOPROXY` (module proxy), and `GONOSUMCHECK` (checksum database bypass). Go programs read config via `os.Getenv()` or the `os.Environ()` slice.

## Go example

The example defines an `EnvVar` struct that holds a key, value, and whether the variable was set. The `GetEnv` function uses `os.LookupEnv` to distinguish between unset variables and variables set to empty string. `main()` checks several common variables (`PATH`, `HOME`, `USER`, `GOPATH`, `GOROOT`, `MY_CUSTOM_VAR`), then lists all environment variables sorted alphabetically using `os.Environ()`.

## Step-by-step execution

1. The OS kernel maintains an environment block for each process: a list of `KEY=VALUE` strings.
2. When a new process is spawned (fork+exec), the parent's environment block is copied into the child's address space.
3. The child process reads variables via language-specific APIs — `os.Getenv` in Go, `getenv` in C, `process.env` in Node.
4. The shell manages environment inheritance: `export VAR=value` marks a variable for export to child processes.
5. Shell initialization files (`~/.bashrc`, `/etc/profile`) pre-populate the environment at login.
6. `GetEnv` uses `os.LookupEnv` which returns the value and a boolean indicating if the variable exists.
7. `os.Environ()` returns all variables as `KEY=VALUE` strings, which the program sorts and prints.

## Common mistakes

| Mistake | Why it happens | Fix |
|---|---|---|
| Hard-coding sensitive values like API keys in source code | Convenience during development leads to accidental commits | Use environment variables for all secrets; never commit `.env` files |
| Assuming environment variables are available in all contexts | Child processes inherit them, but system services or IDEs may not | Check that the variable is set before using it; provide clear error messages |
| Modifying env vars in one terminal and expecting persistence | Variable changes only affect the current shell and its children | Set variables in shell init files (`~/.bashrc`, `~/.zshrc`) for persistence |

## Debugging walkthrough

**Scenario: A developer sets an environment variable in a terminal, but the application running in that terminal does not see it.**
- **Cause:** The variable was set after the application started, or the application was launched in a subshell with a cleared environment.
- **Diagnosis:** Print the variable from within the app using `os.Getenv` and check if it returns empty. Verify with `echo $VAR` in the same terminal.
- **Resolution:** Set the variable before launching the application, or use `export VAR=value` and then start the app in the same shell.

**Scenario: A deployment script fails because a required environment variable is missing in production.**
- **Cause:** The variable is set in the developer's local shell config but not in the production deployment configuration.
- **Diagnosis:** Check the deployment logs for "not set" errors. Compare local and production environment variable lists.
- **Resolution:** Use a `.env` file or deployment secrets manager to explicitly declare all required environment variables for each environment.

## Production notes

Environment variables are the standard mechanism for configuring Twelve-Factor Apps — database URLs, API keys, log levels, feature flags, and deployment-specific settings. Docker, Kubernetes, Heroku, and CI/CD systems all use environment variables for configuration injection.

## Performance implications

- Reading an environment variable is a fast in-memory lookup (~100ns), not a disk or network operation.
- Setting many environment variables (1000+) slightly increases process startup time due to the memory copy at fork.
- Environment variables are limited in size — typical max is 32KB on Windows, 2MB on Linux for the entire environment block.

## Practice task

Write a Go program that reads a `DATABASE_URL` environment variable, prints a connection string, and exits with a clear error if the variable is not set.

## Tests / verification

```bash
go test ./curriculum/modules/01-computers-terminal-git-web/05-environment-variables/
```

## Review questions

1. Set an environment variable in a terminal, run a Go program that reads it with `os.Getenv`, and confirm the value is passed correctly.
2. Why do changes to environment variables in a parent shell not affect already-running child processes?
3. List three common environment variables (`PATH`, `HOME`, `GOOS`) and describe what each controls.
4. What is the difference between `os.Getenv` and `os.LookupEnv` in Go?
5. How does `os.Environ()` differ from reading individual variables?

## NEXT UP

[Lesson 06: Exit codes](../06-exit-codes/README.md)
