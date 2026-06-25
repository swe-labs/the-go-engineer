package main

import "testing"

func TestOffer_TotalFirstYear(t *testing.T) {
	tests := []struct {
		name     string
		offer    Offer
		expected float64
	}{
		{
			name: "simple offer",
			offer: Offer{
				Company: "Test",
				Comp: Compensation{
					BaseSalary:     100000,
					EquityPerYear:  20000,
					SigningBonus:   10000,
					AnnualBonusPct: 10,
				},
			},
			expected: 100000 + 20000 + 10000 + 10000,
		},
		{
			name: "no equity or bonus",
			offer: Offer{
				Company: "Minimal",
				Comp: Compensation{
					BaseSalary: 120000,
				},
			},
			expected: 120000,
		},
		{
			name: "with remote and perks",
			offer: Offer{
				Company: "Perky",
				Comp: Compensation{
					BaseSalary:     90000,
					EquityPerYear:  10000,
					SigningBonus:   5000,
					AnnualBonusPct: 5,
					RemoteBudget:   3000,
					OtherPerks:     2000,
				},
			},
			expected: 90000 + 10000 + 5000 + 4500 + 3000 + 2000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.offer.TotalFirstYear()
			if got != tt.expected {
				t.Errorf("TotalFirstYear() = %.0f, want %.0f", got, tt.expected)
			}
		})
	}
}

func TestOffer_TotalRecurring(t *testing.T) {
	offer := Offer{
		Company: "RecurringCo",
		Comp: Compensation{
			BaseSalary:     150000,
			EquityPerYear:  30000,
			AnnualBonusPct: 10,
			RemoteBudget:   2000,
		},
	}
	expected := 150000.0 + 30000.0 + 15000.0 + 2000.0
	got := offer.TotalRecurring()
	if got != expected {
		t.Errorf("TotalRecurring() = %.0f, want %.0f", got, expected)
	}
}

func TestCompareOffers(t *testing.T) {
	offers := []Offer{
		{
			Company: "A",
			Comp:    Compensation{BaseSalary: 100000},
		},
		{
			Company: "B",
			Comp:    Compensation{BaseSalary: 200000, SigningBonus: 50000, EquityPerYear: 25000, AnnualBonusPct: 10},
		},
	}
	results := CompareOffers(offers)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[1].Offer.Company != "B" {
		t.Errorf("expected company B, got %s", results[1].Offer.Company)
	}
	if results[1].FirstYear <= results[0].FirstYear {
		t.Error("B should have higher first year comp")
	}
}
