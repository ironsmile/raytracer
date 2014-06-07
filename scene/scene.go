package scene

import (
	"fmt"
	"math"

	"github.com/ironsmile/raytracer/common"
)

const (
	HIT = iota
	MISS
	INPRIM
)

const (
	NOTHING = iota
	SPHERE
	PLANE
)

type Material struct {
	Color *common.Color
	Refl  float64
	Diff  float64
}

func (m *Material) GetSpecular() float64 {
	return 1.0 - m.Diff
}

func NewMaterial() *Material {
	mat := new(Material)
	col := common.NewColor(0.2, 0.2, 0.2)
	mat.Color = col
	mat.Refl = 0.0
	mat.Diff = 0.2

	return mat
}

type Primitive interface {
	GetType() int
	Intersect(*common.Ray, float64) (int, float64)
	GetNormal(*common.Vector) *common.Vector
	GetColor() common.Color
	GetMaterial() Material
	IsLight() bool
	GetName() string
}

type BasePrimitive struct {
	Mat   Material
	Light bool
	Name  string
}

func (b *BasePrimitive) GetName() string {
	return b.Name
}

func (p *BasePrimitive) IsLight() bool {
	return p.Light
}

func (b *BasePrimitive) GetColor() common.Color {
	return *b.Mat.Color
}

func (b *BasePrimitive) GetMaterial() Material {
	return b.Mat
}

type Sphere struct {
	BasePrimitive

	Center   *common.Vector
	SqRadius float64
	Radius   float64
	RRadius  float64
}

func (s *Sphere) GetType() int {
	return SPHERE
}

func (s *Sphere) Intersect(ray *common.Ray, dist float64) (int, float64) {
	v := ray.Origin.Minus(s.Center)
	b := -v.Product(ray.Direction)
	det := b*b - v.Product(v) + s.SqRadius

	retdist := dist
	retval := MISS
	if det <= 0 {
		return retval, retdist
	}

	det = math.Sqrt(det)

	i1 := b - det
	i2 := b + det

	if i2 > 0 {
		if i1 < 0 {
			if i2 < dist {
				retdist = i2
				retval = INPRIM
			}
		} else {
			if i1 < dist {
				retdist = i1
				retval = HIT
			}
		}
	}

	return retval, retdist
}

func (s *Sphere) GetNormal(pos *common.Vector) *common.Vector {
	return pos.Minus(s.Center).MultiplyScalar(s.RRadius)
}

func (s *Sphere) String() string {
	return fmt.Sprintf("Sphere<center=%s, radius=%f>", s.Center, s.Radius)
}

func NewSphere(center common.Vector, radius float64) *Sphere {
	s := new(Sphere)
	s.Center = &center
	s.SqRadius = radius * radius
	s.Radius = radius
	s.RRadius = 1.0 / radius
	s.Mat = *NewMaterial()
	return s
}

type PlanePrim struct {
	BasePrimitive

	Plane *common.Plane
}

func (p *PlanePrim) GetType() int {
	return PLANE
}

func NewPlanePrim(normal common.Vector, d float64) *PlanePrim {
	plPrim := new(PlanePrim)
	plPrim.Plane = common.NewPlane(normal, d)
	return plPrim
}

func (p *PlanePrim) GetNormal(_ *common.Vector) *common.Vector {
	return p.Plane.N.Copy()
}

func (p *PlanePrim) GetD() float64 {
	return p.Plane.D
}

func (p *PlanePrim) Intersect(ray *common.Ray, dist float64) (int, float64) {
	d := p.Plane.N.Product(ray.Direction)

	if d == 0 {
		return MISS, dist
	}

	dst := -(p.Plane.N.Product(ray.Origin) + p.Plane.D) / d

	if dst > 0 {
		if dst < dist {
			return HIT, dst
		}
	}

	return MISS, dist
}

func (p *PlanePrim) String() string {
	return fmt.Sprintf("Plane<%s>", p.Name)
}

type Scene struct {
	Primitives []Primitive
	Lights     []Primitive
}

func (s *Scene) GetNrLights() int {
	return len(s.Lights)
}

func (s *Scene) GetLight(index int) Primitive {
	return s.Lights[index]
}

func (s *Scene) GetNrPrimitives() int {
	return len(s.Primitives)
}

func (s *Scene) GetPrimitive(index int) Primitive {
	return s.Primitives[index]
}

func (s *Scene) Intersect(ray *common.Ray) (Primitive, float64) {
	retdist := 1000000.0
	var prim Primitive = nil

	for sInd := 0; sInd < s.GetNrPrimitives(); sInd++ {
		pr := s.GetPrimitive(sInd)

		if pr == nil {
			fmt.Errorf("Primitive with index %d was nil\n", sInd)
		}

		res, resDist := pr.Intersect(ray, retdist)

		if res != MISS {
			prim = pr
			retdist = resDist
		}
	}

	return prim, retdist
}

func (s *Scene) InitScene() {
	s.Primitives = make([]Primitive, 0)
	s.Lights = make([]Primitive, 0)

	plane := NewPlanePrim(*common.NewVector(0, 1, 0), 4)
	plane.Name = "plane"
	plane.Mat.Refl = 0
	plane.Mat.Diff = 1.0
	plane.Mat.Color = common.NewColor(0.4, 0.3, 0.3)

	s.Primitives = append(s.Primitives, plane)

	plane = NewPlanePrim(*common.NewVector(1, 0, 0), 11)
	plane.Name = "plane"
	plane.Mat.Refl = 0
	plane.Mat.Diff = 1.0
	plane.Mat.Color = common.NewColor(0.4, 0.3, 0.3)

	s.Primitives = append(s.Primitives, plane)

	sphere := NewSphere(*common.NewVector(1, -0.8, 3), 2.5)
	sphere.Name = "big sphere"
	sphere.Mat.Refl = 0.8
	sphere.Mat.Diff = 0.9
	sphere.Mat.Color = common.NewColor(1, 0, 0)

	s.Primitives = append(s.Primitives, sphere)

	sphere = NewSphere(*common.NewVector(-5.5, -0.5, 7), 2)
	sphere.Name = "small sphere"
	sphere.Mat.Refl = 0.9
	sphere.Mat.Diff = 0.4
	sphere.Mat.Color = common.NewColor(0.7, 0.7, 1)

	s.Primitives = append(s.Primitives, sphere)

	sphere = NewSphere(*common.NewVector(-6.5, -2.5, 25), 1.5)
	sphere.Name = "small sphere far away"
	sphere.Mat.Refl = 0.9
	sphere.Mat.Diff = 0.4
	sphere.Mat.Color = common.NewColor(0.5, 1, 0)

	s.Primitives = append(s.Primitives, sphere)

	sphere = NewSphere(*common.NewVector(0, 5, 5), 0.1)
	sphere.Name = "Visible light source"
	sphere.Light = true
	sphere.Mat.Color = common.NewColor(0.9, 0.9, 0.9)

	s.Primitives = append(s.Primitives, sphere)
	s.Lights = append(s.Lights, sphere)

	sphere = NewSphere(*common.NewVector(2, 5, 1), 0.1)
	sphere.Name = "Invisible lightsource"
	sphere.Light = true
	sphere.Mat.Color = common.NewColor(0.9, 0.9, 0.9)

	s.Primitives = append(s.Primitives, sphere)
	s.Lights = append(s.Lights, sphere)
}

func NewScene() *Scene {
	scn := new(Scene)
	return scn
}
