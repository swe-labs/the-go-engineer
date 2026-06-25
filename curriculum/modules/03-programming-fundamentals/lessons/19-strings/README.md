# Strings

## Learning objective

Manipulate Go strings using the `strings` and `strconv` packages, explain why strings are immutable `[]byte` under the hood, and handle UTF-8 encoded text correctly without corrupting multi-byte sequences.

## Why this matters

Strings are everywhere: HTTP request bodies, JSON payloads, database queries, log lines, user input, config files. A production Go service may process millions of strings per second. Misunderstanding Go's string model leads to silent data corruption (slicing a multi-byte UTF-8 character in half), memory blowups (concatenating in a loop), or security vulnerabilities (SQL injection from unsanitized string building). Every professional Go engineer must know how strings actually work _in memory_ and how to use the standard library efficiently.

## Mental model

A Go `string` is a **read-only view** into a `[]byte`. The string header is two words: a pointer to the underlying byte array and a length. You can think of a string as a `struct { ptr *byte; len int }`. String operations like `s[3]` (indexing) and `s[3:7]` (slicing) operate on bytes, not characters. This distinction is critical when working with non-ASCII text.

```
String "Hi 👍" in memory:
   ptr ──→ [0x48] 'H'
            [0x69] 'i'
            [0x20] ' '
            [0xF0] '👍' (4-byte UTF-8 sequence: F0 9F 91 8D)
            [0x9F]
            [0x91]
            [0x8D]
   len = 7 (bytes), not 3 (runes/characters)
```

## Core idea

A string in Go is:

- **Immutable**: you cannot modify a string's bytes in place. Any "modification" creates a new string.
- **Not a pointer, not a value**: it is a struct-like header with a pointer and a length, passed by value (the header is copied, but the underlying bytes are shared).
- **Not automatically NUL-terminated**: unlike C, Go strings can contain `\0` bytes.
- **Any sequence of bytes**: there is no requirement that the bytes be valid UTF-8, though idiomatic Go assumes they usually are.

### String literals

```go
s := "hello"          // interpreted literal: \n, \t, \\, \" are escape sequences
raw := `hello\nworld` // raw literal: backslashes are literal, spans newlines
```

### Common operations

| Operation | Syntax | Returns | Notes |
|---|---|---|---|
| Index | `s[i]` | `byte` | Panics if out of bounds |
| Slice | `s[i:j]` | `string` | Shares memory |
| Concatenate | `s + t` | `string` | Allocates a new string |
| Compare | `s == t` | `bool` | Byte-wise, short-circuits |
| Length | `len(s)` | `int` | Byte count, not rune count |

## Under the hood

The `reflect.StringHeader` type (though deprecated, it reveals the layout):

```go
type StringHeader struct {
    Data uintptr  // pointer to the first byte
    Len  int      // number of bytes
}
```

When you write `s := "hello"; t := s`, two headers exist but both point to the same 5 bytes. When you write `t := s + " world"`, Go allocates new memory, copies `s`'s bytes and `" world"`'s bytes into it, and returns a new header.

String concatenation using `+` inside a loop allocates O(n²) memory. The compiler sometimes optimizes small concatenations, but for loops, use `strings.Builder`.

### UTF-8 and strings

Go source code is UTF-8. String literals are UTF-8. The `range` keyword over a string iterates over runes (not bytes):

```go
for i, r := range "Hi 👍" {
    fmt.Printf("%d: %U '%c'\n", i, r, r)
}
// 0: U+0048 'H'
// 1: U+0069 'i'
// 2: U+0020 ' '
// 3: U+1F44D '👍'
```

The index `i` is the byte offset, and `r` is the `rune` (Unicode code point).

## How Go uses it

The `strings` package provides search, replace, splitting, joining, trimming, and case conversion:

```go
strings.Contains(s, substr)
strings.HasPrefix(s, prefix)
strings.Index(s, substr)
strings.Join([]string{"a", "b", "c"}, ",")
strings.ReplaceAll(s, old, new)
strings.Fields(s)        // splits on whitespace
strings.Split(s, sep)
strings.ToLower(s)
strings.TrimSpace(s)
strings.Builder           // efficient concatenation
```

The `strconv` package handles string conversions:

```go
strconv.Itoa(42)              // "42"
strconv.Atoi("42")            // 42, nil
strconv.FormatFloat(3.14, 'f', 2, 64)
strconv.ParseBool("true")
strconv.Quote("hello")        // `"hello"` (with quotes and escapes)
```

## Go example

```go
package main

import (
	"fmt"
	"strings"
	"strconv"
)

func main() {
	s := "The quick brown fox"

	fmt.Println("len (bytes):", len(s))
	fmt.Println("Contains 'fox':", strings.Contains(s, "fox"))
	fmt.Println("HasPrefix 'The':", strings.HasPrefix(s, "The"))
	fmt.Println("Index of 'brown':", strings.Index(s, "brown"))

	words := strings.Fields(s)
	fmt.Println("Words:", words)
	fmt.Println("Joined:", strings.Join(words, "-"))

	fmt.Println("Replace:", strings.ReplaceAll(s, "fox", "cat"))
	fmt.Println("Upper:", strings.ToUpper(s))

	var sb strings.Builder
	for i := 0; i < 5; i++ {
		sb.WriteString(strconv.Itoa(i))
		sb.WriteString(", ")
	}
	fmt.Println("Builder:", sb.String())

	n, _ := strconv.Atoi("42")
	fmt.Println("Atoi 42 + 8 =", n+8)

	quoted := strconv.Quote("hello \"world\"")
	fmt.Println("Quoted:", quoted)
}
```

## Step-by-step execution

For `strings.ReplaceAll(s, "fox", "cat")` where `s = "The quick brown fox"`:

1. `ReplaceAll` calls `Replace(s, "fox", "cat", -1)` (replace all occurrences).
2. Go scans `s` byte-by-byte for the substring `"fox"` using a Boyer-Moore or Rabin-Karp algorithm internally.
3. At byte offset 16 (0-indexed in `"The quick brown "`), `s[16] == 'f'`, `s[17] == 'o'`, `s[18] == 'x'` → match.
4. `Replace` builds a new string: it copies bytes `[0:16]` → `"The quick brown "`, then appends `"cat"`, then appends nothing (no more matches).
5. Result: `"The quick brown cat"`.

For `strconv.Atoi("42")`:

1. `Atoi` calls `ParseInt(s, 10, 0)`.
2. It scans bytes: `'4'` → `4`, `'2'` → `2`.
3. Accumulates: `0 * 10 + 4 = 4`, `4 * 10 + 2 = 42`.
4. Returns `(42, nil)`.

## Common mistakes

- **Slicing multi-byte characters**: `s := "👍"; s[:2]` returns an invalid UTF-8 sequence (first 2 bytes of a 4-byte rune). Use `utf8.ValidString` and `range`, or the `utf8` package, not byte slicing.
- **Using `+` in a loop**: `for i := 0; i < 1000; i++ { s += "x" }` creates O(n²) allocations. Use `strings.Builder`.
- **Equating `len(s)` with character count**: `len("👍")` is 4 bytes, not 1. Use `utf8.RuneCountInString("👍")` for rune count.
- **Assuming strings are mutable**: `s[0] = 'H'` does not compile. Strings are read-only.
- **Confusing escape in raw vs interpreted**: In raw strings `` `\n` ``, `\n` is two characters (backslash + n). In interpreted strings `"\n"`, it's a newline.
- **Forgetting error handling with `strconv.Atoi`**: It returns an error for non-numeric input. Check it.

## Debugging walkthrough

Consider this code:

```go
package main

import "fmt"

func main() {
	name := "José"
	fmt.Println("Name length:", len(name))
	for i := 0; i < len(name); i++ {
		fmt.Printf("byte %d: %c (%x)\n", i, name[i], name[i])
	}
}
```

**Symptom**: Output shows `J o s �` — the `é` is split into two bytes and each is rendered as a corrupted character.

**Explanation**: `len("José")` is 5 bytes, not 4. `é` in UTF-8 is bytes `C3 A9`. Looping by byte renders each byte as if it were a character.

**Fix using `range`**:

```go
for i, r := range name {
    fmt.Printf("rune at byte %d: %c (%U)\n", i, r, r)
}
```

**Fix using `[]rune`**:

```go
runes := []rune(name)
fmt.Println("Rune count:", len(runes))
for _, r := range runes {
    fmt.Printf("%c ", r)
}
```

## Production notes

- `strings.Builder` grows its internal buffer exponentially (like a slice), so it amortizes O(n) allocation across n concatenations.
- `strings.Reader` implements `io.Reader`, `io.Seeker`, and `io.ReaderAt` — useful when you need to treat a string as an I/O source.
- Never concatenate SQL queries using string interpolation. Use parameterized queries (`$1`, `?`) to prevent injection.
- `strconv` functions are allocation-free for parsing and formatting; they operate directly on byte buffers.
- String comparison `==` in Go compares byte-by-byte in O(n) but short-circuits on length differences first. It is safe for Unicode strings because it compares bytes, not grapheme clusters — two strings that are Unicode-equivalent but encoded differently (e.g., `"é"` as precomposed vs decomposed) are not `==`. Use `strings.EqualFold` for case-insensitive comparison or `golang.org/x/text/cases` for Unicode-aware casing.

## Performance implications

| Operation | Cost | Notes |
|---|---|---|
| `len(s)` | O(1) | Stored in the header |
| `s[i]` | O(1) | Single byte access |
| `s[i:j]` | O(1) | No copy — shares memory |
| `s + t` | O(n+m) | Allocates and copies both strings |
| `strings.Join` | O(n) | Pre-allocates exact capacity |
| `strings.Builder` (loop concat) | O(n) amortized | Grows by doubling |
| `range s` | O(n) | Iterates runes, decoding each |

- String slicing does not allocate, but the original string must be kept alive in memory. This can cause memory leaks if you slice a large string and keep only a tiny portion.
- Use `string([]byte)` and `[]byte(string)` conversions judiciously — they allocate. The compiler sometimes optimizes away the allocation in specific contexts, but don't rely on it.
- `strings.TrimPrefix` and `strings.TrimSuffix` are cheaper than general `Trim` functions because they stop after a single check.

## Practice task

Write a function `CleanAndSplit(text string) ([]string, error)` that:

1. Trims leading/trailing whitespace.
2. Converts to lowercase.
3. Replaces all non-alphanumeric characters (except spaces) with a single space.
4. Splits on whitespace into words.
5. Filters out words shorter than 3 characters.
6. Returns the cleaned word list.
7. If the result is empty, returns an error `"empty result"`.

Then write `WordCount(words []string) map[string]int` that counts occurrences.

## Tests / verification

```bash
go run ./curriculum/modules/03-programming-fundamentals/lessons/19-strings
go test ./curriculum/modules/03-programming-fundamentals/lessons/19-strings
```

## Review questions

1. What does `len("café")` return and why?
2. How does `range` over a string differ from a byte index loop `for i := 0; i < len(s); i++`?
3. When would you use `strings.Builder` instead of `+` for concatenation?
4. What is the difference between an interpreted string literal and a raw string literal?
5. What does `strconv.Atoi` return for the input `"12abc"`?

## NEXT UP

Bytes (core-03-20): Work with `[]byte` slices, `bytes.Buffer`, and the `bytes` package for efficient binary data handling.
