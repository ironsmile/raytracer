package transform

import (
	"fmt"
)

type Matrix4x4 struct {
	elements [4][4]float32
}

func (m *Matrix4x4) Get(i, j int) float32 {
	return m.elements[i][j]
}

func (m *Matrix4x4) Multiply(other *Matrix4x4) *Matrix4x4 {
	// // slower:
	// mat := &Matrix4x4{}

	// for i := 0; i < 4; i++ {
	// 	for j := 0; j < 4; j++ {
	// 		mat.elements[i][j] = m.elements[i][0]*other.elements[0][j] +
	// 			m.elements[i][1]*other.elements[1][j] +
	// 			m.elements[i][2]*other.elements[2][j] +
	// 			m.elements[i][3]*other.elements[3][j]
	// 	}
	// }
	// return mat

	// faster:
	return &Matrix4x4{
		[4][4]float32{
			[4]float32{
				m.elements[0][0]*other.elements[0][0] +
					m.elements[0][1]*other.elements[1][0] +
					m.elements[0][2]*other.elements[2][0] +
					m.elements[0][3]*other.elements[3][0],

				m.elements[0][0]*other.elements[0][1] +
					m.elements[0][1]*other.elements[1][1] +
					m.elements[0][2]*other.elements[2][1] +
					m.elements[0][3]*other.elements[3][1],

				m.elements[0][0]*other.elements[0][2] +
					m.elements[0][1]*other.elements[1][2] +
					m.elements[0][2]*other.elements[2][2] +
					m.elements[0][3]*other.elements[3][2],

				m.elements[0][0]*other.elements[0][3] +
					m.elements[0][1]*other.elements[1][3] +
					m.elements[0][2]*other.elements[2][3] +
					m.elements[0][3]*other.elements[3][3],
			},
			[4]float32{
				m.elements[1][0]*other.elements[0][0] +
					m.elements[1][1]*other.elements[1][0] +
					m.elements[1][2]*other.elements[2][0] +
					m.elements[1][3]*other.elements[3][0],

				m.elements[1][0]*other.elements[0][1] +
					m.elements[1][1]*other.elements[1][1] +
					m.elements[1][2]*other.elements[2][1] +
					m.elements[1][3]*other.elements[3][1],

				m.elements[1][0]*other.elements[0][2] +
					m.elements[1][1]*other.elements[1][2] +
					m.elements[1][2]*other.elements[2][2] +
					m.elements[1][3]*other.elements[3][2],

				m.elements[1][0]*other.elements[0][3] +
					m.elements[1][1]*other.elements[1][3] +
					m.elements[1][2]*other.elements[2][3] +
					m.elements[1][3]*other.elements[3][3],
			},
			[4]float32{
				m.elements[2][0]*other.elements[0][0] +
					m.elements[2][1]*other.elements[1][0] +
					m.elements[2][2]*other.elements[2][0] +
					m.elements[2][3]*other.elements[3][0],

				m.elements[2][0]*other.elements[0][1] +
					m.elements[2][1]*other.elements[1][1] +
					m.elements[2][2]*other.elements[2][1] +
					m.elements[2][3]*other.elements[3][1],

				m.elements[2][0]*other.elements[0][2] +
					m.elements[2][1]*other.elements[1][2] +
					m.elements[2][2]*other.elements[2][2] +
					m.elements[2][3]*other.elements[3][2],

				m.elements[2][0]*other.elements[0][3] +
					m.elements[2][1]*other.elements[1][3] +
					m.elements[2][2]*other.elements[2][3] +
					m.elements[2][3]*other.elements[3][3],
			},
			[4]float32{
				m.elements[3][0]*other.elements[0][0] +
					m.elements[3][1]*other.elements[1][0] +
					m.elements[3][2]*other.elements[2][0] +
					m.elements[3][3]*other.elements[3][0],

				m.elements[3][0]*other.elements[0][1] +
					m.elements[3][1]*other.elements[1][1] +
					m.elements[3][2]*other.elements[2][1] +
					m.elements[3][3]*other.elements[3][1],

				m.elements[3][0]*other.elements[0][2] +
					m.elements[3][1]*other.elements[1][2] +
					m.elements[3][2]*other.elements[2][2] +
					m.elements[3][3]*other.elements[3][2],

				m.elements[3][0]*other.elements[0][3] +
					m.elements[3][1]*other.elements[1][3] +
					m.elements[3][2]*other.elements[2][3] +
					m.elements[3][3]*other.elements[3][3],
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
	// 		mat.elements[i][ind] = val
	// 	}
	// }
	// return mat

	// faster:
	return &Matrix4x4{
		[4][4]float32{
			[4]float32{m.elements[0][0], m.elements[1][0], m.elements[2][0], m.elements[3][0]},
			[4]float32{m.elements[0][1], m.elements[1][1], m.elements[2][1], m.elements[3][1]},
			[4]float32{m.elements[0][2], m.elements[1][2], m.elements[2][2], m.elements[3][2]},
			[4]float32{m.elements[0][3], m.elements[1][3], m.elements[2][3], m.elements[3][3]},
		}}
}

func (m *Matrix4x4) GetColumn(index int) [4]float32 {
	return [4]float32{m.elements[0][index], m.elements[1][index], m.elements[2][index],
		m.elements[3][index]}
}

func (m *Matrix4x4) Inverse() (*Matrix4x4, error) {
	return &Matrix4x4{}, nil
}

func (m *Matrix4x4) String() string {

	out := "[\n"

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			out += fmt.Sprintf("%f ", m.elements[i][j])
		}
		out += "\n"
	}

	out += "]"

	return out

	return fmt.Sprintf("[%s %s %s %s]", m.elements[0], m.elements[1], m.elements[2],
		m.elements[3])
}

func NewMatrix(
	a00, a10, a20, a30,
	a01, a11, a21, a31,
	a02, a12, a22, a32,
	a03, a13, a23, a33 float32) *Matrix4x4 {

	return &Matrix4x4{
		[4][4]float32{
			[4]float32{a00, a10, a20, a30},
			[4]float32{a01, a11, a21, a31},
			[4]float32{a02, a12, a22, a32},
			[4]float32{a03, a13, a23, a33},
		}}
}
