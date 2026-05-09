# ST.2 Formatting

## Mission

Learn how `fmt` formats values into readable strings, aligned output, and wrapped errors.

## Prerequisites

- `ST.1` strings

## Mental Model

Formatting is template-driven output.

You give `fmt`:

- a format string
- some values

and `fmt` turns those values into text according to the verbs you chose.

## Visual Model

```mermaid
graph LR
    A["values"] --> B["fmt verbs and templates"]
    B --> C["printed or returned text"]
```

## Machine View

The `fmt` package walks each value, applies the requested verb, and writes formatted bytes to an output destination such as standard output, a string buffer, or an error wrapper.

## Run Instructions

```bash
go run ./06-strings-and-text/2-formatting-string
```

## Code Walkthrough

### General verbs like `%v`, `%+v`, and `%#v`

These show the same value at different levels of detail.

### Type-specific verbs

`%s`, `%d`, `%f`, `%t`, and `%p` format strings, integers, floats, booleans, and pointers.

### Width and alignment

Formatting options like `%-15s` and `%010d` make tables and logs easier to scan.

### `fmt.Sprintf(...)`

This returns a formatted string instead of printing directly.

### `fmt.Errorf(... %w ...)`

This wraps one error with additional context while preserving the underlying cause.

## Try It

1. Change the width values in the table output.
2. Add another struct field and print it with several verbs.
3. Wrap a different error with `fmt.Errorf`.

## ⚠️ In Production

Readable output is operationally important. Logs, CLI tools, diagnostics, and user-facing errors all depend on deliberate formatting choices. Small formatting improvements often make debugging dramatically easier.

## 🤔 Thinking Questions

1. Why is `%+v` often more useful than `%v` while debugging structs?
2. When do you want `Sprintf` instead of `Printf`?
3. Why is error wrapping better than replacing the original error text entirely?

## Next Step

Continue to `ST.3` unicode and runes.
