# TI.5 Stringer

## Mission

Learn how to control how your types are displayed by implementing the fmt.Stringer interface.

## Why This Lesson Exists Now

You have learned how to define structs and methods. The next practical question is: "How do I control what shows when I print my type?"

When you pass a struct to fmt.Println, Go prints the raw fields by default. To control the output, implement the fmt.Stringer interface.

## Prerequisites

- `TI.2` methods
- `TI.3` interfaces

## Mental Model

When you hand someone a business card, the card shows a carefully formatted summary—not raw data. The String() method is your type's business card. Without it, Go prints the raw struct fields. With it, you control exactly how your type is presented.

## Visual Model

```text
fmt.Stringer interface:
┌─────────────────────────┐
│ type Stringer interface │
├─────────────────────────┤
│ String() string         │
└─────────────────────────┘

Implementations:
┌────────────┐  ┌────────────┐
│ HTTPStatus  │  │   Weekday  │
├────────────┤  ├────────────┤
│ String()    │  │  String()  │
└────────────┘  └────────────┘
```

## Machine View

The fmt package automatically checks for the Stringer interface when printing. If your type implements String() string, fmt.Println, fmt.Printf with %s or %v, log.Println, and error messages will all call your method.

## Run Instructions

```bash
go run ./01-foundations/06-types-and-interfaces/5-stringer
```

## Code Walkthrough

### `func (s HTTPStatus) String() string {`

This implements fmt.Stringer for HTTPStatus. Now when you print an HTTPStatus, it shows "HTTP 200: OK" instead of the raw struct.

### Custom types

You can create new types from existing ones: `type Weekday int`. This creates a completely new type—Weekday and int are not interchangeable.

### Stringer with iota

Combine custom types with iota (from Section 02) to create enum-like constants.

## Try It

1. Add a String() method to the Server struct from TI.1 and test it with fmt.Println.
2. Create a custom type based on float64 and implement Stringer.
3. Try printing a value before and after adding the Stringer implementation.

## Common Questions

- Why is Stringer the most commonly implemented interface?
  Because every type needs to be displayed somewhere—logs, errors, user output.

- What is the difference between %v and %s?
  %v uses String() if available, %s specifically calls String().

## Production Relevance

Stringer is essential for logging, debugging, and user-facing output. It makes your types readable in any context where they are printed or logged.

## Next Step

Continue to `TI.5` generics.