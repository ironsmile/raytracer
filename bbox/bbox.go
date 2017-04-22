package bbox

import (
	"math"

	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/utils"
)

// BBox is a structure which defines a box around a object in 3D space
type BBox struct {
	Min geometry.Vector
	Max geometry.Vector
}

// Overlaps returns true if the two bounding boxes overlap
func (b *BBox) Overlaps(other *BBox) bool {
	panic("BBox.Overlaps is not implemented yet")
}

// Inside thells whether a point is inside the bounding box nor not
func (b *BBox) Inside(p geometry.Vector) bool {
	if p.X < b.Min.X || p.X > b.Max.X {
		return false
	}

	if p.Y < b.Min.Y || p.Y > b.Max.Y {
		return false
	}

	if p.Z < b.Min.Z || p.Z > b.Max.Z {
		return false
	}

	return true
}

// Expand modifies the bound box by scaling it with a certain scalar
func (b *BBox) Expand(delta float64) {
	b.Min.Minus(geometry.NewVector(delta, delta, delta))
	b.Max.Plus(geometry.NewVector(delta, delta, delta))
}

// SurfaceArea computes the surface area of the six faces of the box
func (b *BBox) SurfaceArea() float64 {
	d := b.Max.Minus(b.Min)
	return 2 * (d.X*d.Y + d.X*d.Z + d.Y*d.Z)
}

// Volume coputes the inside volume of the bounding box
func (b *BBox) Volume() float64 {
	panic("BBox.Volume is not implemented yet")
}

// MaximumExtend tells the caller which of the three axs is longest:
// 0 -> X
// 1 -> Y
// 2 -> Z
func (b *BBox) MaximumExtend() int {
	x := math.Abs(b.Max.X - b.Min.X)
	y := math.Abs(b.Max.Y - b.Min.Y)
	z := math.Abs(b.Max.Z - b.Min.Z)

	if x > y && x > z {
		return 0
	}

	if y > z {
		return 1
	}

	return 2
}

// Lerp lineary interpolates between the corners of the box by the given amount
func (b *BBox) Lerp(tx, ty, tz float64) geometry.Vector {
	panic("BBox.Lerp is not implemented yet")
}

// Offset returns the position of a point relative to the corners of the box, where
// a position at the minimum corner has offset (0, 0, 0), a point a the maximum corner
// has offset (1, 1, 1)
func (b *BBox) Offset(p *geometry.Vector) *geometry.Vector {
	panic("BBox.Offset is not implemented yet")
}

// Union returns a bounding box which ecompases the the two input boxes
func Union(one, other *BBox) *BBox {
	if one == nil && other == nil {
		return Null()
	}
	if one == nil || one.Min.X == math.MaxFloat64 {
		return other
	}
	if other == nil || other.Min.X == math.MaxFloat64 {
		return one
	}
	union := &BBox{}
	union.Min.X = math.Min(one.Min.X, other.Min.X)
	union.Min.Y = math.Min(one.Min.Y, other.Min.Y)
	union.Min.Z = math.Min(one.Min.Z, other.Min.Z)
	union.Max.X = math.Max(one.Max.X, other.Max.X)
	union.Max.Y = math.Max(one.Max.Y, other.Max.Y)
	union.Max.Z = math.Max(one.Max.Z, other.Max.Z)
	return union
}

// UnionPoint return a new bounding box which includes the original bounding box and a
// point.
func UnionPoint(bb *BBox, p geometry.Vector) *BBox {
	if bb == nil {
		return &BBox{
			Max: p,
			Min: p,
		}
	}
	union := &BBox{}
	union.Min.X = math.Min(bb.Min.X, p.X)
	union.Min.Y = math.Min(bb.Min.Y, p.Y)
	union.Min.Z = math.Min(bb.Min.Z, p.Z)
	union.Max.X = math.Max(bb.Max.X, p.X)
	union.Max.Y = math.Max(bb.Max.Y, p.Y)
	union.Max.Z = math.Max(bb.Max.Z, p.Z)
	return union
}

// IntersectP returns whether a ray intersects this boundng box and at which distances
func (b *BBox) IntersectP(ray geometry.Ray) (bool, float64, float64) {
	var t0 = ray.Mint
	var t1 = ray.Maxt
	var invRayDir, tNear, tFar float64

	invRayDir = 1.0 / ray.Direction.X
	tNear = (b.Min.X - ray.Origin.X) * invRayDir
	tFar = (b.Max.X - ray.Origin.X) * invRayDir
	if tNear > tFar {
		tNear, tFar = tFar, tNear
	}
	if tNear > t0 {
		t0 = tNear
	}
	if tFar < t1 {
		t1 = tFar
	}
	if t0 > t1 {
		return false, t0, t1
	}

	invRayDir = 1.0 / ray.Direction.Y
	tNear = (b.Min.Y - ray.Origin.Y) * invRayDir
	tFar = (b.Max.Y - ray.Origin.Y) * invRayDir
	if tNear > tFar {
		tNear, tFar = tFar, tNear
	}
	if tNear > t0 {
		t0 = tNear
	}
	if tFar < t1 {
		t1 = tFar
	}
	if t0 > t1 {
		return false, t0, t1
	}

	invRayDir = 1.0 / ray.Direction.Z
	tNear = (b.Min.Z - ray.Origin.Z) * invRayDir
	tFar = (b.Max.Z - ray.Origin.Z) * invRayDir
	if tNear > tFar {
		tNear, tFar = tFar, tNear
	}
	if tNear > t0 {
		t0 = tNear
	}
	if tFar < t1 {
		t1 = tFar
	}
	if t0 > t1 {
		return false, t0, t1
	}

	return true, t0, t1
}

// IntersectEdge tells whether a ray intersects any of the edges of this bbox
func (b *BBox) IntersectEdge(ray geometry.Ray) (bool, float64) {
	intersected, t0, t1 := b.IntersectP(ray)

	if !intersected {
		return false, 0
	}

	if t0 > ray.Maxt || t0 < ray.Mint {
		return false, 0
	}

	// Edge size
	bs := .04
	pNear := ray.Origin.Plus(ray.Direction.MultiplyScalar(t0))

	if (utils.EqualFloat64(pNear.Y, b.Min.Y, bs) && utils.EqualFloat64(pNear.Z, b.Min.Z, bs)) ||
		(utils.EqualFloat64(pNear.Y, b.Min.Y, bs) && utils.EqualFloat64(pNear.X, b.Min.X, bs)) ||
		(utils.EqualFloat64(pNear.X, b.Min.X, bs) && utils.EqualFloat64(pNear.Z, b.Min.Z, bs)) ||
		(utils.EqualFloat64(pNear.Y, b.Min.Y, bs) && utils.EqualFloat64(pNear.X, b.Max.X, bs)) ||
		(utils.EqualFloat64(pNear.X, b.Min.X, bs) && utils.EqualFloat64(pNear.Z, b.Max.Z, bs)) ||
		(utils.EqualFloat64(pNear.X, b.Max.X, bs) && utils.EqualFloat64(pNear.Z, b.Min.Z, bs)) ||
		(utils.EqualFloat64(pNear.Z, b.Min.Z, bs) && utils.EqualFloat64(pNear.Y, b.Max.Y, bs)) ||
		(utils.EqualFloat64(pNear.Z, b.Max.Z, bs) && utils.EqualFloat64(pNear.Y, b.Min.Y, bs)) ||
		(utils.EqualFloat64(pNear.Y, b.Max.Y, bs) && utils.EqualFloat64(pNear.Z, b.Max.Z, bs)) ||
		(utils.EqualFloat64(pNear.Y, b.Max.Y, bs) && utils.EqualFloat64(pNear.X, b.Max.X, bs)) ||
		(utils.EqualFloat64(pNear.X, b.Min.X, bs) && utils.EqualFloat64(pNear.Y, b.Max.Y, bs)) ||
		(utils.EqualFloat64(pNear.X, b.Max.X, bs) && utils.EqualFloat64(pNear.Z, b.Max.Z, bs)) {
		return true, t0
	}

	if t1 > ray.Maxt || t1 < ray.Mint {
		return false, 0
	}

	pFar := ray.Origin.Plus(ray.Direction.MultiplyScalar(t1))

	if (utils.EqualFloat64(pFar.Y, b.Min.Y, bs) && utils.EqualFloat64(pFar.Z, b.Min.Z, bs)) ||
		(utils.EqualFloat64(pFar.Y, b.Min.Y, bs) && utils.EqualFloat64(pFar.X, b.Min.X, bs)) ||
		(utils.EqualFloat64(pFar.X, b.Min.X, bs) && utils.EqualFloat64(pFar.Z, b.Min.Z, bs)) ||
		(utils.EqualFloat64(pFar.Y, b.Min.Y, bs) && utils.EqualFloat64(pFar.X, b.Max.X, bs)) ||
		(utils.EqualFloat64(pFar.X, b.Min.X, bs) && utils.EqualFloat64(pFar.Z, b.Max.Z, bs)) ||
		(utils.EqualFloat64(pFar.X, b.Max.X, bs) && utils.EqualFloat64(pFar.Z, b.Min.Z, bs)) ||
		(utils.EqualFloat64(pFar.Z, b.Min.Z, bs) && utils.EqualFloat64(pFar.Y, b.Max.Y, bs)) ||
		(utils.EqualFloat64(pFar.Z, b.Max.Z, bs) && utils.EqualFloat64(pFar.Y, b.Min.Y, bs)) ||
		(utils.EqualFloat64(pFar.Y, b.Max.Y, bs) && utils.EqualFloat64(pFar.Z, b.Max.Z, bs)) ||
		(utils.EqualFloat64(pFar.Y, b.Max.Y, bs) && utils.EqualFloat64(pFar.X, b.Max.X, bs)) ||
		(utils.EqualFloat64(pFar.X, b.Min.X, bs) && utils.EqualFloat64(pFar.Y, b.Max.Y, bs)) ||
		(utils.EqualFloat64(pFar.X, b.Max.X, bs) && utils.EqualFloat64(pFar.Z, b.Max.Z, bs)) {
		return true, t1
	}

	return false, 0
}

// FromPoint returns a new bounding box which bounds around a single post
func FromPoint(p geometry.Vector) *BBox {
	return &BBox{Min: p, Max: p}
}

// Null returns a bounding box which would return false for all intersections and
// any unions with it would result in a bounding box exatly the same as the other bbox or point
func Null() *BBox {
	return FromPoint(geometry.NewVector(math.MaxFloat64, math.MaxFloat64, math.MaxFloat64))
}

// New returns a bounding box defined by two points
func New(p1, p2 geometry.Vector) *BBox {
	return &BBox{
		Min: geometry.NewVector(math.Min(p1.X, p2.X), math.Min(p1.Y, p2.Y), math.Min(p1.Z, p2.Y)),
		Max: geometry.NewVector(math.Max(p1.X, p2.X), math.Max(p1.Y, p2.Y), math.Max(p1.Z, p2.Y)),
	}
}
