package transform

import (
	"testing"

	"github.com/ironsmile/raytracer/geometry"
)

func TestPointTransoformWihtIdentity(t *testing.T) {
	point := geometry.NewPoint(13.0, 14.3, 2.0)

	transform := NewTransformation(NewMatrix(
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1))

	found := transform.Point(point)

	if !found.Equals(point) {
		t.Errorf("Expected %s but got %s", point, found)
	}
}

func TestPointTransform(t *testing.T) {
	transform := NewTransformation(NewMatrix(
		0, 1, 2, 0,
		1, 3, 3, 0,
		2, 0, 0, 3,
		0, 1, 0, 2))

	point := geometry.NewPoint(1, 2, 4)

	found := transform.Point(point)

	expected := geometry.NewPoint(2.5, 4.75, 1.25)

	if !found.Equals(expected) {
		t.Errorf("Expected %s but got %s", expected, found)
	}
}

func TestVectorTransoformWihtIdentity(t *testing.T) {
	vector := geometry.NewVector(13.0, 14.3, 2.0)

	transform := NewTransformation(NewMatrix(
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1))

	found := transform.Vector(vector)

	if !found.Equals(vector) {
		t.Errorf("Expected %s but got %s", vector, found)
	}
}

func TestVectorTransofrm(t *testing.T) {
	transform := NewTransformation(NewMatrix(
		0, 1, 2, 0,
		1, 3, 3, 0,
		2, 0, 0, 3,
		0, 1, 0, 2))

	vector := geometry.NewVector(1, 2, 4)

	found := transform.Vector(vector)

	expected := geometry.NewVector(10, 19, 2)

	if !found.Equals(expected) {
		t.Errorf("Expected %s but got %s", expected, found)
	}
}

func TestTransformationComposition(t *testing.T) {
	t1 := NewTransformation(NewMatrix(
		0, 1, 2, 0,
		1, 3, 3, 0,
		2, 0, 0, 3,
		0, 1, 0, 2))

	t2 := NewTransformation(NewMatrix(
		1, 3, 4, 0,
		2, 1, 0, 3,
		0, 0, 0, 1,
		1, 1, 2, 1))

	expected := NewTransformation(NewMatrix(
		2, 1, 0, 5,
		7, 6, 4, 12,
		5, 9, 14, 3,
		4, 3, 4, 5))

	found := t1.Multiply(t2)

	if !found.Equals(expected) {
		t.Errorf("Expected %s but got %s", expected, found)
	}

	expected = NewTransformation(NewMatrix(
		11, 10, 11, 12,
		1, 8, 7, 6,
		0, 1, 0, 2,
		5, 5, 5, 8))

	found = t2.Multiply(t1)

	if !found.Equals(expected) {
		t.Errorf("Expected %s but got %s", expected, found)
	}
}

/*
   Benchmarks
*/

func BenchmarkTransformationMultiplication(t *testing.B) {
	one := NewMatrix(
		0, 1, 2, 0,
		1, 3, 0, 1,
		2, 3, 0, 0,
		0, 0, 3, 2)

	other := NewMatrix(
		0, 0, 1, 1,
		1, 0, 1, 1,
		1, 1, 0, 0,
		0, 1, 0, 1)

	t1 := NewTransformation(one)
	t2 := NewTransformation(other)

	for i := 0; i < t.N; i++ {
		t1.Multiply(t2)
	}
}
