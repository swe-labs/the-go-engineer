# DB.1 Connecting to SQLite

## Mission

Learn how to use the Go standard library to connect to a SQLite database, understand the connection pool mechanics, and initialize your first table.

## Prerequisites

- `API.9` grpc-service-exercise

## Mental Model

Think of `database/sql` as a **Universal Remote Control**.

1. **The Remote (`database/sql`)**: A single device that has buttons for "Power", "Volume", and "Channel".
2. **The Code (Drivers)**: Each TV (PostgreSQL, MySQL, SQLite) needs a specific code to work with the remote. You "import" this code using a blank import (`_`).
3. **The Setup (`sql.Open`)**: You tell the remote which TV you have. The remote doesn't turn the TV on yet; it just remembers the settings.
4. **The Test (`db.Ping`)**: You press the Power button to make sure the TV actually responds.

## Visual Model

```mermaid
graph LR
    A["Go Code"] --> B["database/sql (Interface)"]
    B --> C["SQLite Driver (Plugin)"]
    C --> D["example.db (File)"]
```

## Machine View

Go's `database/sql` package manages a **Connection Pool** automatically. When you call `sql.Open`, you aren't opening a single socket; you are creating a manager that can open many sockets as needed.
- **Lazy Initialization**: Connections are only created when a query is actually executed.
- **Blank Imports**: Drivers register themselves in their `init()` functions using `sql.Register`. This is why we use `import _ "driver"`.
- **SQLite Special Case**: Unlike Postgres or MySQL, SQLite is "In-Process". The "database" is just a file on your disk, and the driver is a library that reads and writes to that file directly.

## Run Instructions

```bash
go run ./06-backend-db/01-web-and-database/databases/1-connecting-to-db
```

This will create a file named `example.db` in your current directory.

## Code Walkthrough

### Blank Import `_ "modernc.org/sqlite"`
This line is mandatory. It registers the SQLite driver so that `sql.Open("sqlite", ...)` knows which code to use. We use `modernc.org/sqlite` because it is a pure-Go implementation that doesn't require CGO, making it much easier to cross-compile.

### `sql.Open(driver, source)`
The first argument is the driver name. The second is the "Data Source Name" (DSN), which for SQLite is just the path to the file.

### `db.Ping()`
Always call `Ping` after `Open`. It forces the driver to actually open the file and verifies that you have the correct permissions.

### `db.Exec(query)`
Used for "One-Way" queries that don't return rows, such as `CREATE TABLE`, `INSERT`, `UPDATE`, or `DELETE`.

## Try It

1. Change the `dbPath` to `:memory:` and observe that no file is created on your disk (this is great for unit tests!).
2. Try opening a path that you don't have permission to write to and see how `db.Ping` fails.
3. Check the version of the `modernc.org/sqlite` package in `go.mod`.

## In Production
For production databases (like PostgreSQL), you should configure the connection pool limits:
- `db.SetMaxOpenConns(25)`
- `db.SetMaxIdleConns(25)`
- `db.SetConnMaxLifetime(5 * time.Minute)`
Without these limits, your app might try to open thousands of connections during a traffic spike, which can crash your database server.

## Thinking Questions
1. Why doesn't `sql.Open` return a connection error immediately?
2. What is the benefit of a "Blank Import" over a regular import?
3. Why is `db.Close()` deferred? What happens if you forget it in a long-running server?

> [!TIP]
> You have a connection. Now let's put some data in it. In [Lesson 2: Executing Queries (INSERT)](../2-query/README.md), you will learn how to safely insert data using prepared statements to prevent SQL injection.

## Next Step

Next: `DB.2` -> [`06-backend-db/01-web-and-database/databases/2-query`](../2-query/README.md)
