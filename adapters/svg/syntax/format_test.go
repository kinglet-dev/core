package syntax_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/adapters/svg/syntax"
	"github.com/kinglet-dev/core/domain/geom"
)

var _ = Describe("FormatPath writes path data", func() {
	It("writes lines, curves and closed subpaths with absolute commands", func() {
		// Arrange
		p := geom.Path{Subpaths: []geom.Subpath{
			{Start: geom.Point{X: 1, Y: 2}, Segments: []geom.Segment{geom.LineTo(geom.Point{X: 3.5, Y: 2})}, Closed: false},
			{Start: geom.Point{X: 0, Y: 0}, Segments: []geom.Segment{
				geom.CubicTo(geom.Point{X: 0, Y: 1}, geom.Point{X: 1, Y: 2}, geom.Point{X: 2, Y: -0.25}),
			}, Closed: true},
		}}

		// Act
		d := syntax.FormatPath(p)

		// Assert
		Expect(d).To(Equal("M1 2 L3.5 2 M0 0 C0 1 1 2 2 -0.25 Z"))
	})

	It("reads back as the same path", func() {
		// Arrange
		p, err := syntax.ParsePath("M10 10 h20 v5 a5 5 0 0 1 -5 5 q-3 0 -6 -2 z M40 40 L41 41")
		Expect(err).NotTo(HaveOccurred())

		// Act
		back, err := syntax.ParsePath(syntax.FormatPath(p))

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(back.Subpaths).To(HaveLen(len(p.Subpaths)))
		for i, sub := range back.Subpaths {
			Expect(sub.Closed).To(Equal(p.Subpaths[i].Closed))
			Expect(sub.Segments).To(HaveLen(len(p.Subpaths[i].Segments)))
			for j, s := range sub.Segments {
				want := p.Subpaths[i].Segments[j]
				Expect(s.Curved).To(Equal(want.Curved))
				for k, pt := range []geom.Point{s.C1, s.C2, s.To} {
					w := []geom.Point{want.C1, want.C2, want.To}[k]
					Expect(pt.X).To(BeNumerically("~", w.X, 1e-6))
					Expect(pt.Y).To(BeNumerically("~", w.Y, 1e-6))
				}
			}
		}
	})

	It("rounds to a millionth of a unit and never writes -0", func() {
		// Arrange
		p := geom.Path{Subpaths: []geom.Subpath{{Start: geom.Point{X: 1.23456789, Y: -0.0000001}}}}

		// Act
		d := syntax.FormatPath(p)

		// Assert
		Expect(d).To(Equal("M1.234568 0"))
	})
})
