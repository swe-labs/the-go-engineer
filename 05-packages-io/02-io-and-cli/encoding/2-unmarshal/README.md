# EN.2 Unmarshal

## Mission

Learn how to convert JSON byte slices back into Go structs (Unmarshalling) and handle complex or dynamic JSON data shapes.

## Prerequisites

- `EN.1` marshalling

## Mental Model

Think of Unmarshalling as **Unpacking a Delivery**.

You received a "Standard Box" (JSON string) from the "Internet". To use the item inside, you need to:
1. **Identify the Item**: Compare the labels on the box to your "Assembly Guide" (the Go struct).
2. **Unpack**: Carefully move each piece from the box into the correct slot in your project.
3. **Handle Missing Parts**: If the box is missing a piece, you use a "Default Part" (the zero value).
4. **Ignore Junk**: If the box contains extra pieces you don't recognize, you just throw them away.

## Visual Model

```mermaid
graph LR
    A["JSON String: '{\"id\": 1, \"name\": \"Book\"}'"] --> B["json.Unmarshal"]
    B --> C["Go Struct {ID: 1, Name: 'Book'}"]
```

## Machine View

When you call `json.Unmarshal(data, &target)`, you must pass a **pointer** to the target variable. This is because Go functions receive copies of their arguments. If you passed a struct directly, `Unmarshal` would populate a copy of the struct and your original variable would remain empty. By passing a pointer, you give the `json` package the memory address of your variable, allowing it to "reach out" and write the parsed data directly into your memory space.

## Run Instructions

```bash
go run ./05-packages-io/02-io-and-cli/encoding/2-unmarshal
```

## Code Walkthrough

### `json.Unmarshal`
The inverse of `json.Marshal`. It takes a byte slice and a pointer to a Go value. It returns an error if the JSON is malformed or if the data types don't match (e.g., trying to put a JSON string into a Go `int`).

### Nested Structs
If the JSON contains objects within objects, you model this in Go using structs within structs. The `json` package handles this recursively.

### `map[string]any`
When you don't know the keys or the shape of the JSON ahead of time, you can unmarshal into a map. Note that JSON numbers (which don't have a fixed size) are always unmarshalled into `float64` when using a map.

### Zero Values and Ignored Fields
If a JSON field is missing, the corresponding Go field remains its zero value (`0`, `""`, `false`, `nil`). If the JSON contains fields that aren't in your Go struct, they are simply skipped.

## Try It

1. Modify the `weatherJSON` string to include a new field and see it get ignored.
2. Change the value of `humidity` to a string in the JSON and observe the error returned by `Unmarshal`.
3. Try unmarshalling into a struct that is missing the `Forecast` field entirely.

## In Production
Always validate the error returned by `Unmarshal`. In production systems, you should never assume the incoming JSON matches your expectations perfectly. Use "Safe Types" and check for `nil` pointers when dealing with nested optional objects.

## Thinking Questions
1. Why must you pass a pointer (`&target`) to `json.Unmarshal`?
2. What happens to Go fields that don't have a matching key in the incoming JSON?
3. When should you use a `map[string]any` instead of a typed `struct`?

> **Forward Reference:** For small JSON payloads, `Marshal` and `Unmarshal` are perfect. But what if you are reading a 1GB JSON file or a continuous stream of JSON objects from a network? In [Lesson 3: Encoder](../3-encoder/README.md), you will learn how to use the streaming `Encoder` and `Decoder` for high-performance I/O.

## Next Step

Continue to `EN.3` encoder.
