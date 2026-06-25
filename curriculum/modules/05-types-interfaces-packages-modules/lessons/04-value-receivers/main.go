package main

import (
	"fmt"
	"math"
)

type Point struct {
	X, Y float64
}

func (p Point) DistanceToOrigin() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

func (p Point) Translated(dx, dy float64) Point {
	return Point{X: p.X + dx, Y: p.Y + dy}
}

type Temperature struct {
	Celsius float64
}

func (t Temperature) ToFahrenheit() float64 {
	return t.Celsius*1.8 + 32
}

func (t Temperature) ToKelvin() float64 {
	return t.Celsius + 273.15
}

func (t Temperature) IsFreezing() bool {
	return t.Celsius <= 0
}

func main() {
	p := Point{X: 3, Y: 4}
	fmt.Printf("distance: %.2f\n", p.DistanceToOrigin())

	p2 := p.Translated(1, 2)
	fmt.Printf("original: (%.1f, %.1f)\n", p.X, p.Y)
	fmt.Printf("translated: (%.1f, %.1f)\n", p2.X, p2.Y)

	pp := &Point{X: 5, Y: 12}
	fmt.Printf("pointer distance: %.2f\n", pp.DistanceToOrigin())

	temps := []Temperature{
		{Celsius: 100},
		{Celsius: 0},
		{Celsius: -10},
	}
	for _, t := range temps {
		fmt.Printf("%.1f°C = %.1f°F = %.1fK (freezing: %v)\n",
			t.Celsius, t.ToFahrenheit(), t.ToKelvin(), t.IsFreezing())
	}
}
