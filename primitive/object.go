package primitive

import (
	"fmt"

	"github.com/ironsmile/raytracer/shape"
	"github.com/ironsmile/raytracer/transform"
)

type Object struct {
	BasePrimitive

	// Identification string for thins object
	id string
}

func (o *Object) GetType() int {
	return OBJECT
}

func (o *Object) String() string {
	return fmt.Sprintf("Object <%s>", o.id)
}

// NewObject parses an .obj file (`filePath`) and returns an Object, which represents it. It places
// the object at the position, given by its second argument - `center`.
func NewObject(filePath string) (*Object, error) {
	oShape, err := shape.NewObject(filePath)
	if err != nil {
		return nil, err
	}
	obj := &Object{id: filePath}
	obj.shape = oShape
	obj.SetTransform(transform.Identity())
	return obj, nil
}
