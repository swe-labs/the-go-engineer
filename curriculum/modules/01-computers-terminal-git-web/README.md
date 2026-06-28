# Module 01: Computers, Terminal, Git, and the Web

## Learning goal

Make the machine, shell, Git, GitHub, and web basics feel safe before Go syntax.

## Summary

This module builds the foundational mental models every software engineer needs before writing a single line of Go. You will learn what a program actually is — from source code to executable to running process — and how the operating system manages files, processes, and memory. These concepts remove the magic from compilation, debugging, and deployment.

With that foundation in place, you will gain terminal fluency: navigating the filesystem, running commands, piping data, and reading exit codes. You will then master Git — not just memorizing commands but understanding its object model, branching mechanics, and the GitHub pull request workflow that powers modern collaboration. Finally, you will get a preview of the web: how DNS resolves domain names, how TCP connections carry HTTP requests and responses, and how clients and servers communicate across the internet.

## Prerequisites

- Module 00 — Orientation

## Lessons

| # | Lesson | Outcome |
|---|--------|---------|
| 1 | [What is a program?](../01-computers-terminal-git-web/01-what-is-a-program/README.md) | Distinguish source code, compiled binary, and running process; explain the three-part lifecycle of any program. |
| 2 | [Source code, executable, and process](../01-computers-terminal-git-web/02-source-code-executable-and-process/README.md) | Trace the full lifecycle from editing source to terminating a process; understand why editing source does not change a running binary. |
| 3 | [Files, bytes, directories, and paths](../01-computers-terminal-git-web/03-files-bytes-directories-and-paths/README.md) | Navigate the filesystem tree, understand how bytes map to content, and write cross-platform Go file operations. |
| 4 | [Terminal basics](../01-computers-terminal-git-web/04-terminal-basics/README.md) | Run commands, redirect output, pipe data between programs, and explain how the shell forks and manages child processes. |
| 5 | [Environment variables](../01-computers-terminal-git-web/05-environment-variables/README.md) | Read and set environment variables, understand inheritance across processes, and configure Go applications without hard-coded secrets. |
| 6 | [Exit codes](../01-computers-terminal-git-web/06-exit-codes/README.md) | Write programs that signal success or failure via exit codes; chain commands with `&&` and `||` based on exit status. |
| 7 | [How the OS manages processes](../01-computers-terminal-git-web/07-how-the-os-manages-processes/README.md) | Spawn child processes, collect exit codes, and explain scheduling, context switching, and zombie processes. |
| 8 | [Memory preview: stack vs heap](../01-computers-terminal-git-web/08-memory-preview-stack-vs-heap/README.md) | Identify stack- vs heap-allocated variables; use `go build -gcflags=-m` to verify escape analysis decisions. |
| 9 | [Git mental model](../01-computers-terminal-git-web/09-git-mental-model/README.md) | Explain Git's content-addressable object model, the DAG of commits, and why branches are just movable pointers. |
| 10 | [Git basics: status, add, commit](../01-computers-terminal-git-web/10-git-basics-status-add-commit/README.md) | Navigate the working directory, staging area, and commit history; write meaningful commit messages. |
| 11 | [Branching and merging](../01-computers-terminal-git-web/11-branching-and-merging/README.md) | Create branches, merge with fast-forward and three-way merges, and resolve merge conflicts. |
| 12 | [GitHub workflow](../01-computers-terminal-git-web/12-github-workflow/README.md) | Fork, clone, branch, push, and open pull requests; explain the propose-verify-merge cycle. |
| 13 | [Pull requests and code review](../01-computers-terminal-git-web/13-pull-requests-and-code-review/README.md) | Write effective PR descriptions, leave constructive review comments, and respond to feedback professionally. |
| 14 | [Web preview: client, server, DNS, and ports](../01-computers-terminal-git-web/14-web-preview-client-server-dns-and-ports/README.md) | Resolve domain names to IPs, establish TCP connections, and explain how a browser loads a webpage. |
| 15 | [HTTP request and response preview](../01-computers-terminal-git-web/15-http-request-and-response-preview/README.md) | Start an HTTP server, send GET requests, read status codes and headers, and understand the request-response cycle. |

## Project

[**Shell + Git Lab**](./projects/shell-git-lab/README.md) — A hands-on lab where you practice terminal navigation, filesystem operations, and the full Git workflow: initialize a repository, create branches, make commits, merge changes, and resolve a simulated merge conflict. This project reinforces every lesson in the module and produces a terminal + Git workflow you can demonstrate.

## Assessments

[**Shell + Git Lab Assessment**](./assessments/shell-git-lab/README.md) — A 30-minute assessment covering shell commands, Git fundamentals, and web concepts from the module. Includes multiple-choice and short-answer questions with a detailed rubric.

## Estimated time

~750 minutes (15 lessons × 45 min + 60 min project + 30 min assessment)

## NEXT UP

Complete the [Shell + Git Lab](./projects/shell-git-lab/README.md) project, then take the [Module 01 assessment](./assessments/shell-git-lab/README.md) before moving on.
