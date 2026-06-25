package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func Run(r io.Reader, w io.Writer, args []string) int {
	fs := flag.NewFlagSet("wc", flag.ContinueOnError)
	fs.SetOutput(w)
	countLines := fs.Bool("l", false, "count lines (default)")
	countWords := fs.Bool("w", false, "count words")
	countChars := fs.Bool("c", false, "count characters")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	// If no mode flag set, default to lines
	if !*countLines && !*countWords && !*countChars {
		*countLines = true
	}

	scanner := bufio.NewScanner(r)
	var lines, words, chars int

	if *countWords || *countChars {
		for scanner.Scan() {
			line := scanner.Text()
			lines++
			chars += len(line) + 1 // +1 for newline
			words += len(strings.Fields(line))
		}
	} else {
		for scanner.Scan() {
			lines++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(w, "read error:", err)
		return 1
	}

	switch {
	case *countLines && !*countWords && !*countChars:
		fmt.Fprintln(w, lines)
	case *countWords:
		fmt.Fprintln(w, words)
	case *countChars:
		fmt.Fprintln(w, chars)
	default:
		fmt.Fprintln(w, lines)
	}
	return 0
}

func main() {
	os.Exit(Run(os.Stdin, os.Stdout, os.Args[1:]))
}
