# How the OS manages processes

## Learning objective

Understand and apply How the OS manages processes in the context of professional Go software engineering.

## Why this matters

Operating systems need to run multiple programs concurrently while keeping them isolated. The process abstraction provides memory isolation, fair CPU scheduling, and resource accounting.

## Mental model

A process is a running program with its own private memory space, file descriptors, and execution state. The OS kernel acts as a traffic cop — allocating CPU time, memory, and I/O resources among all running processes fairly.

## Core idea

Without OS-managed processes, every program would need its own CPU scheduler, memory isolation mechanism, and I/O multiplexer — essentially reimplementing the OS. Processes are the fundamental unit of workload isolation.

## Under the hood

On Linux, each process has a task_struct containing PID, parent PID, memory map (mm_struct), file descriptor table (files_struct), signal handlers, and scheduling priority. The scheduler on modern Linux (CFS — Completely Fair Scheduler) uses red-black trees to track process runtime and ensures fairness by selecting the process with the least runtime.

## How Go uses it

Go's runtime manages goroutines as user-space threads multiplexed onto OS threads. Go programs typically use one OS thread per CPU core (GOMAXPROCS) and schedule thousands of goroutines across them. The os/exec package spawns external processes, and os.FindProcess looks them up by PID.

## Go example

```go
package main

import (
	"fmt"
	"os"
	"os/exec"
)

// GetPID returns the current process ID.
func GetPID() int {
	return os.Getpid()
}

// GetPPID returns the parent process ID.
func GetPPID() int {
	return os.Getppid()
}

// RunCommand runs a command and returns its exit code.
// Returns -1 if the command could not be started.
func RunCommand(name string, args ...string) (int, error) {
	cmd := exec.Command(name, args...)
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode(), nil
		}
		return -1, err
	}
	return 0, nil
}

func main() {
	fmt.Printf("PID: %d, PPID: %d\n", GetPID(), GetPPID())

	code, err := RunCommand("go", "version")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run command: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("child exit code: %d\n", code)
}
```

## Step-by-step execution

1. The OS boots, initializing the scheduler and the process table (a kernel data structure tracking all processes).
2. When fork() is called, the kernel copies the parent's PCB, creates a new virtual address space (with copy-on-write optimization), and assigns a new PID.
3. The scheduler maintains run queues — when a CPU core is free, it dequeues the highest-priority ready process and performs a context switch.
4. A context switch saves the current process's registers and program counter, then loads the next process's saved state.
5. When a process calls exit(), the kernel closes its file descriptors, frees its memory, stores the exit code, and transitions it to zombie state until the parent calls wait().

## Common mistakes

| Mistake | Why it happens | Fix |
|---------|---------------|-----|
| Thinking a 'process' is the same as a 'program' or a 'thread' | They are distinct OS abstractions — a program is a file, a process is a running instance, a thread is a lightweight execution unit within a process | Learn the distinction: program (disk), process (memory + execution), thread (schedulable unit within a process) |
| Assuming every running program gets its own CPU core | The OS scheduler time-shares cores among all processes | Understand that processes share CPU cores via time-slicing; use parallelism only when you have multiple cores |
| Believing a terminated process disappears instantly | It leaves a zombie entry until the parent collects its exit code | Always call wait() or handle SIGCHLD to reap child processes |

## Debugging walkthrough

### Scenario: A server starts many child processes but never calls wait(), creating zombie processes that exhaust the process table.

**Cause:** Child processes that terminate need their exit code collected; otherwise they remain as zombies consuming a PID entry.

**Fix:** Properly handle SIGCHLD signals or call wait()/waitpid() to reap child processes, or use double-fork technique.

### Scenario: A program spawns too many processes simultaneously, overwhelming the system and causing an Out-of-Memory (OOM) kill.

**Cause:** Each process consumes memory for its address space; the OS has finite resources and may kill processes under memory pressure.

**Fix:** Limit concurrent processes using worker pools, and set resource limits with ulimit or cgroups.

## Production notes

Web servers fork child processes (Apache), containers are isolated process trees (Docker), and shells spawn every command as a child process. Understanding process management is essential for system programming, deployment, and debugging.

## Performance implications

- Creating a new process is expensive (microseconds to milliseconds) — fork+copy-on-write or CreateProcess requires significant kernel work.
- Context switching between processes is slower than between threads because the TLB (translation lookaside buffer) must be flushed.
- The OOM killer activates when memory is overcommitted — processes can be killed arbitrarily if they consume too much memory.

## Practice task

The learner must write a Go program that spawns a child process via os/exec, waits for it to complete, collects the exit code, and print the parent and child PIDs.

## Tests / verification

```bash
go test ./curriculum/modules/01-computers-terminal-git-web/07-how-the-os-manages-processes/
```

## Review questions

1. Learner must run ps (or Task Manager) and identify which processes belong to the Go toolchain vs. the OS kernel.
2. Learner must explain why a CPU-bound process on a single-core machine appears to run 'at the same time' as other processes.
3. Learner must predict what happens to child processes when the parent process is killed.

## NEXT UP

Lesson 08: [Memory preview: stack vs heap](../08-memory-preview-stack-vs-heap/README.md)
