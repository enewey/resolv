package resolv

import "testing"

func TestTriangleOnTriangle(t *testing.T) {
	t1 := NewTriangle(5, 5, 6, 6, 4, 6)
	t2 := NewTriangle(5, 5, 5, 6, 6, 5)

	if !t1.IsColliding(t2) {
		t.Errorf("triangle intersecting with triangle failed")
	}

	t3 := NewTriangle(10, 10, 11, 11, 10, 11)
	if t1.IsColliding(t3) {
		t.Errorf("triangle falsely indicated collision with other triangle")
	}
}

func TestTriangleOnLine(t *testing.T) {
	tr := NewTriangle(5, 5, 10, 10, 0, 10)

	line := NewLine(5, 6, 5, 9)
	if !tr.IsColliding(line) {
		t.Errorf("line inside triangle failed")
	}

	line2 := NewLine(0, 0, 5, 9)
	if !tr.IsColliding(line2) {
		t.Errorf("line intersecting triangle edge failed")
	}

	line3 := NewLine(0, 0, 1, 1)
	if tr.IsColliding(line3) {
		t.Errorf("line falsely indicated collision with triangle")
	}
}

func TestTriangleOnRectangle(t *testing.T) {
	tr := NewTriangle(5, 5, 10, 10, 0, 10)

	rect1 := NewRectangle(5, 7, 5, 5)
	if !tr.IsColliding(rect1) {
		t.Errorf("rectangle interecting triangle failed")
	}

	rect2 := NewRectangle(15, 15, 1, 1)
	if tr.IsColliding(rect2) {
		t.Errorf("rectangle falsely indicated collision with triangle")
	}

	rect3 := NewRectangle(2, 4, 6, 3)
	if !tr.IsColliding(rect3) {
		t.Errorf("rectangle did not indicate tip of triangle inside")
	}
}

func TestTriangleOnCircle(t *testing.T) {
	tr := NewTriangle(10, 10, 18, 18, 2, 18)

	p1C := NewCircle(10, 9, 2)
	if !tr.IsColliding(p1C) {
		t.Errorf("circle encompassing point 1 failed")
	}

	p2C := NewCircle(18, 17, 2)
	if !tr.IsColliding(p2C) {
		t.Errorf("circle encompassing point 2 failed")
	}

	p3C := NewCircle(1, 18, 2)
	if !tr.IsColliding(p3C) {
		t.Errorf("circle encompassing point 3 failed")
	}

	innerC := NewCircle(10, 14, 1)
	if !tr.IsColliding(innerC) {
		t.Errorf("circle inside triangle failed")
	}

	trCCW := NewTriangle(10, 10, 2, 18, 18, 18)
	if !trCCW.IsColliding(innerC) {
		t.Errorf("circle inside counter-clockwise triangle failed")
	}

	e1C := NewCircle(15, 14, 2)
	if !tr.IsColliding(e1C) {
		t.Errorf("circle intersecting first edge of triangle failed")
	}

	e2C := NewCircle(10, 19, 2)
	if !tr.IsColliding(e2C) {
		t.Errorf("circle intersecting second edge of triangle failed")
	}

	e3C := NewCircle(5, 14, 2)
	if !tr.IsColliding(e3C) {
		t.Errorf("circle intersecting third edge of triangle failed")
	}

	badC := NewCircle(30, 30, 2)
	if tr.IsColliding(badC) {
		t.Errorf("circle nowhere near triangle falsely indicated collision")
	}
}
