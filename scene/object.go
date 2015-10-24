package scene

import (
	"fmt"
	"os"

	"github.com/momchil-atanasov/go-data-front/decoder/obj"

	"github.com/ironsmile/raytracer/geometry"
)

type Object struct {
	BasePrimitive

	id    string
	model *obj.Model

	Triangles []Primitive
}

func (o *Object) GetType() int {
	return OBJECT
}

func (o *Object) Intersect(ray *geometry.Ray, dist float64) (int, float64) {
	prim, distance := IntersectPrimitives(o.Triangles, ray)
	if prim == nil {
		return MISS, distance
	}
	return HIT, distance
}

func (o *Object) GetNormal(pos *geometry.Point) *geometry.Vector {
	//!TODO: implement
	return &geometry.Vector{0, 1, 0}
}

func (o *Object) String() string {
	return fmt.Sprintf("Object <%s>", o.id)
}

func NewObject(filePath string) (Primitive, error) {
	decoder := obj.NewDecoder(obj.DefaultLimits())
	objFile, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}
	defer objFile.Close()

	model, err := decoder.Decode(objFile)

	if err != nil {
		return nil, err
	}

	fmt.Printf("model %s has %d models\n", filePath, len(model.Objects))

	o := new(Object)
	o.Triangles = make([]Primitive, 0, len(model.Vertices)/3+1)
	o.id = filePath
	o.model = model

	//!TODO: maybe remove this scale factor?
	scaleFactor := 100.0

	for _, obj := range model.Objects {
		fmt.Printf("object %s has %d meshes\n", obj.Name, len(obj.Meshes))
		for meshIndex, mesh := range obj.Meshes {
			fmt.Printf("mesh %d is from %s and has %d faces\n", meshIndex, mesh.MaterialName,
				len(mesh.Faces))
			for faceIndex, face := range mesh.Faces {
				if len(face.References) != 3 {
					return nil, fmt.Errorf(
						"face %d [mesh: %d, obj: %s] has %d points, don't know how to load it",
						faceIndex, meshIndex, obj.Name, len(face.References))
				}

				a := model.Vertices[face.References[0].VertexIndex]
				b := model.Vertices[face.References[1].VertexIndex]
				c := model.Vertices[face.References[2].VertexIndex]

				triangleVertices := [3]*geometry.Point{
					geometry.NewPoint(a.X/scaleFactor, a.Y/scaleFactor, a.Z/scaleFactor),
					geometry.NewPoint(b.X/scaleFactor, b.Y/scaleFactor, b.Z/scaleFactor),
					geometry.NewPoint(c.X/scaleFactor, c.Y/scaleFactor, c.Z/scaleFactor),
				}

				o.Triangles = append(o.Triangles, NewTriangle(triangleVertices))
			}
		}
	}

	fmt.Printf("%s has %d triangles\n", o, len(o.Triangles))

	return o, nil
}
