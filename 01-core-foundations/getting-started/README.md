# Section 0: Getting Started

Welcome to the **Go Mastery**! This section will get you from zero to running your first Go program. No prior programming experience is required.

## What is Go?

Go (also called Golang) is a programming language created by Google in 2009. It was designed by Robert Griesemer, Rob Pike, and Ken Thompson — the same engineers who built Unix, UTF-8, and the C programming language.

Go is used at Google, Uber, Netflix, Dropbox, Docker, Kubernetes, and thousands of production systems worldwide. It is known for:

- **Simplicity** — Small language, easy to learn, easy to read
- **Speed** — Compiles to native machine code (like C/C++), not interpreted (like Python/JavaScript)
- **Concurrency** — Built-in support for running thousands of tasks simultaneously
- **Reliability** — Strong type system catches bugs at compile time, not at runtime

## Step 1: Install Go

### Linux (Ubuntu/Debian)

```bash
# Option 1: Using snap (easiest, auto-updates)
sudo snap install go --classic

# Option 2: Using apt (may not be the latest version)
sudo apt update && sudo apt install golang-go

# Option 3: Manual install (latest version guaranteed)
# Visit https://go.dev/dl/ and download the Linux tarball
wget https://go.dev/dl/go1.24.2.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.24.2.linux-amd64.tar.gz

# Add Go to your PATH (add this line to ~/.bashrc or ~/.zshrc)
export PATH=$PATH:/usr/local/go/bin
```

### Windows (Native — No WSL)

1. Download the Windows installer from [https://go.dev/dl/](https://go.dev/dl/)
2. Run the `.msi` installer — it will install Go to `C:\Program Files\Go`
3. The installer automatically adds Go to your system `PATH`
4. Open a **new** Command Prompt or PowerShell and verify:

```powershell
go version
```

### Windows (With WSL — Recommended for Developers)

WSL (Windows Subsystem for Linux) gives you a full Linux environment inside Windows. This is the recommended setup for professional Go development.

1. Open PowerShell as Administrator and run:

```powershell
wsl --install
```

1. Restart your computer, then open the Ubuntu terminal
2. Follow the **Linux** instructions above to install Go inside WSL

### macOS

```bash
# Option 1: Using Homebrew (recommended)
brew install go

# Option 2: Download the installer from https://go.dev/dl/
# Download the .pkg file and run it
```

## Step 2: Verify Installation

Open a terminal and run:

```bash
go version
# Expected output: go version go1.24.x linux/amd64 (or similar)
```

If you see a version number, Go is installed correctly. If you see "command not found", revisit the installation steps and make sure Go is in your PATH.

Run this to see your Go environment:

```bash
go env GOPATH GOROOT
```

- **GOROOT** — Where Go itself is installed (you rarely touch this)
- **GOPATH** — Where Go stores downloaded packages (default: `~/go`)

## Step 3: Set Up Your Editor

### Visual Studio Code (Recommended for Beginners)

1. Download from [https://code.visualstudio.com/](https://code.visualstudio.com/)
2. Open VSCode and install the **Go extension**:
   - Press `Ctrl+Shift+X` (or `Cmd+Shift+X` on Mac)
   - Search for "Go" by the Go Team at Google
   - Click **Install**
3. When prompted to install Go tools, click **Install All**

The Go extension provides:

- Auto-completion
- Auto-formatting on save (runs `gofmt` automatically)
- Error highlighting as you type
- Integrated debugging

### GoLand (Professional IDE)

If you prefer a full IDE, JetBrains GoLand is excellent. It has a 30-day free trial and free licenses for students.

## Step 4: Clone This Repository

```bash
git clone https://github.com/rasel9t6/the-go-engineer.git
cd the-go-engineer
```

## Step 5: Run Your First Program

```bash
go run ./00-getting-started/2-hello-world
```

You should see:

```
Hello, World! Welcome to the Go Mastery.
```

🎉 **Congratulations!** You just compiled and executed your first Go program.

## How This Repository Works

This repo uses a single `go.mod` file at the root, making it one Go module with many runnable examples. Each numbered directory is a section, and each sub-directory is a lesson:

```
the-go-engineer/
├── go.mod                    ← Module definition (you don't need to touch this)
├── 00-getting-started/       ← You are here
│   ├── 1-installation/
│   ├── 2-hello-world/        ← Each sub-directory has a main.go you can run
│   ├── 3-how-go-works/
│   └── 4-dev-environment/
├── 01-language-basics/
│   ├── 1-variables/
│   ├── 2-constants/
│   └── ...
└── ...
```

To run any example:

```bash
go run ./SECTION/LESSON
# Example:
go run ./01-language-basics/1-variables
go run ./09-concurrency/3-channels
```

## Common Beginner Mistakes

| Problem | Solution |
| ------- | -------- |
| `command not found: go` | Go is not in your PATH. Re-run the installation steps. |
| `cannot find module` | Make sure you are in the `the-go-engineer/` directory (where `go.mod` lives). |
| `declared and not used` | Go requires you to USE every variable you declare. Remove or use the variable. |
| `imported and not used` | Go requires you to USE every package you import. Remove or use the import. |
| `syntax error: unexpected newline` | Go requires opening braces `{` on the SAME line as the statement. |

## What's Next?

After running the examples in this section, move on to:

➡️ **[01-language-basics](../01-language-basics/)** — Variables, types, constants, and your first real programs.

## Contents

| Directory | Topic | Level |
| --------- | ----- | ----- |
| `1-installation/` | Verify Go installation, print environment info | Beginner |
| `2-hello-world/` | The classic first program, heavily documented | Beginner |
| `3-how-go-works/` | Compilation model, packages, imports explained | Beginner |
| `4-dev-environment/` | Editor setup verification, `go fmt`, `go vet` | Beginner |

## References

- [Official Go Installation Guide](https://go.dev/doc/install)
- [A Tour of Go (Interactive)](https://go.dev/tour/)
- [Go Playground (Run Go in Your Browser)](https://go.dev/play/)
