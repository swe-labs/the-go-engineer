# Section 2: Control Flow

## Beginner → Expert Mapping

| Topic | Level | Importance | Engineering Concept |
|-------|-------|------------|---------------------|
| Conditionals (If) | Beginner | High | Early returns vs Else blocks |
| Loops (For) | Beginner | High | The only looping construct in Go |
| Switch | Intermediate | High | Type switches, empty conditionals |

## Engineering Depth
Go purposefully omitted `while` and `do-while` loops to keep the compiler simple and the language predictable. `for` handles everything.
A critical engineering rule in Go is: **Line of Sight**. You should avoid `else` blocks whenever possible by returning early (guard clauses). This keeps the main "happy path" aligned to the far left margin, which makes code significantly easier to read.

## References
1. **[Effective Go]** [Control structures](https://go.dev/doc/effective_go#control-structures)
2. **[Article]** [Return Early Pattern](https://medium.com/@matryer/line-of-sight-in-code-186dd7cdea88)

---

## 🏗 Exercise: Pricing Calculator (`4-pricing-calculator`)

This project will force you to use switch statements and map lookups effectively.

### Step-by-Step Instructions & Hints
1. **Setup the store:** Create a map mapping item codes (like `TSHIRT`) to float64 prices.
2. **Calculate price func:** Write a function `calculateItemPrice(code string) (float64, bool)`.
   - *Hint:* Return `bool` to indicate if the item was found. This avoids relying on the `0.0` zero value which could be a legitimate free price.
3. **Handle Sales:** If an item ends with `_SALE`, strip the suffix using `strings.TrimSuffix()`, find the base price, and apply a 20% discount.
4. **Use Switch:** Use a switch statement on the boolean `found` flag to print success or error gracefully.


## Learning Path

| ID | Lesson | Concept | Requires |
| --- | --- | --- | --- |
| CF.1 | [for loop](./1-for-loop) | C-style · while-style · infinite · range | 🟢 entry |
| CF.2 | [if / else](./2-if-else) | Guard clauses · if-with-init · comma-ok idiom | CF.1 |
| CF.3 | [switch](./3-switch) | No fall-through · type switch · tagless switch | CF.1, CF.2 |
