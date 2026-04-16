# Track B: Encoding

## Mission

This track teaches you how Go moves structured data into and out of bytes, with special focus on
JSON and stream-oriented decoding.

## Track Map

| ID | Type | Surface | Why It Matters | Requires |
| --- | --- | --- | --- | --- |
| `EN.1` | Lesson | [JSON marshalling](./1-marshalling) | Introduces struct-to-JSON serialization and struct tags. | entry |
| `EN.2` | Lesson | [JSON unmarshalling](./2-unmarshal) | Shows how JSON becomes Go values and where zero values can hide missing fields. | `EN.1` |
| `EN.3` | Lesson | [JSON encoder](./3-encoder) | Teaches stream-oriented writing with `json.NewEncoder`. | `EN.1`, `EN.2` |
| `EN.4` | Lesson | [JSON decoder](./4-decode) | Teaches stream-oriented reading with `json.NewDecoder`. | `EN.2`, `EN.3` |
| `EN.5` | Lesson | [base64](./5-base64_encoding) | Adds binary-to-text transport encoding. | `EN.1` |
| `EN.6` | Exercise | [config parser](./6-config-parser) | Combines file I/O, JSON decoding, and validation in one milestone. | `EN.1`, `EN.2`, `EN.3`, `EN.4` |

## Suggested Order

1. Work through `EN.1` to `EN.5` in order.
2. Complete `EN.6` without copying the finished solution line by line.

## Track Milestone

`EN.6` is the current encoding track milestone.

If you can complete it and explain:

- why decoding directly from an `io.Reader` is usually the right production habit
- why zero values make config validation necessary after parsing
- why JSON and base64 solve different problems even though they both produce text

then the encoding part of Section 09 is doing its job.

## Next Step

After `EN.6`, continue to the [Filesystem track](../filesystem) or back to the
[Section 09 overview](../README.md).
