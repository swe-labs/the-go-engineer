# GT.2 Hello World

## Mission

Learn the smallest useful shape of an executable Go program.

This lesson teaches the learner what a Go program must have in order to run at all.

## Why This Lesson Exists Now

After installation works, the learner needs to see a complete program that still feels small enough
to understand line by line.

That means learning:

- why `package main` exists
- why `func main()` matters
- how `import` gives access to standard-library code
- how printed output reaches the terminal

## Mental Model

A minimal Go program has a clear shape:

1. declare the package
2. import what the file needs
3. define `main`
4. run statements inside `main`

That shape repeats through the whole curriculum.

## Visual Model

```text
package main
import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

```text
program shape:

package -> imports -> main function -> output
```

## Machine View

When you run this lesson, Go does not read one line and immediately print it the way a shell
script might feel.

Instead, the `go` tool:

1. compiles the source file into a runnable program
2. starts execution at `main`
3. executes each statement inside `main`
4. writes output to standard output, which the terminal displays

That is why Go can catch many mistakes before the program ever starts.

## Run Instructions

```bash
go run ./01-foundations/01-getting-started/2-hello-world
```

## Code Walkthrough

### `package main`

This line tells Go the file belongs to the executable package.
Without `package main`, the `go run` flow would not know to build a runnable program from this
file.

### `import "fmt"`

The file needs Go's formatting package so it can print output.
Imports are how one package gains access to code from another package.

### `func main() {`

This is the program entry point.
Execution begins here.
If the file has helper code but no `main` function, the program cannot start as an executable.

### `fmt.Println("Hello, World! Welcome to The Go Engineer.")`

This prints a full line of text.
`Println` adds a newline at the end, so the next output starts on a new line.

### `fmt.Println("Go was created at", "Google", "in", 2009)`

This shows that `Println` can print more than one value.
It inserts spaces between the values automatically.

### `language := "Go"` and `year := 2009`

These two lines store small values in variables.
The section is not formally teaching variables yet.
It is only showing that printed output can come from named values as well as direct text.

### `fmt.Printf("%s was created in %d\n", language, year)`

`Printf` gives more control than `Println`.

This line says:

- `%s` should be replaced by a string
- `%d` should be replaced by an integer
- `\n` should end the line

That helps the learner see that output can be shaped, not only dumped.

## Try It

1. Change the welcome message text and rerun the lesson.
2. Change `year := 2009` to another number and inspect the final line.
3. Add one more `fmt.Println(...)` line below the existing output.

## Common Questions

- Why is `main` special?
  Because executable Go programs start there.

- Why do we need `fmt` just to print?
  Because printing functionality lives in a package, and Go makes package use explicit.

## Next Step

Continue to `GT.3` how Go works.
