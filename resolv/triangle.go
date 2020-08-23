package resolv

import (
	"math"
)

// Triangle represents a triangle.
type Triangle struct {
	BasicShape
	X2, Y2, X3, Y3 int32
}

// NewTriangle returns a pointer to a new Triangle struct
func NewTriangle(x, y, x2, y2, x3, y3 int32) *Triangle {
	//first ensure the points are in clockwise order
	e1 := (x2 - x) * (y2 + y)
	e2 := (x3 - x2) * (y3 + y2)
	e3 := (x - x3) * (y + y3)

	if (e1 + e2 + e3) >= 0 {
		return &Triangle{BasicShape{X: x, Y: y}, x2, y2, x3, y3}
	}

	return &Triangle{BasicShape{X: x, Y: y}, x3, y3, x2, y2}
}

// heron's formula - the area of a triangle made of these three points
func herons(x1, y1, x2, y2, x3, y3 int32) float64 {
	return math.Abs(float64((x2-x1)*(y3-y1) - (x3-x1)*(y2-y1)))
}

// pointCollides determines if a single point would collide with this triangle.
// This is done by comparing the area made of the three triangles that are formed by using the provided point, and then
// two of the other Triangle points. These three triangles are then added together. If the sum of the area of those
// three triangles are the same as the area of the Triangle itself, then the point is contained inside the triangle.
// See http://jeffreythompson.org/collision-detection/tri-point.php for a better explanation.
// However, in order for this to be consistent with the rest of the resolv library, we need to NOT detect a collision
// in the case where a point is exactly on one of the triangle's edges.
func (t *Triangle) pointCollides(x, y int32) bool {
	a := herons(x, y, t.X2, t.Y2, t.X3, t.Y3)
	b := herons(t.X, t.Y, x, y, t.X3, t.Y3)
	c := herons(t.X, t.Y, t.X2, t.Y2, x, y)

	self := t.Area()

	e1, e2, e3 := t.pointOnEdge(x, y)

	return a+b+c == self && !(e1 || e2 || e3)
}

func (t *Triangle) pointOnEdge(x, y int32) (bool, bool, bool) {
	return PointOnLine(x, y, t.X, t.Y, t.X2, t.Y2),
		PointOnLine(x, y, t.X2, t.Y2, t.X3, t.Y3),
		PointOnLine(x, y, t.X3, t.Y3, t.X, t.Y)
}

func (t *Triangle) segmentCollides(px, py, qx, qy int32) bool {
	return LinesIntersect(t.X, t.Y, t.X2, t.Y2, px, py, qx, qy) ||
		LinesIntersect(t.X2, t.Y2, t.X3, t.Y3, px, py, qx, qy) ||
		LinesIntersect(t.X3, t.Y3, t.X, t.Y, px, py, qx, qy)
}

// Area gets the area of the Triangle
func (t *Triangle) Area() float64 {
	return herons(t.X, t.Y, t.X2, t.Y2, t.X3, t.Y3)
}

// SetXY sets the position of base X/Y coord of the Triangle, and both other XY points are set relative to the base XY.
func (t *Triangle) SetXY(x, y int32) {
	dx := x - t.X
	dy := y - t.Y
	t.X, t.Y = x, y
	t.X2 += dx
	t.Y2 += dy
	t.X3 += dx
	t.Y3 += dy
}

// IsColliding returns whether the Triangle is colliding with the specified other Shape or not, including the other Shape
// being wholly contained within the Rectangle.
func (t *Triangle) IsColliding(other Shape) bool {

	switch b := other.(type) {
	case *Triangle:
		return t.segmentCollides(b.X, b.Y, b.X2, b.Y2) || t.segmentCollides(b.X2, b.Y2, b.X3, b.Y3) ||
			t.segmentCollides(b.X3, b.Y3, b.X, b.Y) || t.pointCollides(b.X, b.Y) || t.pointCollides(b.X2, b.Y2) ||
			t.pointCollides(b.X3, b.Y3)
	case *Rectangle:
		rectCollides := b.pointCollides(t.X, t.Y) || b.pointCollides(t.X2, t.Y2) || b.pointCollides(t.X3, t.Y3)
		triCollides := t.pointCollides(b.X, b.Y) || t.pointCollides(b.X+b.W, b.Y) ||
			t.pointCollides(b.X, b.Y+b.H) || t.pointCollides(b.X+b.W, b.Y+b.H)
		segmentsCollide := t.segmentCollides(b.X, b.Y, b.X+b.W, b.Y) || t.segmentCollides(b.X+b.W, b.Y, b.X+b.W, b.Y+b.H) ||
			t.segmentCollides(b.X+b.W, b.Y+b.H, b.X, b.Y+b.H) || t.segmentCollides(b.X, b.Y+b.H, b.X, b.Y)

		return rectCollides || triCollides || segmentsCollide
	case *Line:
		return t.segmentCollides(b.X, b.Y, b.X2, b.Y2) || t.pointCollides(b.X, b.Y) || t.pointCollides(b.X2, b.Y2)
	case *Circle:
		// tests if a circle intersects a triangle.
		// see http://www.phatcode.net/articles.php?id=459

		// test 1 - circle encompasses one of the three triangle points

		c1x := b.X - t.X
		c1y := b.Y - t.Y

		radiusSqr := b.Radius * b.Radius
		c1sqr := (c1x * c1x) + (c1y * c1y) - radiusSqr

		if c1sqr <= 0 {
			return true
		}

		c2x := b.X - t.X2
		c2y := b.Y - t.Y2

		c2sqr := (c2x * c2x) + (c2y * c2y) - radiusSqr

		if c2sqr <= 0 {
			return true
		}

		c3x := b.X - t.X3
		c3y := b.Y - t.Y3

		c3sqr := (c3x * c3x) + (c3y * c3y) - radiusSqr

		if c3sqr <= 0 {
			return true
		}

		// test 2 - circle fully inside of triangle

		// this operation specifically relies on the points being ordered clockwise.
		e1x := t.X2 - t.X
		e1y := t.Y2 - t.Y

		e2x := t.X3 - t.X2
		e2y := t.Y3 - t.Y2

		e3x := t.X - t.X3
		e3y := t.Y - t.Y3

		k := int32((e1y*c1x - e1x*c1y) | (e2y*c2x - e2x*c2y) | (e3y*c3x - e3x*c3y))
		if k >= 0 {
			return true
		}

		// test 3 - circle intersects triangle edge without encompassing a triangle point

		k = c1x*e1x + c1y*e1y

		if k > 0 {
			lng := e1x*e1x + e1y*e1y
			if k < lng {
				if c1sqr*lng <= k*k {
					return true
				}

			}
		}

		k = c2x*e2x + c2y*e2y

		if k > 0 {
			lng := e2x*e2x + e2y*e2y
			if k < lng {
				if c2sqr*lng <= k*k {
					return true
				}

			}
		}

		k = c3x*e3x + c3y*e3y

		if k > 0 {
			lng := e3x*e3x + e3y*e3y
			if k < lng {
				if c3sqr*lng <= k*k {
					return true
				}

			}
		}

		// no collision to report

		return false
	default:
		return b.IsColliding(t)
	}

}

// WouldBeColliding returns whether the Triangle would be colliding with the other Shape if it were to move in the
// specified direction.
func (t *Triangle) WouldBeColliding(other Shape, dx, dy int32) bool {
	t.SetXY(t.X+dx, t.Y+dy)
	isColliding := t.IsColliding(other)
	t.SetXY(t.X-dx, t.Y-dy)
	return isColliding
}
