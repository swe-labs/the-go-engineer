# Module 01 — Computers, Terminal, Git, and the Web

## Mission

This module makes the machine, terminal, Git, GitHub, and web basics feel safe before Go syntax begins.

A learner who skips these foundations often gets stuck later for reasons that are not really about Go: wrong working directory, confusing path, stale process, failed exit code, Git confusion, port mismatch, or misunderstood HTTP request. This module removes those hidden traps early.

## Who this module is for

This module is for:

- beginners who need computing foundations before programming deeply
- learners who can copy terminal commands but do not yet understand what they do
- developers new to Git or GitHub workflow
- anyone who has felt blocked by paths, shells, branches, ports, or HTTP language

## What you will build

You will build a practical mental model of:

- programs and processes
- files, directories, and paths
- terminal commands
- environment variables
- exit codes
- process lifecycle
- memory preview
- Git history
- GitHub collaboration
- client/server web basics
- HTTP request/response vocabulary

The output is not a portfolio project yet. The output is operational confidence.

## Prerequisites

You should have completed:

```text
curriculum/modules/00-orientation/
```

You should be able to:

- open the repository
- find a lesson folder
- run a basic command from the repository root
- explain what zero magic means

## Concept map

```text
program
  ↓
source → executable → process
  ↓
files/directories/paths
  ↓
terminal commands
  ↓
environment + exit code
  ↓
OS process lifecycle + memory preview
  ↓
Git snapshots + branches
  ↓
GitHub collaboration + PR review
  ↓
DNS + ports + HTTP request/response
```

## Lessons

| # | ID | Lesson | Outcome |
|---:|---|---|---|
| 1 | `core-01-01` | [What is a program?](./lessons/01-what-is-a-program/README.md) | Understand a program as a set of instructions that a computer can execute to transform input into output. |
| 2 | `core-01-02` | [Source code, executable, and process](./lessons/02-source-code-executable-and-process/README.md) | Learn the difference between text you edit, the executable artifact produced from it, and the live process running on the machine. |
| 3 | `core-01-03` | [Files, bytes, directories, and paths](./lessons/03-files-bytes-directories-and-paths/README.md) | Understand files as named byte sequences and paths as addresses used to locate them. |
| 4 | `core-01-04` | [Terminal basics](./lessons/04-terminal-basics/README.md) | Become comfortable using the terminal as a text interface for navigating folders, running programs, and reading output. |
| 5 | `core-01-05` | [Environment variables](./lessons/05-environment-variables/README.md) | Understand environment variables as process-level configuration passed from parent process to child process. |
| 6 | `core-01-06` | [Exit codes](./lessons/06-exit-codes/README.md) | Learn how command-line programs report success or failure to shells, scripts, and CI systems. |
| 7 | `core-01-07` | [How the OS manages processes](./lessons/07-how-the-os-manages-processes/README.md) | Preview process lifecycle, scheduling, standard streams, and how the operating system supervises running programs. |
| 8 | `core-01-08` | [Memory preview: stack vs heap](./lessons/08-memory-preview-stack-vs-heap/README.md) | Build a beginner-safe preview of memory as storage that programs use while running. |
| 9 | `core-01-09` | [Git mental model](./lessons/09-git-mental-model/README.md) | Understand Git as a history database for snapshots, not as a mysterious save button. |
| 10 | `core-01-10` | [Git basics: status, add, commit](./lessons/10-git-basics-status-add-commit/README.md) | Learn the safe beginner Git loop: inspect, stage, commit, inspect again. |
| 11 | `core-01-11` | [Branching and merging](./lessons/11-branching-and-merging/README.md) | Understand branches as safe lines of work and merges as ways to combine histories. |
| 12 | `core-01-12` | [GitHub workflow](./lessons/12-github-workflow/README.md) | Learn how local Git work connects to a remote repository and collaboration workflow. |
| 13 | `core-01-13` | [Pull requests and code review](./lessons/13-pull-requests-and-code-review/README.md) | Understand pull requests as structured conversations around a proposed change. |
| 14 | `core-01-14` | [Web preview: client, server, DNS, and ports](./lessons/14-web-preview-client-server-dns-and-ports/README.md) | Preview how a browser or client finds a server and connects to a program listening on a port. |
| 15 | `core-01-15` | [HTTP request and response preview](./lessons/15-http-request-and-response-preview/README.md) | Preview HTTP as a structured request-response conversation before building APIs later. |

## Labs

This module does not need a separate lab because each lesson includes a small runnable proof task.

## Projects

This module has no portfolio project. The concepts are foundations for later portfolio work.

## Assessments

Complete the checkpoint after all lessons:

```text
curriculum/modules/01-computers-terminal-git-web/assessments/checkpoint/
```

## Common failure modes

| Failure | Symptom | Fix |
|---|---|---|
| Wrong working directory | Commands say files do not exist. | Run `pwd`, then navigate to the repository root. |
| Confusing source with process | Edited code does not affect a running program. | Rebuild or rerun the program. |
| Treating Git as a save button | Commits become large and unclear. | Use `status`, `diff`, `add`, `commit` deliberately. |
| Treating ports as URLs | Client cannot connect. | Separate hostname, IP, port, protocol, and path. |
| Ignoring exit codes | Scripts continue after failure. | Check success/failure explicitly. |

## Completion checklist

You are ready for Module 02 when you can:

- [ ] explain source code, executable, and process
- [ ] identify absolute and relative paths
- [ ] run simple terminal commands intentionally
- [ ] explain environment variables and exit codes
- [ ] describe process lifecycle at a beginner level
- [ ] explain Git commits, branches, and merges
- [ ] describe a GitHub pull request workflow
- [ ] explain client/server, DNS, ports, and HTTP request/response
- [ ] pass the Module 01 checkpoint

## Next module

Continue to:

```text
curriculum/modules/02-go-setup-tooling/
```

Module 02 installs and uses Go tooling. Module 01 makes sure the surrounding machine and workflow concepts are already stable.
