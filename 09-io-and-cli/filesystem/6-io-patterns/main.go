// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

// ============================================================================
// Section 09: Filesystem — io.Reader/Writer Patterns
// Level: Advanced
// ============================================================================
//
// RUN: go run ./09-io-and-cli/filesystem/6-io-patterns
// ============================================================================

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// ============================================================================
// Section 09: io.Reader and io.Writer Patterns
// Level: Intermediate → Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - io.Reader and io.Writer — Go's most important interfaces
//   - Composing readers and writers (pipes, tee, multi)
//   - Implementing custom readers
//   - io.Copy — the universal glue
//
// WHY THIS MATTERS:
//   Almost everything in Go implements io.Reader or io.Writer:
//   - os.File, http.Response.Body, bytes.Buffer, strings.Reader,
//     json.Encoder, gzip.Writer, crypto/hash, net.Conn...
//
// THE INTERFACES:
//   type Reader interface { Read(p []byte) (n int, err error) }
//   type Writer interface { Write(p []byte) (n int, err error) }
//
// ENGINEERING DEPTH:
//   The `io.Reader` and `io.Writer` interfaces are the most profound architectural
//   designs in Go. By defining I/O strictly as "fill this byte slice" (`Read`)
//   or "drain this byte slice" (`Write`), Go completely detaches the concept of
//   I/O from physical hardware. A function accepting an `io.Writer` has absolutely
//   no idea if it is writing to a text file, an active TCP socket across the globe,
//   an AWS S3 bucket, or an AES cryptographic cipher. This enables infinite
//   composition: you can pipe a Network Socket through an Unzipper through a
//   JSON Decoder directly into a Struct without ever buffering the 1GB payload
//   into RAM.
// ============================================================================

func main() {
	fmt.Println("=== io.Reader and io.Writer Patterns ===")
	fmt.Println()

	// 1. strings.Reader — turn a string into an io.Reader
	readerFromString()

	// 2. bytes.Buffer — implements BOTH Reader and Writer
	bufferDemo()

	// 3. io.Copy — stream data from Reader to Writer
	copyDemo()

	// 4. io.TeeReader — read and copy simultaneously
	teeReaderDemo()

	// 5. io.MultiReader — concatenate readers
	multiReaderDemo()

	// 6. io.MultiWriter — write to multiple destinations at once
	multiWriterDemo()

	// 7. Custom Reader — implement the interface yourself
	customReaderDemo()
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: FS.7 log search tool")
	fmt.Println("   Current: FS.6 (io.Reader / io.Writer patterns)")
	fmt.Println("---------------------------------------------------")
}

func readerFromString() {
	fmt.Println("--- 1. strings.Reader ---")
	// strings.NewReader creates an io.Reader from a string.
	// Useful when a function expects io.Reader but you have a string.
	reader := strings.NewReader("Hello, Go!")

	// Read into a byte slice
	buf := make([]byte, 5)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			fmt.Printf("  Read %d bytes: %q\n", n, buf[:n])
		}
		if err == io.EOF {
			break
		}
	}
	fmt.Println()
}

func bufferDemo() {
	fmt.Println("--- 2. bytes.Buffer ---")
	// bytes.Buffer implements io.Reader, io.Writer, and io.ByteReader.
	// It's the Swiss Army knife of Go I/O.
	var buf bytes.Buffer

	// Write to buffer (io.Writer)
	buf.WriteString("Hello ")
	buf.WriteString("World")
	fmt.Printf("  Buffer contents: %q\n", buf.String())

	// Read from buffer (io.Reader)
	data := make([]byte, 5)
	buf.Read(data)
	fmt.Printf("  Read 5 bytes: %q\n", data)
	fmt.Printf("  Remaining: %q\n", buf.String())
	fmt.Println()
}

func copyDemo() {
	fmt.Println("--- 3. io.Copy ---")
	// io.Copy streams data from any Reader to any Writer.
	// This is the universal glue that connects all I/O in Go.
	src := strings.NewReader("Streaming data with io.Copy!")
	n, err := io.Copy(os.Stdout, src)
	fmt.Printf("\n  Copied %d bytes (err=%v)\n\n", n, err)
}

func teeReaderDemo() {
	fmt.Println("--- 4. io.TeeReader ---")
	// TeeReader returns a Reader that writes to w what it reads from r.
	// Like the Unix `tee` command — read and copy simultaneously.
	var log bytes.Buffer
	src := strings.NewReader("This gets logged AND processed")

	tee := io.TeeReader(src, &log) // Everything read from tee is also written to log

	// Read through the tee reader
	content, _ := io.ReadAll(tee)
	fmt.Printf("  Processed: %q\n", content)
	fmt.Printf("  Logged:    %q\n", log.String())
	fmt.Println()
}

func multiReaderDemo() {
	fmt.Println("--- 5. io.MultiReader ---")
	// MultiReader concatenates multiple readers into one.
	// Reads from each in sequence — like cat in Unix.
	header := strings.NewReader("[HEADER] ")
	body := strings.NewReader("message body")
	footer := strings.NewReader(" [END]")

	combined := io.MultiReader(header, body, footer)
	result, _ := io.ReadAll(combined)
	fmt.Printf("  Combined: %q\n\n", result)
}

func multiWriterDemo() {
	fmt.Println("--- 6. io.MultiWriter ---")
	// MultiWriter creates a writer that duplicates writes to all provided writers.
	// Like writing to stdout AND a log file simultaneously.
	var logBuf bytes.Buffer
	multi := io.MultiWriter(os.Stdout, &logBuf)

	fmt.Fprint(multi, "  This goes to BOTH stdout and the buffer")
	fmt.Printf("\n  Buffer captured: %q\n\n", logBuf.String())
}

func customReaderDemo() {
	fmt.Println("--- 7. Custom Reader ---")
	// Implementing io.Reader: just implement Read(p []byte) (n int, err error)
	counter := &CountingReader{limit: 5}
	data, _ := io.ReadAll(counter)
	fmt.Printf("  Read from custom reader: %v\n", data)
	fmt.Println()
}

// CountingReader is a custom io.Reader that produces ascending integers.
// This demonstrates how simple it is to implement the io.Reader interface.
type CountingReader struct {
	current int
	limit   int
}

func (r *CountingReader) Read(p []byte) (int, error) {
	if r.current >= r.limit {
		return 0, io.EOF // Signal end of data
	}
	p[0] = byte(r.current)
	r.current++
	return 1, nil
}
