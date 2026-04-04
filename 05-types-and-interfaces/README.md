# Section 5: Types & Interfaces

## Beginner → Expert Mapping

| Topic | Level | Importance | Engineering Concept |
|-------|-------|------------|---------------------|
| Structs | Beginner | High | Data encapsulation without classes |
| Methods | Intermediate | High | Value vs Pointer Receivers |
| Interfaces | Intermediate | **Critical** | Duck typing, dependency inversion |
| Generics | Advanced | Medium | Type parameters, constraints |

## Engineering Depth

Go is not a traditional Object-Oriented language; it lacks inheritance (`extends`). Instead, it uses strictly composition and **implicit interfaces** (duck typing).

- **Memory mechanics:** When defining methods, using a pointer receiver `func (u *User) Login()` passes the memory address (8 bytes). Using a value receiver `func (u User) Login()` creates a copy of the entire struct on the stack. Large structs should always use pointer receivers.
- **Interfaces:** Interfaces in Go are stored as a two-word internal structure combining a pointer to the type information and a pointer to the underlying data.

## References

1. **[Effective Go]** [Interfaces and other types](https://go.dev/doc/effective_go#interfaces_and_types)
2. **[Go Blog]** [An Introduction to Generics](https://go.dev/blog/intro-generics)

---

## 🏗 Exercise: Payroll Processor (`6-payroll-processor`)

This capstone requires leveraging Interfaces to build polymorphic systems without inheritance.

### Step-by-Step Instructions & Hints

1. **Define the interface:** Create an `Employee` interface with a `CalculatePay() float64` method.
2. **Build implementations:** Create `FullTime` (Base Salary) and `Contractor` (Hourly * Hours) structs.
3. **Implement methods:** Attach the `CalculatePay()` method to both structs using value receivers (since the calculation shouldn't modify the struct state).
4. **Process Payroll:** Create a `ProcessPayroll(employees []Employee)` function that iterates the slice, calls the method, and prints a total.
   - *Hint:* Because both structs implement the interface method, the slice can hold mixed types perfectly!


## Learning Path

| ID | Lesson | Concept | Requires |
| --- | --- | --- | --- |
| TI.1 | [structs](./1-struct) | Field layout · constructor pattern · value vs pointer copy | 🟢 entry |
| TI.2 | [methods](./2-methods) | Value receiver vs pointer receiver · method sets | TI.1 |
| TI.3 | [interfaces](./3-interfaces) | Implicit satisfaction · polymorphism · type assertions | TI.1, TI.2 |
| TI.4 | [Stringer](./4-stringer) | fmt.Stringer · named type on int · custom Weekday | TI.2, TI.3 |
| TI.5 | [generics](./5-generics) | [T constraint] · comparable · union constraints · Map/Filter | TI.3, TI.4 |
