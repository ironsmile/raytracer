package accel

import (
	"math"

	"github.com/ironsmile/raytracer/bbox"
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/primitive"
	"github.com/ironsmile/raytracer/utils"
)

// Grid is an accelerator that divides an axis-aligned region of space into equal-sized
// box-shaped chunks (called voxels). Each voxel stores references to the primitives that
// overlap it. Given a ray, the grid steps through each of the voxels that the
// ray passes through in order, checking for intersections with only the primitives in each
// voxel. Useless ray intersection tests are reduced substantially because primitives far
// away from the ray aren’t considered at all. Furthermore, because the voxels are
// considered from near to far along the ray, it is possible to stop performing
// intersection tests once an intersection has been found and it is certain that
// it is not possible for there to be any closer intersections.
type Grid struct {
	Base

	bounds  *bbox.BBox
	nVoxels [3]int

	width    geometry.Vector
	invWidth geometry.Vector

	voxels []*Voxel
}

// NewGrid returns a new grid accelerator for a slice of primitives
func NewGrid(p []primitive.Primitive) *Grid {

	g := &Grid{}

	g.primitives = p

	for i := 0; i < len(p); i++ {
		g.bounds = bbox.Union(g.bounds, p[i].GetWorldBBox())
	}

	delta := g.bounds.Max.Minus(g.bounds.Min)

	maxAxis := g.bounds.MaximumExtend()
	invMaxWidth := 1.0 / delta.ByAxis(maxAxis)
	cubeRoot := 3.0 * math.Pow(float64(len(g.primitives)), 1.0/3.0)
	voxelsPerUnitDist := cubeRoot * invMaxWidth

	for axis := 0; axis < 3; axis++ {
		g.nVoxels[axis] = utils.RoundToInt(delta.ByAxis(axis) * voxelsPerUnitDist)
		g.nVoxels[axis] = utils.ClampInt(g.nVoxels[axis], 1, 64)
	}

	for axis := 0; axis < 3; axis++ {
		g.width.SetByAxis(axis, delta.ByAxis(axis)/float64(g.nVoxels[axis]))

		if g.width.ByAxis(axis) != 0 {
			g.invWidth.SetByAxis(axis, 1/g.width.ByAxis(axis))
		}
	}

	nv := g.nVoxels[0] * g.nVoxels[1] * g.nVoxels[2]
	g.voxels = make([]*Voxel, nv, nv)

	for i := 0; i < len(g.primitives); i++ {
		pb := g.primitives[i].GetWorldBBox()
		var vmin, vmax [3]int
		for axis := 0; axis < 3; axis++ {
			vmin[axis] = g.posToVoxel(pb.Min, axis)
			vmax[axis] = g.posToVoxel(pb.Max, axis)
		}

		for z := vmin[2]; z <= vmax[2]; z++ {
			for y := vmin[1]; y <= vmax[1]; y++ {
				for x := vmin[0]; x <= vmax[0]; x++ {
					o := g.offset(x, y, z)
					if g.voxels[o] == nil {
						g.voxels[o] = NewVoxel()
					}
					g.voxels[o].Add(g.primitives[i])
				}
			}
		}
	}

	return g
}

func (g *Grid) offset(x, y, z int) int {
	return z*g.nVoxels[0]*g.nVoxels[1] + y*g.nVoxels[0] + x
}

func (g *Grid) posToVoxel(p geometry.Vector, axis int) int {
	v := int((p.ByAxis(axis) - g.bounds.Min.ByAxis(axis)) * g.invWidth.ByAxis(axis))
	return utils.ClampInt(v, 0, g.nVoxels[axis]-1)
}

func (g *Grid) voxelToPos(p int, axis int) float64 {
	return g.bounds.Min.ByAxis(axis) + float64(p)*g.width.ByAxis(axis)
}

// Intersect implements the Primitive interface
func (g *Grid) Intersect(ray geometry.Ray) (primitive.Primitive, float64, geometry.Vector) {
	var rayT float64

	if g.bounds.Inside(ray.At(ray.Mint)) {
		rayT = ray.Mint
	} else if intersected, tNear, _ := g.bounds.IntersectP(ray); intersected {
		rayT = tNear
	} else {
		return nil, 0, ray.Direction
	}

	gridIntersect := ray.At(rayT)

	var nextCrossingT, deltaT [3]float64
	var step, out, pos [3]int

	for axis := 0; axis < 3; axis++ {
		pos[axis] = g.posToVoxel(gridIntersect, axis)
		axisComp := ray.Direction.ByAxis(axis)

		if axisComp >= 0 {
			nextCrossingT[axis] = rayT + (g.voxelToPos(pos[axis]+1, axis)-
				gridIntersect.ByAxis(axis))/axisComp
			deltaT[axis] = g.width.ByAxis(axis) / axisComp
			step[axis] = 1
			out[axis] = g.nVoxels[axis]
		} else {
			nextCrossingT[axis] = rayT + (g.voxelToPos(pos[axis], axis)-
				gridIntersect.ByAxis(axis))/axisComp
			deltaT[axis] = -g.width.ByAxis(axis) / axisComp
			step[axis] = -1
			out[axis] = -1
		}
	}

	var ret primitive.Primitive
	var retNormal geometry.Vector

	for {
		voxel := g.voxels[g.offset(pos[0], pos[1], pos[2])]
		if voxel != nil {
			if pr, dist, normal := voxel.Intersect(ray); pr != nil {
				retNormal = normal
				ret = pr
				ray.Maxt = dist
			}
		}

		var stepAxis int

		if nextCrossingT[0] < nextCrossingT[1] && nextCrossingT[0] < nextCrossingT[2] {
			stepAxis = 0
		} else if nextCrossingT[1] < nextCrossingT[2] {
			stepAxis = 1
		} else {
			stepAxis = 2
		}

		if ray.Maxt < nextCrossingT[stepAxis] {
			break
		}

		pos[stepAxis] += step[stepAxis]

		if pos[stepAxis] == out[stepAxis] {
			break
		}

		nextCrossingT[stepAxis] += deltaT[stepAxis]
	}

	return ret, ray.Maxt, retNormal
}

// IntersectP implements the Primitive interface
func (g *Grid) IntersectP(ray geometry.Ray) bool {
	var rayT float64

	if g.bounds.Inside(ray.At(ray.Mint)) {
		rayT = ray.Mint
	} else if intersected, tNear, _ := g.bounds.IntersectP(ray); intersected {
		rayT = tNear
	} else {
		return false
	}

	gridIntersect := ray.At(rayT)

	var nextCrossingT, deltaT [3]float64
	var step, out, pos [3]int

	for axis := 0; axis < 3; axis++ {
		pos[axis] = g.posToVoxel(gridIntersect, axis)
		axisComp := ray.Direction.ByAxis(axis)

		if axisComp >= 0 {
			nextCrossingT[axis] = rayT + (g.voxelToPos(pos[axis]+1, axis)-
				gridIntersect.ByAxis(axis))/axisComp
			deltaT[axis] = g.width.ByAxis(axis) / axisComp
			step[axis] = 1
			out[axis] = g.nVoxels[axis]
		} else {
			nextCrossingT[axis] = rayT + (g.voxelToPos(pos[axis], axis)-
				gridIntersect.ByAxis(axis))/axisComp
			deltaT[axis] = -g.width.ByAxis(axis) / axisComp
			step[axis] = -1
			out[axis] = -1
		}
	}

	for {
		voxel := g.voxels[g.offset(pos[0], pos[1], pos[2])]
		if voxel != nil {
			if intersected := voxel.IntersectP(ray); intersected {
				return true
			}
		}

		var stepAxis int

		if nextCrossingT[0] < nextCrossingT[1] && nextCrossingT[0] < nextCrossingT[2] {
			stepAxis = 0
		} else if nextCrossingT[1] < nextCrossingT[2] {
			stepAxis = 1
		} else {
			stepAxis = 2
		}

		if ray.Maxt < nextCrossingT[stepAxis] {
			break
		}

		pos[stepAxis] += step[stepAxis]

		if pos[stepAxis] == out[stepAxis] {
			break
		}

		nextCrossingT[stepAxis] += deltaT[stepAxis]
	}

	return false
}

// GetWorldBBox implements the Primitive interface
func (g *Grid) GetWorldBBox() *bbox.BBox {
	return g.bounds
}
