# Section 6: Composition & Embedding

## Beginner → Expert Mapping

| Topic | Level | Importance | Engineering Concept |
|-------|-------|------------|---------------------|
| Struct Embedding | Intermediate | High | "Has-a" vs "Is-a" relationship modeling |
| Interface Composition | Intermediate | High | Building larger behavioral contracts |

## Engineering Depth
Go relies entirely on Composition rather than classical Inheritance (No `class extends BaseClass`). Struct embedding (anonymous fields) automatically promotes inner methods and fields to the outer struct, mimicking inheritance. However, the outer struct DOES NOT become the inner struct in type—meaning you cannot natively cast it as you would down a class hierarchy. Time & Space complexity for embedding is $O(1)$ at structural instantiation since memory is allocated in a single contiguous flat block containing both structs.

## References
1. **[Effective Go]** [Embedding](https://go.dev/doc/effective_go#embedding)

---

## 🏗 Exercise: Bank Account System (`3-bank-account`)

### Step-by-Step Instructions & Hints
1. **Define the inner struct:** Create a base `Account` struct holding the core balance logic (`Deposit()`, `Withdraw()`).
2. **Setup the outer struct:** Create a `PremiumAccount` struct, and embed `Account` by passing it as an anonymous field (do not give it a field name!).
3. **Override vs Promote:** 
   - *Hint:* Add an `InterestRate` to `PremiumAccount` and explicitly call `Deposit()`—notice that you don't need to specify `.Account.Deposit()`, it is automatically promoted to the top level!
