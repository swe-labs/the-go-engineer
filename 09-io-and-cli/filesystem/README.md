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


## Learning Path

| ID | Lesson | Concept | Requires |
| --- | --- | --- | --- |
| FS.1 | [files](./1-files) | os.WriteFile · os.ReadFile · bufio.Scanner · os.OpenFile flags | 🟢 entry |
| FS.2 | [paths](./2-paths) | filepath.Join · Base · Dir · Ext · Abs · Clean · Glob | FS.1 |
| FS.3 | [directories](./3-dir) | MkdirAll · ReadDir · WalkDir · RemoveAll · os.Stat | FS.1, FS.2 |
| FS.4 | [temp files](./4-temp) | MkdirTemp · CreateTemp · unique names · immediate defer | FS.1, FS.3 |
| FS.5 | [embed](./5-embed) | //go:embed · string/[]byte · embed.FS · single-binary deployment | FS.3, FS.4 |
| FS.6 | [io.Reader / io.Writer patterns](./6-io-patterns) | strings.Reader · bytes.Buffer · io.Copy · TeeReader · MultiWriter | FS.1, FS.2, FS.3 |
| **FS.8** ⭐ | [fs.FS testing seam](./8-fs-testing-seam) | os.DirFS · fstest.MapFS · zero-disk-IO tests | FS.5, FS.6 |
