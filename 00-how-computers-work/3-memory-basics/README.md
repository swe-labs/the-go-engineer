# HC.3

## Mission

Build a mental model of how a program uses RAM (memory) to store data while it is running.

## Prerequisites

- `HC.2`

## Mental Model

When a program runs, it needs a place to store data temporarily. Your hard drive is too slow. Instead, the OS gives your program a block of super-fast memory called **RAM** (Random Access Memory). 

Inside your program, memory is divided into two main areas:
1. **The Stack:** Fast, organized memory. It stores small, temporary variables that disappear as soon as a function finishes running. It works exactly like a stack of plates.
2. **The Heap:** Slower, messy memory. It stores large or long-living data that needs to survive across different parts of your program.

Every piece of data you create lives at a specific numeric address in memory, much like a house has a street address.

## Visual Model

```mermaid
flowchart TD
    subgraph RAM [System Memory]
        subgraph Program Space
            S[The Stack\n(Fast, Temporary)]
            H[The Heap\n(Slower, Long-Lived)]
        end
    end
    
    A[Function A runs] -->|Puts variables on| S
    A -->|Function finishes| S2[Variables destroyed]
```

## Machine View

When you declare a variable like `score = 100`, the computer doesn't know what "score" is. It asks the OS for an available slot in RAM (e.g., address `0x140000a6018`), and physically alters the silicon at that address to hold the binary representation of `100`.

## Run Instructions

```bash
go run ./00-how-computers-work/3-memory-basics
```

## Code Walkthrough

In Go, we can ask the computer to show us the actual physical memory address where our data is living by placing an `&` in front of the variable name.

1. We create a piece of data called `score`. Because it is small and temporary, Go usually places it on the fast **Stack**.
2. We create another variable using `new(int)`. This tells Go "I want this data to live on the slow, permanent **Heap**."
3. We print out both memory addresses. Notice how the addresses look completely different? That's because the Stack and the Heap are located in entirely different regions of RAM!

## Try It

1. Run the program multiple times. Notice that the memory address changes! This is because the OS assigns your program a new random block of memory every time it starts for security reasons.

## Next Step

We know how code executes and how it stores data. Now, how does the program interact with us while it's running? Move to [HC.4 Terminal Confidence](../4-terminal-confidence/).
