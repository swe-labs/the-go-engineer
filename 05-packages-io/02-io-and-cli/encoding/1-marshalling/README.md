# EN.1 Marshalling

## Mission

Learn how to convert Go structs into JSON byte slices (Marshalling) and use struct tags to control the output format.

## Prerequisites

- `CL.4` file-organizer-project (completion of CLI track)

## Mental Model

Think of Marshalling as **Packaging for Shipping**.

You have a "Product" (your Go struct) that you want to send over the "Internet" (a network connection). To do this, you need to fold it up, put it in a "Standard Box" (JSON), and label the boxes correctly so the person on the other end knows which piece is which.

## Visual Model

```mermaid
graph LR
    A["Go Struct {ID: 1, Name: 'Book'}"] --> B["json.Marshal"]
    B --> C["JSON String: '{\"id\": 1, \"name\": \"Book\"}'"]
```

## Machine View

The `encoding/json` package uses the `reflect` package to inspect your structs at runtime. It looks for exported fields (starting with a capital letter) and reads the metadata attached to them (struct tags). It then builds a buffer and writes the JSON representation of your data into a byte slice. Because reflection is used, only fields that the `json` package can "see" (exported fields) will be included in the output.

## Run Instructions

```bash
go run ./05-packages-io/02-io-and-cli/encoding/1-marshalling
```

## Code Walkthrough

### Struct Tags
The backticks after a field definition (e.g., `` `json:"id"` ``) are called struct tags. They provide instructions to the JSON encoder on how to name the field in the final output.

### `json.Marshal`
Converts a Go value to a compact JSON byte slice. If an error occurs (e.g., trying to marshal a function or a channel), it will return an error.

### `json.MarshalIndent`
Similar to `Marshal`, but adds whitespace and newlines to make the output "pretty" and human-readable. This is extremely useful for debugging or CLI output.

### `omitempty`
A tag option that tells the encoder to skip a field if it contains its "zero value" (like `0`, `""`, `false`, or `nil`).

## Try It

1. Change a field name in the struct tags and see how the JSON output changes.
2. Add a new field to the `Product` struct but keep it unexported (lowercase) and observe that it doesn't appear in the JSON.
3. Use the `json:",string"` tag on the `ID` field and see it get wrapped in quotes in the JSON.

## In Production
JSON marshalling using the standard library is powerful but uses reflection, which can be a bottleneck in extremely high-throughput systems. For most applications, it is more than fast enough. If you are handling millions of requests per second, you might look at code-generation libraries like `easyjson`.

## Thinking Questions
1. Why must a field be exported (capitalized) to be included in the JSON output?
2. When would you use `json.Marshal` versus `json.MarshalIndent`?
3. How do you handle sensitive data (like passwords) using struct tags?

> **Forward Reference:** You have learned how to turn Go data into JSON. But what about the reverse? In [Lesson 2: Unmarshal](../2-unmarshal/README.md), you will learn how to take a JSON string and turn it back into a Go struct.

## Next Step

Next: `EN.2` -> `05-packages-io/02-io-and-cli/encoding/2-unmarshal`

Open `05-packages-io/02-io-and-cli/encoding/2-unmarshal/README.md` to continue.
