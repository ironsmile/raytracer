package transform

type Transform struct {
	mat    *Matrix4x4
	matInv *Matrix4x4
}

func (t *Transform) Inverse() *Transform {
	return NewTransformationWihtInverse(t.matInv, t.mat)
}

func (t *Transform) IsIdentity() bool {
	return *t.mat == *NewMatrix(
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1)
}

func NewTransformation(mat *Matrix4x4) *Transform {
	inv, _ := mat.Inverse()
	return &Transform{mat, inv}
}

func NewTransformationWihtInverse(mat *Matrix4x4, inv *Matrix4x4) *Transform {
	return &Transform{mat, inv}
}
