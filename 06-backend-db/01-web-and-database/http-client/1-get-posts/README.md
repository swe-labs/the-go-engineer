# HCL.1 Basic GET Requests

## Mission

Learn the basics of consuming external APIs in Go using the `net/http` package, while understanding the critical pitfalls of default settings in production environments.

## Prerequisites

- `HS.10` rest-api-exercise

## Mental Model

Think of making an HTTP request as **Sending a Courier to Another Office**.

1. **The Request (`http.Get`)**: You give the courier an address (URL) and tell them to go.
2. **The Handshake (TCP)**: The courier arrives at the other office's front door and establishes a connection.
3. **The Package (`resp.Body`)**: The other office hands the courier a package. The courier doesn't know what's inside; they just bring it back to you.
4. **The Unboxing (`json.NewDecoder`)**: You open the package and sort the items into your own filing system (your structs).
5. **The Safety Lock (`defer resp.Body.Close`)**: Once you've finished unboxing, you must signal to the courier that they can leave. If you don't, they will wait at your desk forever, and eventually, your office will be full of waiting couriers.

## Visual Model

```mermaid
graph LR
    A["Your Code"] -- "http.Get(URL)" --> B["Internet"]
    B -- "io.Reader (Body)" --> C["Response Decoder"]
    C --> D["Your Structs"]
    D -- "Body.Close()" --> E["Finish"]
```

## Machine View

When you call `http.Get`, Go uses a global variable called `http.DefaultClient`. This client is configured to never time out. In a machine environment, this is extremely dangerous. If the remote server is slow to send bytes, your goroutine will block indefinitely. Over time, this "leaks" goroutines and file descriptors until your server runs out of resources and crashes. Additionally, the response body must be manually closed to return the TCP connection to the "Keep-Alive" pool. If you forget to close it, you will eventually hit the "Too many open files" error.

## Run Instructions

```bash
go run ./06-backend-db/01-web-and-database/http-client/1-get-posts
```

The example uses the `dummyjson.com` API to fetch a list of products and map them to a local `RemoteDevice` model.

## Code Walkthrough

### `http.Get(url)`
The simplest way to fetch data. It returns an `*http.Response` and an `error`. Note that the error only indicates a network failure (DNS error, connection refused). It does **not** catch HTTP errors like `404` or `500`.

### `defer resp.Body.Close()`
This is the most important line in any HTTP client code. It ensures the network connection is cleaned up regardless of whether your parsing logic succeeds or fails.

### `json.NewDecoder(resp.Body)`
Since `resp.Body` is an `io.ReadCloser`, we can stream the JSON data directly into our structs. This avoids reading the entire response into a large byte buffer first, saving memory.

### Mapping Models
It is rare that an external API uses the exact same struct field names or structures as your internal domain. We usually define a "Transport Struct" (like `DummyResponse`) and then map those values into our "Domain Model" (like `RemoteDevice`).

## Try It

1. Change the `limit` parameter to `10` and see more results.
2. Intentionally misspell the URL (e.g., `https://dummyjson.xxx`) and observe the error message.
3. Try fetching a URL that returns a `404` (e.g., `https://dummyjson.com/nonexistent`) and check why the `resp.StatusCode` check is necessary.

## In Production
**NEVER use `http.Get` in a production service.** You should always instantiate your own `http.Client` with a reasonable `Timeout` (e.g., 5-10 seconds). This ensures that a single slow dependency cannot bring down your entire application.

```go
client := &http.Client{
    Timeout: 10 * time.Second,
}
resp, err := client.Get(url)
```

## Thinking Questions
1. Why does `http.Get` not return an error for a `404 Not Found` response?
2. What happens if you forget to call `resp.Body.Close()`?
3. How would you handle an API that returns a large list of 100,000 items?

> **Forward Reference:** You can now fetch data from the web. But how do you write tests for this? You don't want your unit tests to actually hit the internet! In [Lesson 2: Refactor for Testability](../2-refactor-for-testability/README.md), you will learn how to use interfaces and `httptest` to mock external APIs.

## Next Step

Continue to `HCL.2` refactor-for-testability.
