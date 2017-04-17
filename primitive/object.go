package primitive

import (
	"github.com/ironsmile/raytracer/shape"
	"github.com/ironsmile/raytracer/transform"
)

// NewObject parses an .obj file (`filePath`) and returns an Object, which represents it. It places
// the object at the position, given by its second argument - `center`.
func NewObject(filePath string) (*BasePrimitive, error) {
	oShape, err := shape.NewObject(filePath)
	if err != nil {
		return nil, err
	}
	obj := &BasePrimitive{shape: oShape}
	obj.SetTransform(transform.Identity())
	return obj, nil
}
