package geom_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/domain/geom"
)

// near is how closely computed points must match, mm.
const near = 1e-9

var _ = Describe("Part takes the stretch of a segment between two parameters", func() {
	It("cuts a straight segment at the points in between", func() {
		// Act
		start, seg := geom.Part(geom.Point{}, geom.LineTo(geom.Point{X: 10}), 0.2, 0.5)

		// Assert
		Expect(start.X).To(BeNumerically("~", 2, near))
		Expect(seg.Curved).To(BeFalse())
		Expect(seg.To.X).To(BeNumerically("~", 5, near))
	})

	It("cuts a curve into a curve that follows the same line exactly", func() {
		// Arrange
		p0 := geom.Point{X: 0, Y: 0}
		c := geom.CubicTo(geom.Point{X: 0, Y: 10}, geom.Point{X: 10, Y: 10}, geom.Point{X: 10, Y: 0})

		// Act
		start, seg := geom.Part(p0, c, 0.25, 0.75)

		// Assert: at every u, the part is where the curve is at 0.25 + u/2.
		Expect(seg.Curved).To(BeTrue())
		for _, u := range []float64{0, 0.1, 0.5, 0.9, 1} {
			got := cubicAt(start, seg.C1, seg.C2, seg.To, u)
			want := cubicAt(p0, c.C1, c.C2, c.To, 0.25+u/2)
			Expect(got.X).To(BeNumerically("~", want.X, near))
			Expect(got.Y).To(BeNumerically("~", want.Y, near))
		}
	})

	It("gives the whole curve back for the whole range", func() {
		// Arrange
		c := geom.CubicTo(geom.Point{X: 1, Y: 5}, geom.Point{X: 7, Y: -2}, geom.Point{X: 9, Y: 3})

		// Act
		start, seg := geom.Part(geom.Point{X: 2, Y: 2}, c, 0, 1)

		// Assert
		Expect(start).To(Equal(geom.Point{X: 2, Y: 2}))
		Expect(seg).To(Equal(c))
	})
})
