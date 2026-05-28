# Module 01 Checkpoint Answer Key

Use this after attempting the checkpoint.

## Part 1

1. Source code is editable text. An executable is the built artifact. A process is a running instance of an executable.
2. The OS gives a process memory, process ID, arguments, environment, file descriptors/standard streams, and scheduled CPU time.
3. A running process is already loaded. Source changes require rebuild/restart or another reload mechanism.

## Part 2

4. An absolute path starts from the filesystem root. A relative path starts from the current working directory.
5. Relative paths and command execution depend on where the shell currently is.
6. Run `pwd`, inspect the command path, then use `ls` to verify the file exists.
7. A shell parses the command, resolves the executable, passes arguments/environment, starts the process, and reports the exit code.

## Part 3

8. An environment variable is a string key/value passed to a process at startup.
9. Environment values can appear in logs, shell history, process inspection, and crash/debug output.
10. Exit code `0` usually means success.
11. CI uses non-zero exit codes to stop unsafe builds, tests, deploys, or scripts.

## Part 4

12. A process ID is the OS identifier for a running process.
13. The stack is orderly call-frame storage; the heap is flexible storage for values with less predictable lifetime.
14. Optimization without measurement often fixes the wrong problem or makes code harder to maintain.

## Part 5

15. The working tree is current files, the staging area is the next snapshot, and a commit is the recorded snapshot.
16. `git status` prevents accidentally staging unrelated or secret files.
17. A branch is a movable label pointing to a commit.
18. A merge conflict means Git needs a human to choose how overlapping changes should combine.
19. `fetch` downloads remote references, `pull` fetches and integrates, and `push` sends local commits to the remote.
20. Scope, reason, changed files, validation, risks, and reviewer guidance.

## Part 6

21. A client initiates a request. A server responds. DNS maps names to IP addresses. A port selects a service on a host.
22. An HTTP request contains method, path, headers, and optional body.
23. An HTTP response contains status code, headers, and optional body.
24. Clients and automation use status codes to decide what happened and what to do next.
25. Check DNS/name, IP/network, port/listener, protocol/TLS, server logs, route/handler, and application logic.

## Part 7

26. Good answers identify one weak concept honestly.
27. Good answers name a command, test, diagram, or explanation task that would create evidence.
