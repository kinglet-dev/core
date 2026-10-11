package geom

// Part returns the stretch of a segment that starts at p0, between the
// parameters t0 and t1 (0 ≤ t0 ≤ t1 ≤ 1): its start point, and a segment of
// the same kind that follows the original exactly. A curve's part comes
// from its blossom (Ramshaw's polar form): the control points of the
// stretch [a, b] are the blossom's values at (a,a,a), (a,a,b), (a,b,b) and
// (b,b,b).
func Part(p0 Point, s Segment, t0, t1 float64) (Point, Segment) {
	if !s.Curved {
		return lerp(p0, s.To, t0), LineTo(lerp(p0, s.To, t1))
	}
	b := func(u, v, w float64) Point { return blossom(p0, s.C1, s.C2, s.To, u, v, w) }
	return b(t0, t0, t0), CubicTo(b(t0, t0, t1), b(t0, t1, t1), b(t1, t1, t1))
}

// blossom evaluates a cubic's blossom: de Casteljau's algorithm with a
// different parameter at each level.
func blossom(p0, p1, p2, p3 Point, u, v, w float64) Point {
	a, b, c := lerp(p0, p1, u), lerp(p1, p2, u), lerp(p2, p3, u)
	a, b = lerp(a, b, v), lerp(b, c, v)
	return lerp(a, b, w)
}

// lerp is the point a fraction t of the way from a to b, written so that
// t = 0 gives a and t = 1 gives b exactly.
func lerp(a, b Point, t float64) Point {
	return Point{X: (1-t)*a.X + t*b.X, Y: (1-t)*a.Y + t*b.Y}
}
