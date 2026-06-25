package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	// Multi-writer: write to stdout and compute SHA-256 simultaneously
	var buf strings.Builder
	hash := sha256.New()
	multi := io.MultiWriter(os.Stdout, &buf, hash)

	writer := bufio.NewWriter(multi)
	writer.WriteString("Hello, io.Writer!\n")
	fmt.Fprintf(writer, "Line %d\n", 2)
	writer.Flush()

	fmt.Printf("\nSHA-256: %s\n", hex.EncodeToString(hash.Sum(nil)))
}
