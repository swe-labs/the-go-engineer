package main

import (
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	data := "Hello, io.Reader!"

	// Chain: base64 encoded gzip data -> base64 decoder -> gzip reader -> read
	var b strings.Builder
	gzWriter := gzip.NewWriter(base64.NewEncoder(base64.StdEncoding, &b))
	gzWriter.Write([]byte(data))
	gzWriter.Close()

	encoded := b.String()
	fmt.Println("Encoded:", encoded[:40], "...")

	// Read it back through chained readers
	base64Reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(encoded))
	gzReader, err := gzip.NewReader(base64Reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating gzip reader: %v\n", err)
		os.Exit(1)
	}
	defer gzReader.Close()

	decoded, err := io.ReadAll(gzReader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Decoded:", string(decoded))
}
