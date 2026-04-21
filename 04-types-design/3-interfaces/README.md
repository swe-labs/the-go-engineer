# TI.3 Interfaces

## Mission

Learn how to define behavior contracts using interfaces and achieve polymorphism without inheritance.

## Why This Lesson Exists Now

You have structs (data) and methods (behavior). The next question is: "How do I write code that works with multiple types that share the same behavior?"

In other languages, you might use inheritance. In Go, you use interfacesâ€”a type that defines what a type can do, not what it is.

## Prerequisites

- `TI.1` structs
- `TI.2` methods

## Mental Model

Think of a power outlet. The outlet defines a contract: "I accept anything with two prongs and a ground pin." A lamp, a phone charger, and a refrigerator all satisfy this contract differently, but the outlet does not care how they work internallyâ€”only that they have the right shape (methods).

## Visual Model

```mermaid
graph TD
    A["data"] --> B["type definition"]
    B --> C["methods or interface behavior"]
```
```text
â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”
â”‚ Shape interface     â”‚
â”œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”¤
â”‚ Area() float64      â”‚
â”‚ Perimeter() float64 â”‚
â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜
         â–²
         â”‚ implements
    â”Œâ”€â”€â”€â”€â”´â”€â”€â”€â”€â”
    â”‚         â”‚
â”Œâ”€â”€â”€â”´â”€â”€â”€â” â”Œâ”€â”€â”´â”€â”€â”€â”€â”
â”‚ Rect  â”‚ â”‚ Circle â”‚
â””â”€â”€â”€â”€â”€â”€â”€â”˜ â””â”€â”€â”€â”€â”€â”€â”€â”˜
```

## Machine View

An interface value is internally a 2-word struct:
- Word 1: pointer to the type descriptor (what concrete type is stored)
- Word 2: pointer to the data (the actual struct value)

This is why interface calls are slightly slower than direct callsâ€”the runtime must look up the method through the type descriptor. This is called "dynamic dispatch."

## Run Instructions

```bash
go run ./04-types-design/3-interfaces
```

## Code Walkthrough

### `type Shape interface { ... }`

This defines an interface. Any type that has both `Area()` and `Perimeter()` methods automatically satisfies this interface.

### Implicit satisfaction

Go has no "implements" keyword. If your type has the right methods, it satisfies the interface. This is called structural typing or duck typing: "if it quacks like a duck, it is a duck."

### `printShapeInfo(s Shape)`

This function accepts any type that satisfies Shape. One function works with Rectangle, Circle, Triangle, and any future type.

### Type assertions

Sometimes you need to extract the concrete type from an interface. Use the comma-ok pattern: `value, ok := s.(Circle)`.

## Try It

1. Add a new shape type (e.g., Square) and see if it works with `printShapeInfo` without any changes.
2. Try a type assertion that fails and observe the ok value.
3. Add a method to one shape type that others do not haveâ€”does it still satisfy Shape?

## Common Questions

- Why no "implements" keyword?
  Go uses structural typing. If the methods match, the contract is satisfied.

- When to use interfaces vs concrete types?
  Use interfaces when you need polymorphism. Use concrete types when you need specificity.

## ⚠️ In Production
Interfaces are Go's primary tool for abstraction and testing. They let you write code that depends on behavior, not concrete types. This is essential for dependency injection, mocking, and flexible API design.

## 🤔 Thinking Questions

1. What problem is this lesson trying to solve?
2. What would change if you removed this idea from the program?
3. Where do you expect to see this pattern again in real Go code?
## Next Step

Continue to `TI.4` interface embedding.
