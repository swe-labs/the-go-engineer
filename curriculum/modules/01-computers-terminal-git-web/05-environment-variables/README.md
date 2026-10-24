# Environment variables

## Learning objective

Read, set, and inherit environment variables across parent and child processes in POSIX and Windows environments.

## Why this matters

Applications need configuration that differs between development, testing, staging, and production environments without changing source code. Environment variables provide a universal, language-agnostic mechanism for injecting runtime configuration, secrets, and system paths into applications.

## Mental model

Environment variables form a string-to-string dictionary (`KEY=VALUE`) stored by the OS kernel in the private address space of every process:
- **Process Scope:** Every running process maintains its own independent environment block in memory.
- **One-Way Inheritance:** When a parent process spawns a child process, the operating system kernel copies a snapshot of the parent's environment block into the child's address space.
- **Strict Isolation:** Changes made to environment variables by the child process cannot alter the parent process or sibling processes. Modifying a variable in a running application does not change the parent shell.
- **The `$PATH` Lookup Variable:** A specialized, system-critical environment variable containing a delimited list of directory paths that the shell searches sequentially from left to right when resolving command names into executable binaries.

## Visual model

```mermaid
flowchart TD
    subgraph EnvironmentInheritance["1. One-Way Environment Inheritance"]
        Parent["Parent Shell Process (bash / pwsh)
ENV: PATH=/usr/local/go/bin:/usr/bin
DB_PORT=5432"]
        Child["Child Process (Go Runtime)
ENV: Inherited snapshot of parent
Modifications here NEVER affect parent"]
        Parent -->|fork + execve (Clones memory block)| Child
    end

    subgraph PathLookupEngine["2. The $PATH Lookup Engine (Resolving 'go')"]
        CMD["User types: 'go version'"] --> DIR1{"Search 1st: /usr/local/bin"}
        DIR1 -->|Not found| DIR2{"Search 2nd: /usr/local/go/bin"}
        DIR2 -->|MATCH FOUND!| EXEC["Execute: /usr/local/go/bin/go"]
        DIR2 -.->|If not found| DIR3["Search 3rd: /usr/bin"]
        DIR3 -.->|If all fail| ERR["command not found: go"]
    end
```

## Core idea

Without environment variables, applications would need hard-coded configuration files baked into images or compile-time constants that require rebuilding the binary for every deployment target. Environment variables provide dynamic, secure, process-scoped configuration that works identically on laptops, bare-metal servers, Docker containers, and Kubernetes pods.

## Under the hood

### 1. The Environment Memory Block
- **Linux / macOS:** The environment is an array of pointers to null-terminated strings (`environ`) placed directly above `argv` at the very top of the process's call stack. When `execve(path, argv, envp)` executes, the kernel copies this block into the new process's address space.
- **Windows:** The OS kernel maintains the environment block within the process's **Process Environment Block (PEB)**.

### 2. The `$PATH` Resolution Engine
When you type a command like `go`, the shell does not search the entire hard drive. It splits the `$PATH` environment variable by its OS delimiter:
- **POSIX (Linux / macOS):** Colon `:` (e.g., `/usr/local/go/bin:/usr/local/bin:/usr/bin`)
- **Windows:** Semicolon `;` (e.g., `C:\Program Files\Go\bin;C:\Windows\System32`)

The shell tests each directory in `$PATH` sequentially from **left to right**:
1. It checks if an executable file named `go` (or `go.exe` on Windows) exists in the first directory.
2. If found, it executes that binary immediately and stops searching (**first match wins**).
3. If not found, it proceeds to the next directory in the list.
4. If all directories are exhausted without a match, the shell returns:
   ```bash
   bash: go: command not found
   ```

### 3. Binary Shadowing & Precedence
Because `$PATH` searches left to right, placing a directory earlier in `$PATH` **shadows** any binary with the same name located in subsequent directories. For example, if you place `/home/user/custom/bin` before `/usr/bin`, any custom `node` or `go` binary in that directory will run instead of the official system installation.

### 4. Why `./app` is Required for Local Binaries
On POSIX operating systems, the current directory (`.`) is intentionally **excluded** from `$PATH` for security. If `.` were in `$PATH`, an attacker could place a malicious executable named `ls` inside a public folder; when you typed `ls`, the shell would run the attacker's malware instead of `/usr/bin/ls`. Forcing the explicit relative prefix `./` ensures you intentionally execute a local binary.

### 5. Persistence Across Shell Sessions
Variables set in an interactive terminal exist only for the lifespan of that specific shell process. To make environment variables permanent:
- **Bash / Zsh:** Add `export KEY="value"` to your user startup file (`~/.bashrc` or `~/.zshrc`).
- **Windows PowerShell:** Use `[System.Environment]::SetEnvironmentVariable('KEY', 'value', 'User')` or configure via Windows System Properties > Environment Variables.

## How Go uses it

Go uses environment variables extensively to configure the compiler, toolchain, and runtime:
- `GOROOT`: The directory where the Go SDK is installed (e.g., `/usr/local/go` or `C:\Program Files\Go`).
- `GOPATH`: The workspace directory for third-party module caches (defaults to `~/go`).
- `GOOS` & `GOARCH`: Target operating system (`linux`, `darwin`, `windows`) and architecture (`amd64`, `arm64`) used for cross-compiling.
- `GOPROXY`: The URL of the module proxy server used to download dependencies (default: `https://proxy.golang.org,direct`).
- In Go code, use `os.Getenv(key)` for basic retrieval, or `os.LookupEnv(key)` to distinguish between an unset variable and a variable explicitly set to an empty string.

## Go example

The example program demonstrates querying specific environment variables, distinguishing unset vs empty values, inspecting the system `$PATH`, and resolving binary paths using `exec.LookPath`:

```go
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// EnvVar holds key, value, and presence status for an environment variable.
type EnvVar struct {
	Key   string
	Value string
	IsSet bool
}

// GetEnv retrieves an environment variable distinguishing between empty and unset.
func GetEnv(key string) EnvVar {
	val, exists := os.LookupEnv(key)
	return EnvVar{
		Key:   key,
		Value: val,
		IsSet: exists,
	}
}

func main() {
	fmt.Println("=== Checking Critical Environment Variables ===")
	targetKeys := []string{"PATH", "HOME", "USER", "GOROOT", "GOPATH", "MY_CONFIG_SECRET"}
	for _, key := range targetKeys {
		v := GetEnv(key)
		if v.IsSet {
			displayVal := v.Value
			if len(displayVal) > 40 {
				displayVal = displayVal[:37] + "..."
			}
			fmt.Printf("✓ %-18s = %q\n", v.Key, displayVal)
		} else {
			fmt.Printf("✗ %-18s = (unset)\n", v.Key)
		}
	}

	fmt.Println("\n=== Deconstructing System $PATH Directories ===")
	pathVar := os.Getenv("PATH")
	delimiter := string(os.PathListSeparator) // ":" on Unix, ";" on Windows
	dirs := filepath.SplitList(pathVar)

	fmt.Printf("OS Path Delimiter: %q | Total PATH Directories: %d\n", delimiter, len(dirs))
	limit := 5
	if len(dirs) < limit {
		limit = len(dirs)
	}
	fmt.Printf("First %d directories searched by shell (in priority order):\n", limit)
	for i := 0; i < limit; i++ {
		fmt.Printf("  [%d] %s\n", i+1, dirs[i])
	}

	fmt.Println("\n=== Binary Resolution via $PATH (exec.LookPath) ===")
	lookupBinary := "go"
	if runtime.GOOS == "windows" {
		lookupBinary = "go.exe"
	}

	resolvedPath, err := exec.LookPath(lookupBinary)
	if err == nil {
		fmt.Printf("✓ Found '%s' on system PATH: %s\n", lookupBinary, resolvedPath)
	} else {
		fmt.Printf("✗ Binary '%s' was NOT found in any $PATH directory!\n", lookupBinary)
	}
}
```

## Step-by-step execution

1. The OS kernel boots and initializes system-wide environment variables from configuration files (`/etc/environment` on Linux, Registry on Windows).
2. The user logs in; the shell reads user startup files (`~/.bashrc`, `~/.zshrc`) and exports user variables (`export PATH="$HOME/go/bin:$PATH"`).
3. The shell executes in memory with its populated environment block.
4. When the user sets a local variable (`FOO="bar"`), it belongs only to the current shell.
5. When the user runs `export FOO="bar"`, the shell marks the variable so that future child processes will inherit a copy.
6. When a command like `go run main.go` executes, the kernel clones the parent shell's environment block into the child's address space.
7. Go's runtime accesses this memory block via `os.Getenv` or `os.LookupEnv`.
8. The child process exits; any modifications it made to its environment vanish with its memory space, leaving the parent shell completely unaffected.

## Common mistakes

| Mistake | Why it happens | Fix |
|---|---|---|
| Hardcoding secrets (API keys, DB passwords) in source | Convenience during development | Always read secrets from environment variables; never commit credentials to Git |
| Forgetting `export` in Bash | Typing `PORT=8080` makes it a shell-local variable, not an exported environment variable | Use `export PORT=8080` so child processes inherit the variable |
| Expecting child process changes to modify the parent shell | Forgetting one-way inheritance and process isolation | Run scripts with `source script.sh` or `. script.sh` if you intend to modify the current shell |
| Misunderstanding `$PATH` delimiters | Using `:` on Windows or `;` on Linux | Use `os.PathListSeparator` or `filepath.SplitList` in Go for cross-platform parsing |

## Debugging walkthrough

**Scenario: A developer installs Go, opens a terminal, runs `go version`, but gets `bash: go: command not found`.**
- **Cause:** The directory containing the `go` binary (e.g., `/usr/local/go/bin` or `C:\Program Files\Go\bin`) is not listed in the `$PATH` environment variable.
- **Diagnosis:** Run `echo $PATH` (POSIX) or `$env:PATH` (PowerShell) and check if the Go installation `bin` folder is present in the output.
- **Resolution:** Add the Go binary path to `$PATH`:
  - Linux/macOS: Add `export PATH="/usr/local/go/bin:$PATH"` to `~/.bashrc` and run `source ~/.bashrc`.
  - Windows: Add `C:\Program Files\Go\bin` to the user or system `Path` variable via System Environment Variables.

**Scenario: A developer runs a script that sets `export ENV=production`, but running `echo $ENV` in another terminal window shows empty.**
- **Cause:** Environment variables are strictly process-scoped. Setting a variable in terminal A does not alter terminal B.
- **Diagnosis:** Understand that terminal windows are separate, independent processes.
- **Resolution:** Add the export statement to your persistent shell configuration file (`~/.bashrc` or `~/.zshrc`), then open a new terminal window.

## Production notes

The **Twelve-Factor App methodology** dictates that configuration must be strictly separated from code and injected via environment variables. In containerized environments (Docker, Kubernetes), container images are immutable. The exact same Docker image runs in local dev, staging, and production; only the injected environment variables (`DATABASE_URL`, `JWT_SECRET`, `LOG_LEVEL`) change between environments.

## Performance implications

- Reading an environment variable with `os.Getenv()` is an in-memory string search (~100 ns) — zero network or disk overhead.
- Huge environment blocks (megabytes of data) increase memory duplication cost during process fork operations.
- Long `$PATH` lists with non-existent network mounted directories can introduce noticeable command lookup latency in interactive shells.

## Tests / verification

1. Open your terminal and set a test environment variable:
   - Bash: `export TEST_TOKEN="supersecret123"`
   - PowerShell: `$env:TEST_TOKEN = "supersecret123"`
2. Verify the variable is set in the current shell:
   - Bash: `echo $TEST_TOKEN`
   - PowerShell: `echo $env:TEST_TOKEN`
3. Launch a child shell and verify inheritance:
   - Bash: `bash -c 'echo "Child sees: $TEST_TOKEN"'`
   - PowerShell: `powershell -Command 'Write-Host "Child sees: $env:TEST_TOKEN"'`
4. Open a completely new terminal window and verify isolation:
   - Notice the variable is unset in the new window, proving process isolation.

## Review questions

1. What does "one-way environment inheritance" mean when a parent shell spawns a child process?
2. If two different directories in `$PATH` both contain an executable named `app`, which one will the shell execute and why?
3. Why is the current directory (`.`) excluded from `$PATH` by default on POSIX operating systems?
4. In Go, what is the crucial behavioral difference between `os.Getenv("VAR")` and `os.LookupEnv("VAR")`?
5. How does the Twelve-Factor App methodology utilize environment variables in production cloud deployments?

## NEXT UP

[Lesson 06: Exit codes](../06-exit-codes/README.md)
