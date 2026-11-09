# Questions

Answer each question with concrete reasoning and evidence from your Pricing Checkout implementation.

## Question 1

Explain your `CalculateTotal` function step by step. What order do you apply discounts and taxes? Why?

## Question 2

Your `FormatPrice` function converts cents to a dollar string. Show how you handle negative amounts. What edge cases did you consider?

## Question 3

What would happen if a country code of `"UK"` (instead of `"GB"`) was passed to `TaxRate`? How does your implementation handle unknown country codes? Is returning 0.0 the right choice, or should it return an error?

## Question 4

The `DiscountPercent` function uses a switch on quantity ranges. What are the boundary values (exactly 10, exactly 50, exactly 100)? Verify your implementation is correct at each boundary.

## Question 5

Why did you choose `int` (cents) for prices instead of `float64`? What problems can floating-point pricing cause in a real checkout system?
