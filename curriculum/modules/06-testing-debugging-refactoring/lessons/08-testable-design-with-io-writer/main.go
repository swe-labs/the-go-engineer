package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func Greet(w io.Writer, name string) {
	fmt.Fprintf(w, "Hello, %s!\n", name)
}

func NewLogger(w io.Writer, prefix string) *log.Logger {
	return log.New(w, prefix, 0)
}

func WriteTable(w io.Writer, rows [][2]string) {
	fmt.Fprintln(w, "Name    | Age")
	fmt.Fprintln(w, "--------|-----")
	for _, row := range rows {
		fmt.Fprintf(w, "%-7s | %s\n", row[0], row[1])
	}
}

func main() {
	Greet(os.Stdout, "Alice")

	logger := NewLogger(os.Stdout, "APP: ")
	logger.Println("started")

	WriteTable(os.Stdout, [][2]string{
		{"Alice", "30"},
		{"Bob", "25"},
	})

	var buf strings.Builder
	WriteTable(&buf, [][2]string{{"Test", "1"}})
	fmt.Print("Captured:\n", buf.String())
}
