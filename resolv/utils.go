package resolv

import (
	"math"
)

// Resolve attempts to move the checking Shape with the specified X and Y values, returning a Collision object
// if it collides with the specified other Shape. The deltaX and deltaY arguments are the movement displacement
// in pixels. For platformers in particular, you would probably want to resolve on the X and Y axes separately.
func Resolve(firstShape Shape, other Shape, deltaX, deltaY int32) Collision {

	out := Collision{}
	out.ResolveX = deltaX
	out.ResolveY = deltaY
	out.ShapeA = firstShape

	if deltaX == 0 && deltaY == 0 {
		return out
	}

	x := float32(deltaX)
	y := float32(deltaY)

	primeX := true
	slope := float32(0)

	if math.Abs(float64(deltaY)) > math.Abs(float64(deltaX)) {
		primeX = false
		if deltaY != 0 && deltaX != 0 {
			slope = float32(deltaX) / float32(deltaY)
		}
	} else if deltaY != 0 && deltaX != 0 {
		slope = float32(deltaY) / float32(deltaX)
	}

	for true {

		if firstShape.WouldBeColliding(other, out.ResolveX, out.ResolveY) {

			if primeX {

				if deltaX > 0 {
					x--
				} else if deltaX < 0 {
					x++
				}

				if deltaY > 0 {
					y -= slope
				} else if deltaY < 0 {
					y += slope
				}

			} else {

				if deltaY > 0 {
					y--
				} else if deltaY < 0 {
					y++
				}

				if deltaX > 0 {
					x -= slope
				} else if deltaX < 0 {
					x += slope
				}

			}

			out.ResolveX = int32(x)
			out.ResolveY = int32(y)
			out.ShapeB = other

		} else {
			break
		}

	}

	if math.Abs(float64(deltaX-out.ResolveX)) > math.Abs(float64(deltaX)*1.5) || math.Abs(float64(deltaY-out.ResolveY)) > math.Abs(float64(deltaY)*1.5) {
		out.Teleporting = true
	}

	return out

}

// Distance returns the distance from one pair of X and Y values to another.
func Distance(x, y, x2, y2 int32) int32 {

	dx := x - x2
	dy := y - y2
	ds := (dx * dx) + (dy * dy)
	return int32(math.Sqrt(math.Abs(float64(ds))))

}

// ColinearPointsOnSegment - Given three colinear points a,b,c, checks if point c lies on the segment ab
func ColinearPointsOnSegment(ax, ay, bx, by, cx, cy int32) bool {
	if cx < max(ax, bx) && cx > min(ax, bx) && cy < max(ay, by) && cy > min(ay, by) {
		return true
	}

	return false
}

// To find orientation of ordered triplet (p, q, r).
// The function returns following values
// 0 --> p, q and r are colinear
// 1 --> Clockwise
// 2 --> Counterclockwise
// int orientation(Point p, Point q, Point r)
// {
//     // See https://www.geeksforgeeks.org/orientation-3-ordered-points/
//     // for details of below formula.
//     int val = (q.y - p.y) * (r.x - q.x) -
//               (q.x - p.x) * (r.y - q.y);

//     if (val == 0) return 0;  // colinear

//     return (val > 0)? 1: 2; // clock or counterclock wise
// }

const (
	colinear = iota
	clockwise
	counterclockwise
)

func orientation(px, py, qx, qy, rx, ry int32) int {
	val := (qy-py)*(rx-qx) - (qx-px)*(ry-qy)
	if val == 0 {
		return colinear
	}
	if val > 0 {
		return clockwise
	}
	return counterclockwise
}

// LinesIntersect - returns true if line segment 'p1q1' and 'p2q2' intersect.
func LinesIntersect(p1x, p1y, q1x, q1y, p2x, p2y, q2x, q2y int32) bool {
	o1, o2, o3, o4 := orientation(p1x, p1y, q1x, q1y, p2x, p2y),
		orientation(p1x, p1y, q1x, q1y, q2x, q2y),
		orientation(p2x, p2y, q2x, q2y, p1x, p1y),
		orientation(p2x, p2y, q2x, q2y, q1x, q1y)

	if o1 != o2 && o3 != o4 {
		// For our purposes, we don't want to consider it an intersection if the point is *on* the line

		if PointOnLine(p1x, p1y, p2x, p2y, q2x, q2y) ||
			PointOnLine(q1x, q1y, p2x, p2y, q2x, q2y) ||
			PointOnLine(p2x, p2y, p1x, p1y, q1x, q1y) ||
			PointOnLine(q2x, q2y, p1x, p1y, q1x, q1y) {
			return false
		}
		// fmt.Printf("points not on line in line intersects %d,%d - %d,%d - %d,%d - %d,%d .. ", p1x, p1y, q1x, q1y, p2x, p2y, q2x, q2y)
		return true
	}

	if o1 == colinear && ColinearPointsOnSegment(p1x, p1y, q1x, q1y, p2x, p2y) {
		return true
	}
	if o2 == colinear && ColinearPointsOnSegment(p1x, p1y, q1x, q1y, q2x, q2y) {
		return true
	}
	if o3 == colinear && ColinearPointsOnSegment(p2x, p2y, q2x, q2y, p1x, p1y) {
		return true
	}
	if o4 == colinear && ColinearPointsOnSegment(p2x, p2y, q2x, q2y, q1x, q1y) {
		return true
	}

	return false
}

// PointOnLine returns true if the point x,y is on the line formed by points a and b
func PointOnLine(x, y, ax, ay, bx, by int32) bool {
	seg1 := fDistance(x, y, ax, ay)
	seg2 := fDistance(x, y, bx, by)
	line := fDistance(ax, ay, bx, by)
	// accuracy to the nearest ten-thousandth should be good enough for int32s!
	return (seg1+seg2)-line <= 0.0001
}

// fDistance returns the distance from one pair of X and Y values to another as a float
func fDistance(x, y, x2, y2 int32) float64 {
	dx := x - x2
	dy := y - y2
	ds := (dx * dx) + (dy * dy)
	return math.Sqrt(math.Abs(float64(ds)))
}

func max(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
