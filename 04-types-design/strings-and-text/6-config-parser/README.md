# ST.6 Config Parser Project

## Mission

Build a small config parser that turns `.env`-style text into structured data and renders a stable summary from that data.

> **Backward Reference:** In [Lesson 5: Text Template](../5-text-template/README.md), you learned how to generate text from data. Now you will combine everything you have learned about strings, formatting, Unicode, Regex, and Templating to build a complete text-processing pipeline.

## Prerequisites

- `ST.1` strings
- `ST.2` formatting
- `ST.4` regex
- `ST.5` text templates

## Mental Model

This milestone is a text-processing pipeline:

- read input line by line
- parse valid key-value lines
- ignore comments and blanks
- store structured data
- render predictable output

It turns the earlier string tools into one end-to-end flow.

## Visual Model

```mermaid
graph LR
    A["raw config text"] --> B["Scanner"]
    B --> C["regex parse"]
    C --> D["map[string]string"]
    D --> E["text/template output"]
```

## Machine View

The parser scans bytes into lines, applies one compiled regular expression to each candidate line, stores results in a map, and then executes a template to produce output. The regex should be compiled once, not rebuilt during each loop iteration.

## Run Instructions

```bash
go run ./04-types-design/strings-and-text/6-config-parser
go run ./04-types-design/strings-and-text/6-config-parser/_starter
go test ./04-types-design/strings-and-text/6-config-parser
```

## Solution Walkthrough

### `bufio.Scanner`

The scanner turns one multi-line input string into a sequence of lines the parser can process incrementally.

### Compiled regex

The regex captures valid key-value lines and keeps the parsing rule centralized.

### Skip comments and blanks

Ignoring non-data lines keeps the parser focused and predictable.

### `map[string]string`

The parsed config values are stored in a lookup-friendly structure.

### `text/template`

Template rendering keeps the final output stable and easier to maintain than scattered print statements.

## Try It

1. Add another valid config line and inspect the rendered output.
2. Add a malformed line and verify the parser reports an error.
3. Change the template and rerun the program.

## Verification Surface

```bash
go run ./04-types-design/strings-and-text/6-config-parser
go run ./04-types-design/strings-and-text/6-config-parser/_starter
go test ./04-types-design/strings-and-text/6-config-parser
```

## In Production
Config parsing is exactly where text bugs become system bugs. Poor parsing rules, unstable output, and ad hoc string handling lead to silent misconfiguration and painful incident debugging.

## Thinking Questions
1. Why is compiling the regex once better than compiling it inside the scan loop?
2. What advantage does a map provide after parsing?
3. Why is template-based rendering safer than manual scattered output?

> **Forward Reference:** You have completed the section on Types and Design, as well as the deep dive into Text. Now it is time to look at how Go programs are organized and how they interact with the outside world. In [Section 05: Packages and I/O](../../../05-packages-io/README.md), you will learn about module boundaries, file systems, and network communication.

## Next Step

Next: `MP.1` -> `05-packages-io/01-modules-and-packages/1-module-basics`

Open `05-packages-io/01-modules-and-packages/1-module-basics/README.md` to continue.
