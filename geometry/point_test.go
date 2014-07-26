package geometry

import (
	"testing"
)

func TestPointEqualsMethod(t *testing.T) {
	p1 := NewPoint(1.00003, 2.000004, 3.000002)
	p2 := NewPoint(1.00003, 2.000004, 3.000002)
	p3 := NewPoint(1.000030001, 2.000004, 3.000002)
	p4 := NewPoint(1.000031, 2.000004, 3.000002)
	p5 := NewPoint(1.00003, 2.000004, 3.00000201)
	p6 := NewPoint(1.00003, 2.000004, 3.000001)

	equal := p1.Equals(p2)

	if !equal {
		t.Errorf("Equals should have returned true but it returned false")
	}

	equal = p1.Equals(p3)

	if !equal {
		t.Errorf("Equals should have returned true but it returned false")
	}

	equal = p1.Equals(p5)

	if !equal {
		t.Errorf("Equals should have returned true but it returned false")
	}

	equal = p1.Equals(p4)

	if equal {
		t.Errorf("Equals should have returned false but it returned true")
	}

	equal = p1.Equals(p6)

	if equal {
		t.Errorf("Equals should have returned false but it returned true")
	}

}
