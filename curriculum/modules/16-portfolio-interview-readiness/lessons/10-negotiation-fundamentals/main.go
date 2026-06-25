package main

import (
	"fmt"
	"strings"
)

type Compensation struct {
	BaseSalary     float64
	EquityPerYear  float64
	SigningBonus   float64
	AnnualBonusPct float64 // percentage of base
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
	return c.BaseSalary + c.EquityPerYear + c.SigningBonus + (c.BaseSalary * c.AnnualBonusPct / 100) + c.RemoteBudget + c.OtherPerks
}

func (o Offer) TotalRecurring() float64 {
	c := o.Comp
	return c.BaseSalary + c.EquityPerYear + (c.BaseSalary * c.AnnualBonusPct / 100) + c.RemoteBudget + c.OtherPerks
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
		fy := o.TotalFirstYear()
		rec := o.TotalRecurring()
		monthly := fy / 12
		results = append(results, ComparisonResult{
			Offer:      o,
			FirstYear:  fy,
			Recurring:  rec,
			MonthlyNet: monthly,
		})
	}
	return results
}

func main() {
	offers := []Offer{
		{
			Company:  "Stripe",
			Title:    "Senior Backend Engineer",
			Location: "Remote US",
			Comp: Compensation{
				BaseSalary:     185000,
				EquityPerYear:  45000,
				SigningBonus:   30000,
				AnnualBonusPct: 10,
				RemoteBudget:   5000,
			},
		},
		{
			Company:  "GitLab",
			Title:    "Senior Go Engineer",
			Location: "Remote Global",
			Comp: Compensation{
				BaseSalary:     165000,
				EquityPerYear:  35000,
				SigningBonus:   15000,
				AnnualBonusPct: 15,
				RemoteBudget:   5000,
				OtherPerks:     3000,
			},
		},
		{
			Company:  "Shopify",
			Title:    "Backend Developer",
			Location: "Remote Canada",
			Comp: Compensation{
				BaseSalary:     175000,
				EquityPerYear:  25000,
				SigningBonus:   10000,
				AnnualBonusPct: 5,
				RemoteBudget:   2000,
			},
		},
	}

	results := CompareOffers(offers)
	fmt.Printf("%-20s %-30s %15s %15s %15s\n", "Company", "Role", "First Year", "Recurring", "Monthly")
	fmt.Println("------" + strings.Repeat("-", 75))
	for _, r := range results {
		fmt.Printf("%-20s %-30s $%12.0f $%12.0f $%12.0f\n",
			r.Offer.Company,
			r.Offer.Title,
			r.FirstYear,
			r.Recurring,
			r.MonthlyNet,
		)
	}
}
