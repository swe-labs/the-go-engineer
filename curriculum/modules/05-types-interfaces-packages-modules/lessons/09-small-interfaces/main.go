package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
)

type UpperWriter struct {
	w io.Writer
}

func (u UpperWriter) Write(p []byte) (int, error) {
	return u.w.Write(bytes.ToUpper(p))
}

type CountingWriter struct {
	w     io.Writer
	count int64
}

func (cw *CountingWriter) Write(p []byte) (int, error) {
	n, err := cw.w.Write(p)
	cw.count += int64(n)
	return n, err
}

func (cw *CountingWriter) BytesWritten() int64 {
	return cw.count
}

func main() {
	r := strings.NewReader("hello world")
	io.Copy(os.Stdout, r)
	fmt.Println()

	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	gw.Write([]byte("compressed data"))
	gw.Close()
	fmt.Printf("gzip wrote %d bytes\n", buf.Len())

	UpperWriter{w: os.Stdout}.Write([]byte("hello\n"))

	cw := &CountingWriter{w: os.Stdout}
	fmt.Fprint(cw, "counting writer test\n")
	fmt.Printf("bytes written: %d\n", cw.BytesWritten())
}
