# HC.1

## Mission

Understand what a computer program actually is at a fundamental level.

## Prerequisites

- None

## Mental Model

A computer only understands electricity: on or off (1s and 0s). A program is simply a structured sequence of these instructions stored on your hard drive. When you run a program, the Operating System loads these instructions into memory and the CPU executes them one by one.

At the end of the day, a program is just data that tells the CPU what to do.

## Visual Model

```mermaid
flowchart LR
    A[Hard Drive\n(Program File)] -->|Loaded by OS| B[RAM\n(Memory)]
    B -->|Fetched by| C[CPU\n(Processor)]
    C -->|Executes Instructions| D[Action\n(Prints to Screen)]
```

## Machine View

When we write code in Go (which is human-readable), the machine cannot run it directly. It must be translated into raw machine code (1s and 0s) that your specific CPU architecture understands. 

Even this simple text file you are reading is stored as binary data under the hood.

## Run Instructions

```bash
go run ./00-how-computers-work/1-what-is-a-program
```

## Code Walkthrough

In the code below, we define a small set of human-readable instructions. 
1. We declare we are making an executable program (`package main`).
2. We define the starting point (`func main()`).
3. We instruct the computer to output text to the screen (`println`).

## Try It

1. Change the text inside the quotes to your name.
2. Run the program again.
3. Observe how the CPU blindly followed your new instruction.

## Next Step

Now that you know a program is just a set of instructions, how do we get from human-readable text to machine code? Move to [HC.2 Code to Execution](../2-code-to-execution/).
