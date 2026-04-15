# GT.3 How Go Works

## Mission

Build a beginner-safe mental model for packages, imports, compilation, and exported names.

This lesson is not trying to teach every compiler detail.
It is trying to make the basic workflow feel explainable.

## Why This Lesson Exists Now

After `hello world`, beginners often know that the code runs but not why it runs.

This lesson answers a few important "behind the scenes" questions:

- What is a package?
- What does `import` really do?
- Why can one file call `strings.ToUpper(...)`?
- What happens when `go run` is executed?

## Mental Model

Go organizes code into packages.
A package groups related code and gives it a name.

This file uses several packages:

- `fmt` for printing
- `strings` for string operations
- `math` for mathematical functions

The program does not own those tools.
It borrows them through imports.

## Visual Model

```text
this file
  |
  +--> fmt.Println(...)
  +--> strings.ToUpper(...)
  +--> strings.Split(...)
  +--> math.Sqrt(...)
```

```text
go run
  |
  +--> read source
  +--> compile program
  +--> run executable
  +--> show terminal output
```

## Machine View

The Go toolchain builds this program before it runs.

That means:

- imported package code is resolved during the build
- function calls are type-checked before execution
- the final program starts only after compilation succeeds

So when the learner sees printed output, they are seeing the result of a compiled executable, not a
line-by-line interpreter loop.

## Run Instructions

```bash
go run ./01-foundations/01-getting-started/3-how-go-works
```

## Code Walkthrough

### `import ("fmt" "math" "strings")`

This lesson imports more than one package because it wants to show that one program can use several
different tools at once.

### `greeting := "hello, go developer!"`

The program begins with a plain string value.
That gives the `strings` package something to work on.

### `strings.ToUpper(greeting)`

This calls a function named `ToUpper` from the `strings` package.

The capital `T` matters.
In Go, names that start with an uppercase letter are exported, which means other packages may use
them.

### `strings.Contains(greeting, "go")`

This asks a yes-or-no question about the string and returns a boolean answer.
The learner does not need full boolean theory yet.
They only need to see that package functions can return useful results, not only print text.

### `strings.Split("one,two,three", ",")`

This turns one string into multiple parts.
It is an early preview that packages can transform data, not just display it.

### `math.Pi`, `math.Sqrt(144)`, and `math.Pow(2, 10)`

These lines show two things:

- packages can expose values like `Pi`
- packages can expose functions like `Sqrt` and `Pow`

That helps the learner see imports as access to organized capabilities.

## Try It

1. Change `greeting` and see how `ToUpper` and `Contains` react.
2. Change `"one,two,three"` to another comma-separated string.
3. Replace `math.Sqrt(144)` with `math.Sqrt(81)` and rerun.

## Common Questions

- Why do package functions use `packageName.FunctionName(...)`?
  Because Go makes package ownership explicit at the call site.

- Why are some names uppercase?
  In Go, uppercase names are exported for use from other packages.

## Next Step

Continue to `GT.4` development environment.
