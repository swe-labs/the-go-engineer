# TM.2 Time Formatting: The Reference Date

## Mission

Master Go's unique approach to time formatting and parsing. Learn the "1-2-3-4-5-6" reference date trick and understand how to convert human-readable strings into precise `time.Time` structs for your backend logic.

## Prerequisites

- `TM.1` time-basics

## Mental Model

Think of Go's Time Formatting as **A Mirror**.

1. **The Reference (`Jan 2 15:04:05 2006`)**: Go uses a specific moment in time as the template.
2. **The Layout (`The Reflection`)**: You don't write `YYYY-MM-DD`. Instead, you write what the Reference Date would look like in your desired format (e.g., `2006-01-02`).
3. **The Mirror**: The `Format` function takes the template, identifies where the month, day, and year are, and reflects the **actual** time using that same layout.

## Visual Model

| Component | Value in Reference Date |
| :--- | :--- |
| **Month** | 01 (January) |
| **Day** | 02 |
| **Hour** | 03 (PM) or 15 (24h) |
| **Minute** | 04 |
| **Second** | 05 |
| **Year** | 06 (2006) |
| **Zone** | 07 (-07:00) |

**Mnemonic**: 1 (Month), 2 (Day), 3 (Hour), 4 (Minute), 5 (Second), 6 (Year), 7 (Zone).

## Machine View

- **Pattern Matching**: When you call `Format`, Go doesn't use regex. It scans your layout string for these specific "Magic Numbers" (like 2006 or 15).
- **Efficiency**: Because it's looking for fixed numbers rather than symbols like `%Y`, the parsing engine is extremely fast and requires no lookup tables.
- **Parsing**: `time.Parse` does the reverse. It uses the layout to extract bits from a string and reconstructs a `time.Time` struct.

## Run Instructions

```bash
go run ./07-concurrency/01-concurrency/time-and-scheduling/2-formatting
```

## Code Walkthrough

### `now.Format(layout)`
To get a string, you call `Format` on a time object.
- Example: `now.Format("Jan 2, 2006")` -> `Apr 29, 2026`.

### `time.Parse(layout, value)`
To get a time object from a string. Note that `time.Parse` returns an error. If the string doesn't match the layout exactly, it will fail.
- Tip: Use `time.RFC3339` for API and Database interactions.

### Built-in Constants
Go provides constants for common standards so you don't have to remember the reference date:
- `time.RFC3339`: `2006-01-02T15:04:05Z07:00` (The Gold Standard).
- `time.Kitchen`: `3:04PM`.

## Try It

1. Create a layout for a "European" date format: `Day/Month/Year`.
2. Parse the string `"2024-02-29"` (Leap year). What happens if you try to parse `"2023-02-29"`?
3. Format the current time to only show the day of the week and the time in 24-hour format.

## Verification Surface

Observe the various formatting outputs:

```text
Current time (default Go format): 2026-04-29 11:30:40...
Formatted as YYYY-MM-DD: 2026-04-29
Formatted as MM/DD/YYYY hh:mm:ss PM: 04/29/2026 11:30:40 AM
Formatted as Day, Month Date, Year: Wed, Apr 29, 2026

Parsed '2025-07-15' -> Year: 2025, Month: July, Day: 15
Parsed RFC3339 (in UTC): 2025-12-25 09:00:00 +0000 UTC
```

## In Production
**Always use `time.ParseInLocation` when parsing local times.**
`time.Parse` assumes UTC unless the string contains timezone info. If you are parsing a user's input from a specific city (like "2024-05-01 10:00"), you must specify the `time.Location` to avoid being off by several hours.

## Thinking Questions
1. Why did Go choose the specific date `January 2, 2006`? (Hint: 1-2-3-4-5-6).
2. How do you format a date with a suffix like "1st", "2nd", or "3rd"? (Hint: Go doesn't support this natively!).
3. What is the difference between `03` and `15` in a layout string?

## Next Step

Next: `TM.3` -> `07-concurrency/01-concurrency/time-and-scheduling/3-timer-and-ticker`

Open `07-concurrency/01-concurrency/time-and-scheduling/3-timer-and-ticker/README.md` to continue.
