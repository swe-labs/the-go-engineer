# TI.7 Advanced Generics

## Mission

Push generic collection helpers further once the core section feels solid. This is a stretch lesson for learners who want more practice with generics.

## Why This Lesson Exists Now

You have completed the core generics lesson (TI.5) and the payroll processor milestone. This stretch lesson provides additional practice with generic patterns like Filter and Map on custom types.

## Prerequisites

- `TI.5` generics

## Run Instructions

```bash
go run ./01-foundations/06-types-and-interfaces/10-advanced-generics
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

## Production Relevance

Generic utilities like Filter and Map are common in real Go codebases. Understanding how to write and use them improves code reusability.

## Next Step

After this stretch lesson, continue to the next section in the foundations path.