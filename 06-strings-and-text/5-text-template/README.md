# ST.5 Text Templates

## Mission

Learn how Go separates presentation from data with `text/template`.

## Prerequisites

- `ST.1` strings
- `ST.2` formatting

## Mental Model

Templates let you describe output structure separately from the data that fills it in.

The flow is:

1. define a template
2. parse it
3. execute it with data
4. collect or print the rendered output

## Visual Model

```mermaid
graph LR
    A["template text"] --> B["Parse"]
    C["data struct"] --> D["Execute"]
    B --> D
    D --> E["rendered output"]
```

## Machine View

`text/template` parses template source into an internal syntax tree. During execution, it uses reflection to read exported fields from the provided data and write rendered output to an `io.Writer`.

## Run Instructions

```bash
go run ./06-strings-and-text/5-text-template
```

## Code Walkthrough

### `type EmailData struct { ... }`

The struct provides the fields the template is allowed to read.

### `template.New(...).Parse(...)`

Parsing turns raw template text into a reusable compiled template object.

### Template actions like `{{if ...}}` and `{{range ...}}`

These actions add conditional and loop behavior inside the template itself.

### `strings.Builder`

The builder captures rendered output efficiently before the lesson prints it.

### `tmpl.Execute(&output, data)`

Execution walks the parsed template and fills in the placeholders from the provided data.

## Try It

1. Add another field to `EmailData` and render it.
2. Change the `if` branch by adjusting `UnreadCount`.
3. Add another item to the `Items` slice and inspect the `range` output.

## ⚠️ In Production

Templates are how teams generate emails, config files, reports, and stable text output without mixing presentation logic into business logic. The separation pays off quickly as systems grow.

## 🤔 Thinking Questions

1. Why is parsing once and executing many times a good design?
2. Why must template-visible struct fields be exported?
3. When is a template clearer than a long chain of manual string concatenation?

## Next Step

Continue to `ST.6` config parser project.
