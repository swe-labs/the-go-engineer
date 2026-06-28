# Answer Key

## Question 1

`CalculateTotal` first sums all item prices, applies a quantity discount (reducing the subtotal), then applies tax (increasing the discounted total). This order matches real-world checkout: discounts are applied before tax, as tax is calculated on the post-discount amount.

```go
func CalculateTotal(items []Item, country string) (int, error) {
    if len(items) == 0 {
        return 0, errors.New("no items")
    }
    sum := 0
    for _, item := range items {
        sum += item.PriceCents
    }
    qty := len(items)
    discount := DiscountPercent(qty)
    if discount > 0 {
        sum -= int(float64(sum) * discount)
    }
    tax := TaxRate(country)
    if tax > 0 {
        sum += int(float64(sum) * tax)
    }
    return sum, nil
}
```

## Question 2

Negative amounts produce `"-$X.XX"` strings:

```go
func FormatPrice(cents int) string {
    if cents < 0 {
        return fmt.Sprintf("-$%d.%02d", -cents/100, -cents%100)
    }
    return fmt.Sprintf("$%d.%02d", cents/100, cents%100)
}
```

Edge cases: zero cents (`$0.00`), single-digit cents (`$0.01`), large values (`$9999.99`), and negative cents (`-$12.34`). The `%02d` format ensures cents always have two digits.

## Question 3

Passing `"UK"` returns `0.0` because the switch only matches `"GB"`. Returning `0.0` for unknown countries is a convention — it silently assumes no tax. A safer alternative would return an error:

```go
func TaxRate(country string) (float64, error) {
    switch country {
    case "US": return 0.07, nil
    case "GB": return 0.20, nil
    // ...
    default: return 0, fmt.Errorf("unknown country: %s", country)
    }
}
```

The current design is simpler but could silently undercharge tax. In production, returning an error is safer.

## Question 4

Boundaries:
- `DiscountPercent(10)` → 0.05 (exactly 10 triggers the `>= 10` case)
- `DiscountPercent(50)` → 0.10 (exactly 50 triggers the `>= 50` case)
- `DiscountPercent(100)` → 0.15 (exactly 100 triggers the `>= 100` case)
- `DiscountPercent(9)` → 0.0 (below 10 falls through to default)

The switch should use `case qty >= 100:` before `case qty >= 50:` before `case qty >= 10:` — order matters because the first matching case wins.

## Question 5

`int` (cents) avoids floating-point rounding errors. `float64` prices can produce `$0.1` instead of `$0.10` due to binary floating-point representation. Real checkout systems always use integer cents or a decimal type to ensure exact monetary calculations. Example: `0.1 + 0.2 == 0.3` is `false` in floating point — unacceptable for payment processing.
