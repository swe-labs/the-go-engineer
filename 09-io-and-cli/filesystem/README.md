# Section 10: Filesystem Operations

## Beginner → Expert Mapping

| Topic | Level | Importance | Engineering Concept |
|-------|-------|------------|---------------------|
| Basic I/O | Beginner | High | `os` File interactions |
| File Traversal | Intermediate | Medium | `filepath.Walk` vs `os.ReadDir` |
| Advanced I/O Patterns | Advanced | **Critical** | `io.Reader` and `io.Writer` composition |
| Embed | Advanced | Medium | Compilation binary embedding (`go:embed`) |

## Engineering Depth

The `io.Reader` and `io.Writer` interfaces are the beating heart of Go's standard library.

- Memory constraint: Direct loading via `os.ReadFile()` loads the entire target to Heap memory ($O(N)$ Space). For large multi-Gigabyte files, this leads to immediate memory exhaustion (OOM).
- To scale cleanly to infinite file sizes ($O(1)$ Space), always read using `bufio.NewScanner()` on the open `*os.File` which acts as a streaming `io.Reader`.

## References

1. **[Go Package]** [io standard library rules](https://pkg.go.dev/io)

---

## 🏗 Exercise: Log Search Tool (`7-log-search`)

### Step-by-Step Instructions & Hints

1. **Initialize the Traversal:** Use `filepath.Walk()` to search a given directory.
2. **Filter criteria:** Skip non-txt/log files based on extensions.
3. **Stream target contents:** When you open a valid file, do not load the whole file. Initialize a `bufio.Scanner` attached to the file reader.
   - *Hint:* By reading line-by-line via `scanner.Scan()`, the program consumes mere kilobytes of RAM, safely navigating logs of any size!
4. **Log the results:** Print the Line Number and Line Content whenever a target regex string matches.
