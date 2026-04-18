# HC.4

## Mission

Demystify the terminal by understanding the three default data streams every program is born with.

## Prerequisites

- `HC.3`

## Mental Model

Every time a program starts, the Operating System automatically attaches three invisible "pipes" to it. These pipes allow the program to communicate with the outside world (the terminal).

1. **Standard Input (`stdin`)**: A pipe where data flows *into* the program (e.g., when you type on your keyboard).
2. **Standard Output (`stdout`)**: A pipe where normal data flows *out* of the program onto your screen.
3. **Standard Error (`stderr`)**: A pipe where error messages flow *out* of the program onto your screen.

When you use `println`, you are literally just pushing text down the `stdout` pipe. The terminal is just a window that displays whatever comes out of those pipes.

## Visual Model

```mermaid
flowchart LR
    Keyboard -->|stdin| P(Your Program)
    P -->|stdout| TerminalScreen[Terminal Screen\n(Normal Text)]
    P -->|stderr| ErrorScreen[Terminal Screen\n(Red/Error Text)]
```

## Machine View

The OS treats these three streams exactly like files. In Linux/Mac, they are literally mapped to file descriptors `0`, `1`, and `2`. Writing to the screen is mechanically identical to writing text to a file on your hard drive.

### The Working Directory
There is one more crucial concept for the terminal: **The Working Directory**. 
When you open a terminal, you are "standing" inside a specific folder on your hard drive. This is your Working Directory.

When we tell the terminal to `go run ./00-how-computers-work/...`, the dot (`.`) means "start from where I am standing right now" (the root of the project). This is called a **relative path**. If you are standing in the wrong folder, the terminal will have no idea what you are talking about.

## Run Instructions

```bash
go run ./00-how-computers-work/4-terminal-confidence
```

## Code Walkthrough

In this code:
1. We ask you to type something.
2. We read data flowing in from the `stdin` pipe until you press Enter.
3. We take that data and push it back out through the `stdout` pipe.

## Try It

1. Run the program. Notice that the program pauses and waits. It is actively listening to the `stdin` pipe.
2. Type your name and press Enter. The data is pulled through the pipe, processed by your CPU, and pushed back out to `stdout`.

## Next Step

We know how a single program runs, uses memory, and talks to the terminal. But your computer runs thousands of programs at once. How does it keep them organized? Move to [HC.5 OS Processes](../5-os-processes/).
