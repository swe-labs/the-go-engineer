# Answer Key

## Question 1

Phone numbers make good map keys because they are (ideally) unique per person and rarely change. Using email as the key would also work, but people change email addresses more frequently than phone numbers. Phone is also shorter and more consistent in format than names (which can have duplicates).

If email were the key, the same logic would apply: `AddContact` would check for duplicate emails instead of duplicate phones. The `Directory` type would be `map[string]Contact` either way — only the key field changes.

## Question 2

Case-insensitive sorting is done using `strings.ToLower` in the `sort.Slice` less function:

```go
sort.Slice(contacts, func(i, j int) bool {
    return strings.ToLower(contacts[i].Name) < strings.ToLower(contacts[j].Name)
})
```

If "alice" and "Alice" both exist, the sort order between them is stable but unspecified — `sort.Slice` is not a stable sort. To guarantee deterministic ordering, add a secondary sort by another field (e.g., email) when names are equal case-insensitively.

## Question 3

Trimming prevents invisible whitespace bugs. A user might copy-paste a name from a spreadsheet and include a trailing newline or leading space. Without trimming, `"Alice "` and `"Alice"` would be treated as different names or phones, causing confusing behaviour like "duplicate phone not detected" or "search misses contact."

## Question 4

Searching only the name field would miss contacts when the user searches by email (e.g., searching "alice@example.com" to find Alice's record). Searching both fields is more convenient for users.

Performance trade-off: for 10,000 contacts, a linear scan through the map is O(n). Searching both fields doubles the string comparison work per contact (two `Contains` calls instead of one). For 10,000 entries this is still fast (~milliseconds), but for millions of records you would want an indexed search (inverted index, full-text search engine).

## Question 5

Rejecting duplicates (returning an error) is safer because it prevents accidental data loss. If the function silently overwrote, a user who calls `AddContact` instead of `UpdateContact` would lose the original contact data without warning.

If the user's intent was to update, they should call `UpdateContact` explicitly. The two operations have different semantics: `AddContact` means "create a new entry" and should fail if the key exists; `UpdateContact` means "modify an existing entry" and should fail if the key doesn't exist.
