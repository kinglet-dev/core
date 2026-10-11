// Package geom holds the plane geometry Kinglet's tools measure: points, affine
// transforms, paths and polylines. Units are whatever the caller uses; the
// reader converts everything to millimeters.
package geom

// Point is a position in the plane.
type Point struct{ X, Y float64 }

// Matrix is an affine transform in SVG order: x' = A·x + C·y + E, y' = B·x + D·y + F.
type Matrix struct{ A, B, C, D, E, F float64 }

// Apply transforms p.
func (m Matrix) Apply(p Point) Point {
	return Point{X: m.A*p.X + m.C*p.Y + m.E, Y: m.B*p.X + m.D*p.Y + m.F}
}

// Identity returns the transform that changes nothing.
func Identity() Matrix { return Matrix{A: 1, D: 1} }

// Then returns the transform that applies m first and n second.
func (m Matrix) Then(n Matrix) Matrix {
	return Matrix{
		A: n.A*m.A + n.C*m.B,
		B: n.B*m.A + n.D*m.B,
		C: n.A*m.C + n.C*m.D,
		D: n.B*m.C + n.D*m.D,
		E: n.A*m.E + n.C*m.F + n.E,
		F: n.B*m.E + n.D*m.F + n.F,
	}
}

// Inverse returns the transform that undoes m, and false when there is
// none: when m flattens the plane onto a line or a point.
func (m Matrix) Inverse() (Matrix, bool) {
	det := m.A*m.D - m.B*m.C
	if det == 0 {
		return Matrix{}, false
	}
	return Matrix{
		A: m.D / det, B: -m.B / det, C: -m.C / det, D: m.A / det,
		E: (m.C*m.F - m.D*m.E) / det, F: (m.B*m.E - m.A*m.F) / det,
	}, true
}
