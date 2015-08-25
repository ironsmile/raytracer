package geometry

import (
	"fmt"
	"math"
)

var COMPARE_PRECISION = 1e-7

type Point struct {
	X, Y, Z float64
}

func (p *Point) Minus(other *Point) *Vector {
	return &Vector{p.X - other.X, p.Y - other.Y, p.Z - other.Z}
}

func (p *Point) MinusVector(vec *Vector) *Point {
	return &Point{p.X - vec.X, p.Y - vec.Y, p.Z - vec.Z}
}

func (p *Point) MinusVectorIP(vec *Vector) *Point {
	p.X, p.Y, p.Z = p.X-vec.X, p.Y-vec.Y, p.Z-vec.Z
	return p
}

func (p *Point) MinusInVector(other *Point, v *Vector) *Vector {
	v.X, v.Y, v.Z = p.X-other.X, p.Y-other.Y, p.Z-other.Z
	return v
}

func (p *Point) PlusVector(vec *Vector) *Point {
	return &Point{p.X + vec.X, p.Y + vec.Y, p.Z + vec.Z}
}

func (p *Point) PlusVectorIP(vec *Vector) *Point {
	p.X, p.Y, p.Z = p.X+vec.X, p.Y+vec.Y, p.Z+vec.Z
	return p
}

func (p *Point) Plus(other *Point) *Point {
	return &Point{p.X + other.X, p.Y + other.Y, p.Z + other.Z}
}

func (p *Point) PlusIP(other *Point) *Point {
	p.X, p.Y, p.Z = p.X+other.X, p.Y+other.Y, p.Z+other.Z
	return p
}

func (p *Point) MultiplyScalar(sclr float64) *Point {
	return &Point{p.X * sclr, p.Y * sclr, p.Z * sclr}
}

func (p *Point) MultiplyScalarIP(sclr float64) *Point {
	p.X, p.Y, p.Z = p.X*sclr, p.Y*sclr, p.Z*sclr
	return p
}

func (p *Point) String() string {
	return fmt.Sprintf("Point<%.8f, %.8f, %.8f>", p.X, p.Y, p.Z)
}

func (p *Point) Equals(other *Point) bool {
	if math.Abs(p.X-other.X) > COMPARE_PRECISION {
		return false
	}
	if math.Abs(p.Y-other.Y) > COMPARE_PRECISION {
		return false
	}
	if math.Abs(p.Z-other.Z) > COMPARE_PRECISION {
		return false
	}
	return true
}

func (p *Point) Vector() *Vector {
	return &Vector{p.X, p.Y, p.Z}
}

func (p *Point) Copy() *Point {
	return &Point{p.X, p.Y, p.Z}
}

func (p *Point) CopyToSelf(o *Point) *Point {
	p.X, p.Y, p.Z = o.X, o.Y, o.Z
	return p
}

func NewPoint(x, y, z float64) *Point {
	return &Point{x, y, z}
}
