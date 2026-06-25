# Runes and Unicode

## Learning objective

Explain the relationship between `rune`, `byte`, Unicode code points, and UTF-8 encoding, and write Go code that correctly handles multi-byte characters.

## Why this matters

The world speaks Unicode. Every modern application handles text containing emoji, CJK characters, accented Latin letters, or right-to-left scripts. Go was designed by Unicode-aware engineers — its string type, rune type, and range loops are all built on Unicode and UTF-8. If you confuse bytes with characters, you will corrupt user text. Understanding runes is essential for any Go developer working with internationalized text.

## Mental model

A `rune` is a single Unicode code point — one conceptual character. A `string` is a sequence of bytes encoded in UTF-8. One rune may take 1, 2, 3, or 4 bytes in the string. The `for range` loop decodes UTF-8 one rune at a time. `len(s)` counts bytes; `utf8.RuneCountInString(s)` counts runes.

Think of runes as the "what you mean" and bytes as the "what's stored". The string `"é"` is one rune (U+00E9) but two bytes (`0xC3 0xA9` in UTF-8). The string `"世界"` is two runes but six bytes.

## Core idea

- **`rune`** is an alias for `int32`. It represents a Unicode code point (a numeric value from 0 to 0x10FFFF).
- **`byte`** is an alias for `uint8`. A byte is one UTF-8 code unit.
- **UTF-8** encodes code points into 1–4 bytes:
  - 0–127 (ASCII): 1 byte.
  - 128–2047: 2 bytes.
  - 2048–65535: 3 bytes.
  - 65536–1114111: 4 bytes.
- **Invalid UTF-8** produces `U+FFFD` (REPLACEMENT CHARACTER, `\ufffd`) when decoded.

## Under the hood

Go strings are immutable byte sequences. They are not null-terminated and may contain arbitrary bytes, including invalid UTF-8. The `for range` loop over a string decodes one UTF-8 sequence per iteration:

```go
for i, r := range "hello" {
	// i is byte offset, r is rune
}
```

The compiler inserts a call to `runtime.decoderune` in the loop. Decoding a rune from a `string` or `[]byte` involves checking the leading byte to determine the sequence length, then extracting the code point bits.

The `unicode/utf8` package provides lower-level functions:
- `DecodeRuneInString(s)` — decodes the first rune.
- `EncodeRune(buf, r)` — writes the UTF-8 encoding of `r`.
- `ValidString(s)` — checks if `s` contains only valid UTF-8.

## How Go uses it

- **`range` over string**: yields `(int, rune)` pairs — byte offset and decoded code point.
- **`unicode` package**: Provides `IsLetter(r)`, `IsDigit(r)`, `ToUpper(r)`, etc.
- **`unicode/utf8`**: Rune counting, validation, encoding/decoding.
- **`unicode/utf16`**: For systems that use UTF-16 (Windows, JavaScript, Java).
- **`text/scanner`**: Tokenizes Go or text using rune-based scanning.
- **`fmt` package**: `%c` formats a rune, `%q` quotes a rune with escaping.

## Go example

```go
package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	s := "Hello, 世界 🚀"
	fmt.Println("string:", s)
	fmt.Println("len(s):", len(s), "bytes")
	fmt.Println("utf8.RuneCountInString(s):", utf8.RuneCountInString(s), "runes")

	// range over string yields runes
	fmt.Println("\nRune iteration:")
	for i, r := range s {
		fmt.Printf("  byte=%d rune=%c (%U) width=%d\n", i, r, r, utf8.RuneLen(r))
	}

	// Manual decoding
	fmt.Println("\nManual decode:")
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		fmt.Printf("  byte=%d rune=%c (%U) size=%d\n", i, r, r, size)
		i += size
	}

	// Invalid UTF-8
	invalid := "hello\xfe\xffworld"
	fmt.Println("\nInvalid UTF-8:")
	fmt.Println("string:", invalid)
	for i, r := range invalid {
		fmt.Printf("  byte=%d rune=%q (%U)\n", i, r, r)
	}
}
```

## Step-by-step execution

For `for i, r := range "Go 🚀"`:

1. String `"Go 🚀"` is bytes: `[0x47, 0x6f, 0x20, 0xF0, 0x9F, 0x9A, 0x80]`.
2. Iteration 1: `i=0`, leading byte `0x47` is ASCII (starts with `0`). Decode 1 byte → rune `'G'` (U+0047).
3. Iteration 2: `i=1`, leading byte `0x6f` is ASCII. Decode 1 byte → rune `'o'` (U+006F).
4. Iteration 3: `i=2`, leading byte `0x20` is ASCII. Decode 1 byte → rune `' '` (U+0020).
5. Iteration 4: `i=3`, leading byte `0xF0` starts with `11110` → 4-byte sequence. Read next 3 bytes: `0x9F, 0x9A, 0x80`. Decode → rune `🚀` (U+1F680).
6. Loop ends: `i=7`, equal to `len(s)`.

For `utf8.RuneCountInString(s)`:

1. Walks the string, counting valid UTF-8 sequences.
2. For `"Go 🚀"`: 4 sequences → returns `4`.

## Common mistakes

- Mistake: Using `len(s)` to count characters.
  - Why it happens: `len(s)` returns byte count. For multi-byte strings, this is larger than the character count.
  - Fix: Use `utf8.RuneCountInString(s)`.

- Mistake: Indexing a string to get the nth character.
  - Why it happens: `s[2]` returns the 3rd byte, not the 3rd character. For strings with multi-byte characters, this may be the middle of a rune.
  - Fix: Convert to `[]rune(s)` and index that, or use `utf8.DecodeRuneInString`.

- Mistake: Not handling invalid UTF-8.
  - Why it happens: External data (files, network) may contain invalid UTF-8. `range` silently produces `U+FFFD` for invalid sequences. The program may propagate corrupted data.
  - Fix: Validate with `utf8.ValidString(s)` or `utf8.Valid(b)` before processing.

- Mistake: Confusing `rune` literals with `string` literals.
  - Why it happens: `'a'` is a rune (`int32`), `"a"` is a string. `'a' + 1` is `98`, `"a" + "b"` is `"ab"`.
  - Fix: Use single quotes for runes, double quotes for strings.

## Debugging walkthrough

```go
package main

import "fmt"

func main() {
	s := "café"
	fmt.Println("len:", len(s))
	for i := 0; i < len(s); i++ {
		fmt.Printf("byte %d: %x\n", i, s[i])
	}
}
```

**Symptom**: Prints `len: 5` but "café" looks like 4 characters. Byte dump: `63 61 66 c3 a9`.

**Investigation**: Add rune iteration:

```go
for i, r := range s {
	fmt.Printf("offset %d: rune %c (%U)\n", i, r, r)
}
```

Output:
```
offset 0: rune c (U+0063)
offset 1: rune a (U+0061)
offset 2: rune f (U+0066)
offset 3: rune é (U+00E9)
```

**Root cause**: `é` is encoded as 2 bytes (`0xC3 0xA9`) in UTF-8. `len(s)=5` because it counts bytes. The `s[3]` direct index gives `0xC3`, the first byte of `é`, not the character.

**Fix**: Use `for range` or `utf8.RuneCountInString` for character-oriented operations. Use `[]rune(s)` for random access by rune index.

## Production notes

- **Always validate UTF-8** from external sources. `utf8.Valid(b)` is fast (a few ns per byte).
- **`[]rune(s)` allocates**: Converting a string to `[]rune` decodes all UTF-8 sequences and allocates a new slice. For large strings, this is expensive.
- **`strings` package functions** work at the byte level on valid UTF-8 strings. They don't split runes, so they are safe for substring searches even in multi-byte text.
- **Normalization**: "é" can be one code point (U+00E9) or two (U+0065 + U+0301). Use `golang.org/x/text/unicode/norm` for NFC/NFD normalization.
- **Emoji and grapheme clusters**: A single visible emoji like "👨‍👩‍👧‍👦" (family) is multiple code points joined by zero-width joiners. `range` iterates over individual code points, not grapheme clusters. For grapheme clusters, use `golang.org/x/text/segment`.

## Performance implications

- `for range` over a string is O(n) in bytes but decodes UTF-8 one rune at a time. It is fast but not as fast as byte iteration.
- `utf8.RuneCountInString(s)` is O(n) — it scans the entire string.
- `[]rune(s)` allocates a `[]int32` of length `runeCount`, plus the copy cost.
- Direct byte access `s[i]` is O(1). Use it when you only need ASCII or raw bytes.
- `unicode.IsLetter(r)`, `unicode.IsDigit(r)` use binary search in Unicode tables — fast but not free.

## Practice task

Write a function `reverseRunes(s string) string` that reverses the runes in `s` (preserving multi-byte characters). Then write a function `countLetters(s string) map[rune]int` that counts occurrences of each Unicode letter in `s`. Letters are runes where `unicode.IsLetter(r)` returns true.

In `main()`, test on `"Hello, 世界 🚀"` and `"café"`.

## Tests / verification

```bash
go run ./curriculum/modules/03-programming-fundamentals/lessons/21-runes-and-unicode
go test ./curriculum/modules/03-programming-fundamentals/lessons/21-runes-and-unicode
```

## Review questions

1. What is the difference between `rune` and `byte` in Go?
2. Why does `len("é")` return 2 instead of 1?
3. What does `for i, r := range s` yield for each iteration?
4. What value does `range` produce for an invalid UTF-8 byte sequence?
5. When would you use `[]rune(s)` instead of `for range s`?

## NEXT UP

Pointers as addresses — Go's reference-by-value pointer model and how to use `&` and `*`.
