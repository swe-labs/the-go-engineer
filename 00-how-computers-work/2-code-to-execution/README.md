# HC.2

## Mission

Understand the difference between source code (text) and machine code (binary), and how compiling bridges the gap.

## Prerequisites

- `HC.1`

## Mental Model

CPUs do not understand English. They do not understand the word `println`. A CPU only understands electrical pulses, represented as numbers (machine code).

To run human-readable code, we need a **Compiler**. The compiler is itself a program that reads your text file and translates it into a binary file containing raw CPU instructions. 

- **Source Code:** For humans to read and write.
- **Compiler:** The translator.
- **Binary / Executable:** For the CPU to run.

## Visual Model

```mermaid
flowchart LR
    A[main.go\n(Human Text)] -->|go build| B(Go Compiler)
    B -->|Translates| C[main.exe\n(Machine Binary)]
    C -->|OS Runs It| D[CPU Executes]
```

## Machine View

When you type `go run main.go`, Go is secretly doing two things:
1. Compiling the code into a temporary binary file.
2. Immediately executing that temporary binary.

If you type `go build main.go`, Go stops after step 1 and hands you the compiled binary file to run yourself.

## Run Instructions

Instead of `go run`, we are going to build the binary manually.
*(Note: If you are on Windows, the output file will be `program.exe`)*

```bash
go build -o program ./00-how-computers-work/2-code-to-execution
./program
```

## Code Walkthrough

This code is identical to the last lesson. The difference is not in the code itself, but in *how* we ask the computer to process it. By using `go build`, we ask the computer to generate a permanent binary translation that we can run without needing the Go compiler ever again.

## Try It

1. Build the program.
2. Notice the new file created in your folder. Try to open it in your text editor. What does it look like? (Hint: It will look like gibberish because text editors expect text, not raw CPU numbers).
3. Delete the binary file when you are done.

## Next Step

Now that we have instructions running on the CPU, where do we store the data those instructions use? Move to [HC.3 Memory Basics](../3-memory-basics/).
