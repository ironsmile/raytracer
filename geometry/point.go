package geometry

import (
	"fmt"
)

type Point struct {
	X, Y, Z float64
}

func (p *Point) Minus(other *Point) *Vector {
	return &Vector{p.X - other.X, p.Y - other.Y, p.Z - other.Z}
}

func (p *Point) MinusVector(vec *Vector) *Point {
	return &Point{p.X - vec.X, p.Y - vec.Y, p.Z - vec.Z}
}

func (p *Point) PlusVector(vec *Vector) *Point {
	return &Point{p.X + vec.X, p.Y + vec.Y, p.Z + vec.Z}
}

func (p *Point) Plus(other *Point) *Point {
	return &Point{p.X + other.X, p.Y + other.Y, p.Z + other.Z}
}

func (p *Point) MultiplyScalar(sclr float64) *Point {
	return &Point{p.X * sclr, p.Y * sclr, p.Z * sclr}
}

func (p *Point) String() string {
	return fmt.Sprintf("Point<%f, %f, %f>", p.X, p.Y, p.Z)
}

func NewPoint(x, y, z float64) *Point {
	return &Point{x, y, z}
}
