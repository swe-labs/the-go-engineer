package main

import "testing"

func TestCelsiusF(t *testing.T) {
	c := Celsius(100.0)
	f := c.F()
	if f != 212.0 {
		t.Errorf("expected 212.0°F, got %.1f", f)
	}
}

func TestCelsiusF_Freezing(t *testing.T) {
	c := Celsius(0.0)
	f := c.F()
	if f != 32.0 {
		t.Errorf("expected 32.0°F, got %.1f", f)
	}
}

func TestRectangleArea(t *testing.T) {
	r := Rectangle{Width: 3, Height: 4}
	if r.Area() != 12.0 {
		t.Errorf("expected 12.0, got %.1f", r.Area())
	}
}

func TestRectanglePerimeter(t *testing.T) {
	r := Rectangle{Width: 3, Height: 4}
	if r.Perimeter() != 14.0 {
		t.Errorf("expected 14.0, got %.1f", r.Perimeter())
	}
}

func TestMethodExpression(t *testing.T) {
	r := Rectangle{Width: 2, Height: 5}
	expr := Rectangle.Area
	result := expr(r)
	if result != 10.0 {
		t.Errorf("expected 10.0 from method expression, got %.1f", result)
	}
}

func TestMilesToKilometers(t *testing.T) {
	m := Miles(100.0)
	km := m.ToKilometers()
	expected := Miles(160.934)
	if km != expected {
		t.Errorf("expected %.3f km, got %.3f", expected, km)
	}
}

func TestTripSpeed(t *testing.T) {
	trip := Trip{Distance: Miles(60), DurationMinutes: 60.0}
	speed := trip.Speed()
	if speed != 60.0 {
		t.Errorf("expected 60.0 mph, got %.1f", speed)
	}
}

func TestTripSpeedHalfHour(t *testing.T) {
	trip := Trip{Distance: Miles(30), DurationMinutes: 30.0}
	speed := trip.Speed()
	if speed != 60.0 {
		t.Errorf("expected 60.0 mph, got %.1f", speed)
	}
}
