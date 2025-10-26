package primitive

import (
	"github.com/ironsmile/raytracer/shape"
	"github.com/ironsmile/raytracer/transform"
)

// NewObject parses an .obj file (`filePath`) and returns an Object, which represents it. It places
// the object at the position, given by its second argument - `center`.
func NewObject(filePath string) ([]*BasePrimitive, error) {

	shapes, err := shape.NewObject(filePath)
	if err != nil {
		return nil, err
	}
	var prims []*BasePrimitive
	for _, shape := range shapes {
		obj := &BasePrimitive{shape: shape}
		obj.SetTransform(transform.Identity())
		obj.id = GetNewID()
		prims = append(prims, obj)
	}
	return prims, nil
}
