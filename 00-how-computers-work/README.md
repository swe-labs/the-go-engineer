# Section 00: How Computers Work

> **Philosophy:** You cannot write great code without understanding what the machine is actually doing. This section builds the mental model that makes later Go code feel explainable instead of magical.

## Mission

Before we teach syntax, we teach the machine.

By the end of this section, the learner should be able to explain:

- what a program really is
- how Go source becomes execution
- how memory is divided and managed
- how the terminal launches programs
- how the OS manages processes
- why cache and I/O shape performance
- why syscalls are the boundary between code and hardware

## Section Map

| ID | Type | Surface | Core Job |
| --- | --- | --- | --- |
| `HC.1` | Lesson | [what-is-a-program](./1-what-is-a-program) | fetch-decode-execute cycle and instruction model |
| `HC.2` | Lesson | [code-to-execution](./2-code-to-execution) | source → tokens → AST → IR → binary |
| `HC.3` | Lesson | [memory-basics](./3-memory-basics) | stack, heap, GC, and escape analysis |
| `HC.4` | Lesson | [terminal-confidence](./4-terminal-confidence) | shell, stdout/stderr, and command flow |
| `HC.5` | Lesson | [os-processes](./5-os-processes) | processes, signals, threads, and file descriptors |
| `HC.6` | Lesson | [cpu-cache-and-performance](./6-cpu-cache-and-performance) | cache hierarchy, locality, and latency |
| `HC.7` | Lesson | [syscalls](./7-syscalls) | user space vs kernel space |
| `HC.8` | Lesson | [blocking-vs-non-blocking-io](./8-blocking-vs-non-blocking-io) | waiting, concurrency, and I/O behavior |

## How To Use This Section

For each lesson:

1. read the `README.md` first
2. run the lesson
3. compare the output with the machine view
4. make one small change in `main.go`
5. rerun and explain what changed

This section is intentionally slow and visual.
Do not rush it just because the programs are small.

## Checkpoint

Before moving to Section 01, you should be able to answer these without looking anything up:

**Core systems:**

- [ ] What is the fetch-decode-execute cycle?
- [ ] What are the six basic CPU operations?
- [ ] What is the difference between compile time and runtime?
- [ ] What is the difference between stack and heap memory?
- [ ] What is escape analysis?
- [ ] What is a process, and how is it different from a program?
- [ ] What are the three default file descriptors?
- [ ] What signal does `Ctrl+C` send?

**Performance and reality:**

- [ ] Why are some programs slow even on fast CPUs?
- [ ] What is the cache hierarchy?
- [ ] Why do cache misses hurt performance?
- [ ] What is a syscall?
- [ ] Why are blocking programs often slow because they wait, not because they compute?

**Terminal fluency:**

- [ ] Run a lesson with `go run`
- [ ] Redirect output to a file
- [ ] Find a process by PID
- [ ] Explain the difference between stdout and stderr

## Next Step

Continue to [GT.1 Installation Verification](../01-getting-started/1-installation).
