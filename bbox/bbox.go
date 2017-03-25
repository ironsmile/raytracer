package bbox

import (
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/shape"
	"github.com/ironsmile/raytracer/utils"
)

// BBox is a structure which defines a box around a object in 3D space
type BBox struct {
	Min geometry.Point
	Max geometry.Point
}

// Overlaps returns true if the two bounding boxes overlap
func (b *BBox) Overlaps(other *BBox) bool {
	panic("BBox.Overlaps is not implemented yet")
}

// Inside thells whether a point is inside the bounding box nor not
func (b *BBox) Inside(point *geometry.Point) bool {
	panic("BBox.Inside is not implemented yet")
}

// Expand modifies the bound box by scaling it with a certain scalar
func (b *BBox) Expand(delta float64) {
	b.Min.MinusVectorIP(geometry.NewVector(delta, delta, delta))
	b.Max.PlusVectorIP(geometry.NewVector(delta, delta, delta))
}

// SurfaceArea computes the surface area of the six faces of the box
func (b *BBox) SurfaceArea() float64 {
	panic("BBox.SurfaceArea is not implemented yet")
}

// Volume coputes the inside volume of the bounding box
func (b *BBox) Volume() float64 {
	panic("BBox.Volume is not implemented yet")
}

// MaximumExtend tells the caller which of the three axs is longest:
// 0 -> X
// 1 -> Y
// 2 -> Z
func (b *BBox) MaximumExtend() int8 {
	panic("BBox.MaximumExtend is not implemented yet")
}

// Lerp lineary interpolates between the corners of the box by the given amount
func (b *BBox) Lerp(tx, ty, tz float64) geometry.Point {
	panic("BBox.Lerp is not implemented yet")
}

// Offset returns the position of a point relative to the corners of the box, where
// a position at the minimum corner has offset (0, 0, 0), a point a the maximum corner
// has offset (1, 1, 1)
func (b *BBox) Offset(p *geometry.Point) *geometry.Vector {
	panic("BBox.Offset is not implemented yet")
}

// BoundingSphere returns a spehre that bounds the bounding box.
func (b *BBox) BoundingSphere() *shape.Sphere {
	panic("BBox.BoundingSphere is not implemted yet")
}

// Union returns a bounding box which ecompases the the two input boxes
func Union(one, other *BBox) *BBox {
	union := &BBox{}
	union.Min.X = utils.Min(one.Min.X, other.Min.X)
	union.Min.Y = utils.Min(one.Min.Y, other.Min.Y)
	union.Min.Z = utils.Min(one.Min.Z, other.Min.Z)
	union.Max.X = utils.Min(one.Max.X, other.Max.X)
	union.Max.Y = utils.Min(one.Max.Y, other.Max.Y)
	union.Max.Z = utils.Min(one.Max.Z, other.Max.Z)
	return union
}

// UnionPoint return a new bounding box which includes the original bounding box and a
// point.
func UnionPoint(bb *BBox, p *geometry.Point) *BBox {
	union := &BBox{}
	union.Min.X = utils.Min(bb.Min.X, p.X)
	union.Min.Y = utils.Min(bb.Min.Y, p.Y)
	union.Min.Z = utils.Min(bb.Min.Z, p.Z)
	union.Max.X = utils.Min(bb.Max.X, p.X)
	union.Max.Y = utils.Min(bb.Max.Y, p.Y)
	union.Max.Z = utils.Min(bb.Max.Z, p.Z)
	return union
}

func (b *BBox) IntersectP(ray geometry.Ray) (bool, float64, float64) {
	var t0, t1, invRayDir, tNear, tFar float64
	for i := 0; i < 3; i++ {
		invRayDir = 1.0 / ray.Direction.ByIndex(i)
		tNear = (b.Min.ByIndex(i) - ray.Origin.ByIndex(i)) * invRayDir
		tFar = (b.Max.ByIndex(i) - ray.Origin.ByIndex(i)) * invRayDir
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
	}
	return true, t0, t1
}

// FromPoint returns a new bounding box which bounds around a single post
func FromPoint(p *geometry.Point) *BBox {
	return &BBox{Min: *p, Max: *p}
}

// New returns a bounding box defined by two points
func New(p1, p2 *geometry.Point) *BBox {
	return &BBox{
		Min: *geometry.NewPoint(utils.Min(p1.X, p2.X), utils.Min(p1.Y, p2.Y), utils.Min(p1.Z, p2.Y)),
		Max: *geometry.NewPoint(utils.Max(p1.X, p2.X), utils.Max(p1.Y, p2.Y), utils.Max(p1.Z, p2.Y)),
	}
}
