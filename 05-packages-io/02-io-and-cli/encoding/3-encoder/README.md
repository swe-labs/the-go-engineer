# EN.3 Encoder

## Mission

Learn how to use the streaming JSON Encoder to write data directly to an `io.Writer`, improving memory efficiency for large datasets and web responses.

## Prerequisites

- `EN.2` unmarshal

## Mental Model

Think of the Encoder as a **Continuous Conveyor Belt**.

Instead of waiting for an entire product to be assembled and boxed (`json.Marshal`), the Encoder takes parts as they come and places them directly onto a conveyor belt (the `io.Writer`) that leads to the shipping dock (the network or disk).

## Visual Model

```mermaid
graph LR
    A["Go Data"] --> B["json.NewEncoder(w)"]
    B -- "Stream" --> C["io.Writer (File/Socket)"]
```

## Machine View

When you use `json.NewEncoder(w).Encode(v)`, Go does not allocate a single large byte slice for the entire JSON payload. Instead, it uses an internal buffer to serialize the data and writes that buffer incrementally to the `io.Writer`. This is particularly important for HTTP servers using `http.ResponseWriter`, as it allows the server to start sending the response headers and the first bytes of the body to the client immediately, rather than waiting for the entire struct to be marshalled in memory.

## Run Instructions

```bash
go run ./05-packages-io/02-io-and-cli/encoding/3-encoder
```

## Code Walkthrough

### `json.NewEncoder(w)`
Initializes a new encoder that will write to the provided `io.Writer` (e.g., `os.Stdout`, an `os.File`, or an `http.ResponseWriter`).

### `enc.SetIndent(prefix, indent)`
Configures the encoder to produce "pretty-printed" JSON. This is similar to `json.MarshalIndent`.

### `enc.Encode(v)`
Serializes the value `v` and writes it to the underlying stream, followed by a newline character (`\n`). This behavior makes the Encoder perfect for generating **JSONL (JSON Lines)** files.

## Try It

1. Change the output of the first example from `os.Stdout` to a buffer using `bytes.Buffer`.
2. Create a loop that encodes 10 different objects and observe how they are separated by newlines in the output.
3. Compare the time taken to marshal a large slice versus encoding it directly to a file (for very large slices, the memory savings will be obvious).

## In Production
For production APIs, `json.NewEncoder(w).Encode(v)` is the standard way to return JSON responses. It keeps your server's memory footprint low and predictable, which is essential for scaling under heavy load.

## Thinking Questions
1. When would you prefer `json.Marshal` over `json.NewEncoder`?
2. What is "JSON Lines" (JSONL) format, and why is it useful for logging?
3. How does the Encoder handle large slices of data?

> [!TIP]
> You have learned how to stream data OUT. Now we will look at the reverse: how to efficiently read a continuous stream of JSON data IN. In [Lesson 4: Decode](../4-decode/README.md), you will learn about the `json.Decoder`.

## Next Step

Next: `EN.4` -> [`05-packages-io/02-io-and-cli/encoding/4-decode`](../4-decode/README.md)
