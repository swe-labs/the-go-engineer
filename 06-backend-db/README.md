# 06 Backend, APIs & Databases

## Mission

This stage teaches HTTP servers, API design, and databases as one application-shaped backend layer. It moves from handling standard library requests to building database-backed services.

By the end of this stage, a learner should be able to:

- build and route HTTP servers using the standard library
- implement HTTP middleware for logging and validation
- design and consume REST and gRPC APIs
- connect to and query SQLite databases
- implement the repository pattern to isolate SQL from business logic

## Stage Map

| Track | Surface | Core Job |
| --- | --- | --- |
| `HS.1-10` | HTTP Servers | Teach request/response boundaries, routing, middleware, and the REST API exercise. |
| `API.1-9` | API Design | Compare REST and gRPC transport choices and implement an API service. |
| `DB.1-8` | Databases | Teach connecting, querying, transactions, the repository pattern, and timeout budgeting. |

## Why This Stage Exists Now

The learner already knows:

- project layout and packaging
- reading and writing configuration
- working with JSON encoding

That is enough to start asking engineering questions like:

- how does Go serve web requests concurrently?
- how do we persist data across service restarts?
- how do multiple services talk to each other safely over the network?

## Suggested Learning Flow

1. Work through the HTTP server track to learn request boundaries and middleware.
2. Use the API lessons to compare REST and gRPC transport choices.
3. Finish the section by extending the database path with queries, repositories, and context timeouts.

## Next Step

After this section, continue to [07 Concurrency](../07-concurrency).
