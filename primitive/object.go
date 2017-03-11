package primitive

import (
	"fmt"

	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/shape"
)

type Object struct {
	BasePrimitive

	// Identification string for thins object
	id string
}

func (o *Object) GetType() int {
	return OBJECT
}

func (o *Object) Intersect(ray *geometry.Ray, dist float64) (int, float64, *geometry.Vector) {
	return o.shape.Intersect(ray, dist)
}

func (o *Object) String() string {
	return fmt.Sprintf("Object <%s>", o.id)
}

// NewObject parses an .obj file (`filePath`) and returns an Object, which represents it. It places
// the object at the position, given by its second argument - `center`.
func NewObject(filePath string, center *geometry.Point) (*Object, error) {
	oShape, err := shape.NewObject(filePath, center)
	if err != nil {
		return nil, err
	}
	obj := &Object{id: filePath}
	obj.shape = oShape
	return obj, nil
}
