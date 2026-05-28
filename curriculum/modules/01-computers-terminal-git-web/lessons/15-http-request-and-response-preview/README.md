# Lesson: HTTP request and response preview

## Mission

Preview HTTP as a structured request-response conversation before building APIs later.

By the end of this lesson, you should be able to explain the concept, run the example, and identify the mistake this lesson is designed to prevent.

## Prerequisites

You should have completed the previous lesson:

- core-01-14 (Web preview: client, server, DNS, and ports)

Understand how clients connect to servers via DNS and ports.

## Mental Model

HTTP is a written exchange: the client asks with a method, path, headers, and optional body; the server replies with a status, headers, and body.

The mental model is deliberately simple. Later modules will add details, but this version is enough to stop common beginner confusion.

## Visual Model

```text
learner action
    ↓
machine boundary
    ↓
observable result
    ↓
proof or debugging signal
```

For this lesson, focus on the boundary: what does the learner ask the machine to do, and what does the machine actually receive or produce?

## Machine View

HTTP bytes travel over a network connection. Servers parse the request, route it to handler logic, and serialize a response.

A professional engineer eventually learns to read the system from the machine's point of view. That does not mean memorizing internals immediately. It means asking what object exists, where it lives, who owns it, and what can observe it.

## Run Instructions

From the repository root:

```bash
go run ./curriculum/modules/01-computers-terminal-git-web/lessons/15-http-request-and-response-preview
go test ./curriculum/modules/01-computers-terminal-git-web/lessons/15-http-request-and-response-preview
```

You may also inspect or try these related commands:

- `curl -i`
- `GET /path HTTP/1.1`
- `HTTP/1.1 200 OK`

## Code Walkthrough

Open `main.go`.

The program models the lesson as a small `conceptCard`:

1. `ID` ties the code back to metadata.
2. `Title` names the concept.
3. `MentalModel` gives the human model.
4. `MachineView` gives the operational model.
5. `CommonMistake` names the trap.
6. `Fix` gives the correction.
7. `Commands` lists commands worth recognizing.

The code is intentionally small because this module is about foundations, not language complexity.

## Try It

Identify method, path, status code, header, and body in a simple HTTP exchange.

Then change one line in `main.go`, rerun the program, and explain what changed.

## Common Mistakes

| Mistake | Why it happens | Correction |
|---|---|---|
| Thinking HTTP status codes are only display messages for users. | The visible surface hides an important machine boundary. | Treat status codes as machine-readable signals between clients and servers. |
| Copying a command without knowing its working directory | The shell accepts the command, but paths resolve somewhere else. | Run `pwd` and explain the current directory first. |
| Treating a word as magic vocabulary | Terms like process, branch, port, or request get memorized but not understood. | Define the term using a concrete example. |

## Debugging Signals

Watch for these signals:

- the command fails only in one directory
- output mentions a missing file or unknown path
- a process keeps running after you expected it to stop
- Git says files are modified but you do not know why
- a network command fails before reaching application code

The first debugging move is to write down the exact command, exact output, and current directory.

## In Production

Backend engineers debug HTTP constantly: status codes, headers, timeouts, payloads, auth, caching, and retries.

The professional habit is to connect vocabulary to operational evidence. If you cannot observe it, test it, log it, inspect it, or explain where it lives, the concept is still too magical.

## Performance Notes

This lesson is not about optimizing code. Its performance value is avoiding wasted time. Clear mental models reduce repeated debugging loops.

## Security Notes

Do not paste secrets into terminal commands, screenshots, issue descriptions, or pull requests. Environment variables, Git history, logs, and shell history can all leak sensitive values.

## Thinking Questions

1. What object is this lesson really about?
2. Where does that object live: source file, filesystem, process, Git history, network, or browser?
3. What command or output would prove you understand it?
4. What mistake would show confusion?
5. How will this idea appear again in later Go work?

## Proof of Understanding

You are complete when you can:

- explain `HTTP request and response preview` in your own words
- run the example
- run the test
- describe the common mistake
- explain the fix without reading this page

## Next Step

Continue to:

```text
module-02 — ../../../02-go-setup-tooling/README.md
```
