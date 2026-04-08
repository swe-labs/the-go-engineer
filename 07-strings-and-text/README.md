# Section 7: Strings & Text Processing

## Beginner → Expert Mapping

| Topic | Level | Importance | Engineering Concept |
|-------|-------|------------|---------------------|
| String Internals | Beginner | High | Immutable byte slices, UTF-8 strings |
| Regexp | Intermediate | Medium | Regular expressions engine |
| Text Templates | Intermediate | High | Dynamic layout generation |

## Engineering Depth
A `string` in Go is a read-only slice of bytes, typically representing UTF-8 encoded text. Since strings are immutable, concatenation using `+=` behaves exceptionally poorly inside loops, operating at $O(N^2)$ due to recursive reallocation. For dynamic text construction, **always use `strings.Builder`** which allocates $O(N)$ memory dynamically just like an integer slice. 

## References
1. **[Go Blog]** [Strings, bytes, runes and characters in Go](https://go.dev/blog/strings)

---

## 🏗 Exercise: Log File Parsing System (`6-log-parser`)

### Step-by-Step Instructions & Hints
1. **Load data:** Read lines from a simulated file or multi-line string.
   - *Hint:* Use `strings.Split()` to tokenize the file.
2. **Apply regex:** Compile a `regexp.MustCompile()` pattern out of the loop.
   - *Hint:* Compiling a regex inside a loop is devastating to performance. Compile it at the package level or before the loop!
3. **Build the output:** Extract matched substrings.
   - *Hint:* Use `strings.Builder` to combine the matching items into a single summary output instead of `+=`.


## Learning Path

| ID | Lesson | Concept | Requires |
| --- | --- | --- | --- |
| ST.1 | [strings](./1-strings) | Immutable byte slice · strings package · strings.Builder | 🟢 entry |
| ST.2 | [formatting](./2-formatting-string) | fmt verbs · width/precision · Sprintf vs Printf vs Fprintf | ST.1 |
| ST.3 | [unicode & runes](./3-unicode) | UTF-8 multi-byte · rune · for-range vs byte iteration | ST.1, ST.2 |
| ST.4 | [regex](./4-regex) | MustCompile · FindAllString · capture groups · ReplaceAllStringFunc | ST.1, ST.3 |
| ST.5 | [text templates](./5-text-template) | text/template · {{range}} · {{if}} · Execute to io.Writer | ST.1, ST.2 |
