# Golden files

## Learning objective

Implement the golden file pattern for testing output, manage `testdata/*.golden` files, support an `-update` flag to regenerate golden files, and know when golden files are appropriate.

## Why this matters

Some functions produce complex output — multi-line strings, formatted reports, serialized data, generated code. Writing explicit assertions for every line is tedious and fragile. Golden files solve this: you store the expected output once, and the test compares the actual output against it. When the expected output changes intentionally, you regenerate the golden file with a flag. This pattern is used in Go's own compiler tests, code generators, and formatters.

## Mental model

A golden file is a snapshot of expected output. The test generates output → reads the golden file → compares them. If they differ, the test fails and shows a diff. When the output is intentionally updated (e.g., you added a field to the serialization), you run the test with `-update` and the golden file is rewritten. This turns output verification into a diff operation.

```
Actual output  ──→  Compare  ←──  testdata/*.golden
                         │
                         ↓
                    Match? → Pass
                    No match → Fail (show diff)
                    -update flag → overwrite golden file
```

## Core idea

The golden file pattern consists of:

1. A `testdata/` directory containing `.golden` files (expected output).
2. A test that generates output and compares it to the golden file.
3. A `-update` flag (via `flag` package) that rewrites golden files instead of comparing.

```go
var update = flag.Bool("update", false, "update golden files")

func TestGolden(t *testing.T) {
    got := GenerateOutput()
    golden := filepath.Join("testdata", tc.name+".golden")
    if *update {
        os.WriteFile(golden, []byte(got), 0644)
    }
    want, _ := os.ReadFile(golden)
    if got != string(want) {
        t.Errorf("output mismatch (-want +got):\n%s", diff)
    }
}
```

Key guidelines:

| Guideline | Reason |
|---|---|
| Use `.golden` extension | Convention, not enforced |
| Check in golden files | They are the expected output |
| Use `flag` for `-update` | Standard Go pattern |
| Always show a diff on failure | Makes debugging fast |
| Golden files per test case | One golden per subtest or behavior |

## Under the hood

## How Go uses it

The standard library uses golden files extensively:

- `go/printer` tests format output against golden files.
- `encoding/xml` tests XML generation.
- `text/template` tests template execution output.
- `cmd/go` tests for expected help text and error messages.

The `flag` package is part of the standard library. When combined with `go test`, the `-update` flag must be passed as `go test -update` (or `go test -args -update` before Go 1.21). The test binary receives the flag.

## Go example

```go
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var update = flag.Bool("update", false, "update golden files")

// FormatPerson returns a formatted string for a person.
func FormatPerson(name string, age int) string {
	return fmt.Sprintf("Name: %s\nAge: %d\n", name, age)
}

func main() {
	fmt.Print(FormatPerson("Alice", 30))
}
```

```go
// main_test.go
package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "update golden files")

func TestFormatPerson(t *testing.T) {
	tests := []struct {
		name string
		person string
		age    int
	}{
		{name: "alice", person: "Alice", age: 30},
		{name: "bob",   person: "Bob", age: 25},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := FormatPerson(tc.person, tc.age)
			golden := filepath.Join("testdata", tc.name+".golden")

			if *update {
				os.MkdirAll("testdata", 0755)
				if err := os.WriteFile(golden, []byte(got), 0644); err != nil {
					t.Fatalf("writing golden file: %v", err)
				}
				return
			}

			wantBytes, err := os.ReadFile(golden)
			if err != nil {
				t.Fatalf("reading golden file: %v", err)
			}
			want := string(wantBytes)
			if got != want {
				t.Errorf("output mismatch\n got: %q\nwant: %q", got, want)
			}
		})
	}
}
```

## Step-by-step execution

First run (no golden files exist):

1. `go test -v` discovers `TestFormatPerson`.
2. Subtests "alice" and "bob" execute.
3. "alice" calls `FormatPerson("Alice", 30)` → `"Name: Alice\nAge: 30\n"`.
4. `*update` is `false`. Tries to read `testdata/alice.golden`. File not found → test fails.
5. You run `go test -v -update` to create golden files.
6. This time, `*update` is `true`. The test writes `got` to `testdata/alice.golden`.
7. Tests pass (no comparison is done when `-update` is set).
8. You run `go test -v` again without `-update`. Golden files exist, comparison succeeds.

If you change `FormatPerson` to add an email field, the test fails on the next run. You update the golden files with `-update`, review the diff, and commit the new golden files alongside the code change.

## Common mistakes

- **Committing golden files without reviewing the diff.** Always verify the new golden content is correct. `git diff testdata/` should be reviewed like any code change.

- **Overwriting golden files on every test run.** The `-update` flag must be explicit. Never update golden files automatically in normal test runs.

- **Binary golden files.** Golden files are text-based by convention. For binary output, compare byte slices or use a hash.

- **Not showing a diff on failure.** Just saying "mismatch" is not helpful. Show what changed. Use `cmp.Diff` from `github.com/google/go-cmp` or `t.Errorf` with both values.

- **Golden files for trivial output.** If the output is `Add(2, 3) = 5`, a golden file is overkill. Use direct assertions.

## Debugging walkthrough

A golden test fails unexpectedly:

```go
func TestGenerateHTML(t *testing.T) {
    got := GenerateHTML()
    golden := "testdata/output.golden"
    want, _ := os.ReadFile(golden)
    if got != string(want) {
        t.Errorf("mismatch")
    }
}
```

**Symptom**: Test fails but the message just says "mismatch" with no details.

**Investigation**: No diff is shown. You add `t.Errorf("got:\n%s\nwant:\n%s", got, want)` to see the difference.

**Root cause**: `GenerateHTML` now includes a CSS class name that changed. The golden file has the old class name.

**Fix**: Review the change. If the new output is correct, run `go test -update` to regenerate the golden file. Verify the diff with `git diff testdata/output.golden`.

## Production notes

- **Always show a diff.** Use `t.Errorf` with `%s` formatting or a proper diff library like `github.com/google/go-cmp/cmp`.
- **One golden per subtest.** Naming: `testdata/TestName/subtest.golden` or `testdata/subtest.golden`.
- **Check in golden files.** They are the source of truth for expected behavior. Excluding them from version control defeats the purpose.
- **Large golden files.** If golden files are large (MB+), consider whether the test is testing the right thing. Large files are hard to review in diffs.
- **`-update` in CI.** Do not run CI with `-update`. It would silently overwrite golden files without review.

## Performance implications

- Reading golden files adds disk I/O per test case. For thousands of cases, batch comparisons or embed golden strings with `//go:embed`.
- `os.ReadFile` of a small golden file is ~microseconds. Negligible.
- The `-update` mode writes files, which is slightly more expensive but only used during development.

## Practice task

Write a function `GenerateReport(items []string) string` that returns a numbered list:
```
1. item1
2. item2
3. item3
```

Create golden files for a report with 0 items, 1 item, and 3 items. Write a golden test that:
- Reads the golden file on normal runs.
- Supports `go test -update` to regenerate.
- Shows a clear diff on mismatch.

## Tests / verification

```bash
go run ./curriculum/modules/06-testing-debugging-refactoring/lessons/07-golden-files
go test -v ./curriculum/modules/06-testing-debugging-refactoring/lessons/07-golden-files
go test -v -update ./curriculum/modules/06-testing-debugging-refactoring/lessons/07-golden-files
```

## Review questions

1. What problem do golden files solve that explicit assertions do not?
2. Why should the `-update` flag be explicit rather than automatic?
3. What is the standard file extension and directory convention for golden files?
4. How would you safely update golden files when the output intentionally changes?
5. What are the risks of using golden files for trivial or unstable output?

## NEXT UP

Testable design with `io.Writer` — injecting dependencies to make output testable.
