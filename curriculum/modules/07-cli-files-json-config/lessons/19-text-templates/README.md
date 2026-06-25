# Text templates

## Learning objective

Generate dynamic text output using Go's `text/template` package with field access (`{{.Field}}`), pipelines, functions, control structures, and understand when to use `html/template` for safe HTML generation.

## Why this matters

Every application generates text: emails, reports, config files, code generation, HTML pages, and notifications. String concatenation and `fmt.Sprintf` break down when the output is complex or needs conditional logic. Templates separate the output structure from code logic, making both easier to maintain and audit. Go's `html/template` automatically escapes output for safe HTML rendering, preventing XSS vulnerabilities.

## Mental model

Think of a template as a mad libs form with blanks to fill in. The template text contains literal parts (the structure) and action markers (`{{...}}`) that are evaluated to insert data. The data (a Go struct, map, or slice) is the "context" passed to the template. Walking through the template, Go executes each action and writes the result to the output writer, leaving the literal parts unchanged.

## Core idea

`text/template` parses template text and executes it with a data context. Key syntax:

- `{{.}}` — the entire context.
- `{{.Field}}` — accesses a struct field or map key.
- `{{.Field.SubField}}` — chained access.
- `{{range .Items}}...{{end}}` — iterates over a slice/map.
- `{{if .Cond}}...{{else}}...{{end}}` — conditional.
- `{{with .Val}}...{{end}}` — sets `.` to `Val` inside the block.
- `{{$var := .Field}}` — variable assignment.
- `{{.Field | fun1 | fun2}}` — pipeline (passs result left-to-right).

Template functions are registered with `template.FuncMap`. Built-in functions include `printf`, `len`, `index`, `and`, `or`, `not`, `eq`, `ne`, `lt`, `le`, `gt`, `ge`.

`template.Must` wraps a template parse call and panics on error — used for template initialization at package level.

## Under the hood

The template compiler parses the template text into a tree of nodes: `TextNode` (literal text), `ActionNode` ({{...}}), `CommandNode`, `PipeNode`, etc. At execution, the tree is walked depth-first. Each node writes to the output writer or modifies the execution context (`.`, variables).

Pipelines evaluate left-to-right: `{{.Name | printf "%q"}}` evaluates `.Name` to get the value, then calls `printf "%q"` with that value as the last argument.

`html/template` is a drop-in replacement for `text/template` that contextually escapes output based on where the template appears (HTML, CSS, JS, URL). It uses the same syntax but adds automatic escaping.

## How Go uses it

- **Code generation**: `go generate` with templates for boilerplate.
- **Email rendering**: HTML or plain-text email bodies from templates.
- **Configuration generation**: Render config files from structured data.
- **CLI output**: `{{range .Results}} {{.Name}} {{end}}` for formatted listing.
- **Web servers**: `html/template.Execute(w, data)` for safe HTML responses.

## Go example

```go
package main

import (
	"os"
	"text/template"
	"time"
)

type Task struct {
	Name     string
	Priority int
	Done     bool
}

type Report struct {
	Title    string
	Tasks    []Task
	GenTime  time.Time
}

const reportTmpl = `
Report: {{.Title}}
Generated: {{.GenTime | date "2006-01-02 15:04"}}
{{if .Tasks}}
Tasks:
{{range .Tasks}}  [{{if .Done}}X{{else}} {{end}}] {{.Name}} (priority {{.Priority}})
{{end}}{{else}}No tasks to report.
{{end}}Total tasks: {{len .Tasks}}
`

// date is a custom function for template use
func date(format string) func(time.Time) string {
	return func(t time.Time) string {
		return t.Format(format)
	}
}

func main() {
	report := Report{
		Title: "Sprint 24 Status",
		Tasks: []Task{
			{Name: "Implement login", Priority: 1, Done: true},
			{Name: "Write tests", Priority: 2, Done: false},
			{Name: "Deploy to staging", Priority: 1, Done: false},
		},
		GenTime: time.Now(),
	}

	tmpl, err := template.New("report").Funcs(template.FuncMap{
		"date": date,
	}).Parse(reportTmpl)
	if err != nil {
		panic(err)
	}

	if err := tmpl.Execute(os.Stdout, report); err != nil {
		panic(err)
	}
}
```

## Step-by-step execution

1. `template.New("report")` creates a new template named "report".
2. `.Funcs(...)` registers the custom `date` function.
3. `.Parse(reportTmpl)` parses the template string into an internal AST.
4. `tmpl.Execute(os.Stdout, report)` executes the template with `report` as data.
5. Template starts with literal `\nReport: `, then `{{.Title}}` → evaluates to `"Sprint 24 Status"`.
6. `{{.GenTime | date "2006-01-02 15:04"}}` — pipeline: `.GenTime` (time.Time) passed to `date` function with format string → formatted time string.
7. `{{if .Tasks}}` — `.Tasks` is non-nil and non-empty → enters `if` branch.
8. `{{range .Tasks}}` — iterates, setting `.` to each `Task`.
9. `[{{if .Done}}X{{else}} {{end}}]` — conditional inside each iteration.
10. `{{.Name}}` and `{{.Priority}}` — field access.
11. After `range`, `{{len .Tasks}}` → `3`.
12. All output written to `os.Stdout`.

## Common mistakes

- Mistake: Calling `Execute` on a nil template or a template that failed to parse — panics.
  - Fix: Use `template.Must(tmpl.Parse(...))` to panic on parse error.

- Mistake: Forgetting to call `Execute` with a pointer to the data — works either way, but `Execute` requires the data to be exported fields accessible via reflection.

- Mistake: Pipeline order confusion — `{{.X | f}}` calls `f(.X)`, not `f` then `.X`.
  - Fix: Think of `|` as "then feed into": "take `.X` then feed it into `f`".

- Mistake: Using `text/template` for HTML — no auto-escaping leads to XSS.
  - Fix: Use `html/template` (same API, different import path).

## Debugging walkthrough

Template executes but output is empty:

```go
tmpl := template.Must(template.New("t").Parse("Hello {{.Name}}!"))
tmpl.Execute(os.Stdout, nil) // nil data
```

**Symptom**: Output is `Hello <no value>!` or empty.

**Root cause**: Passing `nil` as data. The template tries to access `.Name` on nil, which produces `<no value>` in text/template.

**Fix**: Pass a valid data context: `tmpl.Execute(os.Stdout, map[string]string{"Name": "Alice"})`.

## Production notes

- **Use `template.Must` at init time** to ensure templates are valid before the program serves traffic.
- **Cache parsed templates** — parsing is expensive; `Execute` is cheap. Parse once at startup, execute many times.
- **`html/template` escaping is contextual**: it knows if it's in HTML body, attribute, CSS, JS, or URL and escapes appropriately.
- **Template functions should be pure** (no side effects) — they may be called multiple times or not at all depending on template logic.
- **Avoid complex logic in templates**: put computation in Go code, keep templates for presentation only.

## Performance implications

- Template parsing is CPU-intensive (tokenizing, building AST) — always parse at init, not in request handlers.
- `Execute` is fast — it walks the parsed AST and writes to the output writer. For simple templates, it's comparable to `fmt.Sprintf`.
- Function calls in templates have reflection overhead — custom functions called per-row in a loop may be slower than pre-computing values in Go.
- `html/template` is ~10-20% slower than `text/template` due to contextual escaping analysis.

## Practice task

Write a function `RenderTasks(tasks []Task) (string, error)` that renders a template showing each task's name, priority, and done status. Use a custom function `label` that maps priority 1 to `"[HIGH]"`, 2 to `"[MED]"`, 3 to `"[LOW]"`. Output to a `bytes.Buffer` and return the string.

## Tests / verification

```bash
go run ./curriculum/modules/07-cli-files-json-config/lessons/19-text-templates
go test ./curriculum/modules/07-cli-files-json-config/lessons/19-text-templates
```

## Review questions

1. What is the difference between `{{.Field}}` and `{{$var := .Field}}`?
2. How does a pipeline `{{.X | f | g}}` evaluate — what is the order of calls?
3. What does `template.Must` do and why is it useful?
4. When should you use `html/template` instead of `text/template`?
5. Can you use `range` to iterate over a map in Go templates?

## NEXT UP

Base64 as encoding, not encryption — encoding binary data as text with base64.
