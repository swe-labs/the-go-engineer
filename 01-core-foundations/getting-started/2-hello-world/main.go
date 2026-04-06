// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 0: Getting Started — Hello, World!
// Level: Beginner
// ============================================================================
//
// This is the most important file in the entire The Go Engineer curriculum.
// Every single line is explained. Read every comment carefully.
//
// WHAT YOU'LL LEARN:
//   - The structure of every Go program
//   - What "package main" means
//   - What "import" does
//   - What "func main()" is
//   - How to print text to the screen
//
// ENGINEERING DEPTH:
//   Go is a statically typed, compiled language. Unlike Python or JavaScript
//   which are interpreted line-by-line at runtime, Go translates this entire file
//   into raw machine code upfront. This eliminates an entire class of syntax
//   errors from ever reaching production.
//
// RUN: go run ./01-core-foundations/getting-started/2-hello-world
// ============================================================================

// --- LINE-BY-LINE EXPLANATION ---

// LINE 1: "package main"
//
// Every Go file starts with a package declaration. A "package" is Go's way
// of organizing code into reusable groups (like folders for your code).
//
// The package name "main" is SPECIAL. It tells the Go compiler:
//   "This is an executable program, not a library."
//
// If you change "main" to anything else (like "myapp"), Go will refuse
// to run it because it won't know where to start.
//
// RULE: Every executable Go program MUST have exactly one package named "main".

// LINES 3: "import "fmt""
//
// The "import" keyword brings in code from other packages so you can use it.
//
// "fmt" is Go's standard formatting package. The name comes from "format".
// It provides functions for:
//   - Printing text: fmt.Println(), fmt.Printf(), fmt.Print()
//   - Formatting strings: fmt.Sprintf()
//   - Reading input: fmt.Scan(), fmt.Scanln()
//
// Go's standard library has ~150 packages built in. You never need to
// download "fmt" — it ships with Go. Think of it like Python's "print()"
// or JavaScript's "console.log()", but more powerful.
//
// RULE: Go will NOT compile if you import a package and don't use it.
//       This prevents dead code from cluttering your project.

// LINES 53+: "func main()"
//
// "func" declares a function. Functions are reusable blocks of code.
//
// "main" is the name of this function. Like "package main", the function
// named "main" is SPECIAL. It is the ENTRY POINT of your program.
//
// When you run "go run main.go", Go does this:
//  1. Finds "package main"
//  2. Finds "func main()"
//  3. Starts executing the code inside the braces { }
//
// The parentheses () after "main" are for parameters (inputs).
// main() takes no parameters — it always starts with an empty ().
//
// The opening brace { MUST be on the same line as "func main()".
// This is NOT a style choice — Go will not compile if you put { on the next line.
// This eliminates all brace-style debates. Everyone writes the same way.
func main() {

	// fmt.Println() prints text to your terminal, followed by a newline character.
	// "ln" in Println stands for "line" — it adds a line break at the end.
	//
	// The text inside the quotes is called a "string literal".
	// In Go, strings are always enclosed in double quotes "".
	// Single quotes '' are for individual characters (called "runes" in Go).
	fmt.Println("Hello, World! Welcome to The Go Engineer.")

	// You can print multiple values separated by commas.
	// Println automatically adds a space between each value.
	fmt.Println("Go was created at", "Google", "in", 2009)

	// fmt.Printf() gives you more control over formatting.
	// %s = string, %d = integer, %f = float, %v = any value (auto-detect)
	// \n = newline character (Printf does NOT add one automatically)
	language := "Go" // Short variable declaration (covered in Section 01)
	year := 2009
	fmt.Printf("%s was created in %d\n", language, year)

	// WHAT JUST HAPPENED (behind the scenes):
	//
	// 1. You ran:     go run ./01-core-foundations/getting-started/2-hello-world
	// 2. Go compiled: Translated this entire file to native MACHINE CODE
	//                 (the binary language your CPU speaks — 1s and 0s)
	// 3. Go executed: Started at func main() and ran each line top-to-bottom
	// 4. Output:      Println sent text to "stdout" (your terminal)
	//
	// This is fundamentally different from Python or JavaScript:
	//   - Python/JS: Code is read and executed line-by-line (interpreted)
	//   - Go:        Code is compiled to machine code FIRST, then executed
	//
	// Why does this matter?
	//   - Compile-time errors: If you have a typo, Go catches it BEFORE running
	//   - Speed: Compiled code runs 10-100x faster than interpreted code
	//   - Single binary: "go build" produces one file with zero dependencies

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: GS.3 how-go-works")
	fmt.Println("   Current: GS.2 (hello-world)")
	fmt.Println("---------------------------------------------------")
}
