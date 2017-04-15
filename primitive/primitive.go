package primitive

import (
	"github.com/ironsmile/raytracer/bbox"
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/mat"
	"github.com/ironsmile/raytracer/shape"
	"github.com/ironsmile/raytracer/transform"
)

// Primitive is the type which marries the shape to its material. It is resposible for
// the geometry and shading of objects.
type Primitive interface {
	Intersect(geometry.Ray) (pr Primitive, distance float64, normal geometry.Vector)
	IntersectBBoxEdge(geometry.Ray) bool
	GetWorldBBox() *bbox.BBox
	GetColor() *geometry.Color
	GetMaterial() *mat.Material
	IsLight() bool
	GetLightSource() geometry.Vector
	GetName() string
	Shape() shape.Shape
}

// BasePrimitive implements some common methods. This way an actual primitive can be composed
// of one BasicPrimitive which alredy has an implementation of most of the methods.
type BasePrimitive struct {
	Mat         mat.Material
	Light       bool
	LightSource geometry.Vector
	Name        string
	shape       shape.Shape

	objToWorld transform.Transform
	worldToObj transform.Transform

	worldBBox *bbox.BBox
}

// GetLightSource is strange hacky methods. Returns the ligth source point in the world space
// from which the light eliminates for this light primitive.
func (b *BasePrimitive) GetLightSource() geometry.Vector {
	return b.LightSource
}

// GetName returns the name of this primitive as set in the scene
func (b *BasePrimitive) GetName() string {
	return b.Name
}

// IsLight returns true if this primitive is a light source
func (b *BasePrimitive) IsLight() bool {
	return b.Light
}

// GetColor is a hacky method which assumes the whole primitive is from one color and returns it
func (b *BasePrimitive) GetColor() *geometry.Color {
	return b.Mat.Color
}

// GetMaterial returns thie primitive's material
func (b *BasePrimitive) GetMaterial() *mat.Material {
	return &b.Mat
}

// Intersect returns whether a ray intersects this primitive and at what distance from the ray
// origin is this intersection.
func (b *BasePrimitive) Intersect(ray geometry.Ray) (Primitive, float64, geometry.Vector) {

	worldBound := b.GetWorldBBox()
	if worldBound != nil {
		intersected, minT, maxT := worldBound.IntersectP(ray)
		if !intersected {
			return nil, ray.Maxt, ray.Direction
		}
		ray.Mint = minT - geometry.EPSILON
		ray.Maxt = maxT + geometry.EPSILON
	}

	ray = b.worldToObj.Ray(ray)
	res, hitDist, normal := b.shape.Intersect(ray)

	if res != shape.HIT {
		return nil, hitDist, normal
	}

	normal = b.objToWorld.Normal(normal)

	return b, hitDist, normal
}

// IntersectBBoxEdge returns whether a ray intersects this primitive's bounding box
func (b *BasePrimitive) IntersectBBoxEdge(ray geometry.Ray) bool {
	worldBound := b.GetWorldBBox()

	if worldBound == nil {
		return false
	}

	intersected, distance := worldBound.IntersectEdge(ray)

	if !intersected || distance < ray.Mint || distance > ray.Maxt {
		return false
	}

	return true
}

// Shape returns this primitive's shape if there is one
func (b *BasePrimitive) Shape() shape.Shape {
	return b.shape
}

// SetTransform sets the transformation matrices for this primitive's shape
func (b *BasePrimitive) SetTransform(t *transform.Transform) {
	b.objToWorld = *t
	b.worldToObj = *t.Inverse()
	b.refreshWorldBBox()
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
