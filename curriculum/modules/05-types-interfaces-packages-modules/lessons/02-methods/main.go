package main

import "fmt"

type Celsius float64

func (c Celsius) F() float64 {
	return float64(c)*1.8 + 32
}

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

type Miles float64

func (m Miles) ToKilometers() Miles {
	return m * 1.60934
}

type Trip struct {
	Distance        Miles
	DurationMinutes float64
}

func (t Trip) Speed() float64 {
	return float64(t.Distance) / (t.DurationMinutes / 60.0)
}

func main() {
	c := Celsius(100.0)
	fmt.Printf("%.1f°C = %.1f°F\n", c, c.F())

	rect := Rectangle{Width: 3, Height: 4}
	fmt.Printf("area: %.1f, perimeter: %.1f\n", rect.Area(), rect.Perimeter())

	methodExpr := Rectangle.Area
	fmt.Printf("via method expr: %.1f\n", methodExpr(rect))

	miles := Miles(100.0)
	fmt.Printf("%.1f miles = %.1f km\n", miles, miles.ToKilometers())

	trip := Trip{Distance: Miles(60), DurationMinutes: 90.0}
	fmt.Printf("speed: %.1f mph\n", trip.Speed())
}
