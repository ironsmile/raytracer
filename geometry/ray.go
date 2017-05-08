package geometry

import "math"

// Ray represents a straight line with origin and a direction
type Ray struct {
	Origin    Vector
	Direction Vector

	Mint float64
	Maxt float64

	Debug bool
}

// BackToDefaults zeroes out a ray which can then be reused somewhere int the program
func (r *Ray) BackToDefaults() {
	r.Debug = false
	r.Mint = 0
	r.Maxt = math.MaxFloat64
}

// At returns the point which as at distance t from the Origin in Direction
func (r *Ray) At(t float64) Vector {
	return r.Origin.Plus(r.Direction.MultiplyScalar(t))
}

// Refract returns refraction direction between this ray's direction and the normal n
// for ray passing from material with refraction index n1 to one with refraction index n2.
func (r *Ray) Refract(n Vector, n1, n2 float64) (refrDirection Vector, tir bool) {
	nr := n1 / n2
	cosI := -n.Dot(r.Direction)
	sinT2 := nr * nr * (1 - cosI*cosI)
	if sinT2 > 1.0 {
		tir = true
		return
	}
	cosT := math.Sqrt(1 - sinT2)

	refrDirection = r.Direction.MultiplyScalar(nr).Plus(
		n.MultiplyScalar(nr*cosI - cosT))
	return
}

// Intersect returns the intersection point between two rays. Two rays may not always
// intersect so that the second argument says wether there is an intersectoin at all
func (r *Ray) Intersect(o Ray) (Vector, bool) {
	const width = 0.03

	clampInRange := func(p Vector) (Vector, bool) {
		dist := r.Origin.Distance(p)
		if dist < r.Mint || dist > r.Maxt {
			return r.Origin, false
		}

		return p, true
	}

	if r.Origin == o.Origin {
		return r.Origin, true
	}

	d3 := r.Direction.Cross(o.Direction)

	if !d3.Equals(NewVector(0, 0, 0)) {
		matrix := [12]float64{
			r.Direction.X,
			-o.Direction.X,
			d3.X,
			o.Origin.X - r.Origin.X,

			r.Direction.Y,
			-o.Direction.Y,
			d3.Y,
			o.Origin.Y - r.Origin.Y,

			r.Direction.Z,
			-o.Direction.Z,
			d3.Z,
			o.Origin.Z - r.Origin.Z,
		}

		result := solve(matrix, 3, 4)

		a := result[3]
		b := result[7]
		c := result[11]

		if a >= 0 && b >= 0 {
			dist := d3.MultiplyScalar(c)
			if dist.Length() <= width {
				return clampInRange(r.At(a))
			}
			return r.Origin, false
		}
	}

	dP := o.Origin.Multiply(r.Origin)

	a2 := r.Direction.Dot(dP)
	b2 := o.Direction.Dot(dP.Neg())

	if a2 < 0 && b2 < 0 {
		dist := r.Origin.Distance(dP)
		if dP.Length() <= width {
			return clampInRange(r.At(dist))
		}
		return r.Origin, false
	}

	p3a := r.Origin.Plus(r.Direction.MultiplyScalar(a2))
	d3a := o.Origin.Minus(p3a)

	p3b := r.Origin
	d3b := o.Origin.Plus(o.Direction.MultiplyScalar(b2)).Minus(p3b)

	if b2 < 0 {
		if d3a.Length() <= width {
			return clampInRange(p3a)
		}
		return r.Origin, false
	}

	if a2 < 0 {
		if d3b.Length() <= width {
			return clampInRange(p3b)
		}
		return r.Origin, false
	}

	if d3a.Length() <= d3b.Length() {
		if d3a.Length() <= width {
			return clampInRange(p3a)
		}
		return r.Origin, false
	}

	if d3b.Length() <= width {
		return clampInRange(p3b)
	}

	return r.Origin, false
}

// NewRay retursn a new ray with Min zero and Max the maximum float64 value
func NewRay(origin, direction Vector) Ray {
	return Ray{
		Origin:    origin,
		Direction: direction,
		Maxt:      math.MaxFloat64,
	}
}

func solve(matrix [12]float64, rows, cols int) [12]float64 {

	for i := 0; i < cols-1; i++ {
		for j := i; j < rows; j++ {
			if matrix[i+j*cols] != 0 {
				if i != j {
					for k := i; k < cols; k++ {
						temp := matrix[k+j*cols]
						matrix[k+j*cols] = matrix[k+i*cols]
						matrix[k+i*cols] = temp
					}
				}

				j = i

				for v := 0; v < rows; v++ {
					if v == j {
						continue
					} else {
						factor := matrix[i+v*cols] / matrix[i+j*cols]
						matrix[i+v*cols] = 0

						for u := i + 1; u < cols; u++ {
							matrix[u+v*cols] -= factor * matrix[u+j*cols]
							matrix[u+j*cols] /= matrix[i+j*cols]
						}
						matrix[i+j*cols] = 1
					}
				}

				break
			}
		}
	}

	return matrix
}
