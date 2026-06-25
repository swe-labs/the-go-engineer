# Base64 as encoding, not encryption

## Learning objective

Encode binary data as text using `encoding/base64`, use both `StdEncoding` and `URLEncoding`, understand when base64 is appropriate, and distinguish encoding from encryption.

## Why this matters

Base64 is everywhere: JWT tokens, JSON binary fields, email attachments (MIME), data URLs, and TLS certificates. But many developers confuse base64 with encryption — they think base64-encoded data is "hidden" or "secure." This misunderstanding leads to insecure systems where sensitive data is base64-encoded but not encrypted, giving a dangerous false sense of security. Every Go engineer must know what base64 does, when to use it, and — critically — when not to rely on it for protection.

## Mental model

Base64 is a translation system between binary (bytes) and text (ASCII characters). It takes any byte sequence and represents it using only 64 safe characters: A-Z, a-z, 0-9, +, and /. Think of it like Morse code for binary: Morse maps letters to dots/dashes; base64 maps bytes to characters. Both are encoding schemes that anyone can reverse. Encryption is fundamentally different — it uses a key to transform data so that only someone with the key can reverse it.

## Core idea

Base64 encodes 3 bytes of binary data into 4 ASCII characters. Each group of 3 bytes (24 bits) is split into four 6-bit values, and each 6-bit value (0-63) maps to a character in the base64 alphabet.

Variants:
- **StdEncoding**: Uses `+` and `/` as the 63rd and 64th characters, with `=` padding.
- **URLEncoding**: Uses `-` and `_` instead (safe for URLs and filenames), with `=` padding.
- **RawStdEncoding** / **RawURLEncoding**: Same as above but without `=` padding.

Go's `encoding/base64` package provides:
- `base64.StdEncoding.EncodeToString(src []byte) string`
- `base64.StdEncoding.DecodeString(s string) ([]byte, error)`
- Same methods on `URLEncoding`, `RawStdEncoding`, `RawURLEncoding`
- Streaming: `NewEncoder(enc, w) io.WriteCloser` and `NewDecoder(enc, r) io.Reader`

## Under the hood

Encoding works on a 3-byte → 4-character block basis:

```
Input bytes:    [0x4D] [0x61] [0x6E]
Binary:         01001101 01100001 01101110
6-bit chunks:   010011  010110  000101  101110
Decimal:        19      22      5       46
Base64 chars:   T       W       F       u
Result: "TWFu"
```

If the input is not divisible by 3, padding is added: 1 byte → 2 chars + `==`, 2 bytes → 3 chars + `=`. The decoder reverses this: map each character back to its 6-bit value, concatenate bits, split into bytes, discard padding.

Encoding increases data size by ~33% (4 chars for every 3 bytes).

## How Go uses it

- **JWT tokens**: Header and payload are base64-URL-encoded JSON.
- **JSON binary data**: HTTP APIs send binary data as base64-encoded strings.
- **Data URLs**: `data:image/png;base64,...` in HTML/CSS.
- **gRPC/TLS**: Certificate and key PEM files use base64 internally.
- **Cookie values**: Binary session tokens base64-encoded for HTTP cookies.

## Go example

```go
package main

import (
	"encoding/base64"
	"fmt"
)

func main() {
	// Basic encoding/decoding
	data := []byte("Hello, Go! Base64 is encoding, not encryption.")

	encoded := base64.StdEncoding.EncodeToString(data)
	fmt.Println("StdEncoding:")
	fmt.Println(encoded)

	decoded, _ := base64.StdEncoding.DecodeString(encoded)
	fmt.Println("Decoded:", string(decoded))

	// URL-safe encoding (no + or /)
	urlData := []byte("user:password?token=abc123&sig=xyz==")
	urlEncoded := base64.URLEncoding.EncodeToString(urlData)
	fmt.Println("\nURLEncoding:")
	fmt.Println(urlEncoded)

	urlDecoded, _ := base64.URLEncoding.DecodeString(urlEncoded)
	fmt.Println("URL Decoded:", string(urlDecoded))

	// Raw encoding (no padding)
	raw := base64.RawStdEncoding.EncodeToString(data)
	fmt.Println("\nRawStdEncoding (no padding):")
	fmt.Println(raw)

	// Binary data: encoding a small image's bytes
	binaryData := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}
	b64 := base64.StdEncoding.EncodeToString(binaryData)
	fmt.Println("\nBinary data encoded:", b64)

	// Demonstrate size increase
	original := []byte("This is 33 bytes of text!!!!")
	enc := base64.StdEncoding.EncodeToString(original)
	fmt.Printf("\nOriginal: %d bytes\n", len(original))
	fmt.Printf("Encoded: %d chars (%.0f%% increase)\n", len(enc), float64(len(enc)-len(original))/float64(len(original))*100)
}
```

## Step-by-step execution

1. `data` = `"Hello, Go! Base64 is encoding, not encryption."` (52 bytes).
2. `StdEncoding.EncodeToString(data)` loops over bytes in 3-byte groups.
3. For each group: split 24 bits into four 6-bit values, look up each in the alphabet table (A-Z=0-25, a-z=26-51, 0-9=52-61, +=62, /=63).
4. Last group has 1 byte remaining → encoded as 2 characters + `==` padding.
5. Result is a 72-character ASCII string (52 bytes → 72 chars ≈ 38% increase).
6. `DecodeString` reverses: map chars back to 6-bit values, concatenate, split into bytes.
7. `URLEncoding` uses `-` instead of `+` and `_` instead of `/` — safe in URL query strings.

## Common mistakes

- Mistake: Thinking base64 is encryption — "I'll base64-encode the password so it's secure."
  - Why it happens: The encoded output looks like gibberish, creating a false sense of security.
  - Fix: Use actual encryption (AES, ChaCha20) for confidentiality. Base64 is for transport, not protection.

- Mistake: Using `StdEncoding` in URLs — `+` and `/` get URL-encoded to `%2B` and `%2F`.
  - Fix: Use `base64.URLEncoding` for URL-safe output.

- Mistake: Forgetting to handle decoding errors — invalid base64 strings cause an error.
  - Fix: Always check the error from `DecodeString`.

- Mistake: Using base64 for large binary data without considering the 33% size increase.
  - Fix: Consider binary formats (Protobuf, msgpack) or compression before base64.

## Debugging walkthrough

JWT verification fails:

```go
import "encoding/base64"

header := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
decoded, _ := base64.StdEncoding.DecodeString(header)
fmt.Println(string(decoded)) // garbage or error
```

**Symptom**: Decoding produces wrong output or an error.

**Root cause**: JWT uses `base64.RawURLEncoding` (no padding, URL-safe characters). Using `StdEncoding` fails because `_` and `-` are not valid in StdEncoding.

**Fix**:
```go
decoded, _ := base64.RawURLEncoding.DecodeString(header)
fmt.Println(string(decoded)) // {"alg":"HS256","typ":"JWT"}
```

## Production notes

- **Base64 is not encryption, hashing, or obfuscation**. It is encoding for safe transport. Anyone can decode it instantly.
- **Use URLEncoding for query parameters**, filenames, and JWT — avoid `+` and `/`.
- **Use RawEncoding (no padding)** when the length is known or when padding causes issues (JWT, some APIs).
- **Base64 in JSON**: Binary data embedded in JSON must be base64-encoded because JSON has no native binary type.
- **Streaming**: For large data (files > 10 MB), use `base64.NewEncoder` / `NewDecoder` to avoid loading the entire data into memory.

## Performance implications

- Base64 encoding/decoding is CPU-efficient: ~500 MB/s per core on modern hardware.
- The 33% size increase means more network bandwidth and storage for base64-encoded data.
- Decoding is slightly slower than encoding (character lookup + bit manipulation).
- Streaming encoders/decoders add minimal overhead — use them for large payloads to save memory.
- Go's base64 implementation is optimized with lookup tables and SIMD on some architectures.

## Practice task

Write a function `EncodeDecode(original []byte) (encoded string, decoded []byte, err error)` that:
1. Encodes the original bytes using `base64.URLEncoding`.
2. Decodes the encoded string back.
3. Returns the encoded string, decoded bytes, and any error.
Test with both text and binary data (e.g., `[]byte{0, 255, 128, 64}`). In `main`, demonstrate that encoded output is URL-safe (no `+` or `/` characters).

## Tests / verification

```bash
go run ./curriculum/modules/07-cli-files-json-config/lessons/20-base64-as-encoding-not-encryption
go test ./curriculum/modules/07-cli-files-json-config/lessons/20-base64-as-encoding-not-encryption
```

## Review questions

1. What is the purpose of base64 encoding?
2. Why is base64 NOT encryption? What's the key difference?
3. What encoding variant should you use for URL query parameters? Why?
4. By what percentage does base64 increase data size?
5. What does the `=` padding character signify in base64 strings?

## NEXT UP

CLI testability — testing command-line applications in Go.
