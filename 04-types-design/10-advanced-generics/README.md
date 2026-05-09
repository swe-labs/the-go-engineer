# TI.10 Advanced Generics

## Mission

Push generic collection helpers further once the core section feels solid. This is a stretch lesson for learners who want more practice with generics.

## Why This Lesson Exists Now

You have completed the core generics lesson (`TI.9`) and the payroll processor milestone. This stretch lesson provides additional practice with generic patterns like `Filter` and `Map` on custom types.

## Prerequisites

- `TI.9` generics

## Mental Model

Think of this lesson as the "extra tools" drawer after you already know how to use the core workshop safely. The foundational path taught when generics help. This stretch lesson asks how far you can push generic helpers before they become harder to read than the duplication they remove.

## Visual Model

```mermaid
graph TD
    A["data"] --> B["type definition"]
    B --> C["methods or interface behavior"]
```
```text
[]User --Filter--> []User
[]User --Map-----> []string

same helper shape
different concrete types
```

## Machine View

Generic helper calls still compile into concrete versions for the types you use. The extra abstraction cost here is mostly readability and API design, not hidden runtime reflection.

## Run Instructions

```bash
go run ./04-types-design/10-advanced-generics
```

## What You Will Learn

- Using generics to build reusable functional patterns (Map / Filter)
- Defining operations on generic collections
- Understanding type inference and its limits in Go

## Code Walkthrough

### Filter function

The Filter function takes a slice of any type T and a predicate function. It returns a new slice containing only elements where the predicate returns true.

### Map function

The Map function transforms a slice of type T into a slice of type U. Notice that Map requires two type parameters.

### User struct example

The User struct demonstrates how generic functions work with custom types, not just built-in types.

## Try It

1. Add a new Filter condition (e.g., filter by ID > 1).
2. Add a new Map transform (e.g., extract email addresses).
3. Create a different struct type and use Filter/Map with it.

## ⚠️ In Production
Generic utilities like Filter and Map are common in real Go codebases. Understanding how to write and use them improves code reusability.

## 🤔 Thinking Questions

1. What problem is this lesson trying to solve?
2. What would change if you removed this idea from the program?
3. Where do you expect to see this pattern again in real Go code?
## Next Step

Continue to `TI.11` empty interface if you want to keep exploring the optional stretch path.
