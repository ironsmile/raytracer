package primitive

import (
	"github.com/ironsmile/raytracer/bbox"
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/mat"
	"github.com/ironsmile/raytracer/shape"
	"github.com/ironsmile/raytracer/transform"
)

// BasePrimitive implements some common methods. This way an actual primitive can be
// composed of one BasicPrimitive which alredy has an implementation of most
// of the methods.
type BasePrimitive struct {
	Mat         mat.Material
	Light       bool
	LightSource geometry.Vector
	shape       shape.Shape

	objToWorld *transform.Transform
	worldToObj *transform.Transform

	worldBBox *bbox.BBox
}

// CanIntersect returns true when this parimitive can be intersected directly and false
// when it should be refined before intersection.
func (b *BasePrimitive) CanIntersect() bool {
	return b.shape.CanIntersect()
}

// Refine returns the slice of primitives this one primitive is made of. This method
// should be called only when CanIntersect returns false
func (b *BasePrimitive) Refine() []Primitive {
	if b.shape.CanIntersect() {
		panic("Refine should not be called on intersectable primitive: BasePrimitive")
	}
	var prims []Primitive

	for _, objShape := range b.Shape().Refine() {
		pr := FromShape(objShape)
		pr.SetTransform(b.objToWorld)
		prims = append(prims, pr)
	}

	return prims
}

// GetLightSource is strange hacky methods. Returns the ligth source point in the world
// space from which the light eliminates for this light primitive.
func (b *BasePrimitive) GetLightSource() geometry.Vector {
	return b.LightSource
}

// IsLight returns true if this primitive is a light source
func (b *BasePrimitive) IsLight() bool {
	return b.Light
}

// GetColor is a hacky method which assumes the whole primitive is from one color and
// returns it
func (b *BasePrimitive) GetColor() *geometry.Color {
	return b.Mat.Color
}

// GetMaterial returns thie primitive's material
func (b *BasePrimitive) GetMaterial() *mat.Material {
	return &b.Mat
}

// Intersect returns whether a ray intersects this primitive and at what distance from
// the ray origin is this intersection.
func (b *BasePrimitive) Intersect(ray geometry.Ray, in *Intersection) bool {
	if b.IsLight() {
		return false
	}

	worldBound := b.GetWorldBBox()
	intersected, _, _ := worldBound.IntersectP(ray)
	if !intersected {
		return false
	}

	ray = b.worldToObj.Ray(ray)

	if hit := b.shape.Intersect(ray, &in.DfGeometry); !hit {
		return false
	}

	in.Primitive = b
	return true
}

// IntersectP returns whether a ray intersects this primitive and nothing more
func (b *BasePrimitive) IntersectP(ray geometry.Ray) bool {
	if b.IsLight() {
		return false
	}

	worldBound := b.GetWorldBBox()
	intersected, _, _ := worldBound.IntersectP(ray)
	if !intersected {
		return false
	}

	ray = b.worldToObj.Ray(ray)
	return b.shape.IntersectP(ray)
}

// IntersectBBoxEdge returns whether a ray intersects this primitive's bounding box
func (b *BasePrimitive) IntersectBBoxEdge(ray geometry.Ray) bool {
	worldBound := b.GetWorldBBox()

	if worldBound == nil {
		return false
	}

	intersected, _ := worldBound.IntersectEdge(ray)

	return intersected
}

// Shape returns this primitive's shape if there is one
func (b *BasePrimitive) Shape() shape.Shape {
	return b.shape
}

// SetTransform sets the transformation matrices for this primitive's shape. Accepts the
// object-to-world transformation matrix
func (b *BasePrimitive) SetTransform(t *transform.Transform) {
	b.objToWorld = t
	b.worldToObj = t.Inverse()
	b.refreshWorldBBox()
}

// GetTransforms returns the two transformation matrices for this primiitive:
// object-to-world and world-to-object
func (b *BasePrimitive) GetTransforms() (*transform.Transform, *transform.Transform) {
	return b.objToWorld, b.worldToObj
}

// GetWorldBBox returns the bound box around this primitive in world space
func (b *BasePrimitive) GetWorldBBox() *bbox.BBox {
	if b.worldBBox != nil {
		return b.worldBBox
	}
	b.refreshWorldBBox()
	return b.worldBBox
}

func (b *BasePrimitive) refreshWorldBBox() {
	objBBox := b.shape.GetObjectBBox()
	b.worldBBox = bbox.FromPoint(b.objToWorld.Point(objBBox.Min))
	b.worldBBox = bbox.UnionPoint(b.worldBBox, b.objToWorld.Point(objBBox.Max))
}

// FromShape returns a primitive from a given shape
func FromShape(s shape.Shape) *BasePrimitive {
	b := &BasePrimitive{shape: s}
	b.SetTransform(transform.Identity())
	return b
}
