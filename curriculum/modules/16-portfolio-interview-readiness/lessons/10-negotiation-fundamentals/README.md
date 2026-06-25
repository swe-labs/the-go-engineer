# Negotiation fundamentals

## Learning objective

Evaluate job offers by calculating total compensation across salary, equity, bonuses, and benefits, and apply negotiation strategies including BATNA preparation, timing, and value-anchored counter-offers.

## Why this matters

Most engineers focus exclusively on base salary during negotiations, leaving 20-40% of total compensation on the table. Equity, signing bonuses, performance bonuses, remote budgets, education stipends, and other perks can dramatically change the value of an offer. Companies expect candidates to negotiate. A well-prepared negotiation signals seniority and confidence. Engineers who negotiate earn 15-30% more over their careers than those who accept the first offer.

## Mental model

Think of an offer as a multi-dimensional value vector, not a single number. Each component (base, equity, bonus, perks) has different liquidity, tax treatment, and growth potential. Your goal is to maximize the total expected value of the vector, not any single component. Negotiation is collaborative problem-solving: you and the recruiter each want a deal that works. Your leverage comes from having alternatives (BATNA), understanding market rates, and being willing to walk away.

## Core idea

Total compensation is the sum of all monetary and non-monetary benefits an employer provides. The standard components:

| Component | Description | Typical range |
|---|---|---|
| Base salary | Fixed cash, paid bi-weekly or monthly | Market rate ± 10% |
| Equity / RSUs | Company shares, vested over 3-4 years | 10-50% of base |
| Signing bonus | One-time cash, first paycheck | $10K-$100K |
| Annual bonus | Performance-based, % of base | 5-20% of base |
| Remote budget | Home office, internet, co-working | $1K-$5K/year |
| 401k match | Employer retirement contribution | 3-6% of salary |
| Education stipend | Conferences, courses, books | $1K-$10K/year |
| Other perks | Insurance, gym, food, transit | Variable |

BATNA (Best Alternative To a Negotiated Agreement) is your fallback if you don't reach a deal. A strong BATNA — another offer, a current job, or freelance income — dramatically increases your negotiating power.

## Under the hood

Companies have compensation bands determined by level, location, and market data. Recruiters have discretion to move within the band and often have access to a "discretionary pool" for candidates they want to close. The key insight: most negotiable value lives in this discretion pool, not the band itself.

Equity compensation is complex: RSUs (Restricted Stock Units) are taxed as income at vesting, then capital gains at sale. ISO (Incentive Stock Options) have different tax treatment. Early-stage equity may be worthless if the company fails. Public company RSUs are effectively cash-equivalent at the current stock price, discounted by vesting schedule.

## How Go uses it

While negotiation itself doesn't involve Go code, you can build tools in Go to model and compare offers quantitatively. Go's standard library provides everything needed:

- `fmt.Sprintf` for formatting currency amounts.
- `sort` for ordering offers by total value.
- `encoding/json` for reading offer data from configuration files.
- `math` for present value calculations with discount rates.
- `testing` with table-driven tests for verifying compensation logic.

Go's emphasis on correctness and clarity makes it an excellent language for financial modeling where off-by-one errors or rounding mistakes have real dollar consequences.

## Go example

```go
package main

import (
	"fmt"
	"strings"
)

type Compensation struct {
	BaseSalary     float64
	EquityPerYear  float64
	SigningBonus   float64
	AnnualBonusPct float64
	RemoteBudget   float64
	OtherPerks     float64
}

type Offer struct {
	Company  string
	Title    string
	Location string
	Comp     Compensation
}

func (o Offer) TotalFirstYear() float64 {
	c := o.Comp
	return c.BaseSalary + c.EquityPerYear + c.SigningBonus +
		(c.BaseSalary * c.AnnualBonusPct / 100) + c.RemoteBudget + c.OtherPerks
}

type ComparisonResult struct {
	Offer      Offer
	FirstYear  float64
	Recurring  float64
	MonthlyNet float64
}

func CompareOffers(offers []Offer) []ComparisonResult {
	var results []ComparisonResult
	for _, o := range offers {
		results = append(results, ComparisonResult{
			Offer:      o,
			FirstYear:  o.TotalFirstYear(),
			Recurring:  o.TotalRecurring(),
			MonthlyNet: o.TotalFirstYear() / 12,
		})
	}
	return results
}

func main() {
	offers := []Offer{
		{Company: "Stripe", Title: "Senior Backend Engineer", Location: "Remote US",
			Comp: Compensation{BaseSalary: 185000, EquityPerYear: 45000, SigningBonus: 30000, AnnualBonusPct: 10, RemoteBudget: 5000}},
		{Company: "GitLab", Title: "Senior Go Engineer", Location: "Remote Global",
			Comp: Compensation{BaseSalary: 165000, EquityPerYear: 35000, SigningBonus: 15000, AnnualBonusPct: 15, RemoteBudget: 5000, OtherPerks: 3000}},
		{Company: "Shopify", Title: "Backend Developer", Location: "Remote Canada",
			Comp: Compensation{BaseSalary: 175000, EquityPerYear: 25000, SigningBonus: 10000, AnnualBonusPct: 5, RemoteBudget: 2000}},
	}
	results := CompareOffers(offers)
	fmt.Printf("%-20s %15s %15s %15s\n", "Company", "First Year", "Recurring", "Monthly")
	fmt.Println(strings.Repeat("-", 65))
	for _, r := range results {
		fmt.Printf("%-20s $%12.0f $%12.0f $%12.0f\n", r.Offer.Company, r.FirstYear, r.Recurring, r.MonthlyNet)
	}
}
```

Run with `go run .` to see a side-by-side comparison of three real-world Go engineer offers.

## Step-by-step execution

When `CompareOffers` runs:

1. Each `Offer` struct is populated with compensation fields — base, equity, signing bonus, bonus percentage, remote budget, other perks.
2. `TotalFirstYear()` sums base salary, equity per year, signing bonus, annual bonus (base × pct/100), remote budget, and other perks.
3. `CompareOffers` iterates the offers slice, calling `TotalFirstYear()` and `TotalRecurring()` for each.
4. Monthly net is first-year total divided by 12 (a rough take-home proxy, not accounting for taxes).
5. Results are printed in aligned columns with `$` prefixes and zero decimal formatting.

In the example: Stripe first-year = 185000 + 45000 + 30000 + 18500 + 5000 = $283,500. GitLab: 165000 + 35000 + 15000 + 24750 + 5000 + 3000 = $247,750. The numbers reveal that Stripe's higher base and signing bonus make the first year significantly more valuable, but recurring comp (without signing bonus) may shift the comparison.

## Common mistakes

- **Focusing only on base salary.** Base is important but equity and bonuses can add 30-50%. Always calculate total first-year and recurring comp.
- **Not understanding equity.** RSUs at a public company are near-cash. Options at a startup may be worthless. Ask about liquidity, valuation, and dilution. Treat early-stage equity as a lottery ticket, not income.
- **Accepting the first number.** Recruiters expect you to counter. A simple "I was hoping for something closer to $X based on my research and experience" works. The worst they can say is no.
- **Failing to negotiate non-salary items.** If the base is at the ceiling, negotiate signing bonus, remote budget, or education stipend. These are often easier to move.
- **Not preparing a BATNA.** Without a strong alternative, you have no leverage. Maintain your current job search pipeline until you sign. Get multiple offers if possible.
- **Talking first.** Whoever states a number first loses information. Let the recruiter share the range before you anchor.

## Debugging walkthrough

An engineer runs the offer comparison tool and gets unexpected results:

```go
offer := Offer{
    Company: "Co",
    Comp: Compensation{
        BaseSalary:     100000,
        EquityPerYear:  0,
        SigningBonus:   0,
        AnnualBonusPct: 10,
    },
}
fmt.Println(offer.TotalFirstYear()) // Output: 110000
```

**Symptom**: Expected $100,000 (base) + $0 (equity) + $0 (signing) + $10,000 (10% bonus) = $110,000. That seems correct. But why is `TotalRecurring()` returning $110,000 instead of the expected $100,000?

**Investigation**: The engineer realizes `TotalRecurring()` includes annual bonus but not signing bonus. They check the method:

```go
func (o Offer) TotalRecurring() float64 {
    c := o.Comp
    return c.BaseSalary + c.EquityPerYear + (c.BaseSalary * c.AnnualBonusPct / 100)
}
```

**Root cause**: The engineer had personally added `RemoteBudget + OtherPerks` in their mental model but the method didn't include these until they updated the code. Also, bonus is typically not guaranteed year-over-year, so some engineers exclude it from recurring.

**Fix**: Add a comment clarifying what "recurring" includes: base + equity + bonus + remote budget + perks. Optionally provide two recurring numbers: with and without bonus.

## Production notes

- **Have a walk-away number.** Before any negotiation, decide the minimum acceptable total comp. If the offer doesn't meet it, walk away without regret. This prevents emotional decisions.
- **Get everything in writing.** Verbal commitments for equity, bonuses, or promotions are not binding. Request an amended written offer before resigning from your current role.
- **Consider the whole package.** Benefits like unlimited PTO (which often means less PTO taken), insurance quality, parental leave, and sabbatical policies have real monetary value. Add them to your model.
- **Timing matters.** The best time to negotiate is after you have a written offer but before you accept. The second-best time is after a competing offer arrives. Avoid negotiating during initial screening.
- **Use level as leverage.** If base is capped, negotiate for a higher title or level. A senior title now means future offers start from a higher base. Level changes can be worth more than a one-time salary bump.

## Performance implications

- **Equity valuation is the largest uncertainty.** For public companies, multiply RSU count by current stock price, subtract expected tax, and discount for vesting (typically 25% per year). For private companies, reduce the value by 50-80% for risk and illiquidity.
- **Tax impact changes net value.** A $50,000 signing bonus is taxed at your marginal rate (often 35-45%), reducing net to $27,500-$32,500. Factor this into comparisons.
- **Inflation erodes future value.** A dollar today is worth more than a dollar next year. Use a 3-5% discount rate when comparing offers with different timing of payments.
- **Cost of living adjustments.** A $150,000 salary in San Francisco may provide less purchasing power than $120,000 in Austin. Use cost-of-living calculators to normalize.

## Practice task

Extend the offer comparison tool with:

1. A `NetAfterTax(rate float64) float64` method on `Compensation` that returns total first-year comp minus estimated taxes at the given rate.
2. A `PresentValue(rate float64, years int) float64` method that discounts multi-year equity vesting to present value using the formula `FV / (1 + r)^n` where r is the discount rate and n is the year.
3. A `BestOffer` function that takes `[]ComparisonResult` and a preference string ("first_year", "recurring", or "balanced") and returns the index of the best offer based on weighted scoring.
4. `main()` that loads offers from a JSON file (inline or hardcoded), compares them, and prints a ranked list.

## Tests / verification

```bash
go run ./curriculum/modules/16-portfolio-interview-readiness/lessons/10-negotiation-fundamentals
go test ./curriculum/modules/16-portfolio-interview-readiness/lessons/10-negotiation-fundamentals
```

The existing tests verify `TotalFirstYear` (simple, no equity/bonus, with perks), `TotalRecurring`, and `CompareOffers` (length, sorting, ordering). After completing the practice task, add tests for `NetAfterTax` with various rates and `PresentValue` with different discount rates and vesting years.

## Review questions

1. List six components of total compensation beyond base salary. Which two typically have the largest financial impact?
2. What is BATNA and how does it affect your negotiation leverage?
3. Why might "first-year total" and "recurring total" give different rankings when comparing two offers?
4. What is the primary risk difference between RSUs at a public company and options at a startup?
5. Name three things you can negotiate besides base salary when the salary band is maxed out.

## NEXT UP

Open source contribution — learn how to find Go projects, make your first contribution, and build your professional reputation through open source work.
