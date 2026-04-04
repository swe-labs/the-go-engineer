# Section 11: Encoding

## Beginner → Expert Mapping

| Topic | Level | Importance | Engineering Concept |
|-------|-------|------------|---------------------|
| JSON Marshal | Beginner | High | Object serialization via reflection |
| JSON Streams | Intermediate | High | Stream unmarshalling (`json.Decoder`) |
| Custom Marshal | Advanced | Medium | Overriding byte representation |

## Engineering Depth
In Go, `json.Marshal()` operates using Reflection (`reflect` package) to map struct keys to byte arrays at runtime. Because reflection bypasses strict compile-time type safety, it is computationally expensive.
- **Memory Scaling:** `json.Unmarshal(data)` loads the entire `[]byte` representation into memory, allocating massive heap blocks for large arrays. 
- **Production Pattern:** Always use `json.NewDecoder(io.Reader).Decode()` to stream gigabytes of JSON incrementally with $O(1)$ memory footprints.

## References
1. **[Go Blog]** [JSON and Go](https://go.dev/blog/json)

---

## 🏗 Exercise: Config File Parser (`6-config-parser`)

Build a JSON config file parser that reads, decodes, and validates application configuration. Try it yourself with the `_starter/` directory first!

```bash
go run ./11-encoding/6-config-parser/_starter   # Try the exercise
go run ./11-encoding/6-config-parser            # See the solution
```
