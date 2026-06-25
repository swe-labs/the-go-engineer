# Fuzz testing

## Learning objective

Write and run Go fuzz tests using `f.Fuzz`, `f.Add`, and `go test -fuzz` to automatically discover edge cases and crash-inducing inputs.

## Why this matters

Hand-written tests cover the cases you think of. Fuzz testing covers the cases you did not think of — empty strings, negative numbers, Unicode, boundary values, and malicious inputs. The fuzzer generates random inputs and feeds them to your function, looking for panics, crashes, or violated invariants. For security-critical code, parsing, and data validation, fuzz testing is not optional; it is a requirement.

## Mental model

A fuzz test is like a robot that bangs on your function with millions of random inputs, watching for crashes. You give it a seed corpus (known-good inputs) and an invariant (a property that must always hold). The fuzzer mutates the seeds to generate new inputs. If an input causes a crash or violates the invariant, the fuzzer saves it to the testdata corpus and reports it as a finding.

```
seed inputs → [mutate] → random inputs → [your function] → crash? → saved to corpus
```

## Core idea

A fuzz test in Go has this structure:

```go
func FuzzMyFunc(f *testing.F) {
    // Seed the corpus with known-good inputs.
    f.Add("hello")
    f.Add("")

    // Fuzz target: called with random inputs.
    f.Fuzz(func(t *testing.T, input string) {
        result := MyFunc(input)
        // Check an invariant.
        if result != "" && !strings.HasPrefix(result, "prefix") {
            t.Errorf("unexpected result: %s", result)
        }
    })
}
```

Key components:

| Component | Role |
|---|---|
| `f *testing.F` | Fuzz test driver, like `*testing.T` for regular tests |
| `f.Add(inputs...)` | Seeds the fuzzer with known-good inputs |
| `f.Fuzz(fn)` | Registers the fuzz target function that receives random inputs |
| `t *testing.T` inside `f.Fuzz` | Report failures per-input, just like a regular test |
| `go test -fuzz=FuzzMyFunc` | Run the fuzzer (indefinitely until a failure is found) |
| `testdata/fuzz/FuzzMyFunc/*` | Corpus directory where crash inputs are saved |

## Under the hood

Go's fuzzer uses **coverage-guided fuzzing** (similar to libFuzzer). It instruments the binary to track which code paths are executed by each input. Inputs that increase code coverage are retained and mutated further. This means the fuzzer naturally finds edge cases — inputs that trigger uncommon branches, boundary conditions, and error paths.

When a crash is found, the fuzzer:
1. Minimises the input to the smallest version that still crashes.
2. Writes it to `testdata/fuzz/<FuzzTestName>/`.
3. The next `go test` run replays all corpus entries, so the failure becomes a regular test failure.

## How Go uses it

- **Standard library**: Go's standard library uses fuzz testing extensively in `strings`, `net/http`, `encoding/json`, `crypto/*`, and `regexp`.
- **Security audits**: Fuzz testing is a standard part of security vulnerability research (CVEs are often found via fuzzing).
- **Parsing and serialization**: JSON, XML, CSV, protobuf parsers are prime candidates for fuzz testing.
- **Input validation**: Functions that validate emails, URLs, IP addresses, file paths.
- **`go test -fuzz=Fuzz -fuzztime=30s`**: Run the fuzzer for 30 seconds in CI.

## Go example

```go
package main

import (
	"errors"
	"testing"
	"unicode/utf8"
)

// Reverse returns the reversed version of s.
func Reverse(s string) (string, error) {
	if !utf8.ValidString(s) {
		return "", errors.New("invalid utf-8")
	}
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes), nil
}

func FuzzReverse(f *testing.F) {
	seeds := []string{"hello", "world", "12345", "", "a", "abc"}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		// Invariant: reversing twice gives the original.
		first, err := Reverse(s)
		if err != nil {
			return // skip invalid inputs
		}
		second, err := Reverse(first)
		if err != nil {
			t.Errorf("Reverse(Reverse(%q)) returned error: %v", s, err)
		}
		if second != s {
			t.Errorf("Reverse(Reverse(%q)) = %q, want %q", s, second, s)
		}
	})
}
```

Run with:
```bash
go test -fuzz=FuzzReverse -fuzztime=10s
```

## Step-by-step execution

1. `f.Add` seeds the corpus with `"hello"`, `"world"`, `"12345"`, `""`, `"a"`, `"abc"`.
2. `f.Fuzz` registers the target. The fuzzer starts with the seed inputs.
3. For each seed, the target runs. `Reverse("hello")` → `"olleh"`, second reverse → `"hello"`. Passes.
4. The fuzzer mutates `"hello"` → `"hlllo"` (changed 'e' to 'l'). Runs target. Passes.
5. The fuzzer mutates to include multi-byte UTF-8: `"héllo"`. `Reverse` works on runes → `"olléh"`. Passes.
6. The fuzzer generates an invalid UTF-8 byte sequence. `Reverse` returns error. Target returns early. Passes.
7. If the fuzzer finds an input that causes `Reverse(Reverse(s)) != s`, it saves the input and fails the test.

## Common mistakes

- **Fuzzing without an invariant**: The fuzz target must check something. If it just calls the function without assertions, even a crash-only fuzz is useful, but property-based checks are more powerful.
- **Ignoring the seed corpus**: Good seeds guide the fuzzer toward interesting code paths. Start with representative inputs.
- **Fuzzing with side effects**: The fuzz target may be called concurrently by the fuzzer. Avoid mutating shared state inside `f.Fuzz`.
- **Not handling expected errors**: If your function returns errors for invalid inputs, check `err != nil` in the fuzz target and return early — otherwise the fuzzer reports "failure" for expected error paths.
- **Running `go test` without `-fuzz`**: `go test` runs the fuzz corpus entries as regular tests but does not generate new inputs. You must pass `-fuzz=FuzzName` or `-fuzz=.` to start fuzzing.

## Debugging walkthrough

You run `go test -fuzz=FuzzReverse -fuzztime=10s`. The fuzzer finds a crash after 3 seconds:

```
--- FAIL: FuzzReverse (0.00s)
    --- FAIL: FuzzReverse (0.00s)
        main_test.go:30: Reverse(Reverse("¡")) = "¡", want "¡"
```

Wait, that passes? Actually, the output says "reverse of reverse" — let's fix the example to show a real crash. Suppose `Reverse` panics on a multi-byte input:

```
    --- FAIL: FuzzReverse (3.2s)
        main_test.go:28: Reverse("héllo") panicked: runtime error...
```

The crash input is saved in `testdata/fuzz/FuzzReverse/<hash>`. Read it:
```bash
cat testdata/fuzz/FuzzReverse/abcdef123456
go test -run=FuzzReverse/abcdef123456  # replay the exact input
```

Fix the bug (e.g. handle rune slicing correctly), verify the corpus entry passes, then remove it or let the fuzzer confirm.

## Production notes

- **CI fuzzing**: Run `go test -fuzz=Fuzz -fuzztime=30s` in CI on critical packages. It adds 30 seconds to CI but catches regressions.
- **Fuzz time limit**: Always set `-fuzztime` in CI. Without it, the fuzzer runs indefinitely.
- **Corpus in version control**: Check in `testdata/fuzz/` so crash regressions are replayed in CI and by other developers.
- **Fuzzing for security**: Run the fuzzer for hours or days on security-critical code. Use tools like `go-fuzz` (more configurable) or OSS-Fuzz for continuous fuzzing.
- **Race detector with fuzzing**: `go test -fuzz=Fuzz -race` catches data races with generated inputs. Expensive but valuable.

## Performance implications

- Coverage-guided fuzzing is CPU-intensive. A single fuzzing process uses 100% of a CPU core.
- Run `GOMAXPROCS=1` with `go test -fuzz` on a laptop to avoid saturating all cores. Or use `-parallel=1`.
- The fuzzer allocates memory for each input and coverage bitmap. On a long fuzzing run (hours), watch memory usage.
- Fuzzing with `-race` is 5–10x slower but may find races that normal fuzzing misses.

## Practice task

Write a fuzz test for a `ParseEmail(email string) (local, domain string, err error)` function that validates email format. The invariant: if `ParseEmail` succeeds, the result must contain exactly one `@` when concatenated as `local + "@" + domain`. Seed with `"user@example.com"`, `""`, `"@x"`. Run `go test -fuzz=FuzzParseEmail -fuzztime=10s`.

## Tests / verification

```bash
go test ./curriculum/modules/06-testing-debugging-refactoring/lessons/20-fuzz-testing
go test -fuzz=FuzzReverse -fuzztime=5s ./curriculum/modules/06-testing-debugging-refactoring/lessons/20-fuzz-testing
```

(Ensure `testdata/fuzz/` directories are created if crashes are found.)

## Review questions

1. What is the difference between `f.Add` and `f.Fuzz`?
2. What does "coverage-guided fuzzing" mean?
3. Why is it important to check `err != nil` and return early in a fuzz target?
4. What happens when the fuzzer finds a crash input?
5. Why should you set `-fuzztime` in CI?

## NEXT UP

Race detector preview — detecting and preventing data races in concurrent Go programs.
