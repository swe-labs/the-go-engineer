package main

import (
	"math"
	"testing"
)

func TestPointDistanceToOrigin(t *testing.T) {
	p := Point{X: 3, Y: 4}
	d := p.DistanceToOrigin()
	if math.Abs(d-5.0) > 1e-9 {
		t.Errorf("expected 5.0, got %.2f", d)
	}
}

func TestPointDistanceToOriginZero(t *testing.T) {
	p := Point{X: 0, Y: 0}
	d := p.DistanceToOrigin()
	if d != 0 {
		t.Errorf("expected 0, got %.2f", d)
	}
}

func TestPointTranslated(t *testing.T) {
	p := Point{X: 1, Y: 2}
	p2 := p.Translated(3, 4)
	if p2.X != 4 || p2.Y != 6 {
		t.Errorf("expected (4,6), got (%.1f,%.1f)", p2.X, p2.Y)
	}
}

func TestPointOriginalUnchanged(t *testing.T) {
	p := Point{X: 1, Y: 2}
	p.Translated(10, 20)
	if p.X != 1 || p.Y != 2 {
		t.Error("original should be unchanged after Translated")
	}
}

func TestPointPointerDereference(t *testing.T) {
	p := &Point{X: 3, Y: 4}
	d := p.DistanceToOrigin()
	if math.Abs(d-5.0) > 1e-9 {
		t.Errorf("expected 5.0, got %.2f", d)
	}
}

func TestTemperatureToFahrenheit(t *testing.T) {
	tests := []struct {
		c        float64
		expected float64
	}{
		{0, 32},
		{100, 212},
		{-40, -40},
	}
	for _, tc := range tests {
		temp := Temperature{Celsius: tc.c}
		got := temp.ToFahrenheit()
		if got != tc.expected {
			t.Errorf("ToFahrenheit(%.1f) = %.1f, want %.1f", tc.c, got, tc.expected)
		}
	}
}

func TestTemperatureToKelvin(t *testing.T) {
	tests := []struct {
		c        float64
		expected float64
	}{
		{0, 273.15},
		{100, 373.15},
		{-273.15, 0},
	}
	for _, tc := range tests {
		temp := Temperature{Celsius: tc.c}
		got := temp.ToKelvin()
		if got != tc.expected {
			t.Errorf("ToKelvin(%.2f) = %.2f, want %.2f", tc.c, got, tc.expected)
		}
	}
}

func TestTemperatureIsFreezing(t *testing.T) {
	zero := Temperature{Celsius: 0}
	if !zero.IsFreezing() {
		t.Error("0°C should be freezing")
	}
	neg := Temperature{Celsius: -10}
	if !neg.IsFreezing() {
		t.Error("-10°C should be freezing")
	}
	pos := Temperature{Celsius: 1}
	if pos.IsFreezing() {
		t.Error("1°C should not be freezing")
	}
}
