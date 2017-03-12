package transform

import (
	"fmt"
	"math"
)

var COMPARE_PRECISION = 1e-6

type Matrix4x4 struct {
	els [4][4]float64
}

func (m *Matrix4x4) Get(i, j int) float64 {
	return m.els[i][j]
}

func (m *Matrix4x4) Multiply(other *Matrix4x4) *Matrix4x4 {
	// // slower:
	// mat := &Matrix4x4{}

	// for i := 0; i < 4; i++ {
	// 	for j := 0; j < 4; j++ {
	// 		mat.els[i][j] = m.els[i][0]*other.els[0][j] +
	// 			m.els[i][1]*other.els[1][j] +
	// 			m.els[i][2]*other.els[2][j] +
	// 			m.els[i][3]*other.els[3][j]
	// 	}
	// }
	// return mat

	// faster:
	return &Matrix4x4{
		[4][4]float64{
			[4]float64{
				m.els[0][0]*other.els[0][0] +
					m.els[0][1]*other.els[1][0] +
					m.els[0][2]*other.els[2][0] +
					m.els[0][3]*other.els[3][0],

				m.els[0][0]*other.els[0][1] +
					m.els[0][1]*other.els[1][1] +
					m.els[0][2]*other.els[2][1] +
					m.els[0][3]*other.els[3][1],

				m.els[0][0]*other.els[0][2] +
					m.els[0][1]*other.els[1][2] +
					m.els[0][2]*other.els[2][2] +
					m.els[0][3]*other.els[3][2],

				m.els[0][0]*other.els[0][3] +
					m.els[0][1]*other.els[1][3] +
					m.els[0][2]*other.els[2][3] +
					m.els[0][3]*other.els[3][3],
			},
			[4]float64{
				m.els[1][0]*other.els[0][0] +
					m.els[1][1]*other.els[1][0] +
					m.els[1][2]*other.els[2][0] +
					m.els[1][3]*other.els[3][0],

				m.els[1][0]*other.els[0][1] +
					m.els[1][1]*other.els[1][1] +
					m.els[1][2]*other.els[2][1] +
					m.els[1][3]*other.els[3][1],

				m.els[1][0]*other.els[0][2] +
					m.els[1][1]*other.els[1][2] +
					m.els[1][2]*other.els[2][2] +
					m.els[1][3]*other.els[3][2],

				m.els[1][0]*other.els[0][3] +
					m.els[1][1]*other.els[1][3] +
					m.els[1][2]*other.els[2][3] +
					m.els[1][3]*other.els[3][3],
			},
			[4]float64{
				m.els[2][0]*other.els[0][0] +
					m.els[2][1]*other.els[1][0] +
					m.els[2][2]*other.els[2][0] +
					m.els[2][3]*other.els[3][0],

				m.els[2][0]*other.els[0][1] +
					m.els[2][1]*other.els[1][1] +
					m.els[2][2]*other.els[2][1] +
					m.els[2][3]*other.els[3][1],

				m.els[2][0]*other.els[0][2] +
					m.els[2][1]*other.els[1][2] +
					m.els[2][2]*other.els[2][2] +
					m.els[2][3]*other.els[3][2],

				m.els[2][0]*other.els[0][3] +
					m.els[2][1]*other.els[1][3] +
					m.els[2][2]*other.els[2][3] +
					m.els[2][3]*other.els[3][3],
			},
			[4]float64{
				m.els[3][0]*other.els[0][0] +
					m.els[3][1]*other.els[1][0] +
					m.els[3][2]*other.els[2][0] +
					m.els[3][3]*other.els[3][0],

				m.els[3][0]*other.els[0][1] +
					m.els[3][1]*other.els[1][1] +
					m.els[3][2]*other.els[2][1] +
					m.els[3][3]*other.els[3][1],

				m.els[3][0]*other.els[0][2] +
					m.els[3][1]*other.els[1][2] +
					m.els[3][2]*other.els[2][2] +
					m.els[3][3]*other.els[3][2],

				m.els[3][0]*other.els[0][3] +
					m.els[3][1]*other.els[1][3] +
					m.els[3][2]*other.els[2][3] +
					m.els[3][3]*other.els[3][3],
			},
		},
	}
}

func (m *Matrix4x4) Transpose() *Matrix4x4 {

	// // slower:
	// mat := &Matrix4x4{}
	// for i := 0; i < 4; i++ {
	// 	row := m.GetColumn(i)
	// 	for ind, val := range row {
	// 		mat.els[i][ind] = val
	// 	}
	// }
	// return mat

	// faster:
	return &Matrix4x4{
		[4][4]float64{
			[4]float64{m.els[0][0], m.els[1][0], m.els[2][0], m.els[3][0]},
			[4]float64{m.els[0][1], m.els[1][1], m.els[2][1], m.els[3][1]},
			[4]float64{m.els[0][2], m.els[1][2], m.els[2][2], m.els[3][2]},
			[4]float64{m.els[0][3], m.els[1][3], m.els[2][3], m.els[3][3]},
		}}
}

func (m *Matrix4x4) GetColumn(index int) [4]float64 {
	return [4]float64{m.els[0][index], m.els[1][index], m.els[2][index],
		m.els[3][index]}
}

func (m *Matrix4x4) Inverse() (*Matrix4x4, error) {
	indxc, indxr, ipiv := [4]int{}, [4]int{}, [4]int{}
	minv := m.els

	for i := 0; i < 4; i++ {
		irow, icol := -1, -1
		big := 0.0
		for j := 0; j < 4; j++ {
			if ipiv[j] == 1 {
				continue
			}
			for k := 0; k < 4; k++ {
				if ipiv[k] == 0 {
					abs := math.Abs(float64(minv[j][k]))
					if abs >= big {
						big = abs
						irow = j
						icol = k
					}
				} else if ipiv[k] > 1 {
					return nil, fmt.Errorf("Singular matrix in Invert")
				}
			}
		}
		ipiv[icol]++

		if irow != icol {
			for k := 0; k < 4; k++ {
				minv[irow][k], minv[icol][k] = minv[icol][k], minv[irow][k]
			}
		}
		indxr[i] = irow
		indxc[i] = icol
		if minv[icol][icol] == 0.0 {
			return nil, fmt.Errorf("Singular matrix in Invert")
		}

		pivinv := 1.0 / minv[icol][icol]
		minv[icol][icol] = 1.0
		for j := 0; j < 4; j++ {
			minv[icol][j] *= pivinv
		}

		for j := 0; j < 4; j++ {
			if j == icol {
				continue
			}
			save := minv[j][icol]
			minv[j][icol] = 0
			for k := 0; k < 4; k++ {
				minv[j][k] -= minv[icol][k] * save
			}
		}

	}

	for j := 3; j >= 0; j-- {
		if indxr[j] == indxc[j] {
			continue
		}
		for k := 0; k < 4; k++ {
			minv[k][indxr[j]], minv[k][indxc[j]] = minv[k][indxc[j]], minv[k][indxr[j]]
		}
	}

	inverted := &Matrix4x4{}
	inverted.els = minv

	return inverted, nil
}

func (m *Matrix4x4) Equals(other *Matrix4x4) bool {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if math.Abs(float64(m.els[i][j]-other.els[i][j])) > COMPARE_PRECISION {
				return false
			}
		}
	}

	return true
}

func (m *Matrix4x4) String() string {

	out := "[\n"

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			out += fmt.Sprintf("%f ", m.els[i][j])
		}
		out += "\n"
	}

	out += "]"

	return out

	return fmt.Sprintf("[%s %s %s %s]", m.els[0], m.els[1], m.els[2],
		m.els[3])
}

func NewMatrix(
	a00, a10, a20, a30,
	a01, a11, a21, a31,
	a02, a12, a22, a32,
	a03, a13, a23, a33 float64) *Matrix4x4 {

	return &Matrix4x4{
		[4][4]float64{
			[4]float64{a00, a10, a20, a30},
			[4]float64{a01, a11, a21, a31},
			[4]float64{a02, a12, a22, a32},
			[4]float64{a03, a13, a23, a33},
		}}
}
