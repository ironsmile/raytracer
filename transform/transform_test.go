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

}
