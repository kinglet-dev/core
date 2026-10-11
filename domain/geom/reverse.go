package geom

// Reverse returns the subpath drawn the other way: it starts where s ends
// and visits the same points backward; each curve swaps its control points.
func (s Subpath) Reverse() Subpath {
	n := len(s.Segments)
	out := Subpath{Segments: make([]Segment, n), Closed: s.Closed, Start: s.Start}
	if n > 0 {
		out.Start = s.Segments[n-1].To
	}
	for i, seg := range s.Segments {
		from := s.Start
		if i > 0 {
			from = s.Segments[i-1].To
		}
		out.Segments[n-1-i] = Segment{C1: seg.C2, C2: seg.C1, To: from, Curved: seg.Curved}
	}
	return out
}
